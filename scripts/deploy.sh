#!/usr/bin/env bash

set -eo pipefail

aws s3 cp build/main.zip s3://lido-rewards-exporter-source/
aws lambda update-function-code --function-name lido-rewards-exporter --s3-bucket lido-rewards-exporter-source --s3-key main.zip --publish

