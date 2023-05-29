.PHONY: build
build:
	docker-compose run --rm -u "$$(id -u):$$(id -g)" build github.com/at0x0ft/museum/cmd/museum

.PHONY: stat
stat:
	docker-compose run --rm -u "$$(id -u):$$(id -g)" go version -m ./bin/*

.PHONY: clean
clean:
	docker-compose down && \
	rm -rf ./bin/*
