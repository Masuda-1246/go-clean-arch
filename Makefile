.DEFAULT_GOAL := help

### docker ###
DOCKER_COMPOSE                     := docker-compose.develop.yml
DOCKER_EXEC                        := docker-compose exec -it
DOCKER_COMPOSE_PJ_NAME             := app
DOCKER_COMPOSE_DEFAULT_OPTIONS     := -f $(DOCKER_COMPOSE) \
										-p $(DOCKER_COMPOSE_PJ_NAME)
DOCKER_COMPOSE_PRODUCTION_OPTIONS  := -f $(DOCKER_COMPOSE_PRODUCTION) \
										-p $(DOCKER_COMPOSE_PJ_NAME)
DOCKER_COMPOSE_TEST_OPTIONS        := -f $(DOCKER_COMPOSE_TEST) \
										-p $(DOCKER_COMPOSE_PJ_NAME)
API_CONTAINER_NAME                 := api
DB_CONTAINER_NAME                  := postgres

### command ###
RM := rm -rf



.PHONY: up
up: ## [環境構築] docker-compose環境を起動する
	docker-compose $(DOCKER_COMPOSE_DEFAULT_OPTIONS) \
	up -d

.PHONY: up-prod
up-prod: ## [環境構築] docker-compose環境を起動する
	docker compose $(DOCKER_COMPOSE_PRODUCTION_OPTIONS) \
	up -d

.PHONY:build
build: ## [環境構築] docker-compose環境をbuildする(キャッシュ利用)
	docker-compose $(DOCKER_COMPOSE_DEFAULT_OPTIONS) \
	build

.PHONY: clean-build
clean-build: ## [環境構築] docker-compose環境をclean buildする
	docker-compose $(DOCKER_COMPOSE_DEFAULT_OPTIONS) \
	build --no-cache

.PHONY: clean-build-prod
clean-build-prod: ## [環境構築] docker-compose環境をclean buildする
	docker compose $(DOCKER_COMPOSE_PRODUCTION_OPTIONS) \
	build --no-cache

.PHONY: down
down: ## [環境構築] dockerイメージを削除し,docker-compose環境を停止する
	docker-compose $(DOCKER_COMPOSE_DEFAULT_OPTIONS) \
	down \
	--rmi all --volumes --remove-orphans

.PHONY: down-prod
down-prod: ## [環境構築] dockerイメージを削除し,docker-compose環境を停止する
	docker compose $(DOCKER_COMPOSE_PRODUCTION_OPTIONS) \
	down \
	--rmi all --volumes --remove-orphans

.PHONY: fclean
fclean:down ## [環境構築] マウントしたデータを削除,またdockerイメージも削除する

.PHONY: fclean-prod
fclean-prod:down-prod ## [環境構築] マウントしたデータを削除,またdockerイメージも削除する

.PHONY: re
re:fclean clean-build up sleep init-db ## [環境構築] 完全に初期化した状態でdocker環境を立ち上げる

.PHONY: re-prod
re-prod:fclean-prod clean-build-prod up-prod ## [環境構築] 完全に初期化した状態でdocker環境を立ち上げる

.PHONY: log
log: ## [監視] docker-compose環境のログを確認する
	docker-compose $(DOCKER_COMPOSE_DEFAULT_OPTIONS) \
	logs -f

.PHONY: log-prod
log-prod: ## [監視] docker-compose環境のログを確認する
	docker compose $(DOCKER_COMPOSE_PRODUCTION_OPTIONS) \
	logs -f

.PHONY: ps
ps: ## [監視] docker-compose環境のコンテナを確認する
	docker-compose $(DOCKER_COMPOSE_DEFAULT_OPTIONS) \
	ps

.PHONY: attach-api
attach-api: ## [デバッグ] dockerのserverコンテナにアクセスする
	docker-compose $(DOCKER_COMPOSE_DEFAULT_OPTIONS) \
	exec -it $(API_CONTAINER_NAME) \
	/bin/bash

.PHONY: attach-db
attach-db: ## [デバッグ] dockerのdbコンテナにアクセスする
	docker-compose $(DOCKER_COMPOSE_DEFAULT_OPTIONS) \
	exec -it  $(DB_CONTAINER_NAME) \
	psql -Uuser -dcopalettedb

.PHONY: lint
lint: go-lint ## [リンター] プロジェクトのコードを整形する

.PHONY: go-lint
go-lint: ## [リンター] Goのコードを整形する
	gofmt -l -w .

##FIXME: dockerコンテナで動かしたい
.PHONY: migrate-new
migrate-new: ## [マイグレーション] 新しいmigrationファイルを作成する, 引数にcを使う
	cd db/postgres && sql-migrate new -env="develop" $c

.PHONY: migrate-up
migrate-up: ## [マイグレーション] 既存のmigrationsを適用する
	docker-compose $(DOCKER_COMPOSE_DEFAULT_OPTIONS) \
	run --rm migrate \
	bash -c 'cd /db/postgres && sql-migrate up -env="develop"'

.PHONY: migrate-up-prod
migrate-up-prod: ## [マイグレーション] 既存のmigrationsを適用する
	docker exec -it copalette-api bash && cd db/postgres && sql-migrate up -env="production" 
	exit

.PHONY: migrate-down
migrate-down: ## [マイグレーション] migrationsを一つ戻す
	cd db/postgres && sql-migrate down -env="develop" 

.PHONY: migrate-status
migrate-status: ## [マイグレーション] 現在のmigrationsの状況を確認する
	cd db/postgres && sql-migrate status -env="develop" 

.PHONY: insert-mock
insert-mock: ## [モック]　モックデータをinsertする
	docker-compose $(DOCKER_COMPOSE_DEFAULT_OPTIONS) \
	exec -it  $(DB_CONTAINER_NAME) \
	psql -Uuser -dcopalettedb -f ./db/postgres/mock/*.sql

.PHONY: init-db
init-db: migrate-up  ## [初期化]　マイグレーションとmockデータのインサートを行う

## FIXME: API化するのが少し怖いので、cmdとかにしたい. なんかうまくできない。dbが勝手に再起動してしまう
.PHONY: reset-data
reset-data: ## [初期化] データを初期化する
	curl 'http://localhost:8080/v1/develop/data/reset/all'

.PHONY: clean-build-test
clean-build-test: ## [テスト] テスト用のdocker-compose環境をbuildする
	docker-compose $(DOCKER_COMPOSE_TEST_OPTIONS) \
	build --no-cache

.PHONY: test
test: clean-build-test  ## [テスト] docker-compose環境でテストをする
	docker-compose $(DOCKER_COMPOSE_TEST_OPTIONS) up -d
	make sleep
	docker-compose $(DOCKER_COMPOSE_TEST_OPTIONS) run --rm migrate bash -c 'cd /db/postgres && sql-migrate up -env="test"'
	docker-compose $(DOCKER_COMPOSE_TEST_OPTIONS) run --rm api go test -p 1 ./...
	docker-compose $(DOCKER_COMPOSE_TEST_OPTIONS) down -v

#FIXME: dockerizeでうまく監視したい
.PHONY: sleep
sleep:
	@sleep 5

.PHONY: help
help: ## [ヘルプ] コマンドの一覧を標示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
