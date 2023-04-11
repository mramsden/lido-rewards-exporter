#!/usr/bin/env bash

set -eo pipefail

aws cloudformation deploy \
	--stack-name lido-rewards-exporter \
	--template-file infrastructure/app.yml \
	--capabilities CAPABILITY_NAMED_IAM


