#!/bin/bash
[ -e /tmp/edraj ] && rm -rf /tmp/edraj
go build -o edraj && ./edraj
[ -e ./edraj ] && rm ./edraj
