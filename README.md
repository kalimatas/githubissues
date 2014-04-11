## githubissues

Generate HTML for issues from Gitub tracker.

### Install

```
git clone git@github.com:kalimatas/githubissues.git
cd githubissues
git submodule init
git submodule update
export $GOPATH=$(pwd)
go get -u github.com/google/go-querystring/query
go get -u code.google.com/p/goauth2/oauth
go install ./...
```
