import { cmsClient } from '~/server/client/microcms'

export default defineEventHandler(async (event) => {
  const res = await cmsClient
    .getList({
      endpoint: 'volunteer',
    })
    .catch(e => console.error(e))

  console.log(res)

  return res
})
