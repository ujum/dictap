package loader

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

const applicationEnv = "APP_ENV"

// LoadSettings represents how to load config to config struct.
type LoadSettings struct {
	ConfigDir      string // path to config directory
	FileNamePrefix string // config file prefix
	ConfigType     string // config file type (see viper.SupportedExts)
	// SetEnvPrefix defines a prefix that ENVIRONMENT variables will use.
	// See viper.SetEnvPrefix
	EnvPrefix string
	// Flag to load environment variables
	LoadSysEnv bool
}

// Load config from files and system env to app config struct.
// Use 'APP_ENV' sys env key to denote app profile config.
func Load(configStruct interface{}, loadSettings *LoadSettings) error {
	viper.SetConfigType(loadSettings.ConfigType)
	viper.AddConfigPath(loadSettings.ConfigDir)
	mergeCommonConfig(loadSettings)
	mergeAppEnvConfig(loadSettings)
	if loadSettings.LoadSysEnv {
		mergeSystemEnvConfig(loadSettings.EnvPrefix)
	}
	if err := viper.Unmarshal(configStruct); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		return err
	}
	return nil
}

// merge common(base) config file to config struct
func mergeCommonConfig(loadSettings *LoadSettings) {
	const baseConfigName = "base"
	viper.SetConfigName(resolveFileName(loadSettings.FileNamePrefix, baseConfigName))
	if err := viper.MergeInConfig(); err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("Loaded %s config (%s)", baseConfigName, viper.ConfigFileUsed())
	}
}

// merge app profile config file to config struct
func mergeAppEnvConfig(loadSettings *LoadSettings) {
	var envKey = applicationEnv
	if loadSettings.EnvPrefix != "" {
		envKey = strings.ToUpper(loadSettings.EnvPrefix) + "_" + applicationEnv
	}
	profile := os.Getenv(envKey)
	if profile == "" {
		profile = viper.GetString("app.profile")
	}
	if profile != "" {
		viper.SetConfigName(resolveFileName(loadSettings.FileNamePrefix, profile))
		if err := viper.MergeInConfig(); err != nil {
			if err, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Printf("Config for %s profile not found, %v", profile, err)
			}
		} else {
			log.Printf("Loaded %v profile config (%s)", profile, viper.ConfigFileUsed())
		}
	}
}

// merge system environment variables to config struct
func mergeSystemEnvConfig(envPrefix string) {
	viper.SetEnvPrefix(envPrefix)
	viper.AllowEmptyEnv(false)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func resolveFileName(filePrefix, name string) string {
	if filePrefix == "" {
		return name
	}
	return filePrefix + "_" + name
}
