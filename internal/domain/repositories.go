package domain

import "context"

type FileRepository interface {
	Exists(ctx context.Context, filePath string) (bool, error)
	Write(ctx context.Context, filePath string, content string) error
}

type TemplateRepository interface {
	Generate(ctx context.Context, config ComposeConfig) (string, error)
}

type UserInputRepository interface {
	AskProjectName(ctx context.Context) (string, error)
	AskFramework(ctx context.Context) (Framework, error)
	AskDatabase(ctx context.Context) (DatabaseType, error)
	AskDatabaseVersion(ctx context.Context, dbType DatabaseType) (string, error)
	AskDatabasePort(ctx context.Context, dbType DatabaseType) (string, error)
	AskOverwrite(ctx context.Context, filePath string) (bool, error)
}