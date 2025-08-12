package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"compose-gen/internal/application"
	"compose-gen/internal/infrastructure"
)

type CLI struct {
	useCase *application.GenerateComposeUseCase
}

func NewCLI() *CLI {
	fileRepo := infrastructure.NewFileRepository()
	templateRepo := infrastructure.NewTemplateRepository()
	userInputRepo := infrastructure.NewUserInputRepository()

	useCase := application.NewGenerateComposeUseCase(
		fileRepo,
		templateRepo,
		userInputRepo,
	)

	return &CLI{
		useCase: useCase,
	}
}

func (c *CLI) RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compose-gen",
		Short: "docker-compose.yml自動生成ツール",
		Long: `compose-gen - docker-compose.yml自動生成ツール

docker-compose.ymlファイルを、対話形式の簡単な質問に答えるだけで自動生成します。

対応フレームワーク:
  - Next.js
  - Golang

対応データベース:
  - MySQL
  - MariaDB  
  - PostgreSQL
  - なし`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.runInteractiveMode(cmd.Context()); err != nil {
				fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return cmd
}

func (c *CLI) runInteractiveMode(ctx context.Context) error {
	fmt.Println("compose-gen - docker-compose.yml自動生成ツール")
	fmt.Println("対話形式でdocker-compose.ymlを生成します。")
	fmt.Println()

	input := application.GenerateComposeInput{
		OutputPath: "docker-compose.yml",
	}

	output, err := c.useCase.Execute(ctx, input)
	if err != nil {
		return err
	}

	fmt.Printf(" docker-compose.ymlファイルが生成されました: %s\n", output.FilePath)
	fmt.Printf(" プロジェクト名: %s\n", output.Config.ProjectName)
	fmt.Printf(" フレームワーク: %s\n", output.Config.Framework.DisplayName())

	if output.Config.Database.IsEnabled() {
		fmt.Printf(" データベース: %s %s (ポート: %s)\n",
			output.Config.Database.Type.DisplayName(),
			output.Config.Database.Version,
			output.Config.Database.Port,
		)
	} else {
		fmt.Printf(" データベース: なし\n")
	}

	fmt.Println()
	fmt.Println("次のステップ:")
	fmt.Println("  docker-compose up -d")
	fmt.Println()

	return nil
}
