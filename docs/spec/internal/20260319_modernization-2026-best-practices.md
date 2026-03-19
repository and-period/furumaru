# 2026年ベストプラクティスに基づくモダナイゼーション

| 項目 | 内容 |
|----|---|
| 機能 | フロントエンド・インフラのツールチェーン・依存ライブラリを2026年現在のベストプラクティスに更新する |

## 仕様

### 背景・課題

リポジトリ全体を調査した結果、Go バックエンド（Go 1.26, slog, golangci-lint v2, errgroup 等）は2026年時点で十分モダンな構成であることを確認した。一方、フロントエンドおよびインフラ周辺に以下の技術的負債・改善機会が特定された。

### 調査結果サマリ

#### Go バックエンド（対応不要）

現時点で十分モダンな構成であり、大規模な改修は不要と判断。

| 項目 | 現状 | 評価 |
|------|------|------|
| Go バージョン | 1.26 | 最新 |
| HTTP フレームワーク | Gin v1.11.0 | 安定・実績あり。stdlib routing (Go 1.22+) への移行は可能だが ROI が低い |
| ORM | GORM v1.31.1 | 安定。sqlc 等への移行は大規模リファクタが必要で ROI が低い |
| ログ | log/slog + sentry-slog + slog-multi | 2026年標準に準拠 |
| テスト | testify + go.uber.org/mock + gotestsum | モダンなパターン |
| Lint | golangci-lint v2 | 最新 |
| エラーハンドリング | errors.Is/As + sentinel errors | 標準的 |
| DI | コンストラクタ + Options パターン | シンプルで十分 |

#### フロントエンド（対応内容）

| # | 項目 | Before | After | 優先度 |
|---|------|--------|-------|--------|
| F1 | user app TypeScript | 4.9.5 | 5.9.3 | **Critical** |
| F2 | user app Vitest | 0.26.3 + coverage-c8 | 3.2.4 + coverage-v8 | **Critical** |
| F3 | user app happy-dom | 8.1.1 | 18.0.1 | **Critical** |
| F4 | user app vue-tsc | 1.8.27 | 2.2.12 | **Critical** |
| F5 | user/admin Sentry | @sentry/vue + @sentry/node v7 | @sentry/nuxt v9 | **High** |
| F6 | LIFF app Nuxt | 3.18.1 | ^4.1.3 | **Medium** |
| F7 | パッケージマネージャ | Yarn v1 (1.22.22) | pnpm 10.x | **Medium** |
| F8 | user app ESLint deps | レガシー deps 残存 | 不要 deps 削除 | **Low** |

※ F8: ESLint は既に Flat Config（eslint.config.mjs）に移行済み。package.json に残っていた `@babel/eslint-parser`, `@nuxtjs/eslint-config-typescript`, `vue-eslint-parser`, `@typescript-eslint/eslint-plugin` 等の不要パッケージを削除。

#### インフラ（対応内容）

| # | 項目 | Before | After | 優先度 |
|---|------|--------|-------|--------|
| I1 | Docker Node.js (開発用) | 16.14.2 / 20.11.1 | 22-alpine | **Critical** |
| I2 | .tool-versions Node.js | 20.19.6 / 20.11.1 | 22.22.1 | **Critical** |
| I3 | Serverless Runtime | nodejs18.x | nodejs22.x | **High** |
| I4 | コンテナセキュリティ | なし | Trivy スキャン in CI | **Low** |
| I5 | Observability 統合 | New Relic + Sentry + CloudWatch（分散） | OpenTelemetry Collector 統合 | **Medium**（将来課題） |

## 設計概要

独立してリリース可能な単位で段階的に改修する。

## 設計詳細

### 変更ファイル一覧

#### Node.js バージョン更新

| ファイル | 変更内容 |
|---------|---------|
| `web/user/.tool-versions` | 20.19.6 → 22.22.1 |
| `web/admin/.tool-versions` | 20.19.6 → 22.22.1 |
| `web/liff/.tool-versions` | 20.19.6 → 22.22.1 |
| `web/shared/.tool-versions` | 20.19.6 → 22.22.1 |
| `func/.tool-versions` | 20.11.1 → 22.22.1 |
| `infra/serverless/.tool-versions` | 20.11.1 → 22.22.1 |
| `infra/docker/web/user/Dockerfile.development` | node:16.14.2-alpine → node:22-alpine |
| `infra/docker/web/admin/Dockerfile.development` | node:20.11.1-alpine → node:22-alpine |
| `infra/docker/func/node/Dockerfile` | node:20.11.1-alpine → node:22-alpine |
| `infra/serverless/serverless.yml` | nodejs18.x → nodejs22.x |

#### user app ツールチェーン更新

| ファイル | 変更内容 |
|---------|---------|
| `web/user/package.json` | TS 5.9.3, Vitest 3.2.4, coverage-v8, happy-dom 18.0.1, vue-tsc 2.2.12 |
| `web/user/package.json` | レガシー ESLint deps 削除（@babel/eslint-parser, @nuxtjs/eslint-config-typescript 等） |
| `web/user/vitest.config.ts` | `/// <reference types="vitest" />` 削除（Vitest 3.x では不要） |

