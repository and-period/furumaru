# Architecture - システム設計意図・構造

システムの設計思想、意思決定の経緯、システム構造を文書化します。

## 📐 設計ドキュメント

### 🎯 [design-decisions.md](./design-decisions.md)
アーキテクチャ設計決定の記録
- マイクロサービス採用の背景と理由
- 技術選択の決定経緯
- トレードオフと今後の課題

### 📊 [overview.md](./overview.md)
システム全体の構造
- サービス構成とデータフロー
- 認証・認可アーキテクチャ
- ビジネスドメイン構造
- データベース設計

### 🔧 [api/README.md](./api/README.md)
バックエンドAPI設計思想
- モジュラーモノリス + レイヤードアーキテクチャ
- エンジニア実装のしやすさを重視

### 🌐 Web アプリ構造 (`web/`)
フロントエンドアプリケーションの構造
- **[web/README.md](./web/README.md)**: フロントエンドアーキテクチャ概要
- **[web/components.md](./web/components.md)**: コンポーネント設計・分類
- **[web/state-management.md](./web/state-management.md)**: Pinia状態管理
- **[web/api-integration.md](./web/api-integration.md)**: API連携アーキテクチャ

### 🤝 [shared/README.md](./shared/README.md)
共通設計仕様・ガイドライン

## 設計思想

### 実装しやすさの重視
エンジニアが理解しやすく実装しやすいパターンを採用

### モジュール分離
機能ごとに明確に分離し、責務を明確化

### 長期保守性
運用・保守のしやすさを考慮した設計
