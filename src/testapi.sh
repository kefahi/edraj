#!/bin/bash -x

curl -s http://127.0.0.1:5533/api/entry/ -d @delete.json | jq
curl -s http://127.0.0.1:5533/api/entry/ -d @create.json | jq
curl -s http://127.0.0.1:5533/api/entry/ -d @get.json | jq
curl -s http://127.0.0.1:5533/api/entry/ -d @delete.json | jq
