#!/bin/bash
config=$1
chapdir=$2
userdir=$3
#echo "+++ config $config"
#echo "+++ chapdir $chapdir"
#echo "+++ userdir $userdir"
export PYTHONPATH=$chapdir/lib/python*/site-packages/:$userdir
#python -V
#python -c "from users import *; print(UserProcessor)"
#echo "#### start CHAP with $PYTHONPATH"
$chapdir/bin/CHAP --config $config 2>&1
