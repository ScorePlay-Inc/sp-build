#!/bin/sh -l

packages=$(/sp-build dependencies list --json)
services=$(/sp-build dependencies list --services --json)

echo "packages=$packages" >> $GITHUB_OUTPUT
echo "services=$services" >> $GITHUB_OUTPUT