#!/usr/bin/env bash

set -e

aws dynamodb create-table \
  --table-name grpc-gateway-example-user \
  --attribute-definitions \
  AttributeName=user_id,AttributeType=S \
  --key-schema \
  AttributeName=user_id,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=10,WriteCapacityUnits=10 \
  --endpoint-url http://localhost:54000 \
  --region ap-northeast-2 || true | cat