#!/bin/bash
printenv | sed 's/^\(.*\)$/export \1/g' > /env_vars.sh
chmod +x /env_vars.sh