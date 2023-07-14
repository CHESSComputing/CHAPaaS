#!/bin/bash
#
# check user input
if [ $# -ne 4 ]; then
    echo "Not enough arguments, usage: publish.sh <user-repo> <release-tag> <release-notes> <token>"
    exit 1;
fi

dir=$1   # e.g /path/CHAPUsers
tag=$2   # e.g. v0.0.0
notes=$3 # e.g. "some release notes"
token=$4 # github access token string

if [ -z "$token" ]; then
    echo "Please provide proper github access token for CHAPUsers/CHAPBook repo"
    exit 1;
fi

repo="CHAPUsers"
owner="CHESSComputing"


echo "Create release:"
echo "Directory     : $dir"
echo "Tag           : $tag"
echo "Notes         : $notes"
payload=$(printf '{"tag_name": "%s","target_commitish": "main","name": "Auto-generated release %s","body": "%s","draft": false,"prerelease": false}' $tag $tag "$notes")
echo "payload       : $payload"

# if tag==-1 then we'll grab the last available tag and increment it by 1
if [ "$tag" == "-1" ]; then
    # capture last available tag
    lastTag=`curl -ks -H "Authorization: Bearer $CHAPUSERS_TOKEN" https://api.github.com/repos/CHESSComputing/CHAPUsers/releases | grep tag_name | awk '{print $2}' | sort | tail -1 | sed -e "s,\",,g" -e  "s#,##g"`
    echo "Last tag      : $lastTag"
    mainVersion=`echo v0.0.2 | awk '{split($1,a,"."); print a[1]}'`
    majorNumber=`echo v0.0.2 | awk '{split($1,a,"."); print a[2]}'`
    minorNumber=`echo v0.0.2 | awk '{split($1,a,"."); print a[3]}'`
    newMinorNumber=$((minorNumber+1))
    tag="${mainVersion}.${majorNumber}.${newMinorNumber}"
    echo "New tag      : $tag"
fi

# see: https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28

# post new release tag
curl -k -s -L \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer $token"\
  -H "X-GitHub-Api-Version: 2022-11-28" \
  -d "$payload" \
  https://api.github.com/repos/$owner/$repo/releases
