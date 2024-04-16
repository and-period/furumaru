.PHONY: build push deploy

build:
	./bin/build.sh

push:
	aws s3 cp ./dist/app.zip s3://${BUCKET_NAME}/functions/furumaru/cloud-origin-response.zip

deploy:
	aws lambda update-function-code --region us-east-1 --function-name ${FUNCTION_NAME} --s3-bucket ${BUCKET_NAME} --s3-key functions/furumaru/cloud-origin-response.zip
	aws lambda publish-version --region us-east-1 --function-name ${FUNCTION_NAME} --description $(shell date +%F_%T)
