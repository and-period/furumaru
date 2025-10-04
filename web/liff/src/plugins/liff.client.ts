import liff from '@line/liff';

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig();

  // Expose a single init promise for consumers to await
  const liffReady = (async () => {
    try {
      const liffId = config.public?.LIFF_ID as string | undefined;
      if (!liffId) {
        console.warn('[LIFF] Missing public.LIFF_ID. Skipping liff.init');
        return;
      }
      await liff.init({ liffId });
    }
    catch (err) {
      console.error('[LIFF] init failed:', err);
    }
  })();

  return {
    provide: {
      liff,
      liffReady,
    },
  };
});
