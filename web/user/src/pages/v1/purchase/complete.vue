<script setup lang="ts">
const router = useRouter()

const handleBackTopPageButton = () => {
  router.push('/')
}
</script>
<template>
  <div class="text-main">
    <div class="hidden md:block">
      <div class="flex items-center justify-center gap-x-[40px]">
        <div class="text-xl font-bold tracking-[2px]">
          <p>ご注文ありがとうございます！</p>
          <p>ご購入手続きが完了しました。</p>
        </div>
        <img src="/img/purchase/complete.svg" />
      </div>
    </div>
    <div class="mt-[40px] block md:hidden">
      <div class="flex justify-center">
        <img src="/img/purchase/complete.svg" />
      </div>
      <div class="text-[16px] font-bold tracking-[1.6px]">
        <p class="mt-[40px] flex justify-center">ご注文ありがとうございます！</p>
        <p class="mt-2 flex justify-center">ご購入手続きが完了しました。</p>
      </div>
    </div>
    <div class="mt-[40px] flex flex-col items-center justify-center md:mt-20">
      <p class="text-[14px] font-medium tracking-[2px] md:text-xl">
        お客様の注文番号は「1111111」です
      </p>
      <div
        class="flex flex-col items-center justify-center px-[15px] text-[14px] font-medium md:px-0 md:text-[16px]"
      >
        <p class="mt-[40px] md:mt-10">
          ご登録いただいたメールアドレスに確認メールをお送りしましたのでご確認ください。
        </p>
        <p>
          メールが届かない場合は上記の注文番号をお書き添えの上、こちらからお問い合わせください。
        </p>
        <p>
          注文の詳細は、アカウントページの注文履歴からもご確認いただけます。
        </p>
      </div>
      <div class="mb-4 mt-[40px] flex w-full justify-center px-16 md:mt-20 md:px-0">
        <button
          class="w-[400px] bg-main py-2 tracking-[1.6px] text-white"
          @click="handleBackTopPageButton"
        >
          TOPページへ
        </button>
      </div>
    </div>
  </div>
</template>
