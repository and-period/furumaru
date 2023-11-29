const AWS = require('aws-sdk');
const cloudfront = new AWS.CloudFront();

exports.disableDistribution = async (event) => {
    console.log(event);
    const distributionConfig = event.GetDistributionConfigResult.DistributionConfig;
    const etag = event.GetDistributionConfigResult.ETag;
    distributionConfig.Enabled = false;
    distributionConfig.Origins.Items[0].CustomOriginConfig.HTTPPort = 80;
    distributionConfig.Origins.Items[0].CustomOriginConfig.HTTPSPort = 443;

    delete distributionConfig.Origins.Items[0].CustomOriginConfig.HttpPort;
    delete distributionConfig.Origins.Items[0].CustomOriginConfig.HttpsPort;
    delete distributionConfig.ViewerCertificate.SslSupportMethod;

    const params = {
        DistributionConfig: distributionConfig,
        Id: event.DistributionId,
        IfMatch: etag
    };
    console.log(params);
    const result = await cloudfront.updateDistribution(params).promise();
    console.log(result);
    return event.Id;
}
