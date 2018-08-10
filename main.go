package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	// We require a github token to make API calls due to rate limiting
	if os.Getenv("GITHUB_TOKEN") == "" {
		fmt.Println("Error: Please supply an environment variable named GITHUB_TOKEN with a valid token")
		os.Exit(1)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	c := github.NewClient(tc)

	// Query for issues with said label
	sc := &github.IssueListByRepoOptions{
		Labels: []string{"wontfix"},
	}
	sc.PerPage = 100

	issues, response, err := c.Issues.ListByRepo(ctx, "helm", "charts", sc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print info about the response
	fmt.Println("Response Status: ", response.Status)

	// Iterate over issues
	for _, issue := range issues {

		// Display existing issue from the query
		fmt.Printf("> Issue: %d\n", issue.ID)
		fmt.Printf(">>> State: %s\n", github.Stringify(issue.State))
		fmt.Printf(">>> Labels: %v\n\n", issue.Labels)

		// Add label we want. Note, all of the labels are changed so we need to
		// copy the existing labels over.
		req := &github.IssueRequest{
			Labels: &[]string{"lifecycle/stale"},
		}

		for _, l := range issue.Labels {
			*req.Labels = append(*req.Labels, *l.Name)
		}

		// Saving change. Being conservative so as not remove any labels including
		// the wontfix one
		ctx2 := context.Background()
		_, _, err := c.Issues.Edit(ctx2, "helm", "charts", *issue.Number, req)
		if err != nil {
			fmt.Printf("There was an error updating %s: %s", github.Stringify(issue.Number), err)
		}
	}
}
