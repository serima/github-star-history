package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	slice := strings.Split(os.Args[1], "/")
	owner := slice[0]
	repo := slice[1]

	token := os.Getenv("GITHUB_API_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	ctime, _ := getCreatedAtFromRepo(ctx, client, owner, repo)
	months := iterateMonth(ctime)
	summary := map[string]int{}
	initSummary(summary, months)

	// 1ページ目の読み込み
	stargazers, resp, err := client.Activity.ListStargazers(ctx, owner, repo, nil)
	if err != nil {
		panic(err)
	}
	for _, stargazer := range stargazers {
		//printStargazers(stargazer)
		tallyStargazers(stargazer, summary)
	}

	// 2ページ目以降の読み込み
	p := resp.NextPage
	for p != 0 {
		opt := &github.ListOptions{Page: p}
		stargazers, resp, err := client.Activity.ListStargazers(ctx, owner, repo, opt)
		if err != nil {
			panic(err)
		}
		for _, stargazer := range stargazers {
			//printStargazers(stargazer)
			tallyStargazers(stargazer, summary)
		}
		p = resp.NextPage
	}

	printTallySummary(summary)
}

func initSummary(summary map[string]int, months []string) {
	for _, value := range months {
		summary[value] = 0
	}
}

// get CreatedAt from repo
func getCreatedAtFromRepo(ctx context.Context, client *github.Client, owner string, repo string) (createdAt time.Time, err error) {
	repoinfo, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		fmt.Println(err)
		return
	}
	var shortForm = "2006-01-02 15:04:05 -0700 UTC"
	ctime, _ := time.Parse(shortForm, fmt.Sprintf("%s", repoinfo.CreatedAt))

	return ctime, nil
}

func iterateMonth(ctime time.Time) []string {
	months := []string{}
	now := time.Now()
	for d := ctime; now.After(d); d = d.AddDate(0, 1, 0) {
		months = append(months, fmt.Sprintf("%s", d)[0:7])
	}

	return months
}

func printStargazers(stargazer *github.Stargazer) {
	user := stargazer.User
	fmt.Printf("starred_at:%v\tuser_login:%v\n", stargazer.StarredAt, *user.Login)
}

func printTallySummary(summary map[string]int) {
	for key, value := range summary {
		fmt.Printf("%v\t%v\n", key, value)
	}
}

func tallyStargazers(stargazer *github.Stargazer, summary map[string]int) {
	key := stargazer.StarredAt.String()[0:7]
	summary[key]++
}
