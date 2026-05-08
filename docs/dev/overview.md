# YUNO プロジェクト概要

## プロジェクトについて

YUNOは、家族の思い出を共有するためのフォトシェアリングアプリケーションです。AIによる顔認識機能を搭載し、家族・グループ単位での写真管理をサポートします。Webアプリ・モバイルアプリの両方に対応しており、クロスプラットフォームで利用可能です。

---

## リポジトリ構成

```
yuno/
├── front/       # Webフロントエンド (React + TypeScript)
├── mobile/      # モバイルアプリ (React Native)
├── server/      # バックエンドAPI (Go + Echo)
├── ai/          # AI顔認識サービス (Python + FastAPI)
├── infra/       # インフラ構成 (Terraform / AWS)
└── docs/        # プロジェクトドキュメント
```

---

## 主要機能

| 機能 | 説明 |
|------|------|
| 写真・動画のアップロード | Web/モバイルから家族アルバムへアップロード |
| AI顔認識 | InsightFace (ArcFace) による顔検出・同一人物識別 |
| 家族・グループ管理 | ファミリー単位でのグループ分け、招待リンク発行 |
| アルバム管理 | グループごとのアクセス権限付きアルバム |
| OAuthログイン | LINE / Kakao のソーシャルログイン対応 |
| 多言語対応 | 韓国語・日本語の国際化 (i18n) |
| PWA対応 | オフラインサポート (Webアプリ) |

---

## 技術スタック一覧

| レイヤー | 技術 | バージョン | 用途 |
|----------|------|-----------|------|
| **Webフロントエンド** | React | 18.2.0 | UIフレームワーク |
| | TypeScript | 5.0.2 | 型安全性 |
| | Vite | 7.0.6 | ビルドツール |
| | Tailwind CSS | 4.1.11 | スタイリング |
| | React Router | 6.18.0 | クライアントルーティング |
| | TanStack Query | 5.7.2 | サーバー状態管理 |
| | shadcn/ui | - | UIコンポーネント |
| | Firebase | 10.5.2 | 認証補助 |
| | i18next | 25.8.11 | 多言語対応 |
| **モバイル** | React Native | 0.73.5 | クロスプラットフォームUI |
| | TypeScript | 5.0.4 | 型安全性 |
| | React Navigation | 6.1+ | モバイルルーティング |
| | NativeWind | 2.0.11 | Tailwind CSS for RN |
| | TanStack Query | 5.24.8 | サーバー状態管理 |
| **バックエンド** | Go | 1.24.0 | サーバー言語 |
| | Echo | 4.11.4 | WebフレームワークE |
| | sqlc | - | 型安全なSQLコード生成 |
| | golang-migrate | 4.19.0 | DBマイグレーション |
| | Swagger (swaggo) | - | APIドキュメント自動生成 |
| **データベース** | PostgreSQL | 17 | メインDB |
| | pgvector | - | ベクトル類似検索 |
| **AI/ML** | FastAPI | 0.121.1 | AIサービスフレームワーク |
| | InsightFace | 0.7.3 | 顔検出・埋め込み生成 |
| | OpenCV | 4.12.0 | 画像処理 |
| **ストレージ** | MinIO | latest | S3互換オブジェクトストレージ (ローカル開発) |

---

## 認証フロー

```
ユーザー
  │
  ├─ OAuth (LINE / Kakao)
  │     └─ コールバック → サーバーでJWT発行
  │
  └─ JWTトークン
        ├─ Webアプリ: Session Storage に保存
        └─ モバイル: AsyncStorage に保存
```

---

## 写真アップロードフロー

```
1. ユーザーが写真をアップロード (Web / Mobile)
       ↓
2. バックエンドサーバーが受信
   - サムネイル・閲覧用画像を生成
   - MinIO / S3 に保存
       ↓
3. AIサービスに顔検出リクエスト
   - ArcFaceモデルで512次元の埋め込みを生成
       ↓
4. サーバーが埋め込みをpgvectorに保存
   - IVFFlat インデックスで高速検索
       ↓
5. 類似度検索で同一人物を識別
```

---

## ローカル開発環境

ローカル開発では `docker-compose` を使用します。

```yaml
サービス構成:
  - app      : Goサーバー         (port: 1323)
  - ai       : PythonAIサービス   (port: 8000)
  - postgres : PostgreSQL 17      (port: 5433)
  - minio    : S3互換ストレージ   (port: 9000 / 9090)
```

詳細は [development.md](./development.md) を参照してください。

---

## ドキュメント一覧

| ファイル | 説明 |
|----------|------|
| [overview.md](./overview.md) | プロジェクト概要（このファイル） |
| [frontend.md](./frontend.md) | Webフロントエンド詳細 |
| [mobile.md](./mobile.md) | モバイルアプリ詳細 |
| [server.md](./server.md) | バックエンドサーバー詳細 |
| [ai.md](./ai.md) | AIサービス詳細 |
| [development.md](./development.md) | ローカル開発環境セットアップ |
