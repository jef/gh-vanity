package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shurcooL/graphql"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
)

type userNode struct {
	Company       graphql.String
	Email         graphql.String
	IsEmployee    graphql.Boolean
	Name          graphql.String
	Login         graphql.String
	Url           graphql.String
	Organizations struct {
		Edges []struct {
			Node struct {
				Name graphql.String
			}
		}
	} `graphql:"organizations(first: 10)"`
}

var query struct {
	Repository struct {
		Stargazers struct {
			TotalCount graphql.Int
			Edges      []struct {
				Node userNode
			}
			PageInfo struct {
				EndCursor   graphql.String
				HasNextPage graphql.Boolean
			}
		} `graphql:"stargazers(first: 100, after: $stargazerCursor)"`
	} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

type arguments struct {
	company         string
	companyFilter   []string
	employee        bool
	repositoryOwner string
	repositoryName  string
}

func getFlags() (*arguments, error) {
	a := new(arguments)

	flag.StringVar(&a.company, "company", "",
		"Filter stargazers by company name(s). Can be comma separated.\n"+
			"If no names are given, then all stargazers will output.")
	flag.BoolVar(&a.employee, "employee", false, "Filter stargazers that are GitHub employees.")
	flag.StringVar(&a.repositoryName, "repo", "", "(Required) The name of the repository.")
	flag.StringVar(&a.repositoryOwner, "owner", "", "(Required) The owner or organization of the repository.")

	flag.Parse()

	required := []string{"owner", "repo"}
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			return nil, errors.New(fmt.Sprintf("missing required -%s flag\n", req))
		}
	}

	a.companyFilter = strings.Split(a.company, ",")

	return a, nil
}

func filterStargazers(a *arguments, u userNode) bool {
	company := strings.ReplaceAll(string(u.Company), "@", "")
	company = strings.TrimSpace(company)
	company = strings.ToLower(company)

	for _, v := range a.companyFilter {
		if v == company {
			return true
		}
	}

	if a.employee {
		return bool(u.IsEmployee)
	}

	return false
}

func getGitHubClient() (*graphql.Client, error) {
	pat, found := os.LookupEnv("GITHUB_PAT")

	if !found {
		fmt.Print("Enter GitHub PAT: ")
		pw, err := terminal.ReadPassword(0)
		fmt.Println()
		if err != nil {
			return nil, err
		}
		pat = string(pw)
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)

	httpClient := oauth2.NewClient(context.Background(), src)

	client := graphql.NewClient("https://api.github.com/graphql", httpClient)

	return client, nil
}

func getOrganizations(user userNode) []string {
	var organizations []string

	if len(user.Organizations.Edges) > 0 {
		for _, o := range user.Organizations.Edges {
			organizations = append(organizations, fmt.Sprintf("%s", o.Node.Name))
		}
	}

	return organizations
}

func getStargazers(a *arguments, ghc *graphql.Client) {
	variables := map[string]interface{}{
		"repositoryName":  graphql.String(a.repositoryName),
		"repositoryOwner": graphql.String(a.repositoryOwner),
		"stargazerCursor": (*graphql.String)(nil),
	}

	for {
		err := ghc.Query(context.Background(), &query, variables)
		if err != nil {
			log.Fatal(err)
		}

		stargazers := query.Repository.Stargazers
		if !stargazers.PageInfo.HasNextPage {
			break
		}

		variables["stargazerCursor"] = stargazers.PageInfo.EndCursor
		for _, v := range stargazers.Edges {
			user := v.Node

			if !filterStargazers(a, user) {
				continue
			}

			organizations := strings.Join(getOrganizations(user), ", ")
			fmt.Printf("[%s] %s (%s) <%s>: %s (%s)\n", user.Company, user.Name, user.Login, user.Email, user.Url, organizations)
		}
	}
}

func main() {
	a, err := getFlags()
	if err != nil {
		log.Fatal(err)
	}

	ghc, err := getGitHubClient()
	if err != nil {
		log.Fatal(err)
	}

	getStargazers(a, ghc)
}
