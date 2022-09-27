# gh-vanity [![Release](https://img.shields.io/github/workflow/status/jef/gh-vanity/Release)](https://github.com/jef/gh-vanity/actions/workflows/release.yaml)

> **star·gaz·er** `/ˈstärˌɡāzər/` Someone that has starred a repository.
>
> **van·i·ty** `/ˈvanədē/` Excessive pride in or admiration of one's own appearance or achievements.

I created this out of pure vanity, hence the name. I was curious as to who has starred my repositories (and others) and what companies they worked for.

This allows programs lets a user understand that data without parsing through _many_ pages of stargazers.

## Installation

Download the specific binary from [GitHub Releases](https://github.com/jef/gh-vanity/releases) or use `go get`:

```shell
go get -u github.com/jef/gh-vanity
```

## Usage

It is required that you use a GitHub Personal Access Token (PAT). You can generate one [here](https://github.com/settings/tokens/new). The required scopes are `['read:org', 'user:email', 'read:user']`. Set your PAT to environment variable `GITHUB_PAT`. If `GITHUB_PAT` isn't set, you will be prompted for your PAT in the beginning of startup.

```
Usage of ./gh-vanity:
  -company string
    	Filter stargazers by company name(s). Can be comma separated.
    	If no names are given, then all stargazers will output.
  -employee
    	Filter stargazers that are GitHub employees.
  -repo string
    	(Required) The name of the repository.
  -owner string
    	(Required) The owner or organization of the repository.
```

### Examples

Filter Amazon, Google, and GitHub employees for [cli/cli](https://github.com/cli/cli):

```
./gh-vanity -company=amazon,google -employee -owner=cli -repo=cli
```

Filter Nvidia employees for [jef/streetmerchant](https://github.com/jef/streetmerchant):

```
./gh-vanity -company=nvidia -owner=jef -repo=streetmerchant
```

## Development

- `make build`: Builds source
- `make clean`: Cleans executable
- `make dist`: Cross-compilation for distribution
