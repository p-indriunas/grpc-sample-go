#!/bin/bash

#
# Based on tutorial:
# https://github.com/grpc-ecosystem/grpc-gateway#usage
#

mkdir -p ./google

function download() {
  SOURCE_URL=$1
  TARGET_DIR=$2

  if [ -z $TARGET_DIR ]; then
    TARGET_DIR="."
  fi

  TARGET_FILE="${TARGET_DIR}/$(basename "$SOURCE_URL")"
  if [ ! -f $TARGET_FILE ]; then
    echo "Downloading ${SOURCE_URL}"

    wget -q $SOURCE_URL -P $TARGET_DIR
    if [ $? -ne 0 ]; then
      echo "Download failed. Terminating."
    fi
  fi
  return $?
}

download "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto" ./google/api
download "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto" ./google/api
download "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto" ./google/api
download "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/httpbody.proto" ./google/api

if [[ ! -x "$(command -v protoc-gen-go)" ]]; then
    echo "INSTALLING protoc-gen-go"
    go install google.golang.org/protobuf/cmd/protoc-gen-go
fi

if [[ ! -x "$(command -v protoc-gen-grpc-gateway)" ]]; then
    echo "INSTALLING protoc-gen-grpc-gateway"
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
fi

mkdir -p ../gen/go

protoc -I . \
    --go_out=../gen/go \
    --go-grpc_out=require_unimplemented_servers=false:../gen/go \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    EchoService.proto

protoc -I . \
    --grpc-gateway_out ../gen/go \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    EchoService.proto