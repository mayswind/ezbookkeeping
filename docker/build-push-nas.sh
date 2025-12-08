#!/bin/bash

# Check if an image name is provided
if [ -z "$1" ]; then
  echo "Error: Please provide a Docker Hub image name (e.g., username/ezbookkeeping)."
  echo "Usage: $0 <image_name>"
  exit 1
fi

IMAGE_NAME=$1

# Ensure docker buildx is available and use it
if ! docker buildx version > /dev/null 2>&1; then
  echo "Error: docker buildx is not installed or enabled."
  echo "Please install Docker Desktop or enable buildx."
  exit 1
fi

echo "Building and pushing image: $IMAGE_NAME for platform linux/amd64..."

# Create a new builder instance if one doesn't exist (optional but recommended for multi-arch)
# docker buildx create --use

# Build and push the image
docker buildx build --platform linux/amd64 -t "$IMAGE_NAME" --push .

if [ $? -eq 0 ]; then
  echo "Successfully built and pushed $IMAGE_NAME"
else
  echo "Failed to build or push the image."
  exit 1
fi
