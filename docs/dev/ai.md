# AIサービス（顔認識）

**ディレクトリ:** `ai/`

---

## 概要

ArcFaceモデルを使った顔検出・顔埋め込みベクトル生成サービスです。バックエンドサーバーから呼び出され、アップロードされた写真から顔を検出し、512次元の埋め込みベクトルを生成します。このベクトルをpgvectorで検索することで、同一人物の写真を自動的にグルーピングします。

---

## 使用技術

| カテゴリ       | ライブラリ        | バージョン |
| -------------- | ----------------- | ---------- |
| フレームワーク | FastAPI           | 0.121.1    |
| サーバー       | Uvicorn           | 0.38.0     |
| 顔認識モデル   | InsightFace       | 0.7.3      |
| 画像処理       | OpenCV (headless) | 4.12.0     |
| 画像処理       | Pillow            | 11.3.0     |
| データ検証     | Pydantic          | 2.12.4     |
| 数値計算       | NumPy             | 2.0.2      |
| 機械学習       | scikit-learn      | 1.6.1      |
| 科学計算       | SciPy             | 1.13.1     |
| 画像分析       | scikit-image      | 0.24.0     |

---

## ディレクトリ構成

```
ai/
├── src/
│   ├── main.py              # FastAPIエンドポイント定義
│   ├── service/
│   │   └── face.py          # 顔検出ロジック
│   └── utils/               # ユーティリティ関数
├── models/                  # 事前学習済みモデルファイル
├── requirements.txt         # Python依存関係
└── docker/                  # Dockerコンテナ設定
```

---

## APIエンドポイント

### `GET /`

ヘルスチェック

**レスポンス:**

```json
{ "status": "ok" }
```

---

### `POST /face-detection/file`

アップロードされた画像ファイルから顔を検出します。

**リクエスト:** `multipart/form-data`

- `file`: 画像ファイル (JPEG, PNG等)

**レスポンス:**

```json
{
  "faces": [
    {
      "bbox": [x, y, width, height],
      "embedding": [0.12, -0.34, ...]  // 512次元ベクトル
    }
  ]
}
```

---

### `POST /face-detection/url`

URLで指定した画像から顔を検出します。

**リクエスト:**

```json
{
  "url": "https://example.com/image.jpg"
}
```

**レスポンス:** `/face-detection/file` と同形式

---

## 顔認識モデル

### InsightFace (ArcFace)

- **モデル:** `buffalo_s`
- **出力:** 512次元の顔埋め込みベクトル
- **特徴:** 正規化された埋め込みにより、コサイン類似度で人物の同一性を判定

### 処理フロー

```
入力画像
    ↓
OpenCV / Pillow で前処理
    ↓
InsightFace (ArcFace) で顔検出
    ↓
バウンディングボックス座標を抽出
    ↓
各顔の512次元埋め込みを生成
    ↓
正規化してレスポンス
```

---

## バックエンドとの連携

```
写真アップロード
    ↓ (サーバーが呼び出し)
POST /face-detection/file または /url
    ↓
[bbox, embedding] を返却
    ↓
サーバーが face_detections テーブルに保存
    ↓
pgvector の IVFFlat インデックスで類似検索
    ↓
同一人物 (identity) にグルーピング
```

---

## 開発コマンド

```bash
# 依存関係インストール
pip install -r requirements.txt

# 開発サーバー起動
uvicorn src.main:app --reload --host 0.0.0.0 --port 8000

# Dockerで起動 (docker-composeから)
docker-compose up ai
```

---

## 注意事項

- **GPU対応:** InsightFaceはCPU/GPU両対応ですが、ローカル開発ではCPUモードで動作します。
- **モデルファイル:** `models/` ディレクトリに事前学習済みモデルが必要です。初回起動時に自動ダウンロードされる場合があります。
- **メモリ:** 顔認識モデルはメモリを多く使用するため、コンテナのメモリ制限に注意してください。
