PROJECT:=fast-api

linuxBash=GOPROXY=https://goproxy.cn,direct GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ./bin/${PROJECT}-${@} ./

macBash=GOPROXY=https://goproxy.cn,direct GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/${PROJECT}-${@} ./

winBash=GOPROXY=https://goproxy.cn,direct GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CGO_ENABLED=0  go build -ldflags "-s -w" -o ./bin/${PROJECT}-${@} ./

ifeq ($(OS),Windows_NT)
 PLATFORM="windows"
 autoBash=$(winBash)
else
 ifeq ($(shell uname),Darwin)
  PLATFORM="mac"
  autoBash=$(macBash)
 else
  PLATFORM="linux"
  autoBash=$(linuxBash)
 endif
endif

$(PLATFORM):
	@echo 当前系统是$(PLATFORM)
	$(autoBash)

linux:
	$(linuxBash)

mac:
	$(macBash)

win:
	$(winBash)

clean:
	rm -rf ./bin/${PROJECT}-*




PWD := $(shell pwd)
VERSION := $(shell git log --pretty=format:"%h" | head -n 1)

build-docker-gobase:
	@if [ ! $(shell docker image ls --format "{{.Repository}}"|grep gobase) ]; then docker build -t gobase:latest -f ./DockerfileBase .; fi

# make build-linux
build-docker:
	@echo $(VERSION)
	@docker build -t fast-api:${VERSION} -t fast-api:latest .
	@echo "build successful"


# make run
run:
    # delete fast-api-api container
	@if [ $(shell docker ps -aq --filter name=fast-api) ]; then docker rm -f fast-api; fi
	@if [ $(shell docker ps -aq --filter name=fast-mq) ]; then docker rm -f fast-asynq; fi
	@if [ $(shell docker ps -aq --filter name=fast-cron) ]; then docker rm -f fast-cron; fi


	@docker run -d --name fast-api --privileged  --network web -v ${PWD}/config.yaml/:config.yaml -v ${PWD}/static/:/go/static/ -v /var/log/fast-api/:/go/temp/ fast-api:${VERSION} ./main server -c config.yaml
	@docker run -d --name fast-asynq --privileged  --network web -v ${PWD}/config.yaml/:config.yaml -v /var/log/fast-api-asynq/:/go/temp/ fast-api:${VERSION} ./main mq -c config.yaml
	@docker run -d --name fast-cron --privileged  --network web -v ${PWD}/config.yaml/:config.yaml -v /var/log/fast-api-cron/:/go/temp/ fast-api:${VERSION} ./main cron -c config.yaml

	@echo "fastApi service is running..."

	# delete Tag=<none> 的镜像
	@docker image prune -f
	@docker images fast-api --format "{{.Repository}}:{{.Tag}}" | grep -v 'fast-api:latest' | grep -v 'fast-api:${VERSION}' | xargs -r docker image rm
	@docker ps -a | grep "fast"


# make deploy
deploy-dev:
	make build-docker-gobase
	make build-docker
	make run

swag-api:
	@swag init --parseDependency --parseInternal --parseGoList=false   --parseDepth=6 -o ./docs -d .