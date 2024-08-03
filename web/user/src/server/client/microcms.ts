import { createClient } from 'microcms-js-sdk'

export const cmsClient = (serviceDomain: string, apiKey: string) => {
  return createClient({
    serviceDomain,
    apiKey,
  })
}
