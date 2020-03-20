package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	var accessToken, secret string
	var ok bool

	if accessToken, ok = os.LookupEnv("ACCESS_TOKEN"); !ok {
		log.Fatal("ACCESS_TOKEN not set")
	}

	if secret, ok = os.LookupEnv("SECRET"); !ok {
		log.Fatal("SECRET not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		payload, err := github.ValidatePayload(r, []byte(secret))
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, err)
			return
		}

		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprint(w, err)
			return
		}

		switch event := event.(type) {
		case *github.PushEvent:
			input := &github.RepoStatus{
				State:       github.String("pending"),
				TargetURL:   github.String("https://github.com/kunzese"),
				Description: github.String("My description"),
				Context:     github.String("Example/Golang"),
			}

			_, _, err := githubClient.Repositories.CreateStatus(ctx, event.GetRepo().Owner.GetName(), event.GetRepo().GetName(), event.GetAfter(), input)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
				return
			}
		}

		fmt.Fprintf(w, "%+v\n", event)
	})

	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
