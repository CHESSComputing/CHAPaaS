#!/bin/bash

# delete previous processes
cpid=`ps auxww | grep "chapaas -config" | grep -v grep | awk '{print $2}'`
if [ -n "$cpid" ]; then
    echo "kill previous chapaas process $cpid"
    kill -9 $cpid
fi
jpid=`ps auxww | grep "jupyter-notebook --config" | grep -v grep | awk '{print $2}'`
if [ -n "$jpid" ]; then
    echo "kill previous jupyter-notebook process $jpid"
    kill -9 $jpid
fi


# create proper environment
idir=/home/chess_chapaas/chess/CHAPaaS
jdir=$idir/jupyter
cdir=$idir/chap
unset PYTHONPATH
source venv/bin/activate
mkdir -p $jdir/logs
mkdir -p $cdir/logs

# start jupyter server
mkdir -p $jdir/logs
cd $jdir
echo "starting Jupyter notebook server in $PWD"
nohup jupyter notebook --config $idir/scripts/jupyter_config.py \
    2>&1 1>& $jdir/logs/jupyter.log < /dev/null & \
    echo $! > $jdir/logs/jupyter.pid
echo "Jupyter notebook PID=`cat $jdir/logs/jupyter.pid`"
echo "Jupyter notebook logs $jdir/logs/jupyter.log"
cd -

# start chapaas server
cd $cdir
echo "starting chapaas server in $PWD"
nohup $idir/chapaas -config $idir/config-http.json \
    2>&1 1>& $cdir/logs/chapaas.log < /dev/null & \
    echo $! > $cdir/logs/chapaas.pid
echo "chapaas PID=`cat $cdir/logs/chapaas.pid`"
echo "chapaas logs $cdir/logs/chapaas.log"
cd $idir

echo "Current chap processes:"
ps axu | grep --color=auto -v grep | grep --color=auto "chap" -i --color=auto
sleep 3
echo
echo "Jupyter tail $jdir/logs/jupyter.log"
tail $jdir/logs/jupyter.log
echo
echo "CHAPaaS tail $cdir/logs/chapaas.log"
tail $cdir/logs/chapaas.log
