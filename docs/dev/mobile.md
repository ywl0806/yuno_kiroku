# モバイルアプリ

**ディレクトリ:** `mobile/`

---

## 使用技術

| カテゴリ           | ライブラリ                   | バージョン |
| ------------------ | ---------------------------- | ---------- |
| フレームワーク     | React Native                 | 0.73.5     |
| 言語               | TypeScript                   | 5.0.4      |
| バンドラー         | Metro                        | -          |
| スタイリング       | NativeWind (Tailwind for RN) | 2.0.11     |
| ナビゲーション     | React Navigation             | 6.1+       |
| サーバー状態管理   | TanStack React Query         | 5.24.8     |
| HTTPクライアント   | Axios                        | 1.6.8      |
| ローカルストレージ | AsyncStorage                 | 2.2.0      |
| メディアアクセス   | Camera Roll                  | 7.4.2      |
| ファイルシステム   | React Native FS              | 2.20.0     |
| 動画再生           | React Native Video           | 5.2.1      |
| UIアイコン         | React Native Vector Icons    | -          |
| 写真一覧           | Masonry List                 | 1.4.2      |

---

## ディレクトリ構成

```
mobile/
├── src/
│   ├── screens/          # 画面コンポーネント
│   │   ├── HomeScreen    # タイムライン・写真一覧
│   │   ├── LoginScreen   # ログイン
│   │   ├── UploadScreen  # 写真アップロード
│   │   └── SettingsScreen# 設定
│   ├── navigation/       # ナビゲーション設定
│   │   └── カスタムタブバー
│   ├── services/         # APIサービス
│   │   ├── getPhotos.ts       # サーバーから写真取得
│   │   ├── uploadPhoto.ts     # 写真アップロード
│   │   └── getLocalPhotos.ts  # 端末内写真取得
│   ├── context/          # グローバル状態
│   │   ├── AuthContext   # 認証状態管理
│   │   └── TabContext    # タブナビゲーション状態
│   ├── components/       # 再利用可能コンポーネント
│   ├── hooks/            # カスタムフック
│   ├── icons/            # アイコン素材
│   ├── types/            # TypeScript型定義
│   ├── utils/            # ユーティリティ関数
│   └── lib/              # ヘルパー関数
├── ios/                  # iOSネイティブコード
├── android/              # Androidネイティブコード
└── package.json
```

---

## ナビゲーション構成

```
RootStack (Native Stack)
├── LoginScreen          # 未認証時
└── MainTabs (Bottom Tab Navigator)
    ├── Home             # ホームタブ
    ├── Upload           # アップロードタブ
    └── Settings         # 設定タブ
```

- ボトムタブにカスタムタブバーを使用
- 認証状態に応じてルート切り替え

---

## 認証

- **OAuth:** LINE / Kakao ソーシャルログイン
- **トークン保存:** `AsyncStorage`（永続化）
- **グローバル状態:** `AuthContext` で管理
- **Axiosインターセプター:** リクエスト時にJWTを自動付与

---

## 状態管理

### サーバー状態: TanStack React Query

- 写真一覧などのAPIデータをキャッシュ・管理

### クライアント状態: Context API

- `AuthContext`: 認証トークン・ユーザー情報
- `TabContext`: タブナビゲーション状態

---

## メディア機能

### 端末内写真取得

- `Camera Roll` で端末のカメラロールにアクセス
- 写真のメタデータ（撮影日時・GPS情報）取得

### 写真アップロード

1. カメラロールから写真を選択
2. バックグラウンドでサーバーへアップロード
3. React Query でキャッシュを無効化して画面を更新

### 動画再生

- `React Native Video` でサーバー上の動画を再生

---

## 開発コマンド

```bash
# 依存関係インストール
npm install

# iOSポッドインストール
cd ios && pod install

# Metro開発サーバー起動
npm start

# iOS実行
npm run ios

# Android実行
npm run android
```

---

## 対応プラットフォーム

| プラットフォーム | 対応状況 |
| ---------------- | -------- |
| iOS              | ✅ 対応  |
| Android          | ✅ 対応  |
