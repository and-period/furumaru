##################################################
# Container Commands - Run All
##################################################
.PHONY: setup build install start stop down remove logs

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

start: proto migrate
	docker-compose up --remove-orphans

stop:
	docker-compose stop

down:
	docker-compose down

remove:
	docker-compose run mysql_test bash -c "echo 'DROP DATABASE migrations;' | mysql -u root -p12345678 "
	docker-compose down --rmi all --volumes --remove-orphans
	rm -r ./tmp/** && touch ./tmp/.keep

logs:
	docker-compose logs

##################################################
# Container Commands - Run Container Group
##################################################
.PHONY: start-web start-api start-swagger start-test

start-web:
	docker-compose up user_web admin_web

start-api: migrate
	docker-compose up user_gateway admin_gateway mysql_test

start-swagger:
	docker-compose up swagger_generator swagger_user swagger_admin

start-test:
	docker-compose up mysql_test

##################################################
# Container Commands - Single
##################################################
.PHONY: proto swagger migrate

swagger:
	docker-compose run --rm swagger_generator yarn generate
	docker-compose run --rm user_web yarn lintfix
	docker-compose run --rm admin_web yarn lintfix

migrate:
	docker-compose up -d mysql_test
	docker-compose exec mysql_test bash -c "until mysqladmin ping -u root -p12345678 2> /dev/null; do echo 'waiting for ping response..'; sleep 3; done; echo 'mysql_test is ready!'"
	$(MAKE) proto
	docker-compose run --rm executor sh -c "cd ./hack/database-migrate; go run ./main.go -db-host=mysql_test -db-port=3306"
	docker-compose down mysql_test
