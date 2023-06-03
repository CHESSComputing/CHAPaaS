#!/bin/bash
pyver=`python -V | awk '{split($2,a,"."); print ""a[1]"."a[2]""}'`
config=$1
chapdir=$2
userdir=$3
#echo "+++ config $config"
#echo "+++ chapdir $chapdir"
#echo "+++ userdir $userdir"
export PYTHONPATH=$chapdir/lib/python$pyver/site-packages:$chapdir/lib64/python$pyver/site-packages:$userdir
export PYTHONPATH=$PYTHONPATH:$chapdir/venv/lib/python$pyver/site-packages
#python -V
#python -c "from users import *; print(UserProcessor)"
#echo "#### start CHAP with $PYTHONPATH"
$chapdir/bin/CHAP --config $config 2>&1
