#!/bin/bash -x

curl -s http://127.0.0.1:5533/api/entry/ -d @sample_api/delete.json | jq .Status
curl -s http://127.0.0.1:5533/api/entry/ -d @sample_api/delete2.json | jq .Status
curl -s http://127.0.0.1:5533/api/entry/ -d @sample_api/create.json | jq .Status
curl -s http://127.0.0.1:5533/api/entry/ -d @sample_api/create2.json | jq .Status
curl -s http://127.0.0.1:5533/api/entry/ -d @sample_api/query.json | jq .Entries[].Content.ID
curl -s http://127.0.0.1:5533/api/entry/ -d @sample_api/get.json | jq .Entries[].Content.ID
curl -s http://127.0.0.1:5533/api/entry/ -d @sample_api/delete.json | jq .Status
curl -s http://127.0.0.1:5533/api/entry/ -d @sample_api/delete2.json | jq .Status
