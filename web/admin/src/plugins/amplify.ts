import { Amplify, Analytics } from 'aws-amplify'

const enabledAnalytics: boolean = process.env.AWS_PINPOINT_APP_ID !== ''

Amplify.configure({
  Auth: {
    identityPoolId: process.env.AWS_COGNITO_IDENTITY_POOL_ID,
    region: process.env.AWS_REGION,
  },
  Analytics: {
    disabled: !enabledAnalytics,
    autoSessionRecord: true,
    AWSPinpoint: {
      appId: process.env.AWS_PINPOINT_APP_ID,
      region: process.env.AWS_REGION,
      mandatorySignIn: false,
      flushInterval: 5000, // 5s
      resendLimit: 5,
    },
  },
})

if (enabledAnalytics) {
  Analytics.record(process.env.APP_NAME || '')
}
