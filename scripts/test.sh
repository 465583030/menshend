#!/usr/bin/env bash
set -e
go get github.com/mitchellh/gox
go get github.com/wadey/gocovmerge
go test -v -coverpkg=$(go list ./... | grep -v /vendor/)
export PKGS=$(go list ./... | /bin/grep -v /vendor/ | /bin/grep -v apis/client | /bin/grep -v pfclient)
export PKGS_DELIM=$(echo "$PKGS" | paste -sd "," -)
go list -f '{{if or (len .TestGoFiles) (len .XTestGoFiles)}}go test -covermode count -coverprofile {{.Name}}_{{len .Imports}}_{{len .Deps}}.coverprofile -coverpkg $PKGS_DELIM {{.ImportPath}}{{end}}' $PKGS | xargs -I {} bash -c {}
gocovmerge `ls *.coverprofile` >  coverage.txt
bash <(curl -s https://codecov.io/bash)
