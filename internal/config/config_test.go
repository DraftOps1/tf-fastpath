package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Parallel()

	cfg, err := Load(
		func(string) string { return "" },
		func() (string, error) { return "/workspace/tf-fastpath", nil },
	)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.WorkingDir != "/workspace/tf-fastpath" {
		t.Fatalf("WorkingDir = %q, want %q", cfg.WorkingDir, "/workspace/tf-fastpath")
	}
	if cfg.DataDir != "/workspace/tf-fastpath/.tf-fastpath" {
		t.Fatalf("DataDir = %q", cfg.DataDir)
	}
	if cfg.SQLitePath != "/workspace/tf-fastpath/.tf-fastpath/tf-fastpath.sqlite3" {
		t.Fatalf("SQLitePath = %q", cfg.SQLitePath)
	}
	if cfg.TerraformBin != "terraform" {
		t.Fatalf("TerraformBin = %q", cfg.TerraformBin)
	}
	if cfg.OpenTofuBin != "tofu" {
		t.Fatalf("OpenTofuBin = %q", cfg.OpenTofuBin)
	}
	if cfg.DefaultRuntime != "terraform" {
		t.Fatalf("DefaultRuntime = %q", cfg.DefaultRuntime)
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	t.Parallel()

	env := map[string]string{
		"TFFASTPATH_WORKING_DIR":     "/repo/envdir",
		"TFFASTPATH_TERRAFORM_BIN":   "/usr/local/bin/terraform",
		"TFFASTPATH_OPENTOFU_BIN":    "/usr/local/bin/tofu",
		"TFFASTPATH_GIT_BIN":         "/usr/bin/git",
		"TFFASTPATH_DEFAULT_RUNTIME": "opentofu",
	}

	cfg, err := Load(
		func(key string) string { return env[key] },
		func() (string, error) { return "/workspace/tf-fastpath", nil },
	)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.WorkingDir != "/repo/envdir" {
		t.Fatalf("WorkingDir = %q", cfg.WorkingDir)
	}
	if cfg.TerraformBin != "/usr/local/bin/terraform" {
		t.Fatalf("TerraformBin = %q", cfg.TerraformBin)
	}
	if cfg.OpenTofuBin != "/usr/local/bin/tofu" {
		t.Fatalf("OpenTofuBin = %q", cfg.OpenTofuBin)
	}
	if cfg.GitBin != "/usr/bin/git" {
		t.Fatalf("GitBin = %q", cfg.GitBin)
	}
	if cfg.DefaultRuntime != "opentofu" {
		t.Fatalf("DefaultRuntime = %q", cfg.DefaultRuntime)
	}
}
