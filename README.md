# gh-vanity [![Release](https://github.com/jef/gh-vanity/actions/workflows/release.yaml/badge.svg)](https://github.com/jef/gh-vanity/actions/workflows/release.yaml)

I created this out of pure vanity, hence the name. I was curious as to who has starred my repositories (and others) and what companies they worked for.

This allows programs lets a user understand that data without parsing through _many_ pages of stargazers.

## Installation

1. Install the `gh` cli - see the [installation](https://github.com/cli/cli#installation)

   _Installation requires a minimum version (2.0.0) of the GitHub CLI that supports extensions._

2. Install this extension:

   ```shell
   gh extension install jef/gh-vanity
   ```

<details>
<summary><strong>Manual Installation</strong></summary>

Requirements: `cli/cli` and `go`.

1. Clone the repository

   ```shell
   # git
   git clone git@github.com:jef/gh-vanity.git

   # GitHub CLI
   gh repo clone jef/gh-vanity
   ```

2. `cd` into it

   ```shell
   cd gh-vanity
   ```

3. Build it

   ```shell
   make build
   ```

4. Install it locally

   ```shell
   gh extension install .
   ```
</details>

## Usage

To run:

```shell
gh vanity
```

To upgrade:

```sh
gh extension upgrade vanity
```

### Examples

Filter Amazon, Google, and GitHub employees for [jef/streetmerchant](https://github.com/jef/streetmerchant):

```shell
gh vanity --company=amazon,google --employee --owner=jef --name=streetmerchant
```
