package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/go-github/v55/github"
)

type Repo struct {
	FullName     string         `json:"full_name"`
	Owner        string         `json:"owner"`
	Repository   string         `json:"repository"`
	LanguagesURL string         `json:"languages_url"`
	Languages    map[string]int `json:"languages"`
}

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
		// TODO: make better error here
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	repos := []*Repo{}
	for _, ghRepo := range results.Repositories {
		// TODO: check if not nil
		repos = append(repos, &Repo{
			FullName:     *ghRepo.FullName,
			Owner:        *ghRepo.Owner.Login,
			Repository:   *ghRepo.Name,
			LanguagesURL: *ghRepo.LanguagesURL,
		})
	}

	rhc.populateLanguages(repos)

	jsonData, err := json.Marshal(repos)
	if err != nil {
		// TODO: make better error here
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

// populateLanguages add languages for each repo.
// cancel every goroutines if we cannot add languages for a single repo.
func (rhc reposHandlerConfig) populateLanguages(repos []*Repo) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(rhc.ctx)
	defer cancel()

	for _, repo := range repos {
		wg.Add(1)
		repo := repo
		go func() {
			defer wg.Done()
			res, _, err := rhc.ghClient.Repositories.ListLanguages(ctx, repo.Owner, repo.Repository)
			if err != nil {
				cancel()
			}
			repo.Languages = res
		}()
	}

	wg.Wait()
}
