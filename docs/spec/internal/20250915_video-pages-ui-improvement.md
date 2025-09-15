# 動画管理ページUI/UX改善

| 項目 | 内容 |
|----|------|
| 機能 | 管理者ポータルの動画管理ページのUI/UX改善 |
| 対象ページ | `/videos`、`/videos/new`、`/videos/[id]` |

## 仕様

動画管理ページのUI/UXをスケジュール管理ページの設計パターンと統一し、管理者の操作性を向上させる。

## 設計概要

既存の動画管理ページを以下の観点で改善：

1. **統一されたデザインパターン**: スケジュールページと同様のレイアウト・コンポーネント構成
2. **テンプレートコンポーネント分離**: ページロジックとUIロジックの責任分離
3. **向上したユーザビリティ**: 視認性とアクセシビリティの改善
4. **一貫したページネーション**: 統一されたページング処理

## 設計詳細

### Web

#### 変更対象ファイル

##### 一覧ページ
- **新規作成**: `web/admin/src/components/templates/VideoList.vue`
- **変更**: `web/admin/src/pages/videos/index.vue`

##### 新規作成・編集ページ
- **新規作成**: `web/admin/src/components/templates/VideoNew.vue`
- **変更**: `web/admin/src/components/templates/VideoEdit.vue`
- **変更**: `web/admin/src/pages/videos/new.vue`
- **変更**: `web/admin/src/pages/videos/[id]/index.vue`

#### UI改善内容

##### 1. VideoListテンプレートコンポーネント作成

**主な機能**:
- スケジュールページと統一されたヘッダーデザイン
- 動画専用のアイコン（PlayCircle）とタイトル
- サムネイル画像の改善（16:9比率、再生アイコンオーバーレイ）
- 向上した削除確認ダイアログ

**Props**:
```typescript
interface VideoListProps {
  loading: boolean
  adminType: AdminType
  deleteDialog: boolean
  isAlert: boolean
  alertType: AlertType
  alertText: string
  videos: Video[]
  tableItemsPerPage: number
  tableItemsTotal: number
}
```

**Emits**:
```typescript
interface VideoListEmits {
  'click:update-page': (page: number) => void
  'click:update-items-per-page': (page: number) => void
  'click:row': (videoId: string) => void
  'click:add': () => void
  'click:delete': (videoId: string) => void
  'update:delete-dialog': (v: boolean) => void
}
```

##### 2. データテーブル構造改善

| カラム | 変更前 | 変更後 | 改善内容 |
|--------|--------|--------|----------|
| サムネイル | 50px固定 | 100px（16:9） | 再生アイコンオーバーレイ追加、最左配置 |
| タイトル | テキストのみ | タイトル+説明 | 説明テキスト追加（50文字まで） |
| 公開日時 | 基本表示 | 未設定対応 | "未設定"表示改善 |
| ステータス | 基本チップ | 小サイズチップ | サイズ統一 |

##### 3. 削除ダイアログ改善

**変更前**:
```vue
<v-card-text class="text-h7">
  {{ selectedItem?.title }}を本当に削除しますか？
</v-card-text>
```

**変更後**:
```vue
<v-card-title class="text-h6 py-4">
  動画削除の確認
</v-card-title>
<v-card-text class="pb-4">
  <div class="text-body-1">
    「{{ selectedItem?.title || "" }}」を削除しますか？
  </div>
  <div class="text-body-2 text-medium-emphasis mt-2">
    この操作は取り消せません。
  </div>
</v-card-text>
```

#### ページロジック改善

##### 1. 統一されたデータフェッチパターン

**変更前**: 直接的な`useAsyncData`使用
**変更後**: スケジュールページと統一されたフェッチパターン

```typescript
const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchVideos()
})

watch(pagination.itemsPerPage, (): void => {
  fetchVideos()
})
```

##### 2. エラーハンドリング統一

**統一されたエラーハンドリング**:
```typescript
const handleClickDelete = async (videoId: string): Promise<void> => {
  try {
    loading.value = true
    const video = videoResponse.value?.videos.find((video: Video): boolean => {
      return video.id === videoId
    })
    if (!video) {
      throw new Error(`failed to find video. videoId=${videoId}`)
    }
    await videoStore.deleteVideo(videoId)
    commonStore.addSnackbar({
      message: `${video.title}を削除しました。`,
      color: 'info',
    })
    deleteDialog.value = false
    fetchState.execute()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}
```

## チェックリスト

### 実装開始前

* [x] スケジュールページのデザインパターン分析完了
* [x] VideoListテンプレートコンポーネント設計完了
* [x] 既存動画ストアのページネーション対応確認

### 動作確認

* [ ] 動画一覧表示の正常動作確認
* [ ] ページネーション動作確認
* [ ] 削除機能の正常動作確認
* [ ] レスポンシブデザインの確認
* [ ] アクセシビリティの確認
* [ ] エラーハンドリングの確認

## リリース時確認事項

### リリース順

1. Webのみのリリース（API変更なし）

### リリース制御

特になし

### インフラ設定

特になし

### パフォーマンスチェック

* 画像読み込みのパフォーマンスに影響がないことを確認
* サムネイル表示の最適化（resized images使用）

### その他

* 既存の動画詳細・新規作成ページは変更対象外
* 動画ストアのページネーション対応は既存実装を活用

#### 新規作成・編集ページ改善内容

##### 1. VideoNewテンプレートコンポーネント作成

**主な機能**:
- モダンなカード型レイアウト
- セクション分けされた入力フォーム（基本情報、メディア、関連コンテンツ、公開設定）
- 向上したファイルアップロードUI
- 商品・体験の検索と紐付けダイアログ

##### 2. VideoEditテンプレートコンポーネント改善

**主な変更**:
- タブUIの改善（アイコン付き、見やすい配置）
- コメント管理タブの追加
- ダイアログUIの統一
- フッターボタンのデザイン改善

##### 3. 共通改善点

- **アイコン使用**: 各セクションに適切なアイコンを配置
- **カラースキーム**: プライマリカラーを効果的に使用
- **レスポンシブ対応**: 画面サイズに応じた最適な表示
- **エラーハンドリング**: 統一されたエラー表示
- **ローディング状態**: 明確なローディング表示

## 関連リンク

- [スケジュール管理ページUI改善](./20250115_schedule-experience-ui-improvement.md)
- [商品管理ページUI改善](./20250915_product-pages-ui-improvement.md)
- [生産者管理ページUI改善](./20250915_producer-pages-ui-improvement.md)