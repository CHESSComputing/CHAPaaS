#!/bin/bash
### for more info see: https://wiki.classe.cornell.edu/Computing/GridEngine
#
# check user input
if [ $# -ne 4 ]; then
    echo "Not enough arguments, usage: batch.sh <workflow> <config> <worflows-dir> <user-dir>"
    exit 1;
fi
pyver=`python -V | awk '{split($2,a,"."); print ""a[1]"."a[2]""}'`
workflow=$1
config=$2
wdir=$3
udir=$4

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

# finally, cd to user directory, copy all necessary intput files and run CHAP job
mkdir -p $udir/$workflow
cd $udir/$workflow
cp -f -r $wdir/* .i

# create tarball of user workflow
tar cf /tmp/batch.tar .

# create new batch_submit.sh script which will send tarball to worker node
cat > /tmp/batch_submit.sh << EOF
#!/bin/bash
#$ -q all.q
#$ -S /bin/bash
#$ -l mem_free=8G
#
CHAP $config 2>&1 1>& chap.log
cat chap.log
EOF

# TODO: implement batch sybmission for CHAP workflow
qsub -q all.q /tmp/batch_submit.sh
