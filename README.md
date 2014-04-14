## githubissues

Generate HTML for issues from Gitub tracker.

### Install

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
