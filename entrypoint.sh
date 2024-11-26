#!/bin/sh -l

go_mod_path=$@
if [[ $go_mod_path ]]; then
  packages=$(sp-build dependencies list --module $go_mod_path --json)
  services=$(sp-build dependencies list --module $go_mod_path --services --json)
else
  packages=$(sp-build dependencies list --json)
  services=$(sp-build dependencies list --services --json)
fi

echo "packages=$packages" >> $GITHUB_OUTPUT
echo "services=$services" >> $GITHUB_OUTPUT