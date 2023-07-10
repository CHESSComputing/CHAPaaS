#!/bin/bash
#
# check user input
if [ $# -ne 3 ]; then
    echo "Not enough arguments, usage: publish.sh <notebook> <user-dir> <user-repo>"
    exit 1;
fi

notebook=$1
userDir=$2 # e.g /path/chap/Users
userRepo=$3 # e.g /path/CHAPUsers
echo "Publish $notebook"
cd $userDir
# rsync commands
rsync --progress -avuz --numeric-ids --exclude=__pycache__ "$userDir" "$userRepo"
cd $userRepo
# git commands
#git commit -m "CHAPBook publish update"
echo "git commit -m \"CHAPBook publish update\" ."
cd -
