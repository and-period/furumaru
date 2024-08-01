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
    .catch((e) => {
      console.log('エラーです')
      // if (e.response.status === 404) {
      throw createError({
        statusCode: 404,
        statusMessage: 'Not found.',
      })
      // } else {
      //   throw createError({
      //     statusCode: 500,
      //     statusMessage: "Internal Server Error.",
      //   });
      // }
    })

  return res
})
