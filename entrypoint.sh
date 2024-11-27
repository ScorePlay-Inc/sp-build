#!/bin/sh

sh -c "git config --global --add safe.directory $PWD"

packages=$(/sp-build dependencies list --working-directory $1 --json)
services=$(/sp-build dependencies list --working-directory $1 --services --json)

echo "packages=$packages" >> $GITHUB_OUTPUT
echo "services=$services" >> $GITHUB_OUTPUT