#!/bin/bash

# Update the Lambda@Edge function in the CloudFront distribution
# - ORIGIN_RESPONSE_ARN: オリジンレスポンス用のLambda@Edge関数のARN
# - CLOUDFRONT_DISTRIBUTION_ID: CloudFrontディストリビューションID

### Fetch the latest version of the Lambda@Edge function
# For origin response
ORIGIN_RESPONSE_LATEST_ARN=$(aws lambda list-versions-by-function --function-name ${ORIGIN_RESPONSE_ARN} | jq -r '.Versions[-1].FunctionArn')

# Get CloudFront distribution configuration
CLOUDFRONT_ETAG=$(aws cloudfront get-distribution-config --id ${CLOUDFRONT_DISTRIBUTION_ID} | jq -r '.ETag')

# Edit the CloudFront distribution configuration to include the new Lambda@Edge function
aws cloudfront get-distribution-config --id ${CLOUDFRONT_DISTRIBUTION_ID} | \
  jq '.DistributionConfig' | \
  jq "(.DefaultCacheBehavior.LambdaFunctionAssociations.Items[] | select(.EventType == \"origin-response\") | .LambdaFunctionARN) |= \"${ORIGIN_RESPONSE_LATEST_ARN}\"" | \
  > ./config.json

# Update the CloudFront distribution with the new Lambda@Edge function
aws cloudfront update-distribution --id ${CLOUDFRONT_DISTRIBUTION_ID} --if-match ${CLOUDFRONT_ETAG} --distribution-config file://config.json
