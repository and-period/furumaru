const AWS = require('aws-sdk')
const medialive = new AWS.MediaLive();
const defaultSettings = require('./default.json')
const uuid = require('uuid');

module.exports.createChannelWithRTMPOutputs = async (event) => {
  console.log(event);
  const settings = event;
  console.log(settings)

  defaultSettings.Name = settings.Name;
  defaultSettings.RoleArn = settings.RoleArn;
  defaultSettings.InputAttachments = settings.InputAttachments;
  defaultSettings.Destinations = settings.Destinations;
  defaultSettings.EncoderSettings.GrobalConfiguration = settings.GrobalConfiguration;

  defaultSettings.EncoderSettings.OutputGroups.forEach((outputGroup) => {
    outputGroup.Outputs.forEach((output) => {
      output.OutputName = uuid.v4();
    });
  });
  console.log(defaultSettings)
  try {
    const data = await medialive.createChannel(defaultSettings).promise();
    console.log(data);
    return data;
  } catch (err) {
    console.log(err);
    return err;
  }
}