#!/bin/bash

if [ -z "${CA}" ]; then
    echo "CA is unset or set to the empty string"
fi
if [ -z "${TLS_CERT}" ]; then
    echo "TLS_KEY_FILE is unset or set to the empty string"
fi
if [ -z "${TLS_KEY}" ]; then
    echo "TLS_KEY is unset or set to the empty string"
fi

sed "s/__CA__/${CA}/" deployments/validating-webhook-configuration.yaml | kubectl apply -f -
sed "s/__CA__/${CA}/" deployments/admission-webhook.yaml | sed "s/__TLS_CERT__/${TLS_CERT}/" | sed "s/__TLS_KEY__/${TLS_KEY}/" | kubectl apply -f -