package domain

import "errors"

var (
	ErrProjectNameRequired = errors.New("プロジェクト名は必須です")
	ErrFrameworkRequired   = errors.New("フレームワークの選択は必須です")
	ErrFileAlreadyExists   = errors.New("ファイルが既に存在します")
	ErrTemplateNotFound    = errors.New("テンプレートファイルが見つかりません")
	ErrInvalidConfig       = errors.New("設定が無効です")
)