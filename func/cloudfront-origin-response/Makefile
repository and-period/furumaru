.PHONY: build push

build: build
	yarn build

push:
	aws s3 cp ./app.zip s3://${BUCKET_NAME}/functions/furumaru/image-converter.zip

deploy: build push
	aws lambda update-function-code --region us-east-1 --function-name ${FUNCTION_NAME} --s3-bucket ${BUCKET_NAME} --s3-key functions/furumaru/image-converter.zip
	aws lambda publish-version --region us-east-1 --function-name ${FUNCTION_NAME} --description $(shell date +%F_%T)
