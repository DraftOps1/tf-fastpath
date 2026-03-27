package config

import (
	"fmt"
	"path/filepath"
)

// Config holds the shared runtime settings used by all commands.
type Config struct {
	WorkingDir     string
	DataDir        string
	SQLitePath     string
	TerraformBin   string
	OpenTofuBin    string
	GitBin         string
	DefaultRuntime string
}

// Load resolves configuration from environment variables and working directory.
func Load(lookupEnv func(string) string, getwd func() (string, error)) (Config, error) {
	if lookupEnv == nil {
		return Config{}, fmt.Errorf("lookupEnv is required")
	}
	if getwd == nil {
		return Config{}, fmt.Errorf("getwd is required")
	}

	workingDir, err := getwd()
	if err != nil {
		return Config{}, fmt.Errorf("resolve working directory: %w", err)
	}

	cfg := Config{
		WorkingDir:     envOr(lookupEnv, "TFFASTPATH_WORKING_DIR", workingDir),
		TerraformBin:   envOr(lookupEnv, "TFFASTPATH_TERRAFORM_BIN", "terraform"),
		OpenTofuBin:    envOr(lookupEnv, "TFFASTPATH_OPENTOFU_BIN", "tofu"),
		GitBin:         envOr(lookupEnv, "TFFASTPATH_GIT_BIN", "git"),
		DefaultRuntime: envOr(lookupEnv, "TFFASTPATH_DEFAULT_RUNTIME", "terraform"),
	}
	cfg.DataDir = envOr(lookupEnv, "TFFASTPATH_DATA_DIR", filepath.Join(cfg.WorkingDir, ".tf-fastpath"))
	cfg.SQLitePath = envOr(lookupEnv, "TFFASTPATH_SQLITE_PATH", filepath.Join(cfg.DataDir, "tf-fastpath.sqlite3"))

	return cfg, nil
}

func envOr(lookupEnv func(string) string, key string, fallback string) string {
	if value := lookupEnv(key); value != "" {
		return value
	}
	return fallback
}
