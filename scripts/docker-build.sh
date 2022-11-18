#!/bin/sh
HOST_NAME="tedvo.dev"
PROJECT_ID="ilt-poker-club-bot"
IMAGE_NAME="golang"
IMAGE_VERSION="v1.0.0"

TAG="$HOST_NAME/$PROJECT_ID/$IMAGE_NAME:$IMAGE_VERSION"

docker build --pull --rm -f $FILE_PATH -t $TAG ./ --no-cache
