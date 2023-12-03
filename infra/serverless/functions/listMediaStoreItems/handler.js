
const { MediaStoreDataClient, ListItemsCommand } = require("@aws-sdk/client-mediastore-data"); // CommonJS import

exports.listMediaStoreItems = async (event) => {
  console.log(event);
  const client = new MediaStoreDataClient({
    region: 'ap-northeast-1',
    endpoint: event.CheckContainerExistsResult.Container.Endpoint,
  });
  const input = {
    Path: '/',
    MaxResults: Number(1),
  };
  const command = new ListItemsCommand(input);
  const response = await client.send(command);

  return response;
}
