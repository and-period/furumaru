exports.getExistsChannelIndex = async (event) => {
  console.log(event);
  const channelNames = event.ChannelNames;
  const channelName = event.ChannelName;
  const index = channelNames.indexOf(channelName);

  return index;
}
