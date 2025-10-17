#!/bin/bash

node ../ssc-cli/cli/deployTCCTest.js
cd ../ssc-cli/cli-py
python ./simulate_tcc.py
cd -

notify-send "TCC Test" "TCC test completed successfully."