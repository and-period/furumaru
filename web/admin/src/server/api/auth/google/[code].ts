import { Buffer } from 'buffer'

export default defineEventHandler(async (event): Promise<void> => {
  console.log('debug', event)

  const token = Buffer.from(`${process.env.COGNITO_CLIENT_ID}:`).toString('base64')

  const res = await $fetch(
    `https://${process.env.COGNITO_AUTH_DOMAIN}/oauth2/token`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
        'Authorization': `Basic ${token}`,
      },
      body: {
        grant_type: 'authorization_code',
        client_id: process.env.COGNITO_CLIENT_ID,
        redirect_uri: process.env.OAUTH_GOOGLE_REDIRECT_URI,
        code: event.context.params?.code || '',
      },
    },
  ).catch((err) => {
    console.error('error', err)
  })

  console.log('debug', res)
})
