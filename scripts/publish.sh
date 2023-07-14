#!/bin/bash
#
# check user input
if [ $# -ne 4 ]; then
    echo "Not enough arguments, usage: publish.sh <user-repo> <token> <release-tag> <release-notes>"
    exit 1;
fi

dir=$1   # e.g /path/CHAPUsers/CHAPBook
token=$2 # github access token string
tag=$3   # e.g. v0.0.0
notes=$4 # e.g. "some release notes"

if [ -z "$token" ]; then
    echo "Please provide proper github access token for CHAPUsers/CHAPBook repo"
    exit 1;
fi

repo="CHAPBook"
owner="CHAPUsers"
url=https://api.github.com/repos/$owner/$repo/releases

echo "Create release:"
echo "Directory     : $dir"
echo "Tag           : $tag"
echo "Notes         : $notes"
echo "URL           : $url"

# if tag==0 then we'll grab the last available tag and increment it by 1
if [ "$tag" == "0" ]; then
    # capture last available tag
    lastTag=`curl -ks -H "Accept: application/vnd.github+json" -H "Authorization: Bearer $token" $url | grep tag_name | awk '{print $2}' | sort | tail -1 | sed -e "s,\",,g" -e  "s#,##g"`
    echo "Last tag      : $lastTag"
    mainVersion=`echo $lastTag | awk '{split($1,a,"."); print a[1]}'`
    majorNumber=`echo $lastTag | awk '{split($1,a,"."); print a[2]}'`
    minorNumber=`echo $lastTag | awk '{split($1,a,"."); print a[3]}'`
    newMinorNumber=$((minorNumber+1))
    tag="${mainVersion}.${majorNumber}.${newMinorNumber}"
    echo "New tag       : $tag"
    # change notes as well
    notes=`echo $notes | sed -e "s,0,$tag,g"`
fi
payload=$(printf '{"tag_name": "%s","target_commitish": "main","name": "%s","body": "%s","draft": false,"prerelease": false}' $tag "$notes" "$notes")
echo "payload       : $payload"

# see: https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28

# post new release tag
curl -k -s -L \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer $token"\
  -H "X-GitHub-Api-Version: 2022-11-28" \
  -d "$payload" \
  $url
