.PHONY: build push deploy

S3_OBJECT_KEY = functions/furumaru/cognito-migrate-user-pool.zip

build:
	./bin/build.sh

push:
	aws s3 cp ./dist/app.zip s3://${S3_BUCKET_NAME}/${S3_OBJECT_KEY}

deploy:
	aws lambda update-function-code --region ap-northeast-1 --function-name ${FUNCTION_NAME} --s3-bucket ${S3_BUCKET_NAME} --s3-key ${S3_OBJECT_KEY} | jq .
	aws lambda wait function-updated --region ap-northeast-1 --function-name ${FUNCTION_NAME}
	aws lambda publish-version --region ap-northeast-1 --function-name ${FUNCTION_NAME} --description $(shell date +%F_%T) | jq .
