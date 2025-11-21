package config

// AppConfig represents the full config structure for vmdeploy v2.x
type AppConfig struct {
	AppName      string `json:"app_name" yaml:"app_name"`
	RepoURL      string `json:"repo_url" yaml:"repo_url"`
	Branch       string `json:"branch" yaml:"branch"`
	ArtifactName string `json:"artifact_name" yaml:"artifact_name"`

	DBType string `json:"db_type" yaml:"db_type"`
	DBHost string `json:"db_host" yaml:"db_host"`
	DBPort int    `json:"db_port" yaml:"db_port"`
	DBName string `json:"db_name" yaml:"db_name"`

	ServiceName string `json:"service_name" yaml:"service_name"`
	AppPort     int    `json:"app_port" yaml:"app_port"`
}

// VM settings (where to deploy)
type VMConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
}
