#!/bin/bash
#
# check user input
if [ $# -ne 2 ]; then
    echo "Not enough arguments, usage: galaxy.sh <user-repo> <chap-repo>"
    exit 1;
fi

udir=$1 # user-repo, e.g. /path/users/user-name
cdir=$2 # chap-repo, e.g. /path/ChessAnalysisPipeline

# get tag of CHAP repo
chapVersion=`grep EASY /Users/vk/tmp/ChessAnalysisPipeline/install/bin/CHAP | sed -e "s#,# #g" -e "s,',,g" -e "s,ChessAnalysisPipeline==,,g" | awk '{print $3}'`
pythonVersion=`python -V | awk '{print $2}'`
profile="21.05"
year=`date +%Y`
user=`echo $udir | awk '{z=split($1,a,"/"); print a[z]}'`

# generate galaxy XML file
cid="CHAPBook_$user"
cat > /tmp/chapbook-galaxy.xml << EOF
<tool id="$cid" name="CHESS Analysis Pipeline" version="$chapVersion" python_template_version="pythonVersion" profile="$profile">
    <requirements>
        <requirement type="package" version="$chapVersion">ChessAnalysisPipeline</requirement>
    </requirements>
    <command detect_errors="exit_code"><![CDATA[
        CHAP --config '$config'
    ]]></command>
    <inputs>
        <param type="data" name="config" format="yaml" />
    </inputs>
    <outputs>
    </outputs>
    <tests>
        <test>
            <param name="config" value="config.yaml"/>
        </test>
    </tests>
    <help><![CDATA[
        CHESS Analysis Pipeline (CHAP):

        To run it on command line you'll use:
        CHAP --config CONFIG

        To run it within galaxy you'll only need to upload your
        required configuration pipeline and necessary data.
    ]]></help>
    <citations>
        <citation type="bibtex">
@misc{githubChessAnalysisPipeline,
  author = {github, $user},
  year = {$year},
  title = {CHAPBook analysis},
  publisher = {GitHub},
  journal = {GitHub repository},
  url = {https://github.com/CHAPUsers/CHAPBook},
}</citation>
    </citations>
</tool>
EOF
