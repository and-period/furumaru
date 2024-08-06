import { cmsClient } from '~/server/client/microcms'
import type { VolunteerBlogItemResponse } from '~/types/cms/volunteer'

export default defineEventHandler(async (event) => {
  const { MICRO_CMS_DOMAIN, MICRO_CMS_API_KEY } = useRuntimeConfig()

  const id = event.context.params?.id || ''

  if (!id) {
    throw createError({
      statusCode: 404,
      statusMessage: 'Param id is required.',
    })
  }

  const res = await cmsClient(MICRO_CMS_DOMAIN, MICRO_CMS_API_KEY)
    .get<VolunteerBlogItemResponse>({
      endpoint: 'volunteer',
      contentId: id,
    })
    .catch((e) => {
      console.log(e)
      throw createError({
        statusCode: 500,
        statusMessage: 'Internal Server Error',
      })
    })

  return res
})
