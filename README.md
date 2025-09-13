# Furumaru - 全国ふるさとマルシェ

地域特産品を扱う日本のECマーケットプレイスプラットフォーム（ライブコマース機能付き）

## プロジェクト概要

FurumaruはGoマイクロサービスとNuxtフロントエンドによるモノレポ構成のECプラットフォームです。
- **バックエンド**: モジュラーモノリス（Go）+ レイヤードアーキテクチャ
- **フロントエンド**: Nuxt マルチアプリケーション（管理者/購入者/LINE）
- **データベース**: MySQL/TiDB
- **認証**: AWS Cognito, OAuth2.0（Google/LINE）

## クイックスタート

```bash
# 初回セットアップ
make setup

# 環境設定
cp .env.temp .env    # .envファイルを編集

# サービス起動
make start

# アクセスURL
# ユーザー側: http://localhost:3000
# 管理者側: http://localhost:3010
```

## 開発コマンド

```bash
# API開発
cd api && make test          # APIテスト実行
cd api && make lint          # リンター実行

# フロントエンド開発
cd web/user && yarn dev      # ユーザー側開発サーバー
cd web/admin && yarn dev     # 管理者側開発サーバー
```

## プロジェクト構造

```
furumaru/
├── api/                     # バックエンドAPI（Go）
│   ├── cmd/                 # エントリポイント
│   └── internal/            # ビジネスモジュール
│       ├── gateway/         # APIゲートウェイ
│       ├── user/            # ユーザー管理
│       ├── store/           # EC機能
│       ├── media/           # メディア管理
│       └── messenger/       # 通知機能
├── web/                     # フロントエンド
│   ├── admin/               # 管理者ポータル（Nuxt）
│   ├── user/                # 購入者ポータル（Nuxt）
│   ├── liff/                # LINEアプリ（Nuxt）
│   └── shared/              # 共通コンポーネント
├── infra/                   # インフラ設定
├── docs/                    # ドキュメント
│   ├── agents/              # AIエージェント向けガイド
│   ├── architecture/        # アーキテクチャ設計
│   ├── knowledge/           # 実装パターン・知見
│   ├── rules/               # 開発規約
│   └── spec/                # 仕様書・API定義
└── Makefile                 # ビルド・運用コマンド
```

## ドキュメント

### 開発者向け
- [アーキテクチャ概要](docs/architecture/overview.md)
- [開発規約](docs/rules/coding-standards.md)
- [よく使うコマンド](docs/knowledge/commands.md)
- [トラブルシューティング](docs/knowledge/troubleshooting.md)
- [仕様書・API定義](docs/spec/README.md)

### AIエージェント向け
- [エージェント運用ガイド](AGENTS.md)
- [基本方針・運用規約](docs/agents/charter.md)
- [開発環境セットアップ](docs/agents/quick-start.md)

## 技術スタック

### バックエンド
- Go
- Echo (Web Framework)
- GORM (ORM)
- MySQL / TiDB

### フロントエンド
- Nuxt
- Vue
- Vuetify (管理画面)
- Tailwind CSS (ユーザー画面)

### インフラ・ツール
- Docker / Docker Compose
- AWS (Cognito, S3, CloudFront)
- GitHub Actions (CI/CD)

## ライセンス

Proprietary - All Rights Reserved
