package config

import (
	"database/sql"
	_ "embed"
	"errors"
	"github.com/BurntSushi/toml"
	"io/fs"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

const (
	DatabaseFile = "cornhole/cornhole.db"
	ConfFile     = "cornhole/config.toml"
)

//go:embed "config.toml"
var defaultConf []byte

type Config struct {
	Server   server
	Database *sql.DB
}

type configFile struct {
	Server server
	Database database
}

type server struct {
	Port int
}

type database struct {
	Path *string
}

func New() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(configDir, ConfFile)

	err = initializeConfig(configPath)
	if err != nil {
		return nil, err
	}

	configFile, err := loadFromFile(configPath)
	if err != nil {
		return nil, err
	}

	var dbPath string
	if configFile.Database.Path != nil {
		dbPath = *configFile.Database.Path
	} else {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return nil, err
		}
		dbPath = filepath.Join(cacheDir, DatabaseFile)
	}

	db, err := initializeDatabase(dbPath)
	if err != nil {
		return nil, err
	}

	config := Config{}
	config.Server = configFile.Server
	config.Database = db

	return &config, nil
}

func initializeConfig(configPath string) error {
	if _, err := os.Stat(configPath); errors.Is(err, fs.ErrNotExist) {
		if err = os.MkdirAll(filepath.Dir(configPath), 0744); err != nil {
			return err
		}
		if err = os.WriteFile(configPath, defaultConf, 0644); err != nil {
			return err
		}
	}

	return nil
}

func initializeDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func loadFromFile(filename string) (*configFile, error) {
	var config configFile
	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
