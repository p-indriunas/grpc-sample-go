#!/bin/bash

#
# Based on:
# https://github.com/grpc-ecosystem/grpc-gateway#usage
#

mkdir -p ./google

function download() {
  source_url=$1
  target_dir=$2

  if [ -z "$target_dir" ]; then
    target_dir="."
  fi

  target_file="${target_dir}/$(basename "$source_url")"
  if [ ! -f "$target_file" ]; then
    echo "Downloading ${source_url}"

    wget -q "$source_url" -P $target_dir

    ret=$?
    if [ $ret -ne 0 ]; then
      echo "Download failed."
    fi
  fi
  return $ret
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

# Generate GRPC service contract
protoc -I . \
    --go_out=../gen/go \
    --go-grpc_out=require_unimplemented_servers=false:../gen/go \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    EchoService.proto

# Generate GRPC gateway
protoc -I . \
    --grpc-gateway_out ../gen/go \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    EchoService.proto