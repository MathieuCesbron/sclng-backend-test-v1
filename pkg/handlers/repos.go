package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/google/go-github/v55/github"
)

// Repo only keep useful fields from a github repo.
type Repo struct {
	FullName   string `json:"full_name"`
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
	// LanguagesURL string         `json:"languages_url"`
	Languages map[string]int `json:"languages"`
}

// reposHandlerConfig is the config for reposHandler
type reposHandlerConfig struct {
	ctx      context.Context
	log      *log.Logger
	ghClient *github.Client
}

// NewReposHandler returns the repos handler function with his config.
func NewReposHandler(log *log.Logger) func(http.ResponseWriter, *http.Request) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		panic("GITHUB_TOKEN env variable should be set.")
	}

	ghClient := github.NewClient(nil).WithAuthToken(token)
	c := reposHandlerConfig{ctx: context.Background(), log: log, ghClient: ghClient}

	return c.reposHandler
}

// reposHandler returns the 100 latest updated public github repos.
func (rhc reposHandlerConfig) reposHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// https://api.github.com/search/repositories?q=Q&sort=updated
	results, _, err := rhc.ghClient.Search.Repositories(
		context.TODO(),
		"is:public",
		&github.SearchOptions{
			Sort:        "updated",
			ListOptions: github.ListOptions{Page: 1, PerPage: 100},
		},
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed accessing recent github repos: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	repos := []*Repo{}
	for _, ghRepo := range results.Repositories {
		if ghRepo == nil || ghRepo.FullName == nil || ghRepo.Owner == nil ||
			ghRepo.Owner.Login == nil || ghRepo.Name == nil || ghRepo.LanguagesURL == nil {
			http.Error(w, "got github repo with wrong format", http.StatusInternalServerError)
			return
		}

		repos = append(repos, &Repo{
			FullName:   *ghRepo.FullName,
			Owner:      *ghRepo.Owner.Login,
			Repository: *ghRepo.Name,
		})
	}

	err = rhc.populateLanguages(repos)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed populating languages to repos: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	repos, err = filterRepos(repos, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed populating languages to repos: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed marshalling repos: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed writing to the conection: %s", err.Error()), http.StatusInternalServerError)
	}
}

// populateLanguages add languages for each repo.
// cancel every goroutines if we cannot add languages for a single repo.
func (rhc reposHandlerConfig) populateLanguages(repos []*Repo) (err error) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(rhc.ctx)
	defer cancel()

	for _, repo := range repos {
		wg.Add(1)
		repo := repo
		go func() {
			defer wg.Done()
			var res map[string]int
			res, _, err = rhc.ghClient.Repositories.ListLanguages(ctx, repo.Owner, repo.Repository)
			if err != nil {
				cancel()
			}
			repo.Languages = res
		}()
	}

	wg.Wait()

	return err
}

func filterRepos(repos []*Repo, r *http.Request) ([]*Repo, error) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return repos, err
	}
	languagesToFilter := strings.Split(params.Get("languages"), ",")

	repos = slices.DeleteFunc(repos, func(r *Repo) bool {
		for language := range r.Languages {
			for _, languageToFilter := range languagesToFilter {
				if language == languageToFilter {
					return false
				}
			}
		}
		return true
	})

	return repos, nil
}
