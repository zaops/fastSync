package database

import (
	"testing"
)

// TestDatabaseConfig tests the database configuration structure
func TestDatabaseConfig(t *testing.T) {
	config := Config{
		Type:     "sqlite",
		DSN:      ":memory:",
		LogLevel: "info",
	}

	if config.Type != "sqlite" {
		t.Errorf("Expected type 'sqlite', got '%s'", config.Type)
	}

	if config.DSN != ":memory:" {
		t.Errorf("Expected DSN ':memory:', got '%s'", config.DSN)
	}

	if config.LogLevel != "info" {
		t.Errorf("Expected log level 'info', got '%s'", config.LogLevel)
	}
}

// TestUnsupportedDatabaseType tests error handling for unsupported database types
func TestUnsupportedDatabaseType_Mock(t *testing.T) {
	config := Config{
		Type:     "mysql",
		DSN:      "test",
		LogLevel: "silent",
	}

	err := Initialize(config)
	if err == nil {
		t.Error("Should return error for unsupported database type")
	}

	expectedError := "unsupported database type: mysql"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestAutoMigrateWithoutDB tests auto migrate without database initialization
func TestAutoMigrateWithoutDB_Mock(t *testing.T) {
	// 保存原始DB状态
	originalDB := DB
	DB = nil
	defer func() {
		DB = originalDB
	}()

	err := AutoMigrate()
	if err == nil {
		t.Error("Should return error when database is not initialized")
	}

	expectedError := "database not initialized"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestGetDBWithoutInitialization tests GetDB without initialization
func TestGetDBWithoutInitialization(t *testing.T) {
	// 保存原始DB状态
	originalDB := DB
	DB = nil
	defer func() {
		DB = originalDB
	}()

	db := GetDB()
	if db != nil {
		t.Error("GetDB should return nil when database is not initialized")
	}
}

// TestCloseWithoutDB tests Close without database
func TestCloseWithoutDB(t *testing.T) {
	// 保存原始DB状态
	originalDB := DB
	DB = nil
	defer func() {
		DB = originalDB
	}()

	err := Close()
	if err != nil {
		t.Errorf("Close should not return error when DB is nil, got: %v", err)
	}
}