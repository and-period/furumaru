import { cmsClient } from '~/server/client/microcms'

export default defineEventHandler(async (event) => {
  const id = event.context.params.id

  if (!id) {
    throw createError({
      statusCode: 404,
      statusMessage: 'Param id is required.',
    })
  }

  const res = await cmsClient
    .get({
      endpoint: 'volunteer',
      contentId: id,
    })
    .catch(e => console.error(e))

  return res
})
