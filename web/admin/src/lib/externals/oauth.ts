interface FetchTokenResponse {
  access_token: string
  refresh_token: string
  id_token: string
  token_type: string
  expires_in: number
}

export interface OAuth {
  accessToken: string
  refreshToken: string
  idToken: string
  tokenType: string
  expiresIn: number
}

const config = useRuntimeConfig()

/**
 * トークン発行のための非同期関数
 */
export async function fetchOAuthToken(code: string, redirectUri: string): Promise<OAuth> {
  if (code === '' || redirectUri === '') {
    throw new Error('code or redirectUri is empty.')
  }

  const params = new URLSearchParams({
    grant_type: 'authorization_code',
    client_id: config.public.COGNITO_CLIENT_ID || '',
    redirect_uri: redirectUri,
    code,
  })

  const out = await $fetch<FetchTokenResponse>(
    `https://${config.public.COGNITO_AUTH_DOMAIN}/oauth2/token`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: params.toString(),
    },
  )

  return {
    accessToken: out.access_token,
    refreshToken: out.refresh_token,
    idToken: out.id_token,
    tokenType: out.token_type,
    expiresIn: out.expires_in,
  }
}
