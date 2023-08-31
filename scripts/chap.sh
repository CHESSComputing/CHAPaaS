#!/bin/bash
#
# check user input
if [ $# -ne 3 ]; then
    echo "Not enough arguments, usage: chap.sh <config> <worflows-dir> <user-dir>"
    exit 1;
fi
pyver=`python -V | awk '{split($2,a,"."); print ""a[1]"."a[2]""}'`
config=$1
wdir=$2
udir=$3

# obtain appropriate conda env for work workflow
cenv=`cat $wdir/conda.yml | grep ^name | awk '{print $2}'`
echo "conda environment: $cenv"

# initialize workflow conda environment
testEnv=`conda env list | grep ^$cenv`
if [ -z "$testEnv" ]; then
    echo "Unable to identify $cenv in conda environment"
    conda env list
    exit 1
fi
eval "$(conda shell.bash hook)"
conda activate $cenv

# copy workflow files to user area
cp -f -r $wdir/* $udir

# finally, cd to user directory to run CHAP job
cd $udir
CHAP $config 2>&1
