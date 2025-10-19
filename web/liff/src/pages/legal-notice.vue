<script setup lang="ts">
const router = useRouter();

definePageMeta({
  title: '特定商取引法に基づく表記',
  layout: 'init',
  public: true,
});
</script>

<template>
  <div class="min-h-screen bg-white relative mb-[24px] overflow-auto">
    <!-- Navigation header for mobile -->
    <div class="top-0 z-10 bg-white border-b border-gray-200 px-4 py-3 fixed w-full">
      <div class="flex items-center">
        <button
          class="flex items-center text-gray-600 hover:text-gray-800 cursor-pointer"
          @click="router.back()"
        >
          <svg
            class="w-6 h-6 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15 19l-7-7 7-7"
            />
          </svg>
          戻る
        </button>
      </div>
    </div>

    <div class="mt-12 py-2 px-4">
      <h2 class="font-semibold font-inter text-center w-full">
        特定商取引法に基づく表記
      </h2>
      <!-- 特定商取引法に基づく表記本文 -->
      <div class="mt-4 text-gray-700 space-y-6">
        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            販売業者
          </h3>
          <p class="text-sm leading-relaxed">
            株式会社アンドピリオド
          </p>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            販売責任者
          </h3>
          <p class="text-sm leading-relaxed">
            藤中拓弥
          </p>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            ホームページ
          </h3>
          <p class="text-sm leading-relaxed">
            <a
              href="https://www.and-period.co.jp/"
              class="text-blue-600 hover:text-blue-800 underline"
              target="_blank"
              rel="noopener noreferrer"
            >
              https://www.and-period.co.jp/
            </a>
          </p>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            所在地
          </h3>
          <p class="text-sm leading-relaxed">
            広島県豊田郡大崎上島町中野5748-2
          </p>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            電話番号
          </h3>
          <p class="text-sm leading-relaxed">
            050-3700-7904
          </p>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            メールアドレス
          </h3>
          <p class="text-sm leading-relaxed">
            <a
              href="mailto:support@and-period.co.jp"
              class="text-blue-600 hover:text-blue-800 underline"
            >
              support@and-period.co.jp
            </a>
          </p>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            URL
          </h3>
          <p class="text-sm leading-relaxed">
            <a
              href="https://www.furumaru.and-period.co.jp/"
              class="text-blue-600 hover:text-blue-800 underline"
              target="_blank"
              rel="noopener noreferrer"
            >
              https://www.furumaru.and-period.co.jp/
            </a>
          </p>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            引き渡し時期
          </h3>
          <p class="text-sm leading-relaxed">
            ご注文から7日以内に発送いたします。（商品によっては収穫時期により異なる）
          </p>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            お支払い方法
          </h3>
          <ul class="text-sm leading-relaxed list-disc list-inside space-y-1">
            <li>クレジットカード</li>
            <li>QRコード決済</li>
          </ul>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            返金・交換について
          </h3>
          <p class="text-sm leading-relaxed">
            本サイトで取り扱う商品はすべて食品となるため、返金・交換は一切承っておりません。ただし、お届けした商品に破損、汚損等あった場合や、注文品と違う商品が届いた場合はその限りではありません。
          </p>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            返金・交換の対象となる商品
          </h3>
          <ul class="text-sm leading-relaxed list-disc list-inside space-y-1">
            <li>商品が腐っている</li>
            <li>商品が傷んでいる</li>
            <li>商品に破損がある</li>
            <li>注文品と違う商品が届いた</li>
          </ul>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            返金・交換対応の流れ
          </h3>
          <div class="text-sm leading-relaxed space-y-2">
            <ol class="list-decimal list-inside space-y-1">
              <li>商品の状態がわかるお写真を撮影してください。</li>
              <li>
                お問い合わせフォームより以下を記載の上、商品のお写真を添付して送信してください
                <ul class="list-disc list-inside ml-4 space-y-1 mt-1">
                  <li>注文番号</li>
                  <li>商品の状態</li>
                </ul>
              </li>
              <li>ふるまる運営事務局より今後の対応について1営業日以内に返信いたします。</li>
            </ol>
            <p class="mt-2">
              以上で手続きは完了となります。
            </p>
            <p class="text-xs text-gray-500 mt-2">
              ※宅配ボックスへのお届けや置き配など、一部適用外となる場合もあります。
            </p>
            <div class="bg-gray-50 p-3 rounded-md mt-3">
              <p class="text-xs text-gray-600">
                ＜返金・交換対応になった商品について＞<br>
                商品については対応が確定するまでは破棄せずに保管してください。<br>
                ご迷惑をおかけし申し訳ありませんが何卒ご了承ください。
              </p>
            </div>
          </div>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            商品代金以外に必要な費用
          </h3>
          <ul class="text-sm leading-relaxed list-disc list-inside">
            <li>商品により配送費用</li>
          </ul>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            送料に関して
          </h3>
          <div class="text-sm leading-relaxed space-y-2">
            <p>ふるマルのご注文には一部商品を除き、送料（税込）が箱ごとにかかる場合がございます。また、お届け先によって料金が異なります。</p>
          </div>
        </div>

        <div>
          <h3 class="font-semibold text-gray-900 mb-2">
            返品送料
          </h3>
          <p class="text-sm leading-relaxed">
            返金・交換対象に当てはまる場合は当社負担で配送させて頂きます。
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
