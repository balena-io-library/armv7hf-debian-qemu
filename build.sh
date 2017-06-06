#!/bin/bash
set -e
set -o pipefail

TAR_FILE=resin-xbuild$VERSION.tar.gz

rm -rf resin-xbuild

go build -ldflags "-w -s" resin-xbuild.go

tar -cvzf $TAR_FILE resin-xbuild

curl -SLO "http://resin-packages.s3.amazonaws.com/SHASUMS256.txt"
sha256sum $TAR_FILE >> SHASUMS256.txt

# Upload to S3 (using AWS CLI)
printf "$ACCESS_KEY\n$SECRET_KEY\n$REGION_NAME\n\n" | aws configure
aws s3 cp $TAR_FILE s3://$BUCKET_NAME/resin-xbuild/v$VERSION/
aws s3 cp SHASUMS256.txt s3://$BUCKET_NAME/
