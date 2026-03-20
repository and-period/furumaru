export default defineSitemapEventHandler(() => {
  return [
    {
      loc: '/',
      // automatically creates: /en/about-us, /fr/about-us, etc.
      _i18nTransform: true,
    },
  ]
})
