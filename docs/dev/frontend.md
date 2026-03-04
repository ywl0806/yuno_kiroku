# Webフロントエンド

**ディレクトリ:** `front/`

---

## 使用技術

| カテゴリ         | ライブラリ              | バージョン       |
| ---------------- | ----------------------- | ---------------- |
| UIフレームワーク | React                   | 18.2.0           |
| 言語             | TypeScript              | 5.0.2            |
| ビルドツール     | Vite + SWC              | 7.0.6            |
| スタイリング     | Tailwind CSS            | 4.1.11           |
| CSS-in-JS        | Emotion                 | -                |
| UIコンポーネント | shadcn/ui (Radix UI)    | -                |
| ルーティング     | React Router DOM        | 6.18.0           |
| サーバー状態管理 | TanStack React Query    | 5.7.2            |
| フォーム管理     | React Hook Form         | 7.61.1           |
| バリデーション   | Zod                     | 4.0.10           |
| HTTPクライアント | Axios                   | 1.6.0            |
| 認証補助         | Firebase                | 10.5.2           |
| 多言語対応       | i18next + react-i18next | 25.8.11 / 16.5.4 |
| アイコン         | Lucide React, MUI Icons | -                |
| カルーセル       | Swiper                  | 11.1.0           |
| フォトアルバム   | react-photo-album       | 2.3.1            |
| PWA              | vite-plugin-pwa         | 1.0.2            |

---

## ディレクトリ構成

```
front/
├── src/
│   ├── feature/          # 機能モジュール
│   │   ├── auth/         # 認証関連
│   │   ├── home/         # ホーム画面
│   │   ├── setting/      # 設定画面
│   │   └── upload/       # 写真アップロード
│   ├── page/             # ページコンポーネント
│   ├── components/       # 再利用可能なUIコンポーネント
│   ├── service/          # APIサービス層
│   ├── lib/              # ユーティリティ
│   │   ├── axios.ts      # Axiosラッパー (認証トークン自動付与)
│   │   ├── query-client.ts # React Query設定
│   │   └── session-storage.ts # セッションストレージ管理
│   ├── config/           # Firebase等の設定
│   ├── i18n/             # 国際化リソース
│   │   ├── kr.json       # 韓国語
│   │   └── jp.json       # 日本語
│   ├── hooks/            # カスタムReactフック
│   ├── types/            # TypeScript型定義
│   ├── providers/        # Contextプロバイダー
│   │   └── upload-photo-provider.tsx
│   ├── router.tsx        # ルーティング設定
│   └── app.tsx           # ルートコンポーネント
├── public/               # 静的ファイル
├── package.json
├── tsconfig.json
└── vite.config.ts
```

---

## ルーティング構成

```
/                   # ルート (ログイン済みならリダイレクト)
├── /login          # ログインページ
└── /               # DefaultLayout (認証必須)
    ├── /home       # ホーム (写真タイムライン)
    ├── /upload     # 写真アップロード
    ├── /album      # アルバム一覧・詳細
    ├── /family     # 家族・グループ管理
    └── /settings   # 設定ページ
```

---

## 状態管理

### サーバー状態: TanStack React Query

- APIのデータ取得・キャッシュ・同期に使用
- カスタムフックでクエリをカプセル化

```typescript
// 例: 写真一覧取得
const { data, isLoading } = useQuery({
  queryKey: ['mediaItems'],
  queryFn: () => mediaService.getMediaItems(),
})
```

### クライアント状態: Context API

- アップロード処理の状態管理に `UploadPhotoProvider` を使用

---

## 認証

- **OAuth:** LINE / Kakao ソーシャルログイン
- **トークン保存:** Session Storage
- **Axiosラッパー:** リクエストヘッダーにJWTを自動付与
- **Firebase:** 認証補助に使用

---

## フォーム・バリデーション

React Hook Form + Zod の組み合わせ。

```typescript
const schema = z.object({
  name: z.string().min(1),
})

const { register, handleSubmit } = useForm({
  resolver: zodResolver(schema),
})
```

---

## 多言語対応 (i18n)

- 対応言語: 韓国語 (`kr`) / 日本語 (`jp`)
- `i18next` + `react-i18next` を使用
- 翻訳ファイル: `src/i18n/kr.json`, `src/i18n/jp.json`

---

## PWA (Progressive Web App)

- `vite-plugin-pwa` によるサービスワーカー生成
- オフラインキャッシュサポート
- ホーム画面への追加が可能

---

## 開発コマンド

```bash
# 依存関係インストール
npm install

# 開発サーバー起動 (Vite HMR)
npm run dev

# 本番ビルド
npm run build

# 型チェック
npm run type-check
```
