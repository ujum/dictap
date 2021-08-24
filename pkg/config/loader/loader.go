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
	v := viper.New()
	v.SetConfigType(loadSettings.ConfigType)
	v.AddConfigPath(loadSettings.ConfigDir)
	mergeCommonConfig(v, loadSettings)
	mergeAppEnvConfig(v, loadSettings)
	if loadSettings.LoadSysEnv {
		mergeSystemEnvConfig(v, loadSettings.EnvPrefix)
	}
	if err := v.Unmarshal(configStruct); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
		return err
	}
	return nil
}

// merge common(base) config file to config struct
func mergeCommonConfig(v *viper.Viper, loadSettings *LoadSettings) {
	const baseConfigName = "base"
	v.SetConfigName(resolveFileName(loadSettings.FileNamePrefix, baseConfigName))
	if err := v.MergeInConfig(); err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("loaded %s config (%s)", baseConfigName, v.ConfigFileUsed())
}

// merge app profile config file to config struct
func mergeAppEnvConfig(v *viper.Viper, loadSettings *LoadSettings) {
	var envKey = applicationEnv
	if loadSettings.EnvPrefix != "" {
		envKey = strings.ToUpper(loadSettings.EnvPrefix) + "_" + applicationEnv
	}
	profile := os.Getenv(envKey)
	if profile == "" {
		profile = v.GetString("app.profile")
	}
	if profile != "" {
		v.SetConfigName(resolveFileName(loadSettings.FileNamePrefix, profile))
		if err := v.MergeInConfig(); err != nil {
			if err, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Printf("config for %s profile not found, %v", profile, err)
			} else {
				log.Printf("%v", err)
			}
			return
		}
		log.Printf("loaded %v profile config (%s)", profile, v.ConfigFileUsed())
	}
}

// merge system environment variables to config struct
func mergeSystemEnvConfig(v *viper.Viper, envPrefix string) {
	v.SetEnvPrefix(envPrefix)
	v.AllowEmptyEnv(false)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}

func resolveFileName(filePrefix, name string) string {
	if filePrefix == "" {
		return name
	}
	return filePrefix + "_" + name
}
