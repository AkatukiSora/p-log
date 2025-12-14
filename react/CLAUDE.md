# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

```bash
npm run dev      # 開発サーバー起動 (Vite)
npm run build    # プロダクションビルド
npm run lint     # ESLint実行
npm run preview  # ビルド後のプレビュー
```

## Architecture

React 19 + Vite + React Router v7 を使用したSPA。目標管理・ソーシャル機能を持つアプリケーション (p-log)。

### ディレクトリ構成

```
src/
├── main.jsx           # エントリーポイント
├── App.jsx            # ルーティング定義 (BrowserRouter)
├── index.css          # グローバルスタイル
├── pages/             # ページコンポーネント (各ページにCSS Modules)
│   ├── pathData/path.js  # ルート定義定数 (ROUTES)
│   ├── myGoalsList/      # 自分の目標一覧
│   ├── myProfile/        # マイプロフィール
│   ├── friendList/       # フレンド一覧
│   ├── friendProfile/    # フレンドプロフィール
│   ├── friendGoalsList/  # フレンドの目標一覧
│   ├── timeLine/         # タイムライン
│   ├── settingGoal/      # 目標設定
│   ├── AddPost/          # 投稿追加
│   └── goals/            # 目標ページ (現在のルートパス "/")
└── components/        # 再利用可能コンポーネント
    ├── button/           # ボタン系 (ButtonTopMyOption, ButtonTopFriendOption, ButtonBottomOption)
    └── goalsData/        # 目標データ表示 (GoalData, GoalsData)
```

### スタイリング

- CSS Modules を使用 (`*.module.css`)
- グローバルスタイルは `src/index.css`

### ルーティング

`src/pages/pathData/path.js` の `ROUTES` オブジェクトでパスを一元管理。日本語パスを含む。

## Tech Stack

- **React 19** + **react-dom 19**
- **Vite 6** (ビルドツール)
- **React Router DOM v7** (ルーティング)
- **react-icons** (アイコン)
- **ESLint** (eslint-plugin-react, react-hooks, react-refresh)
