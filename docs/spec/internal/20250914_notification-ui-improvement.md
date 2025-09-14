# Notification UI Improvement

| 項目 | 値 |
|----|--|
| 機能 | お知らせ管理UIの改善・統一化 |

## 仕様

お知らせの新規作成・編集コンポーネントの設計改善と共通化を行う。

## 設計概要

現在のNotificationNew.vueとNotificationEdit.vueは以下の問題を抱えている：
- コード重複（共通ロジック、定数定義）
- 責務の混在（UIコンポーネント内のビジネスロジック）
- 一貫性の欠如（実装パターンの相違）

これらを解決するため、以下のリファクタリングを実施する：
1. 共通ロジックのComposable化
2. 定数定義の共通化
3. プロモーション表示コンポーネントの分離
4. フォーム検証ロジックの統一

## 設計詳細

### Web

#### ファイル構成

```
web/admin/src/
├── components/
│   ├── molecules/
│   │   ├── NotificationPromotionDisplay.vue   # NEW: プロモーション情報表示
│   │   └── NotificationFormFields.vue         # NEW: フォーム共通フィールド
│   └── templates/
│       ├── NotificationNew.vue               # REFACTOR
│       └── NotificationEdit.vue              # REFACTOR
├── composables/
│   ├── useNotificationForm.ts                # NEW: フォーム共通ロジック
│   └── useNotificationDisplay.ts             # NEW: 表示用ヘルパー
└── constants/
    └── notification.ts                       # NEW: 定数定義
```

#### 実装方針

**1. 共通定数の分離 (`constants/notification.ts`)**
```typescript
export const NOTIFICATION_TYPES = [
  { title: 'システム関連', value: NotificationType.NotificationTypeSystem },
  // ...
]

export const NOTIFICATION_TARGETS = [
  { title: 'ユーザー', value: NotificationTarget.NotificationTargetUsers },
  // ...
]
```

**2. 表示ロジックのComposable化 (`useNotificationDisplay.ts`)**
```typescript
export const useNotificationDisplay = () => {
  const getDateTime = (unixTime: number): string => { /* ... */ }
  const getPromotionTerm = (promotion: Promotion): string => { /* ... */ }
  const getPromotionDiscount = (promotion: Promotion): string => { /* ... */ }
  
  return {
    getDateTime,
    getPromotionTerm,
    getPromotionDiscount,
  }
}
```

**3. フォームロジックのComposable化 (`useNotificationForm.ts`)**
```typescript
export const useNotificationForm = (
  formData: Ref<CreateNotificationRequest | UpdateNotificationRequest>
) => {
  const timeDataValue = computed({ /* ... */ })
  const onChangePublishedAt = () => { /* ... */ }
  
  return {
    timeDataValue,
    onChangePublishedAt,
  }
}
```

**4. プロモーション表示の分離 (`NotificationPromotionDisplay.vue`)**
- プロモーション情報表示の専用コンポーネント
- 表示モード（readonly/editable）の対応
- 統一されたスタイリング

**5. フォームフィールドの共通化 (`NotificationFormFields.vue`)**
- 共通フィールド（公開範囲、投稿日時、本文、備考）
- バリデーション表示の統一
- 条件分岐の最小化

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

* [ ] お知らせ新規作成機能の動作確認
* [ ] お知らせ編集機能の動作確認
* [ ] プロモーション情報表示の確認
* [ ] バリデーション機能の確認
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

- [Web Admin Notification Components](../../web/admin/src/components/templates/)
- [Vuetify Design Guidelines](https://vuetifyjs.com/en/introduction/why-vuetify/)