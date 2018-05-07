#!/bin/bash

rm -rf pki/example
mkdir -p pki/example
cfssl gencert -initca pki/ca-csr.json | cfssljson -bare pki/example/ca
cfssl gencert \
  -ca=pki/example/ca.pem \
  -ca-key=pki/example/ca-key.pem \
  -config=pki/cert-config.json \
  -profile=default \
  pki/admission-webhook-csr.json | cfssljson -bare pki/example/admission-webhook