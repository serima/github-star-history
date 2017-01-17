package main

import "fmt"

import "github.com/google/go-github/github"

const (
	owner = "kevincobain2000"
	repo  = "ionic-photo-browser"
)

func main() {
	summary := map[string]int{}
	client := github.NewClient(nil)
	stargazers, resp, err := client.Activity.ListStargazers(owner, repo, nil)
	if err != nil {
		panic(err)
	}
	for _, stargazer := range stargazers {
		printStargazers(stargazer)
		tallyStargazers(stargazer, summary)
	}

	p := resp.NextPage
	for p != 0 {
		opt := &github.ListOptions{Page: p}
		stargazers, resp, err := client.Activity.ListStargazers(owner, repo, opt)
		if err != nil {
			panic(err)
		}
		for _, stargazer := range stargazers {
			printStargazers(stargazer)
			tallyStargazers(stargazer, summary)
		}
		p = resp.NextPage
	}
	fmt.Println(summary)
}

func printStargazers(stargazer *github.Stargazer) {
	user := stargazer.User
	fmt.Printf("starred_at:%v\tuser_login:%v\n", stargazer.StarredAt, *user.Login)
}

func tallyStargazers(stargazer *github.Stargazer, summary map[string]int) {
	key := stargazer.StarredAt.String()[0:7]
	summary[key]++
}
