#!/bin/sh
HOST_NAME="asia.gcr.io"
PROJECT_ID="meepo-vn"
IMAGE_NAME="ilt-poker-club-bot"
IMAGE_VERSION="latest"

TAG="$HOST_NAME/$PROJECT_ID/$IMAGE_NAME:$IMAGE_VERSION"

docker build --pull --rm -t $TAG ./ --no-cache

docker push $TAG
