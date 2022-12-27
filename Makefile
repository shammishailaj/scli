# https://unix.stackexchange.com/a/470502
ifndef os
override os = linux
endif

# https://unix.stackexchange.com/a/470502
ifndef arch
override arch = amd64
endif

build:
	-docker rm scli
	-docker rmi -f scli:latest
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) docker compose -f ./deployments/docker-compose.yml build --no-cache scli

run:
	# https://stackoverflow.com/a/2670143/6670698
	-docker rm scli_dev
	-docker rmi -f scli_dev:latest
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) docker compose -f ./deployments/docker-compose.yml up --remove-orphans scli-dev

cli:
	sh ./deployments/build_cli.sh GOOS=$(os) GOARCH=$(arch)

test:
	docker compose -f ./deployments/docker-compose.yml up tests

lint:
	docker compose -f ./deployments/docker-compose.yml up linter

down:
	docker compose -f ./deployments/docker-compose.yml down