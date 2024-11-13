<script setup lang="ts">
const config = useRuntimeConfig()

useHead({
  script: [
    {
      key: 'newRelic',
      src: `/new-relic/${config.public.ENVIRONMENT}.js`,
      defer: true,
      type: 'text/javascript',
    },
    {
      key: 'clarity',
      src: `/clarity/${config.public.ENVIRONMENT}.js`,
      defer: true,
      type: 'text/javascript',
    },
    // 本番環境にだけ Meta Pixel Code と Google Tag Manager を仕込む
    ...(config.public.ENVIRONMENT === 'prd'
      ? [
          {
            key: 'meta-picel',
            src: '/meta/pixel-code.js',
            defer: true,
            type: 'text/javascript',
          },
          {
            key: 'google-tag-manager',
            src: '/gtm/index.js',
            defer: true,
            type: 'text/javascript',
          },
        ]
      : []),
  ],
  // @ts-ignore
  noscript: [
    // 本番環境にだけ Meta Pixel Code と Google Tag Manager を仕込む
    ...(config.public.ENVIRONMENT === 'prd'
      ? [
          {
            key: 'meta-picel-noscirpt',
            children:
              '<img height="1" width="1" style="display:none" src="https://www.facebook.com/tr?id=610594923400694&ev=PageView&noscript=1"/>',
            tagDuplicateStrategy: 'merge',
          },
          {
            name: 'google-tag-manager-noscript',
            children:
              '<iframe src="https://www.googletagmanager.com/ns.html?id=GTM-5SBNWRJT" height="0" width="0" style="display:none;visibility:hidden"></iframe>',
            tagDuplicateStrategy: 'merge',
            tagPosition: 'bodyOpen',
          },
        ]
      : []),
  ],
})
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
