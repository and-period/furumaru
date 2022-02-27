##################################################
# Container Commands - Run All
##################################################
.PHONY: setup build install start stop down remove logs

setup: build install proto migrate
	if [ ! -f $(PWD)/.env ]; then \
		cp $(PWD)/.env.temp $(PWD)/.env; \
	fi

install:
	echo 'wip...'

build:
	docker-compose build --parallel

start: proto migrate
	docker-compose up --remove-orphans

stop:
	docker-compose stop

down:
	docker-compose down

remove:
	docker-compose down --rmi all --volumes --remove-orphans

logs:
	docker-compose logs

##################################################
# Container Commands - Run Container Group
##################################################
.PHONY: start-api start-test

start-api: proto migrate
	docker-compose up user_gateway user_api

start-test:
	docker-compose up mysql_test

##################################################
# Container Commands - Single
##################################################
.PHONY: proto migrate

proto:
	docker-compose run --rm proto bash -c "cd ./api; make install; make protoc"

migrate:
	docker-compose up -d mysql mysql_test
	docker-compose exec mysql bash -c "until mysqladmin ping -u root -p12345678 2> /dev/null; do echo 'waiting for ping response..'; sleep 3; done; echo 'mysql is ready!'"
	docker-compose exec mysql_test bash -c "until mysqladmin ping -u root -p12345678 2> /dev/null; do echo 'waiting for ping response..'; sleep 3; done; echo 'mysql_test is ready!'"
	$(MAKE) proto
	docker-compose run --rm executor sh -c "cd ./hack/database-migrate; go run ./main.go -db-host=mysql -db-port=3306"
	docker-compose run --rm executor sh -c "cd ./hack/database-migrate; go run ./main.go -db-host=mysql_test -db-port=3306"
	docker-compose down mysql mysql_test
