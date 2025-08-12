package application

import (
	"context"
	"compose-gen/internal/domain"
)

type GenerateComposeUseCase struct {
	fileRepo      domain.FileRepository
	templateRepo  domain.TemplateRepository
	userInputRepo domain.UserInputRepository
}

func NewGenerateComposeUseCase(
	fileRepo domain.FileRepository,
	templateRepo domain.TemplateRepository,
	userInputRepo domain.UserInputRepository,
) *GenerateComposeUseCase {
	return &GenerateComposeUseCase{
		fileRepo:      fileRepo,
		templateRepo:  templateRepo,
		userInputRepo: userInputRepo,
	}
}

type GenerateComposeInput struct {
	OutputPath string
}

type GenerateComposeOutput struct {
	FilePath string
	Config   domain.ProjectConfig
}

func (uc *GenerateComposeUseCase) Execute(ctx context.Context, input GenerateComposeInput) (*GenerateComposeOutput, error) {
	projectName, err := uc.userInputRepo.AskProjectName(ctx)
	if err != nil {
		return nil, err
	}

	framework, err := uc.userInputRepo.AskFramework(ctx)
	if err != nil {
		return nil, err
	}

	dbType, err := uc.userInputRepo.AskDatabase(ctx)
	if err != nil {
		return nil, err
	}

	var database domain.Database
	if dbType != domain.DatabaseNone {
		version, err := uc.userInputRepo.AskDatabaseVersion(ctx, dbType)
		if err != nil {
			return nil, err
		}

		port, err := uc.userInputRepo.AskDatabasePort(ctx, dbType)
		if err != nil {
			return nil, err
		}

		database = domain.NewDatabase(dbType, version, port)
	} else {
		database = domain.NewDatabase(domain.DatabaseNone, "", "")
	}

	config := domain.NewProjectConfig(projectName, framework, database)
	if err := config.Validate(); err != nil {
		return nil, err
	}

	outputPath := input.OutputPath
	if outputPath == "" {
		outputPath = "docker-compose.yml"
	}

	exists, err := uc.fileRepo.Exists(ctx, outputPath)
	if err != nil {
		return nil, err
	}

	if exists {
		overwrite, err := uc.userInputRepo.AskOverwrite(ctx, outputPath)
		if err != nil {
			return nil, err
		}
		if !overwrite {
			return nil, domain.ErrFileAlreadyExists
		}
	}

	composeConfig := domain.NewComposeConfig(config)
	content, err := uc.templateRepo.Generate(ctx, composeConfig)
	if err != nil {
		return nil, err
	}

	if err := uc.fileRepo.Write(ctx, outputPath, content); err != nil {
		return nil, err
	}

	return &GenerateComposeOutput{
		FilePath: outputPath,
		Config:   config,
	}, nil
}