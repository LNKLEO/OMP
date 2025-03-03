package config

import (
	"bytes"
	"fmt"
	"os"
	stdOS "os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LNKLEO/OMP/log"
	"github.com/LNKLEO/OMP/runtime/path"
	"github.com/gookit/goutil/jsonutil"

	json "github.com/goccy/go-json"
	yaml "github.com/goccy/go-yaml"
	toml "github.com/pelletier/go-toml/v2"
)

// LoadConfig returns the default configuration including possible user overrides
func Load(configFile, sh string) *Config {
	defer log.Trace(time.Now())

	cfg := loadConfig(configFile)
	return cfg
}

func Path(config string) string {
	defer log.Trace(time.Now())

	// if the config flag is set, we'll use that over OMP_THEME
	// in our internal shell logic, we'll always use the OMP_THEME
	// due to not using --config to set the configuration
	hasConfig := len(config) > 0

	if OMPTheme := os.Getenv("OMP_THEME"); len(OMPTheme) > 0 && !hasConfig {
		log.Debug("config set using OMP_THEME:", OMPTheme)
		return OMPTheme
	}

	if len(config) == 0 {
		return ""
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
