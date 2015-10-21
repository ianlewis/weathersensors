#!/bin/sh

# Builds a docuker image for webfront

set -e

# Build a really static binary.
CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo .

VERSION=`./aggre_mod -version`

docker build -t aggremod .

docker tag aggremod asia.gcr.io/ianlewis-org/aggremod:${VERSION}

gcloud docker push asia.gcr.io/ianlewis-org/aggremod:${VERSION}
