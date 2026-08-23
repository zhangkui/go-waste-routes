#!/usr/bin/env bash
set -euo pipefail
DOCKER_PLATFORM="${DOCKER_PLATFORM:-linux/amd64}"
IMAGE_NAME="${IMAGE_NAME:-go-waste-routes}"
docker build --platform "$DOCKER_PLATFORM" -f benzhi.Dockerfile -t "$IMAGE_NAME" .
