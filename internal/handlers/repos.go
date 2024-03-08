package handlers

import (
	"context"
	"encoding/json"
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
	ghClient *github.Client
}

func NewReposHandler() func(http.ResponseWriter, *http.Request) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		panic("GITHUB_TOKEN env variable should be set.")
	}

	ghClient := github.NewClient(nil).WithAuthToken(token)
	c := reposHandlerConfig{ghClient: ghClient}

	return c.reposHandler
}

// ReposHandler returns the 100 latest updated public github repos.
func (rhc reposHandlerConfig) reposHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// https://api.github.com/search/repositories?q=Q&sort=updated
	results, _, err := rhc.ghClient.Search.Repositories(
		context.Background(),
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

	err = rhc.PopulateLanguages(repos)
	if err != nil {
		// TODO: make better error here
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(repos)
	if err != nil {
		// TODO: make better error here
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func (rhc reposHandlerConfig) PopulateLanguages(repos []*Repo) error {
	var wg sync.WaitGroup

	for _, repo := range repos {
		wg.Add(1)
		repo := repo
		go func() {
			defer wg.Done()
			res, _, err := rhc.ghClient.Repositories.ListLanguages(context.TODO(), repo.Owner, repo.Repository)
			if err != nil {
				// TODO: make better error handling
			}
			repo.Languages = res
		}()
	}

	wg.Wait()

	return nil
}
