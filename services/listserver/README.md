# List Server

Simple server to provide some "keys" for the Mixer Adapter listchecker adapter.

## Build

docker build -t .
docker tag listserver:latest bladedancer/listserver:1.0.0
docker push bladedancer/listserver:1.0.0
