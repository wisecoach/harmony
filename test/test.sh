#!/bin/bash

node ../cli/deployTCCTest.js
cd ../cli-py
python ../cli-py/simulate_tcc.py
cd -

notify-send "TCC Test" "TCC test completed successfully."