package handlers

import (
	"context"
	"encoding/json"
	"net/http"
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

// ReposHandler returns the latest public github repos.
func ReposHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// https://api.github.com/search/repositories?q=Q&sort=updated
	client := github.NewClient(nil).WithAuthToken("token")
	results, _, err := client.Search.Repositories(
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

	err = PopulateLanguages(client, repos)
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

func PopulateLanguages(client *github.Client, repos []*Repo) error {
	var wg sync.WaitGroup

	for _, repo := range repos {
		wg.Add(1)
		repo := repo
		go func() {
			defer wg.Done()
			res, _, err := client.Repositories.ListLanguages(context.TODO(), repo.Owner, repo.Repository)
			if err != nil {
				// TODO: make better error handling
			}
			repo.Languages = res
		}()
	}

	wg.Wait()

	return nil
}
