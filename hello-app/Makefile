
# Makefile for hello-app

version=v2.0

.PHONY: all
all: go-build docker-build

.PHONY: set-app-version
set-app-version:
	sed -i.tmp 's/SET_TO_APP_VERSION/${version}/' hello.go
	rm hello.go.tmp

.PHONY: reset-app-version
reset-app-version:
	sed -i.tmp 's/${version}/SET_TO_APP_VERSION/' hello.go
	rm hello.go.tmp

.PHONY: go-build-as-is
go-build-as-is:
	go build hello.go

.PHONY: go-build
go-build: set-app-version go-build-as-is reset-app-version

.PHONY: go-run
go-run:
	./hello

.PHONY: docker-build-as-is
docker-build-as-is:
	docker build -t hello:${version} .

.PHONY: docker-build
docker-build: set-app-version docker-build-as-is reset-app-version

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 -d hello:${version}

.PHONY: check-registry-env-vars
check-registry-env-vars:
	@if [[ -z "${IRHOSTNAME}" ]] || [[ -z "${IRNAMESPACE}" ]] || [[ -z "${IRUSER}" ]] || [[ -z "${IRPASSWORD}" ]]; then \
		printf "\nERROR: Image Registry env vars not set. \nRequired:\n"; \
		printf "  - IRHOSTNAME - full hostname of image registry\n"; \
		printf "  - IRNAMESPACE - namespace in image registry\n"; \
		printf "  - IRUSER - user for image registry\n"; \
		printf "  - IRPASSWORD - password for image registry\n"; \
	fi

.PHONY: docker-tag-and-push
docker-tag-and-push: check-registry-env-vars 
	docker login -u ${IRUSER} -p ${IRPASSWORD} ${IRHOSTNAME}
	docker tag hello:${version} ${IRHOSTNAME}/${IRNAMESPACE}/hello:${version}
	docker push ${IRHOSTNAME}/${IRNAMESPACE}/hello:${version}
