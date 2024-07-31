import { cmsClient } from '~/server/client/microcms'
import type { VolunteerBlogListResponse } from '~/types/cms/volunteer'

export default defineEventHandler(async (event) => {
  const query = getQuery(event)

  const offset = query.offset ? Number(query.offset) : 0
  const limit = query.limit ? Number(query.limit) : 20

  const res = await cmsClient
    .getList<VolunteerBlogListResponse>({
      endpoint: 'volunteer',
      queries: {
        limit,
        offset,
      },
    })
    .catch(e => console.error(e))

  return res
})
