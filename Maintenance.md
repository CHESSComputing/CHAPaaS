### CHAPBook service maintenance
```
# login to Chapbook node
ssh -S none -t chessdata_svc@lnx15 'cd /home/chessdata_svc/CHAPBook/src/CHAPaaS/

# initialize conda environment
. /home/chessdata_svc/miniconda3/etc/profile.d/conda.sh

# check status of CHAPBook services
./scripts/manage status

# start CHAPBook services
./scripts/manage start

# stop CHAPBook services
./scripts/manage stop

# restart CHAPBook services
./scripts/manage restart
```
