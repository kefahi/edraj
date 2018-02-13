#!/bin/env python3

# pip3 install --user openapi-spec-validator

from openapi_spec_validator import validate_spec_url
validate_spec_url('file:///home/kefah/Development/edraj/api-spec/edraj.json')
