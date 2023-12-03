const { MediaStoreDataClient, DeleteObjectCommand, ListItemsCommand } = require("@aws-sdk/client-mediastore-data");

exports.deleteMediaStoreObjects = async (event) => {
  console.log(event);
  const client = new MediaStoreDataClient({
    region: 'ap-northeast-1',
    endpoint: event.CheckContainerExistsResult.Container.Endpoint,
  });

  const listCommand = new ListItemsCommand({
    Path: `/test`,
  });
  const listResponse = await client.send(listCommand);
  console.log(listResponse);

  for (const item of listResponse.Items) {
    const deleteInput = {
      Path: `/test/${item.Name}`,
    };
    console.log(deleteInput);
    const deleteCommand = new DeleteObjectCommand(deleteInput);
    await client.send(deleteCommand);
  }

  return { message: 'All items deleted' };
}
