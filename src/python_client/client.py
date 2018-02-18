#!/bin/env python3
"""The Python implementation of the GRPC edraj.EntryService client."""

from __future__ import print_function

import grpc

import edraj_pb2
import edraj_pb2_grpc

def run():

    with open('../../certs/edrajRootCA.crt') as f:
        rootca_crt = bytes(f.read(), 'utf8')

    with open('../../certs/kefah.key') as f:
        kefah_key = bytes(f.read(), 'utf8')

    with open('../../certs/kefah.crt') as f:
        kefah_crt = bytes(f.read(), 'utf8')


    credentials = grpc.ssl_channel_credentials(root_certificates=rootca_crt, private_key=kefah_key, certificate_chain=kefah_crt )
    channel = grpc.secure_channel('localhost:50051', credentials)
    stub = edraj_pb2_grpc.EntryServiceStub(channel)
    response = stub.Query(edraj_pb2.QueryRequest(query=edraj_pb2.Query(entry_type='CONTENT')), metadata=[('edraj-foo', 'bar')])

    for entry in response.entries:
        #print(entry.content.tags)
        print(entry.content)


if __name__ == '__main__':
    run()
