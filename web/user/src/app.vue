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
    // 本番環境にだけ Meta Pixel Code を仕込む
    ...(config.public.ENVIRONMENT === 'prd'
      ? [
          {
            key: 'meta-picel',
            src: '/meta/pixel-code.js',
            defer: true,
            type: 'text/javascript',
          },
        ]
      : []),
  ],
  // @ts-ignore
  noscript: [
    // 本番環境にだけ Meta Pixel Code を仕込む
    ...(config.public.ENVIRONMENT === 'prd'
      ? [
          {
            key: 'meta-picel-noscirpt',
            innerHTML:
              '<img height="1" width="1" style="display:none" src="https://www.facebook.com/tr?id=610594923400694&ev=PageView&noscript=1"/>',
            tagDuplicateStrategy: 'merge',
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
