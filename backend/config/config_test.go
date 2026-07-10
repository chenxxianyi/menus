package config

import "testing"

func TestValidateProductionRejectsUnsafeSettings(t *testing.T) {
	cfg := &Config{
		App:      AppConfig{Environment: "production", Debug: false},
		Database: DatabaseConfig{Password: "database-password", AutoMigrate: false},
		JWT:      JWTConfig{SecretKey: "your-secret-key-change-in-production"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() accepted a default production JWT secret")
	}

	cfg.JWT.SecretKey = "a-very-long-random-production-secret-that-is-not-a-default"
	cfg.Database.AutoMigrate = true
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() accepted production auto_migrate")
	}
}

func TestLoadConfigOverridesDatabaseFromEnvironment(t *testing.T) {
	t.Setenv("DATABASE_HOST", "127.0.0.1")
	t.Setenv("DATABASE_PORT", "3306")
	t.Setenv("DATABASE_NAME", "ai_workbench")
	t.Setenv("DATABASE_USER", "root")
	t.Setenv("DATABASE_PASSWORD", "123456")
	t.Setenv("DATABASE_CHARSET", "utf8mb4")
	t.Setenv("DATABASE_LOC", "UTC")
	t.Setenv("DATABASE_MAX_OPEN_CONNS", "20")
	t.Setenv("DATABASE_MAX_IDLE_CONNS", "10")
	t.Setenv("DATABASE_CONN_MAX_LIFETIME_MINUTES", "30")
	t.Setenv("DB_AUTO_MIGRATE", "true")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Database.Host != "127.0.0.1" {
		t.Fatalf("Database.Host = %q, want %q", cfg.Database.Host, "127.0.0.1")
	}
	if cfg.Database.Port != 3306 {
		t.Fatalf("Database.Port = %d, want %d", cfg.Database.Port, 3306)
	}
	if cfg.Database.DBName != "ai_workbench" {
		t.Fatalf("Database.DBName = %q, want %q", cfg.Database.DBName, "ai_workbench")
	}
	if cfg.Database.User != "root" {
		t.Fatalf("Database.User = %q, want %q", cfg.Database.User, "root")
	}
	if cfg.Database.Password != "123456" {
		t.Fatalf("Database.Password = %q, want %q", cfg.Database.Password, "123456")
	}
	if cfg.Database.Charset != "utf8mb4" {
		t.Fatalf("Database.Charset = %q, want %q", cfg.Database.Charset, "utf8mb4")
	}
	if cfg.Database.MaxOpenConns != 20 {
		t.Fatalf("Database.MaxOpenConns = %d, want %d", cfg.Database.MaxOpenConns, 20)
	}
	if cfg.Database.MaxIdleConns != 10 {
		t.Fatalf("Database.MaxIdleConns = %d, want %d", cfg.Database.MaxIdleConns, 10)
	}
	if cfg.Database.Loc != "UTC" {
		t.Fatalf("Database.Loc = %q, want %q", cfg.Database.Loc, "UTC")
	}
	if cfg.Database.ConnMaxLifetimeMinutes != 30 {
		t.Fatalf("Database.ConnMaxLifetimeMinutes = %d, want %d", cfg.Database.ConnMaxLifetimeMinutes, 30)
	}
	if !cfg.Database.AutoMigrate {
		t.Fatal("Database.AutoMigrate = false, want true")
	}
}
