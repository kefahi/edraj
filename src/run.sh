#!/bin/bash
reset
[ -e /tmp/edraj ] && rm -rf /tmp/edraj
go build -o edraj 
(netstat -nptl | grep 27017) > /dev/null 2>&1
if [ $? -eq 0 ]; then
	./edraj
else
	echo "Local Mongod service doesn't seem to be runnig (port 27017)"
fi

[ -e ./edraj ] && rm -f ./edraj
