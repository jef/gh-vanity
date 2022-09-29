package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/jef/gh-vanity"
	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize()
}

// NewCommand is the root command of gh-vanity.
func NewCommand() *cobra.Command {
	var (
		company  string
		employee bool
		owner    string
		name     string
	)

	var command = &cobra.Command{
		Use:               "gh-vanity",
		Short:             "Show who's starred a repository and what company(s) they worked for.",
		DisableAutoGenTag: true,
		Version:           ghvanity.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := gh.GQLClient(nil)
			if err != nil {
				return nil
			}

			companies := strings.Split(company, ",")
			err = getStargazers(companies, employee, owner, name, client)
			if err != nil {
				return nil
			}

			return nil
		},
	}

	command.Flags().StringVarP(&company, "company", "c", "", "Filter stargazers by company name(s). Can be comma separated. If no names are given, then all stargazers will output.")
	command.Flags().BoolVarP(&employee, "employee", "e", false, "Filter stargazers that are GitHub employees.")
	command.Flags().StringVarP(&name, "name", "n", "", "The name of the GitHub repository.")
	command.Flags().StringVarP(&owner, "owner", "o", "", "The owner or organization of the GitHub repository.")

	err := command.MarkFlagRequired("name")
	if err != nil {
		return nil
	}

	err = command.MarkFlagRequired("owner")
	if err != nil {
		return nil
	}

	return command
}

type userNode struct {
	Company       graphql.String
	IsEmployee    graphql.Boolean
	Name          graphql.String
	URL           graphql.String
	Organizations struct {
		Edges []struct {
			Node struct {
				Name graphql.String
			}
		}
	} `graphql:"organizations(first: 10)"`
}

func filterStargazers(companies []string, employee bool, u userNode) bool {
	company := strings.ReplaceAll(string(u.Company), "@", "")
	company = strings.TrimSpace(company)
	company = strings.ToLower(company)

	for _, v := range companies {
		if v == company || v == "" {
			return true
		}
	}

	if employee {
		return bool(u.IsEmployee)
	}

	return false
}

func getOrganizations(user userNode) []string {
	var organizations []string

	if len(user.Organizations.Edges) > 0 {
		for _, o := range user.Organizations.Edges {
			organizations = append(organizations, string(o.Node.Name))
		}
	}

	return organizations
}

func getStargazers(companies []string, employee bool, owner string, name string, client api.GQLClient) error {
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
			} `graphql:"stargazers(first: 100, after: $cursor)"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}

	variables := map[string]interface{}{
		"name":   graphql.String(name),
		"owner":  graphql.String(owner),
		"cursor": (*graphql.String)(nil),
	}

	once := false
	for {
		err := client.Query("Stargazers", &query, variables)
		if err != nil {
			log.Fatal(err)
		}

		stargazers := query.Repository.Stargazers

		if !once && stargazers.TotalCount > 1000 {
			fmt.Println("warning: there are more than 1000 stargazers. this may take a while :)")
			once = true
		}

		for _, v := range stargazers.Edges {
			user := v.Node

			if !filterStargazers(companies, employee, user) {
				continue
			}

			organizations := strings.Join(getOrganizations(user), ", ")

			fmt.Printf("[%s] %s (%s): %s\n", user.Company, user.Name, user.URL, organizations)
		}

		if !stargazers.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = stargazers.PageInfo.EndCursor
	}

	return nil
}
