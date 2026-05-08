# ローカル開発環境セットアップ

---

## 前提条件

以下のツールが事前にインストールされている必要があります。

| ツール                  | バージョン | 用途                     |
| ----------------------- | ---------- | ------------------------ |
| Docker + Docker Compose | 最新版     | ローカルインフラ         |
| Go                      | 1.24.0+    | バックエンド開発         |
| Node.js                 | 18+        | フロントエンド・モバイル |
| Python                  | 3.10+      | AIサービス               |
| air                     | 最新版     | Goホットリロード         |
| sqlc                    | 最新版     | SQLコード生成            |
| golang-migrate          | 最新版     | DBマイグレーション       |

---

## 1. Dockerサービスの起動

`server/` ディレクトリでDockerサービスを起動します。

```bash
cd server
docker-compose up -d
```

### 起動されるサービス

| サービス   | ポート      | 説明                         |
| ---------- | ----------- | ---------------------------- |
| `postgres` | 5433        | PostgreSQL 17 (pgvector)     |
| `minio`    | 9000 / 9090 | S3互換オブジェクトストレージ |
| `ai`       | 8000        | Python顔認識サービス         |
| `app`      | 1323        | Goバックエンドサーバー       |

> **Note:** `app` と `ai` は個別に起動する場合は除外できます。
>
> ```bash
> docker-compose up -d postgres minio
> ```

---

## 2. バックエンドサーバーのセットアップ

```bash
cd server

# 環境変数設定
cp .env.example .env
# .env を編集して各値を設定

# DBマイグレーション実行
make migrate-up

# 開発サーバー起動 (ホットリロード)
air
```

サーバーは `http://localhost:1323` で起動します。

### Swaggerドキュメント

```
http://localhost:1323/swagger/index.html
```

---

## 3. Webフロントエンドのセットアップ

```bash
cd front

# 依存関係インストール
npm install

# 開発サーバー起動
npm run dev
```

Webアプリは `http://localhost:5173` で起動します。

---

## 4. モバイルアプリのセットアップ

```bash
cd mobile

# 依存関係インストール
npm install

# iOSの場合: Podインストール
cd ios && pod install && cd ..

# Metro開発サーバー起動
npm start

# 別ターミナルでiOS起動
npm run ios

# またはAndroid起動
npm run android
```

---

## 5. AIサービスのセットアップ (Dockerを使わない場合)

```bash
cd ai

# 仮想環境作成
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate

# 依存関係インストール
pip install -r requirements.txt

# サーバー起動
uvicorn src.main:app --reload --host 0.0.0.0 --port 8000
```

---

## MinIO の設定

MinIOが起動したら、初回はバケットを作成する必要があります。

**MinIOコンソール:** `http://localhost:9090`

- ユーザー名: `minioadmin`
- パスワード: `minioadmin`

以下のバケットを作成してください:

- `thumbnails`
- `originals`

> docker-compose の `create-buckets` サービスが自動作成します。

---

## 環境変数一覧

### server/.env

```env
# データベース
DATABASE_URL=postgres://postgres:postgres@localhost:5433/yuno

# JWT
JWT_SECRET=your-secret-key

# OAuth - LINE
LINE_CLIENT_ID=your-line-client-id
LINE_CLIENT_SECRET=your-line-client-secret
LINE_REDIRECT_URL=http://localhost:1323/auth/callback/line

# OAuth - Kakao
KAKAO_CLIENT_ID=your-kakao-client-id
KAKAO_REDIRECT_URL=http://localhost:1323/auth/callback/kakao

# ストレージ (MinIO)
STORAGE_ENDPOINT=localhost:9001
STORAGE_ACCESS_KEY=minioadmin
STORAGE_SECRET_KEY=minioadmin
STORAGE_USE_SSL=false
```

---

## よく使う Makefile コマンド

```bash
# DBマイグレーション
make migrate-up       # マイグレーション適用
make migrate-down     # 最後のマイグレーションを戻す

# コード生成
make sqlc             # sqlc でDBコードを生成
make swag             # Swagger ドキュメントを生成

# データ投入
make seed             # テストデータ投入
```

---

## トラブルシューティング

### PostgreSQLに接続できない

```bash
# コンテナの状態確認
docker-compose ps

# ログ確認
docker-compose logs postgres
```

### pgvector拡張が有効でない

```sql
-- PostgreSQLに接続して実行
CREATE EXTENSION IF NOT EXISTS vector;
```

### AIサービスが起動しない

モデルファイルが `ai/models/` に存在するか確認してください。存在しない場合は初回起動時に自動ダウンロードが実行されます（ネットワーク接続が必要）。

### ポートが競合する

`.env` ファイルでポート番号を変更してください。docker-compose.yml と合わせて修正が必要です。
