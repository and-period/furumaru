import type { Liff } from '@line/liff';
import { useAuthStore } from '~/stores/auth';
import { useShoppingCartStore } from '~/stores/shopping';
import { useUserStore } from '~/stores/user';
import { ResponseError } from '~/types/api/facility';

export default defineNuxtRouteMiddleware(async (to) => {
  if (import.meta.server) {
    // サーバーサイドでは処理しない
    return;
  }

  // 認証はデフォルト必須。明示的に public: true または requiresAuth: false のときだけスキップ
  const isPublic = to.meta?.public === true;
  if (isPublic) {
    return;
  }

  const nuxtApp = useNuxtApp();
  const liff = nuxtApp.$liff as Liff;
  await nuxtApp.$liffReady;

  if (!liff.isLoggedIn()) {
    // LINEログインを強制（リダイレクト先は現URL）
    liff.login({ redirectUri: window.location.href });
    return;
  }

  const authStore = useAuthStore();
  const shoppingCartStore = useShoppingCartStore();
  const userStore = useUserStore();

  const facilityId = String(to.params.facilityId ?? '');

  if (!facilityId) {
    // No facility context, skip bootstrap in middleware
    return;
  }

  if (authStore.isAuthenticated) {
    return;
  }

  const idToken = liff.getIDToken();
  if (!idToken) {
    return;
  }

  try {
    await authStore.signIn(idToken);
    await shoppingCartStore.getCart();
    const accessToken = authStore.token!.accessToken;
    await userStore.fetchMe(facilityId, accessToken);
  }
  catch (err) {
    if (err instanceof ResponseError) {
      if (err.response.status === 404) {
        const path = `/${facilityId}/checkin/new`;
        return navigateTo(path);
      }
      if (err.response.status === 401) {
        liff.logout();
        return;
      }
    }

    console.error('[auth.global] bootstrap failed:', err);
  }
});
