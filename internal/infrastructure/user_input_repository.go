package infrastructure

import (
	"context"
	"fmt"

	"compose-gen/internal/domain"

	"github.com/AlecAivazis/survey/v2"
)

type UserInputRepositoryImpl struct{}

func NewUserInputRepository() *UserInputRepositoryImpl {
	return &UserInputRepositoryImpl{}
}

func (r *UserInputRepositoryImpl) AskProjectName(ctx context.Context) (string, error) {
	var projectName string
	prompt := &survey.Input{
		Message: "プロジェクト名を入力してください:",
	}

	if err := survey.AskOne(prompt, &projectName, survey.WithValidator(survey.Required)); err != nil {
		return "", err
	}

	return projectName, nil
}

func (r *UserInputRepositoryImpl) AskFramework(ctx context.Context) (domain.Framework, error) {
	frameworks := []string{
		domain.FrameworkNextJS.DisplayName(),
		domain.FrameworkGolang.DisplayName(),
	}

	var selected string
	prompt := &survey.Select{
		Message: "アプリケーションのフレームワークを選択してください:",
		Options: frameworks,
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return "", err
	}

	switch selected {
	case domain.FrameworkNextJS.DisplayName():
		return domain.FrameworkNextJS, nil
	case domain.FrameworkGolang.DisplayName():
		return domain.FrameworkGolang, nil
	default:
		return "", domain.ErrInvalidConfig
	}
}

func (r *UserInputRepositoryImpl) AskDatabase(ctx context.Context) (domain.DatabaseType, error) {
	databases := []string{
		domain.DatabaseNone.DisplayName(),
		domain.DatabaseMySQL.DisplayName(),
		domain.DatabaseMariaDB.DisplayName(),
		domain.DatabasePostgreSQL.DisplayName(),
	}

	var selected string
	prompt := &survey.Select{
		Message: "データベースの種類を選択してください:",
		Options: databases,
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return "", err
	}

	switch selected {
	case domain.DatabaseNone.DisplayName():
		return domain.DatabaseNone, nil
	case domain.DatabaseMySQL.DisplayName():
		return domain.DatabaseMySQL, nil
	case domain.DatabaseMariaDB.DisplayName():
		return domain.DatabaseMariaDB, nil
	case domain.DatabasePostgreSQL.DisplayName():
		return domain.DatabasePostgreSQL, nil
	default:
		return "", domain.ErrInvalidConfig
	}
}

func (r *UserInputRepositoryImpl) AskDatabaseVersion(ctx context.Context, dbType domain.DatabaseType) (string, error) {
	defaultVersion := dbType.DefaultVersion()

	var version string
	prompt := &survey.Input{
		Message: fmt.Sprintf("%sのバージョンを入力してください:", dbType.DisplayName()),
		Default: defaultVersion,
	}

	if err := survey.AskOne(prompt, &version); err != nil {
		return "", err
	}

	if version == "" {
		version = defaultVersion
	}

	return version, nil
}

func (r *UserInputRepositoryImpl) AskDatabasePort(ctx context.Context, dbType domain.DatabaseType) (string, error) {
	defaultPort := dbType.DefaultPort()

	var port string
	prompt := &survey.Input{
		Message: fmt.Sprintf("%sのポート番号を入力してください:", dbType.DisplayName()),
		Default: defaultPort,
	}

	if err := survey.AskOne(prompt, &port); err != nil {
		return "", err
	}

	if port == "" {
		port = defaultPort
	}

	return port, nil
}

func (r *UserInputRepositoryImpl) AskOverwrite(ctx context.Context, filePath string) (bool, error) {
	var overwrite bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("ファイル '%s' が既に存在します。上書きしますか?", filePath),
		Default: false,
	}

	if err := survey.AskOne(prompt, &overwrite); err != nil {
		return false, err
	}

	return overwrite, nil
}
