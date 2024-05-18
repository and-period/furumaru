const AWS = require('aws-sdk')
const medialive = new AWS.MediaLive();
const defaultSettings = require('./default.json')
const uuid = require('uuid');

module.exports.createChannelWithRTMPOutputs = async (event) => {
  const params = event;
  console.log(params)

  const settings = JSON.parse(JSON.stringify(defaultSettings));
  settings.Name = params.Name;
  settings.RoleArn = params.RoleArn;
  settings.InputAttachments = params.InputAttachments;
  settings.Destinations = params.Destinations;
  settings.EncoderSettings.GlobalConfiguration = params.GlobalConfiguration;

  settings.EncoderSettings.OutputGroups.forEach((outputGroup) => {
    outputGroup.Outputs.forEach((output) => {
      output.OutputName = uuid.v4();
    });
  });

  const RTMPOutputs = params.RTMPOutputs;
  const RTMPOutputGroup = require('./RTMPOutputGroup.json');
  if (RTMPOutputs.length > 0) {
    RTMPOutputs.forEach((output, i) => {
      settings.Destinations.push({
        "Id": `RTMP${i}`,
        "Settings": [
          {
            "StreamName": output.StreamKey,
            "Url": output.StreamUrl
          }
        ]
      });
    });
    RTMPOutputs.forEach((output, i) => {
      RTMPOutputGroup.Outputs.push(
        {
          "AudioDescriptionNames": [
            "オーディオ"
          ],
          "CaptionDescriptionNames": [],
          "OutputName": output.Name,
          "OutputSettings": {
            "RtmpOutputSettings": {
              "CertificateMode": "VERIFY_AUTHENTICITY",
              "ConnectionRetryInterval": 2,
              "Destination": {
                "DestinationRefId": `RTMP${i}`
              },
              "NumRetries": 10
            }
          },
          "VideoDescriptionName": "動画"
        }
      );
    });
  }
  settings.EncoderSettings.OutputGroups.push(RTMPOutputGroup);
  console.log(settings);
  try {
    const data = await medialive.createChannel(settings).promise();
    console.log(data);
    return data;
  } catch (err) {
    console.error(`Error creating channel: ${err.message}`, err);
    return err;
  }
}