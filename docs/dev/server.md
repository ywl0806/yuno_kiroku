# バックエンドサーバー

**ディレクトリ:** `server/`

---

## 使用技術

| カテゴリ | ライブラリ | バージョン |
|----------|-----------|-----------|
| 言語 | Go | 1.24.0 |
| Webフレームワーク | Echo | 4.11.4 |
| DB | PostgreSQL | 17 |
| ベクトル検索 | pgvector | - |
| SQLコード生成 | sqlc | - |
| DBマイグレーション | golang-migrate | 4.19.0 |
| バリデーション | go-playground/validator | v10 |
| JWT | golang-jwt | v5.3.0 |
| AWSクライアント | aws-sdk-go-v2 | - |
| 画像処理 | govips (libvips) | v2 |
| UUID | google/uuid | - |
| 設定管理 | viper | - |
| EXIFデータ | goexif | - |
| APIドキュメント | swaggo (Swagger) | - |
| ホットリロード | air | - |

---

## ディレクトリ構成

```
server/
├── internal/
│   ├── handlers/       # HTTPリクエストハンドラー
│   ├── routers/        # APIルート定義
│   ├── services/       # ビジネスロジック
│   ├── store/          # データアクセス層
│   ├── db/             # DB操作 (sqlc生成コード)
│   ├── middlewares/    # HTTPミドルウェア
│   ├── oauth/          # OAuthプロバイダー (LINE, Kakao)
│   ├── apperr/         # カスタムエラー型
│   ├── consts/         # 定数
│   ├── enums/          # 列挙型
│   ├── i18n/           # 国際化 (エラーメッセージ等)
│   └── validator/      # カスタムバリデーター
├── db/
│   ├── schema.sql      # DBスキーマ定義
│   ├── queries/        # sqlc用SQLクエリ
│   └── migrations/     # マイグレーションファイル
├── docker-compose.yml  # ローカル開発環境
├── .air.toml           # ホットリロード設定
├── Makefile            # 開発タスク
└── go.mod
```

---

## アーキテクチャ

レイヤードアーキテクチャを採用しています。

```
HTTPリクエスト
    ↓
Middleware (CORS / 認証 / ログ)
    ↓
Handler (リクエスト解析・レスポンス整形)
    ↓
Service (ビジネスロジック)
    ↓
Store (データアクセス)
    ↓
DB (PostgreSQL / sqlc生成コード)
```

---

## APIエンドポイント一覧

### 認証

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/auth/login` | ログイン |
| GET | `/auth/callback/line` | LINEコールバック |
| GET | `/auth/callback/kakao` | Kakaoコールバック |
| POST | `/auth/refresh` | トークンリフレッシュ |

### ユーザー

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/users/me` | プロフィール取得 |
| PUT | `/users/me` | プロフィール更新 |

### メディア

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/media-items` | 写真・動画アップロード |
| GET | `/media-items` | 日付範囲で取得 |
| POST | `/media-items/batch` | バッチアップロード |

### アルバム

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/albums` | アルバム一覧 |
| POST | `/albums` | アルバム作成 |
| GET | `/albums/:id` | アルバム詳細 |
| PUT | `/albums/:id` | アルバム更新 |

### 顔認識・ID管理

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/identities` | 識別済み人物一覧 |
| POST | `/identities` | 人物作成 |
| GET | `/face-detections` | 顔検出データ取得 |

### 招待

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/invites` | 招待トークン生成 |
| GET | `/invites/:token` | 招待トークン検証 |

---

## データベーススキーマ

### 主要テーブル

| テーブル名 | 説明 |
|-----------|------|
| `families` | 家族組織 |
| `groups` | 家族内のグループ（父方・母方など） |
| `users` | ユーザーアカウント（OAuthサポート） |
| `albums` | 写真アルバム（グループ権限付き） |
| `media_items` | 写真・動画アイテム |
| `media_files` | ファイル実体（オリジナル・サムネイル・閲覧用） |
| `identities` | 顔認識で識別された人物 |
| `face_detections` | 顔検出データ（512次元埋め込みベクトル） |
| `invite_tokens` | 招待リンクトークン |
| `refresh_tokens` | JWTリフレッシュトークン |

### pgvector による顔認識

```sql
-- 顔埋め込みベクトルを保存
embedding vector(512)

-- IVFFlat インデックスでコサイン類似度検索を高速化
CREATE INDEX ON face_detections
  USING ivfflat (embedding vector_cosine_ops);
```

### media_files のロール

| ロール | 説明 |
|--------|------|
| `original` | 元の高解像度ファイル |
| `thumbnail` | 一覧表示用サムネイル |
| `view` | ビュー用に最適化されたファイル |
| `live` | Live Photoなど動的コンテンツ |

---

## 認証・認可

- **JWT:** アクセストークン + リフレッシュトークン
- **OAuthプロバイダー:** LINE / Kakao
- **ミドルウェア:** 認証が必要なルートに `AuthMiddleware` を適用
- **リフレッシュトークン:** DBに保存して管理

---

## ファイルストレージ

- **本番:** AWS S3
- **ローカル開発:** MinIO（S3互換）
- **バケット構成:**
  - `thumbnails` バケット
  - `originals` バケット

---

## 画像処理

`govips`（libvipsのGoバインディング）を使用して高速な画像処理を実現。

- リサイズ・サムネイル生成
- EXIFデータ（撮影日時・GPS）の読み取り

---

## 開発コマンド

```bash
# 開発サーバー起動 (ホットリロード)
air

# DBマイグレーション
make migrate-up
make migrate-down

# sqlc コード生成
make sqlc

# Swaggerドキュメント生成
make swag

# テスト実行
go test ./...
```

---

## 環境変数

`.env.example` を参照して `.env` ファイルを作成してください。

主要な環境変数:

| 変数名 | 説明 |
|--------|------|
| `DATABASE_URL` | PostgreSQL接続URL |
| `AI_API_URL` | AIサービスのURL |
| `LINE_CLIENT_ID` | LINE OAuthクライアントID |
| `LINE_CLIENT_SECRET` | LINE OAuthクライアントシークレット |
| `KAKAO_CLIENT_ID` | Kakao OAuthクライアントID |
| `S3_BUCKET` | S3バケット名 |
| `JWT_SECRET` | JWTシークレットキー |
