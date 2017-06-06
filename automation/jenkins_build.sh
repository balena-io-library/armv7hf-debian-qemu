#!/bin/bash

# Jenkins build steps

docker build -t xbuild-builder .

docker run --rm -e VERSION=$VERSION \
				-e ACCESS_KEY=$ACCESS_KEY \
				-e SECRET_KEY=$SECRET_KEY \
				-e BUCKET_NAME=$BUCKET_NAME xbuild-builder bash -ex build.sh

# Clean up builder image after every run
docker rmi -f xbuild-builder
