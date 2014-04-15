## githubissues

Generate HTML for issues from Gitub tracker.

To get access to private repositories you have to generate a personal API [access token](https://help.github.com/articles/creating-an-access-token-for-command-line-use).

### Install

Binary is available for download from [releases](https://github.com/kalimatas/githubissues/releases/tag/v1.0) page.

Or you can compile it yourself:


```
git clone git@github.com:kalimatas/githubissues.git
cd githubissues
git submodule init
git submodule update
export GOPATH=$(pwd)
go get -u github.com/google/go-querystring/query
go get -u code.google.com/p/goauth2/oauth
go install github.com/kalimatas/githubissues
```

### Usage

```
Usage of githubissues:
  -access-token="": access token for authenticated access
  -issues="": comma-separated issue numbers list
  -milestone=0: milestone number
  -owner="": repository owner
  -repository="": repository name
```
