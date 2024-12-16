package config

import (
	"bytes"
	"fmt"
	"os"
	stdOS "os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gookit/goutil/jsonutil"

	"github.com/LNKLEO/OMP/cache"
	"github.com/LNKLEO/OMP/log"
	"github.com/LNKLEO/OMP/runtime/path"

	json "github.com/goccy/go-json"
	yaml "github.com/goccy/go-yaml"
	toml "github.com/pelletier/go-toml/v2"
)

// LoadConfig returns the default configuration including possible user overrides
func Load(configFile, sh string, migrate bool) *Config {
	defer log.Trace(time.Now())

	cfg := loadConfig(configFile)

	// only migrate automatically when the switch isn't set
	if !migrate && cfg.Version < Version {
		cfg.BackupAndMigrate()
	}

	if !cfg.ShellIntegration {
		return cfg
	}

	return cfg
}

func Path(config string) string {
	defer log.Trace(time.Now())

	// if the config flag is set, we'll use that over OMP_THEME
	// in our internal shell logic, we'll always use the OMP_THEME
	// due to not using --config to set the configuration
	hasConfig := len(config) > 0

	if poshTheme := os.Getenv("OMP_THEME"); len(poshTheme) > 0 && !hasConfig {
		log.Debug("config set using OMP_THEME:", poshTheme)
		return poshTheme
	}

	if len(config) == 0 {
		return ""
	}

	if strings.HasPrefix(config, "https://") {
		filePath, err := Download(cache.Path(), config)
		if err != nil {
			log.Error(err)
			return ""
		}

		return filePath
	}

	isCygwin := func() bool {
		return runtime.GOOS == "windows" && len(os.Getenv("OSTYPE")) > 0
	}

	// Cygwin path always needs the full path as we're on Windows but not really.
	// Doing filepath actions will convert it to a Windows path and break the init script.
	if isCygwin() {
		log.Debug("cygwin detected, using full path for config")
		return config
	}

	configFile := path.ReplaceTildePrefixWithHomeDir(config)

	abs, err := filepath.Abs(configFile)
	if err != nil {
		log.Error(err)
		return filepath.Clean(configFile)
	}

	return abs
}

func loadConfig(configFile string) *Config {
	defer log.Trace(time.Now())

	if len(configFile) == 0 {
		log.Debug("no config file specified, using default")
		return Default(false)
	}

	var cfg Config
	cfg.origin = configFile
	cfg.Format = strings.TrimPrefix(filepath.Ext(configFile), ".")

	data, err := stdOS.ReadFile(configFile)
	if err != nil {
		log.Error(err)
		return Default(true)
	}

	switch cfg.Format {
	case "yml", "yaml":
		cfg.Format = YAML
		err = yaml.Unmarshal(data, &cfg)
	case "jsonc", "json":
		cfg.Format = JSON

		str := jsonutil.StripComments(string(data))
		data = []byte(str)

		decoder := json.NewDecoder(bytes.NewReader(data))
		err = decoder.Decode(&cfg)
	case "toml", "tml":
		cfg.Format = TOML
		err = toml.Unmarshal(data, &cfg)
	default:
		err = fmt.Errorf("unsupported config file format: %s", cfg.Format)
	}

	if err != nil {
		log.Error(err)
		return Default(true)
	}

	return &cfg
}