#### Sentry v7 → v9 移行

| ファイル | 変更内容 |
|---------|---------|
| `web/user/package.json` | @sentry/vue + @sentry/node → @sentry/nuxt ^9.0.0 |
| `web/admin/package.json` | @sentry/vue → @sentry/nuxt ^9.0.0 |
| `web/user/src/plugins/sentry.client.ts` | import を @sentry/vue → @sentry/nuxt に変更 |
| `web/user/src/server/plugins/sentry.ts` | import を @sentry/node → @sentry/nuxt に変更 |
| `web/admin/src/plugins/sentry.client.ts` | import を @sentry/vue → @sentry/nuxt に変更 |

#### LIFF Nuxt 4 移行

| ファイル | 変更内容 |
|---------|---------|
| `web/liff/package.json` | nuxt 3.18.1 → ^4.1.3 |
| `web/liff/nuxt.config.ts` | `compatibilityDate: 'latest'` 追加 |

#### pnpm 移行

| ファイル | 変更内容 |
|---------|---------|
| `pnpm-workspace.yaml` | 新規作成（web/admin, user, liff, shared） |
| `web/*/package.json` | packageManager: yarn → pnpm, link: → workspace:* |
| `.github/actions/setup-node/action.yaml` | pnpm/action-setup 追加、yarn cache → pnpm cache |
| `.github/workflows/ci-web-*.yaml` | yarn → pnpm コマンド |
| `.github/workflows/cd-serverless-for-main.yaml` | yarn → pnpm |
| `Makefile` | yarn → pnpm |
| `docker-compose.yaml` | yarn → pnpm |

#### Trivy コンテナスキャン

| ファイル | 変更内容 |
|---------|---------|
| `.github/workflows/_build_and_push.yaml` | Trivy スキャンステップ追加（CRITICAL, HIGH） |

### 将来課題: Observability 統合（OpenTelemetry Collector）

現状は New Relic（APM）+ Sentry（エラートラッキング）+ CloudWatch（ログ）が個別に動作しており、相関分析が困難。インフラ構成変更（OTel Collector デプロイ）が必要なため、本 PR のスコープ外とする。

**目標構成**:

```
Go サービス / Nuxt SSR
  │
  ├── OTLP (traces + metrics)
  │   └── OpenTelemetry Collector
  │       ├── → New Relic (APM/traces)
  │       ├── → CloudWatch (metrics)
  │       └── → Sentry (errors, via Sentry v9 OTel integration)
  │
  └── slog (logs)
      └── stdout → CloudWatch Logs (ECS)
```

## チェックリスト

### 実装開始前

* [x] 各変更の影響範囲を特定
* [x] Amplify で利用可能な Node.js LTS バージョンを確認（22.x）
* [ ] pnpm 移行時の暗黙的依存の洗い出し（`pnpm install` 実行時に確認）

### 動作確認

* [ ] `pnpm install` が全 web app で成功すること
* [ ] user app: `pnpm dev` / `pnpm build` / `pnpm test` が正常動作
* [ ] admin app: `pnpm dev` / `pnpm build` が正常動作
* [ ] liff app: `pnpm dev` / `pnpm build` が正常動作（Nuxt 4）
* [ ] shared lib: `pnpm build` が正常動作
* [ ] Docker 開発環境が Node 22 で正常起動（`make build && make start`）
* [ ] Sentry にエラーレポートが送信されること（staging で確認）
* [ ] CI ワークフローが pnpm で正常実行
* [ ] Trivy スキャンが正常動作

## リリース時確認事項

### リリース順

1. pnpm-lock.yaml 生成（`pnpm install` 実行）
2. 全 yarn.lock を削除
3. CI が全て通ることを確認
4. ステージングで動作確認後、本番リリース

### リリース制御

- ステージングで十分な検証後に本番リリース

### インフラ設定

- Docker イメージの再ビルド（Node 22 ベース）
- Lambda ランタイム変更（Serverless Framework 経由で自動）
- Amplify のビルドイメージが Node 22 をサポートしていることを確認

### パフォーマンスチェック

- pnpm 移行後の CI 所要時間比較
- Node 22 移行後の Nuxt ビルド時間・SSR レスポンス時間比較

### その他

- 旧 yarn.lock ファイルは pnpm 移行完了後に削除
- Sentry v9 移行後の動作確認はステージングで実施

## 関連リンク

- [Node.js 22.22.1 LTS](https://nodejs.org/en/blog/release/v22.22.1)
- [AWS Amplify Supported Node.js Versions](https://docs.aws.amazon.com/amplify/latest/userguide/ssr-supported-features.html)
- [Sentry Nuxt SDK](https://docs.sentry.io/platforms/javascript/guides/nuxt/)
- [pnpm Workspaces](https://pnpm.io/workspaces)
- [Vitest Migration Guide](https://vitest.dev/guide/migration)
- [Nuxt 4 Upgrade Guide](https://nuxt.com/docs/getting-started/upgrade)
- [Trivy Container Scanning](https://aquasecurity.github.io/trivy/)
