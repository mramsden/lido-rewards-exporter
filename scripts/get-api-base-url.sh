#!/usr/bin/env bash

aws cloudformation describe-stacks --stack-name lido-rewards-exporter --query 'Stacks[0].Outputs[?OutputKey==`ApiGatewayUrl`] | [0].OutputValue' --output text
