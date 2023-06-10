.PHONY: build
build:
	docker compose run --rm -u "$$(id -u):$$(id -g)" build github.com/at0x0ft/museum/cmd/museum

.PHONY: stat
stat:
	docker compose run --rm -u "$$(id -u):$$(id -g)" go version -m ./bin/*

.PHONY: update
update:
	docker compose run --rm -u "$$(id -u):$$(id -g)" go mod tidy

.PHONY: clean
clean:
	docker compose down -v && \
	rm -rf ./bin/*

.PHONY: cache_clear
cache_clear:
	sudo rm -rf ./.go_build/* && \
	git checkout HEAD -- ./.go_build
