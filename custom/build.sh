#!/bin/sh

# 1 - Setup Environment Variables
. ./env.sh

# 2 - Get Istio
mkdir -p $GOPATH/src/istio.io/
cd $GOPATH/src/istio.io/
git clone https://github.com/istio/istio

# 3 - Build mixer server & client binary
cd $ISTIO/istio && make mixs
cd $ISTIO/istio && make mixc

# 4 - Copy
mkdir -p $MIXER_REPO/adapter/onetimeadapter/config
cp $ROOT_FOLDER/onetimeadapter/onetimeadapter_impl.go $MIXER_REPO/adapter/onetimeadapter/onetimeadapter_impl.go
cp -rf $ROOT_FOLDER/onetimeadapter/cmd $MIXER_REPO/adapter/onetimeadapter/cmd
cp $ROOT_FOLDER/onetimeadapter/config/config.proto $MIXER_REPO/adapter/onetimeadapter/config/config.proto

# 5 - Build and generate
cd $MIXER_REPO/adapter/onetimeadapter
go generate ./...
go build ./...

# 6 - Copy generated
cp $MIXER_REPO/adapter/onetimeadapter/config/onetimeadapter.yaml $ROOT_FOLDER/generated
cp $MIXER_REPO/testdata/config/attributes.yaml $ROOT_FOLDER/generated
cp $MIXER_REPO/template/authorization/template.yaml $ROOT_FOLDER/generated
cp $ROOT_FOLDER/onetimeadapter/testdata/sample_operator_cfg.yaml $ROOT_FOLDER/generated/sample_operator_cfg.yaml

# 7 - To Test

## Start Adapter
echo "STARTING ADAPTER"
cd $MIXER_REPO/adapter/onetimeadapter
go run cmd/main.go 44225 &
ADAPTER_PID=$!
sleep 5

## Start Local mixer server
echo "STARTING MIXS"
$GOPATH/out/linux_amd64/release/mixs server --configStoreURL=fs://$ROOT_FOLDER/generated --log_output_level=attributes:debug &
MIXS_PID=$!
sleep 5

## Send test request
echo "Sending test request"
$GOPATH/out/linux_amd64/release/mixc check -s destination.service="svc.cluster.local" --stringmap_attributes "request.headers=x-custom-token:abc"

kill -9 $ADAPTER_PID
kill -9 $MIXS_PID
