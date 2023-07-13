#!/bin/bash
#
# check user input
if [ $# -ne 2 ]; then
    echo "Not enough arguments, usage: publish.sh <user-dir> <user-repo>"
    exit 1;
fi

userDir=$1 # e.g /path/chap/Users/user-name
userRepo=$2 # e.g /path/CHAPUsers
echo "Publish:"
echo "$userDir"
echo "to github repo:"
echo "$userRepo"
cd $userDir

# rsync commands
rsync --progress -avuz --numeric-ids --exclude=__pycache__ "$userDir" "$userRepo"
cd $userRepo
# git commands, commit only if we have git access
user=`echo $userDir | awk '{split($1,a,"/"); print a[length(a)]}'`
gitAccess=`cat .git/config | grep url | awk '{print $1,$2,$3}' | grep ^url | grep "git@github"`
echo "Update code for user: $user"
tstamp=`date -u`
if [ -n "$gitAccess" ]; then
    echo "commit to github $user"
    git add $user
    git commit -m "CHAPBook publish update for $user on $tstamp" $user
    git push -f
else
    echo "skip github commit due to lack of protocol permission"
fi
