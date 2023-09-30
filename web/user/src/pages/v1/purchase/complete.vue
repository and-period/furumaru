<script lang="ts" setup>

</script>

<template>
  <div class="text-main">
  <div class="flex items-center justify-center gap-x-[40px]">
    <div class="text-xl font-bold tracking-[2px]">
      <p> ご注文ありがとうございます！ </p>
      <p> ご購入手続きが完了しました。</p>
    </div>
    <div>
      <img src="/img/purchase/complete.svg">
    </div>
  </div>
  <div class="mt-20 flex flex-col items-center justify-center">
      <p class="text-xl font-medium tracking-[2px]"> お客様の注文番号は「1111111」です </p>
      <div class="flex flex-col items-center justify-center text-[16px] font-medium">
        <p class="mt-10">ご登録いただいたメールアドレスに確認メールをお送りしましたのでご確認ください。</p>
        <p>メールが届かない場合は上記の注文番号をお書き添えの上、こちらからお問い合わせください。</p>
        <p>注文の詳細は、アカウントページの注文履歴からもご確認いただけます。</p>
      </div>
      <div class="mb-4 mt-20 flex w-full justify-center">
          <button onclick="location.href='/'" class="w-[400px] bg-main py-2 tracking-[1.6px] text-white">TOPページへ</button>
        </div>
  </div>
  </div>
</template>
