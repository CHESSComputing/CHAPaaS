#!/bin/bash
config=$1
export PYTHONPATH=/Users/vk/tmp/ChessAnalysisPipeline/install/lib/python3.11/site-packages/:/tmp:/tmp/test
/Users/vk/tmp/ChessAnalysisPipeline/install/bin/CHAP --config $config 2>&1
