LAMBDA_ROLE :=  arn:aws:iam::386661535629:role/service-role/furumaru-stg-lambda

create: build
	aws lambda create-function \
		--function-name ${FUNC_NAME} \
		--zip-file fileb://function.zip \
		--handler main \
		--runtime go1.x \
		--role ${LAMBDA_ROLE}

update: build
	aws lambda update-function-code \
		--function-name ${FUNC_NAME} \
		--zip-file fileb://function.zip

build:
	GOOS=linux go build -o main ${FUNC_DIR}/main.go
	zip function.zip main
