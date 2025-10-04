import type { Liff } from '@line/liff';

declare module '#app' {
  interface NuxtApp {
    $liff: Liff;
    $liffReady: Promise<void>;
  }
}

declare module 'vue' {
  interface ComponentCustomProperties {
    $liff: Liff;
    $liffReady: Promise<void>;
  }
}

export {};
