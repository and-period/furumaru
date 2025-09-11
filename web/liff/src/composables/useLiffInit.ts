import liff from '@line/liff';

export const useLiffInit = () => {
  const init = async (liffId: string | undefined) => {
    if (!liffId) {
      console.error('Please set LIFF_ID in .env file');
      return;
    }

    try {
      console.log(liffId);
      await liff.init({ liffId });
    }
    catch (error) {
      console.error('LIFF init failed', error);
    }

    if (!liff.isLoggedIn()) {
      const redirectUri = window.location.href;
      liff.login({ redirectUri });
      return;
    }
  };

  return { init };
};
