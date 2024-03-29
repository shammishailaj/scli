#!/bin/sh
printf "Building inside a docker container...\n"
OS=$1
ARCH=$2
# See https://docs.docker.com/compose/compose-file/
# See https://docs.docker.com/develop/develop-images/multistage-build/#before-multi-stage-builds
docker rmi -f scli_cli:latest
GOOS=$OS GOARCH=$ARCH CGO_ENABLED=0 COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./deployments/docker-compose.yml build scli-cli
printf "Extracting from the built docker container image\n"
printf "Creating container...\n"
docker container create --name extractor scli_cli:latest
printf "Copying file to local filesystem...\n"

if [ ! -d ./bin ]
then
echo "Executables Directory not found. Creating..."
mkdir ./bin
fi
VERSION=$(cat VERSION)
docker container cp extractor:/bin/scli ./bin/scli-"$VERSION"-"$OS"-"$ARCH"
printf "Removing container...\n"
docker container rm -f extractor
docker rmi -f scli_cli:latest
