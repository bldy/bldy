#!/bin/bash	
echo $TRAVIS_BRANCH
echo $TRAVIS_PULL_REQUEST
if [ "$TRAVIS_PULL_REQUEST" == false ]; then
	if [[ $TRAVIS_BRANCH == "master" ]]; then
		gox -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}"
		ghr --username bldy --token $GITHUB_TOKEN --replace --prerelease --debug `git describe --always`  dist/
	fi
fi