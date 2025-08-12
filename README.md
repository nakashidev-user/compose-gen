# compose-gen

> `docker-compose.yml` 自動生成CLIツール

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 概要

`compose-gen`は、開発の初期段階で必要となる`docker-compose.yml`ファイルを、対話形式の簡単な質問に答えるだけで自動生成するCLIツールです。

開発者は複雑なDockerやdocker-composeの記法を覚えることなく、直感的な質問に答えるだけで、構築できます。

## 特徴

- **対話式**: 複雑なコマンドオプションを覚える必要なし
- **テンプレートベース**: ベストプラクティスに基づいた設定
- **日本語対応**: わかりやすい日本語プロンプト
- **高速起動**: 単一バイナリでクロスプラットフォーム対応
- **拡張性**: 新しいフレームワークやデータベースを簡単に追加可能

## 対応技術

### フレームワーク
- **Next.js** - React フレームワーク
- **Golang** - Go言語Webアプリケーション

### データベース
- **MySQL** - 人気のリレーショナルデータベース
- **MariaDB** - MySQLの高性能フォーク
- **PostgreSQL** - 高機能オープンソースデータベース
- **なし** - データベースを使用しない構成

## インストール

### ソースからビルド

```bash
# リポジトリをクローン
git clone https://github.com/nakashidev-user/compose-gen.git
cd compose-gen

# ビルド
go build -o compose-gen ./cmd

# 実行
./compose-gen
```

## 使い方

### 基本的な使用方法

```bash
# 対話モードで実行
./compose-gen
```

### ヘルプの表示

```bash
./compose-gen --help
```

## 使用例

### 1. Next.js + MySQL の構成

```bash
$ ./compose-gen
compose-gen - docker-compose.yml自動生成ツール
対話形式でdocker-compose.ymlを生成します。

? プロジェクト名を入力してください: my-nextjs-app
? アプリケーションのフレームワークを選択してください: Next.js
? データベースの種類を選択してください: MySQL
? MySQLのバージョンを入力してください: 8.0
? MySQLのポート番号を入力してください: 3306

docker-compose.ymlファイルが生成されました: docker-compose.yml
プロジェクト名: my-nextjs-app
フレームワーク: Next.js
データベース: MySQL 8.0 (ポート: 3306)

次のステップ:
  docker-compose up -d
```

### 2. Golang + PostgreSQL の構成

```bash
$ ./compose-gen
? プロジェクト名を入力してください: my-go-api
? アプリケーションのフレームワークを選択してください: Golang
? データベースの種類を選択してください: PostgreSQL
? PostgreSQLのバージョンを入力してください: 15
? PostgreSQLのポート番号を入力してください: 5432

docker-compose.ymlファイルが生成されました: docker-compose.yml
プロジェクト名: my-go-api
フレームワーク: Golang
データベース: PostgreSQL 15 (ポート: 5432)
```

### 3. 生成される docker-compose.yml の例

**Next.js + MySQL の場合:**

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
    ports:
      - "3000:3000"
    environment:
      DATABASE_URL: mysql://root:password@db:3306/my-nextjs-app
    volumes:
      - .:/app
    depends_on:
      - db

  db:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: my-nextjs-app
    volumes:
      - db-data:/var/lib/mysql

volumes:
  db-data:
```

## プロジェクト構造

本プロジェクトは**クリーンアーキテクチャ**に基づいて設計されています：

```
compose-gen/
├── cmd/                    # エントリーポイント
├── internal/
│   ├── domain/            # Domain層（ビジネスロジック）
│   ├── application/       # Application層（ユースケース）
│   ├── infrastructure/    # Infrastructure層（外部ライブラリ連携）
│   └── interface/         # Interface層（CLI）
├── logs/                  # 実装ログ
├── docs/                  # ドキュメント
└── README.md
```

## 開発

### 要件

- Go 1.21以上
- Git

### セットアップ

```bash
# リポジトリをクローン
git clone https://github.com/nakashidev-user/compose-gen.git
cd compose-gen

# 依存関係をインストール
go mod tidy

# 開発版を実行
go run .
```

### ビルド

```bash
# ローカル環境用
go build -o compose-gen .

# クロスプラットフォーム用
GOOS=linux GOARCH=amd64 go build -o compose-gen-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o compose-gen-windows-amd64.exe .
GOOS=darwin GOARCH=amd64 go build -o compose-gen-darwin-amd64 .
```

## ライセンス

MIT License - 詳細は [LICENSE](LICENSE) ファイルを参照してください。
## リンク

- [リポジトリ](https://github.com/nakashidev-user/compose-gen)
- [Issue報告](https://github.com/nakashidev-user/compose-gen/issues)
- [リリース](https://github.com/nakashidev-user/compose-gen/releases)

# compose-gen
