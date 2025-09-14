# Promotion UI Improvement

| 項目 | 値 |
|----|--|
| 機能 | プロモーション管理UIの改善・統一化 |

## 仕様

プロモーションの新規作成・編集コンポーネントの設計改善と共通化を行う。

## 設計概要

現在のPromotionNew.vueとPromotionEdit.vueは以下の問題を抱えている：
- コード重複（共通ロジック、定数定義、バリデーション）
- 責務の混在（UIコンポーネント内のビジネスロジック）
- デザインの一貫性不足（variant統一、日時UI、必須表示）

これらを解決するため、以下のリファクタリングを実施する：
1. 共通ロジックのComposable化
2. 定数定義とバリデーションの共通化
3. 日時入力UIの改善
4. フォーム一貫性の向上

## 設計詳細

### Web

#### ファイル構成

```
web/admin/src/
├── components/
│   ├── molecules/
│   │   └── PromotionPeriodInput.vue           # NEW: 期間入力コンポーネント
│   └── templates/
│       ├── PromotionNew.vue                   # REFACTOR
│       └── PromotionEdit.vue                  # REFACTOR
├── composables/
│   ├── usePromotionForm.ts                    # NEW: フォーム共通ロジック
│   └── usePromotionValidation.ts              # NEW: バリデーション共通ロジック
└── constants/
    └── promotion.ts                           # NEW: 定数定義
```

#### 実装方針

**1. 共通定数の分離 (`constants/promotion.ts`)**
```typescript
export const DISCOUNT_METHODS = [
  { method: '円', value: DiscountType.DiscountTypeAmount },
  { method: '%', value: DiscountType.DiscountTypeRate },
  { method: '送料無料', value: DiscountType.DiscountTypeFreeShipping },
]

export const PROMOTION_STATUS_OPTIONS = [
  { status: '有効', value: true },
  { status: '無効', value: false },
]
```

**2. フォームロジックのComposable化 (`usePromotionForm.ts`)**
```typescript
export const usePromotionForm = (formData: Ref<CreatePromotionRequest | UpdatePromotionRequest>) => {
  const startTimeDataValue = computed({ /* ... */ })
  const endTimeDataValue = computed({ /* ... */ })
  const onChangeStartAt = () => { /* ... */ }
  const onChangeEndAt = () => { /* ... */ }
  const generateRandomCode = () => { /* ... */ }
  
  return {
    startTimeDataValue,
    endTimeDataValue,
    onChangeStartAt,
    onChangeEndAt,
    generateRandomCode,
  }
}
```

**3. バリデーション共通化 (`usePromotionValidation.ts`)**
```typescript
export const usePromotionValidation = (formData: Ref<any>) => {
  const getDiscountErrorMessage = (): string => { /* ... */ }
  
  return {
    getDiscountErrorMessage,
  }
}
```

**4. 期間入力コンポーネント化 (`PromotionPeriodInput.vue`)**
- 開始・終了日時の入力を統一したコンポーネント
- エラー表示の統一
- レスポンシブ対応

**5. フォームデザインの改善**
- 全フィールドを`variant="outlined"`に統一
- `ga-4`での間隔統一
- 必須フィールドに`*`表示追加
- 割引コード生成UIの改善

### API

#### エンドポイント

APIの変更は不要。

#### シーケンス

既存のフローを維持。

## チェックリスト

### 実装開始前

* [x] 既存コンポーネントの問題点分析完了
* [x] リファクタリング方針策定完了
* [ ] 共通化対象の特定完了
* [ ] 破壊的変更がないことを確認

### 動作確認

* [ ] プロモーション新規作成機能の動作確認
* [ ] プロモーション編集機能の動作確認
* [ ] 割引コード自動生成の確認
* [ ] 日時入力・バリデーションの確認
* [ ] レスポンシブデザインの確認

## リリース時確認事項

### リリース順

Web側のみの変更のため、特別な順序指定なし。

### リリース制御

特になし。

### インフラ設定

特になし。

### パフォーマンスチェック

* [ ] バンドルサイズへの影響確認
* [ ] コンポーネント描画パフォーマンス確認

### その他

* [ ] 既存機能の回帰テスト実施
* [ ] TypeScript型安全性の確認

## 関連リンク

- [Web Admin Promotion Components](../../web/admin/src/components/templates/)
- [Vuetify Design Guidelines](https://vuetifyjs.com/en/introduction/why-vuetify/)