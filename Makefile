.PHONY: help

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

##################################################
# Container Commands - Run All
##################################################
.PHONY: setup install build start stop down remove logs

setup: build install swagger ## 初回の環境構築用
	if [ ! -f $(PWD)/.env ]; then \
		cp $(PWD)/.env.temp $(PWD)/.env; \
	fi

install: migrate ## ライブラリインストール/アップデート用
	docker compose run --rm swagger_generator yarn
	docker compose run --rm user_web yarn
	docker compose run --rm admin_web yarn

build: ## コンテナイメージのビルド用
	docker compose build --parallel

start: migrate ## 全てのシステムを起動
	docker compose up --remove-orphans

stop: ## 全てのシステムを停止
	docker compose stop

down: ## 全てのシステムを停止/削除
	docker compose down

remove: ## 全てのシステムを停止/削除(データも含めて)
	docker compose down --rmi all --volumes --remove-orphans
	rm -rf ./tmp/data ./tmp/logs

logs: ## コンテナログを確認
	docker compose logs

##################################################
# Container Commands - Run Container Group
##################################################
.PHONY: start-user start-admin start-web start-api start-swagger

start-user: ## 購入者関連システムの起動
	docker compose up user_web user_gateway mysql

start-admin: ## 管理者関連システムの起動
	docker compose up admin_web admin_gateway mysql

start-web: ## フロントエンド関連の起動
	docker compose up user_web admin_web

start-api: ## バックエンド関連の起動
	docker compose up user_gateway admin_gateway mysql mysql_test

start-swagger: ## API仕様書用システムの起動
	docker compose up swagger_generator swagger_user swagger_admin

##################################################
# Container Commands - Single
##################################################
.PHONY: swagger migrate

swagger: ## API仕様書の生成
	docker-compose run --rm swagger_generator make build
	cd ./web/admin && yarn format
	cd ./web/user && yarn format

migrate: ## データベースにDDLを適用
	docker compose up -d mysql mysql_test
	docker compose exec mysql bash -c "until mysqladmin ping -u root -p12345678 2> /dev/null; do echo 'waiting for ping response..'; sleep 3; done; echo 'mysql_test is ready!'"
	docker compose exec mysql_test bash -c "until mysqladmin ping -u root -p12345678 2> /dev/null; do echo 'waiting for ping response..'; sleep 3; done; echo 'mysql_test is ready!'"
	docker compose run --rm executor sh -c "cd ./hack/database-migrate; go run ./main.go -db-host=mysql -db-port=3306"
	docker compose run --rm executor sh -c "cd ./hack/database-migrate; go run ./main.go -db-host=mysql_test -db-port=3306"
	docker compose stop mysql mysql_test
