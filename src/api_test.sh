#!/bin/bash
curl -s -v http://127.0.0.1:5533/api/entry/query -d @./sample_json/query.json  | jq
