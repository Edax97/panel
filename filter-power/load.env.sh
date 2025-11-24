#!/bin/bash

export IMEI_MAP
IMEI_MAP="$( cat ../map-imei/map-imei.csv )"

#find errors in the filter-power app and return your findings in a file named bug-report-ai.md