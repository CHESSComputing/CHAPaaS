#!/bin/bash
#
# check user input
if [ $# -ne 4 ]; then
    echo "Not enough arguments, usage: chap.sh <config> <chap-dir> <user-dir> <profile>"
    exit 1;
fi
pyver=`python -V | awk '{split($2,a,"."); print ""a[1]"."a[2]""}'`
config=$1
chapdir=$2
userdir=$3
profile=$4
#echo "+++ config $config"
#echo "+++ chapdir $chapdir"
#echo "+++ userdir $userdir"
#echo "+++ profile $profile"
export PYTHONPATH=$chapdir/lib/python$pyver/site-packages:$chapdir/lib64/python$pyver/site-packages:$userdir
export PYTHONPATH=$PYTHONPATH:$chapdir/venv/lib/python$pyver/site-packages
#python -V
#python -c "from users import *; print(UserProcessor)"
#echo "#### start CHAP with $PYTHONPATH"
if [ -n "$profile" ]; then
    echo "run profile code"
    $chapdir/bin/CHAP --profile --config $config 2>&1
else
    echo "run non-profiled code"
    $chapdir/bin/CHAP --config $config 2>&1
fi
