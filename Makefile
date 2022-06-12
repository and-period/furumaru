##################################################
# Container Commands - Run All
##################################################
.PHONY: setup install build start stop down remove logs

setup: build install swagger
	if [ ! -f $(PWD)/.env ]; then \
		cp $(PWD)/.env.temp $(PWD)/.env; \
	fi

install: migrate
	docker-compose run --rm swagger_generator yarn
	docker-compose run --rm user_web yarn
	docker-compose run --rm admin_web yarn

build:
	docker-compose build --parallel

start: migrate
	docker-compose up --remove-orphans

stop:
	docker-compose stop

down:
	docker-compose down

remove:
	docker-compose down --rmi all --volumes --remove-orphans
	rm -rf ./tmp/data ./tmp/logs

logs:
	docker-compose logs

##################################################
# Container Commands - Run Container Group
##################################################
.PHONY: start-user start-admin start-web start-api start-swagger

start-user:
	docker-compose up user_web user_gateway mysql

start-admin:
	docker-compose up admin_web admin_gateway mysql

start-web:
	docker-compose up user_web admin_web

start-api:
	docker-compose up user_gateway admin_gateway mysql mysql_test

start-swagger:
	docker-compose up swagger_generator swagger_user swagger_admin

##################################################
# Container Commands - Single
##################################################
.PHONY: swagger migrate

swagger:
	docker-compose run --rm swagger_generator yarn generate
	docker-compose run --rm user_web yarn lintfix
	docker-compose run --rm admin_web yarn lintfix

migrate:
	docker-compose up -d mysql mysql_test
	docker-compose exec mysql bash -c "until mysqladmin ping -u root -p12345678 2> /dev/null; do echo 'waiting for ping response..'; sleep 3; done; echo 'mysql_test is ready!'"
	docker-compose exec mysql_test bash -c "until mysqladmin ping -u root -p12345678 2> /dev/null; do echo 'waiting for ping response..'; sleep 3; done; echo 'mysql_test is ready!'"
	docker-compose run --rm executor sh -c "cd ./hack/database-migrate; go run ./main.go -db-host=mysql -db-port=3306"
	docker-compose run --rm executor sh -c "cd ./hack/database-migrate; go run ./main.go -db-host=mysql_test -db-port=3306"
	docker-compose stop mysql mysql_test
