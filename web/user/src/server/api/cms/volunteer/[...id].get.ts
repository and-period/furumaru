import { cmsClient } from '~/server/client/microcms'
import type { VolunteerBlogItemResponse } from '~/types/cms/volunteer'

export default defineEventHandler(async (event) => {
  const id = event.context.params?.id || ''

  if (!id) {
    throw createError({
      statusCode: 404,
      statusMessage: 'Param id is required.',
    })
  }

  const res = await cmsClient
    .get<VolunteerBlogItemResponse>({
      endpoint: 'volunteer',
      contentId: id,
    })
    .catch(e => console.error(e))

  return res
})
