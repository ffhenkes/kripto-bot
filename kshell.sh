#!/bin/bash

set -e

export KRIPTO_SERVER_ENDPOINT=https://localhost:20443/v1
export KRIPTO_USERNAME=ffhenkes
export KRIPTO_PASSWORD=test
export KRIPTO_INSECURE_TLS=true
export KRIPTO_APP=sample_app

./kripto-bot

source ${KRIPTO_CLIENT_APP}'.env'

env | grep SAMPLE
rm -f ${KRIPTO_CLIENT_APP}'.env'
