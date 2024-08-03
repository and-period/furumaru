import { cmsClient } from '~/server/client/microcms'
import type { VolunteerBlogListResponse } from '~/types/cms/volunteer'

export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  const { MICRO_CMS_DOMAIN, MICRO_CMS_API_KEY } = useRuntimeConfig()

  const offset = query.offset ? Number(query.offset) : 0
  const limit = query.limit ? Number(query.limit) : 20

  const res = await cmsClient(MICRO_CMS_DOMAIN, MICRO_CMS_API_KEY)
    .getList<VolunteerBlogListResponse>({
      endpoint: 'volunteer',
      queries: {
        limit,
        offset,
      },
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
