package domain

type Framework string

const (
	FrameworkNextJS Framework = "nextjs"
	FrameworkGolang Framework = "golang"
)

func (f Framework) String() string {
	return string(f)
}

func (f Framework) DisplayName() string {
	switch f {
	case FrameworkNextJS:
		return "Next.js"
	case FrameworkGolang:
		return "Golang"
	default:
		return string(f)
	}
}

func (f Framework) DefaultPort() string {
	switch f {
	case FrameworkNextJS:
		return "3000"
	case FrameworkGolang:
		return "8080"
	default:
		return "8080"
	}
}

type DatabaseType string

const (
	DatabaseNone       DatabaseType = "none"
	DatabaseMySQL      DatabaseType = "mysql"
	DatabaseMariaDB    DatabaseType = "mariadb"
	DatabasePostgreSQL DatabaseType = "postgresql"
)

func (d DatabaseType) String() string {
	return string(d)
}

func (d DatabaseType) DisplayName() string {
	switch d {
	case DatabaseNone:
		return "なし"
	case DatabaseMySQL:
		return "MySQL"
	case DatabaseMariaDB:
		return "MariaDB"
	case DatabasePostgreSQL:
		return "PostgreSQL"
	default:
		return string(d)
	}
}

func (d DatabaseType) DefaultVersion() string {
	switch d {
	case DatabaseMySQL:
		return "8.0"
	case DatabaseMariaDB:
		return "10.9"
	case DatabasePostgreSQL:
		return "15"
	default:
		return ""
	}
}

func (d DatabaseType) DefaultPort() string {
	switch d {
	case DatabaseMySQL, DatabaseMariaDB:
		return "3306"
	case DatabasePostgreSQL:
		return "5432"
	default:
		return ""
	}
}

func (d DatabaseType) InternalPort() string {
	return d.DefaultPort()
}

type Database struct {
	Type    DatabaseType
	Version string
	Port    string
}

func NewDatabase(dbType DatabaseType, version, port string) Database {
	db := Database{
		Type:    dbType,
		Version: version,
		Port:    port,
	}

	if db.Version == "" {
		db.Version = dbType.DefaultVersion()
	}
	if db.Port == "" {
		db.Port = dbType.DefaultPort()
	}

	return db
}

func (d Database) IsEnabled() bool {
	return d.Type != DatabaseNone
}

type ProjectConfig struct {
	ProjectName string
	Framework   Framework
	Database    Database
}

func NewProjectConfig(projectName string, framework Framework, database Database) ProjectConfig {
	return ProjectConfig{
		ProjectName: projectName,
		Framework:   framework,
		Database:    database,
	}
}

func (p ProjectConfig) Validate() error {
	if p.ProjectName == "" {
		return ErrProjectNameRequired
	}
	if p.Framework == "" {
		return ErrFrameworkRequired
	}
	return nil
}

type ComposeService struct {
	Name        string
	Image       string
	Build       *BuildConfig
	Ports       []string
	Environment map[string]string
	Volumes     []string
	DependsOn   []string
}

type BuildConfig struct {
	Context string
}

func NewAppService(config ProjectConfig) ComposeService {
	service := ComposeService{
		Name: "app",
		Build: &BuildConfig{
			Context: ".",
		},
		Ports: []string{
			config.Framework.DefaultPort() + ":" + config.Framework.DefaultPort(),
		},
		Volumes: []string{
			".:/app",
		},
		Environment: make(map[string]string),
	}

	if config.Database.IsEnabled() {
		service.DependsOn = []string{"db"}
		service.Environment["DATABASE_URL"] = config.Database.ConnectionString(config.ProjectName)
	}

	return service
}

func NewDatabaseService(config ProjectConfig) *ComposeService {
	if !config.Database.IsEnabled() {
		return nil
	}

	service := ComposeService{
		Name:        "db",
		Image:       config.Database.ImageName(),
		Ports:       []string{config.Database.Port + ":" + config.Database.Type.InternalPort()},
		Environment: config.Database.EnvironmentVars(config.ProjectName),
		Volumes:     []string{"db-data:" + config.Database.DataPath()},
	}

	return &service
}

func (d Database) ImageName() string {
	switch d.Type {
	case DatabaseMySQL:
		return "mysql:" + d.Version
	case DatabaseMariaDB:
		return "mariadb:" + d.Version
	case DatabasePostgreSQL:
		return "postgres:" + d.Version
	default:
		return ""
	}
}

func (d Database) ConnectionString(projectName string) string {
	switch d.Type {
	case DatabaseMySQL, DatabaseMariaDB:
		return "mysql://root:password@db:3306/" + projectName
	case DatabasePostgreSQL:
		return "postgresql://postgres:password@db:5432/" + projectName
	default:
		return ""
	}
}

func (d Database) EnvironmentVars(projectName string) map[string]string {
	env := make(map[string]string)

	switch d.Type {
	case DatabaseMySQL, DatabaseMariaDB:
		env["MYSQL_ROOT_PASSWORD"] = "password"
		env["MYSQL_DATABASE"] = projectName
	case DatabasePostgreSQL:
		env["POSTGRES_PASSWORD"] = "password"
		env["POSTGRES_DB"] = projectName
	}

	return env
}

func (d Database) DataPath() string {
	switch d.Type {
	case DatabaseMySQL, DatabaseMariaDB:
		return "/var/lib/mysql"
	case DatabasePostgreSQL:
		return "/var/lib/postgresql/data"
	default:
		return ""
	}
}

type ComposeConfig struct {
	Services map[string]ComposeService
	Volumes  map[string]interface{}
}

func NewComposeConfig(config ProjectConfig) ComposeConfig {
	compose := ComposeConfig{
		Services: make(map[string]ComposeService),
		Volumes:  make(map[string]interface{}),
	}

	appService := NewAppService(config)
	compose.Services["app"] = appService

	if dbService := NewDatabaseService(config); dbService != nil {
		compose.Services["db"] = *dbService
		compose.Volumes["db-data"] = nil
	}

	return compose
}