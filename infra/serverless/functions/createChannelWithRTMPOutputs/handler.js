const AWS = require('aws-sdk')
const medialive = new AWS.MediaLive();
const defaultSettings = require('./default.json')
const uuid = require('uuid');

module.exports.createChannelWithRTMPOutputs = async (event) => {
  const params = event;
  console.log(params)

  // Defaultの設定を作成
  defaultSettings.Name = params.Name;
  defaultSettings.RoleArn = params.RoleArn;
  defaultSettings.InputAttachments = params.InputAttachments;
  defaultSettings.Destinations = params.Destinations;
  defaultSettings.EncoderSettings.GlobalConfiguration = params.GlobalConfiguration;

  defaultSettings.EncoderSettings.OutputGroups.forEach((outputGroup) => {
    outputGroup.Outputs.forEach((output) => {
      output.OutputName = uuid.v4();
    });
  });

  const RTMPOutputs = params.RTMPOutputs;
  const RTMPOutputGroup = require('./RTMPOutputGroup.json');
  if (RTMPOutputs.length > 0) {
    for (let i = 0; i < RTMPOutputs.length; i++) {
      defaultSettings.Destinations.push({
        "Id": `RTMP${i}`,
        "Settings": [
          {
            "StreamName": RTMPOutputs[i].StreamKey,
            "Url": RTMPOutputs[i].StreamUrl
          }
        ]
      });
    }
    for (let i = 0; i < RTMPOutputs.length; i++) {
      RTMPOutputGroup.Outputs.push(
        {
          "AudioDescriptionNames": [
            "オーディオ"
          ],
          "CaptionDescriptionNames": [],
          "OutputName": RTMPOutputs[i].Name,
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
    }
  }
  defaultSettings.EncoderSettings.OutputGroups.push(RTMPOutputGroup);
  console.log(defaultSettings);
  try {
    const data = await medialive.createChannel(defaultSettings).promise();
    console.log(data);
    return data;
  } catch (err) {
    console.error(`Error creating channel: ${err.message}`, err);
    return err;
  }
}