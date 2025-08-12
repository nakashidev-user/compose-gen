# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

compose-genは、docker-compose.ymlファイルを対話形式で自動生成するCLIツールです。クリーンアーキテクチャに基づいて設計され、Go言語で実装されています。

## Commands

### Development & Build Commands
```bash
# 開発版の実行
go run ./cmd
# または
go run .

# ビルド
go build -o compose-gen ./cmd
# または
go build -o compose-gen .

# 依存関係のインストール・更新
go mod download
go mod tidy

# 実行例（対話モード）
./compose-gen
./compose-gen --help
```

### Cross-Platform Build
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o compose-gen-linux-amd64 ./cmd

# Windows  
GOOS=windows GOARCH=amd64 go build -o compose-gen-windows-amd64.exe ./cmd

# macOS
GOOS=darwin GOARCH=amd64 go build -o compose-gen-darwin-amd64 ./cmd
```

## Architecture

このプロジェクトはクリーンアーキテクチャの4層構造で設計されています：

### Domain Layer (`internal/domain/`)
- `entities.go`: コアとなるエンティティとビジネスルールを定義
- `errors.go`: ドメイン固有のエラーを定義  
- `repositories.go`: リポジトリインターフェースを定義
- Framework（Next.js, Golang）、DatabaseType（MySQL, MariaDB, PostgreSQL）の型定義
- ProjectConfig、ComposeConfig等のドメインオブジェクト

### Application Layer (`internal/application/`)  
- `usecase.go`: GenerateComposeUseCaseでアプリケーションのユースケースを実装
- ドメインロジックを組み合わせてビジネス要件を満たす処理フロー

### Infrastructure Layer (`internal/infrastructure/`)
- `file_repository.go`: ファイル操作の実装
- `template_repository.go`: テンプレートエンジンの実装  
- `user_input_repository.go`: ユーザー入力処理の実装（survey/v2使用）
- `templates/docker-compose.yml.tmpl`: docker-compose.ymlのテンプレート

### Interface Layer (`internal/interface/cli/`)
- `root.go`: CobraによるCLIインターフェースの実装
- 対話モードでユーザーとのやり取りを管理

## Key Dependencies

- **CLI Framework**: `github.com/spf13/cobra` v1.9.1 - コマンドライン構造とサブコマンド
- **Interactive Prompts**: `github.com/AlecAivazis/survey/v2` v2.3.7 - 対話型プロンプト（選択、入力）
- **Template Engine**: Go標準の`text/template` - docker-compose.yml動的生成
- **Go Version**: 1.24.1 - 最低要求バージョン

## Dependency Injection Pattern

プロジェクトのエントリーポイント（`cmd/main.go`）で以下の順序で依存関係を注入：

1. `infrastructure.NewFileRepository()` - ファイル操作
2. `infrastructure.NewTemplateRepository()` - テンプレート処理  
3. `infrastructure.NewUserInputRepository()` - ユーザー入力
4. `application.NewGenerateComposeUseCase()` - 上記3つを注入してユースケース構築
5. `cli.NewCLI()` - ユースケースを注入してCLI構築

## Adding New Frameworks

1. `internal/domain/entities.go`に新しいFramework定数を追加
2. `DisplayName()`、`DefaultPort()`メソッドに対応するケースを追加
3. 必要に応じてテンプレートを更新

## Adding New Database Types

1. `internal/domain/entities.go`に新しいDatabaseType定数を追加
2. `DisplayName()`、`DefaultVersion()`、`DefaultPort()`メソッドに対応するケースを追加
3. `ImageName()`、`ConnectionString()`、`EnvironmentVars()`、`DataPath()`メソッドに実装を追加
4. テンプレートの更新（必要に応じて）

## Template System

- `internal/infrastructure/templates/docker-compose.yml.tmpl`がメインテンプレート
- Go templateエンジンを使用してComposeConfigから動的生成
- サービス定義、ポート、環境変数、ボリューム、依存関係を動的に構成

## File Generation Flow

1. ユーザー入力収集（プロジェクト名、フレームワーク、データベース設定）
2. ドメインオブジェクト構築（ProjectConfig）
3. ComposeConfig生成（サービス定義）
4. テンプレート適用
5. ファイル出力（上書き確認付き）

## Project Entry Point

- メインエントリーポイント: `cmd/main.go`
- CLIルートコマンド: `internal/interface/cli/root.go`
- DI構成: main.goでリポジトリとユースケースの依存関係注入を実行