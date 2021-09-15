package loader

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type Config struct {
	Number int
	Part   *ConfigPart
}

type ConfigPart struct {
	Str string
}

type ConfigTestParam struct {
	name         string
	resultConfig interface{}
	loadSettings *LoadSettings
	env          map[string]string
	want         *Config
	wantErr      bool
}

func TestErrorCases(t *testing.T) {
	runTests(t, []ConfigTestParam{
		{
			name:         "emptyConfigDir",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{ConfigFile: &ConfigFileSettings{
				ConfigDir: "",
			}},
			wantErr: true,
		},
		{
			name:         "nilConfigFileSettings",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{},
			wantErr:      true,
		},
		{
			name:         "nilLoadSettings",
			resultConfig: &Config{},
			wantErr:      true,
		},
		{
			name: "nonPointerConfigStruct",
			loadSettings: &LoadSettings{ConfigFile: &ConfigFileSettings{
				ConfigDir: "./testdata/noprofiles",
			}},
			resultConfig: Config{},
			wantErr:      true,
		},
		{
			name: "nilConfigStruct",
			loadSettings: &LoadSettings{ConfigFile: &ConfigFileSettings{
				ConfigDir: "./testdata/noprofiles",
			}},
			resultConfig: nil,
			wantErr:      true,
		},
		{
			name: "baseConfigFileNotExist",
			loadSettings: &LoadSettings{ConfigFile: &ConfigFileSettings{
				ConfigDir: "/notexistsdir",
			}},
			wantErr: true,
		},
		{
			name: "profileConfigFileNotExist",
			loadSettings: &LoadSettings{
				ConfigFile: &ConfigFileSettings{
					ConfigDir: "./testdata/noprofiles",
				},
			},
			env:     map[string]string{"APP_PROFILE": "stage"},
			wantErr: true,
		},
		{
			name: "invalidBaseConfig",
			loadSettings: &LoadSettings{
				ConfigFile: &ConfigFileSettings{
					ConfigType: "toml",
					ConfigDir:  "./testdata/invalidbaseconfig",
				},
			},
			wantErr: true,
		},
		{
			name: "invalidProfileConfig",
			loadSettings: &LoadSettings{
				ConfigFile: &ConfigFileSettings{
					ConfigType: "toml",
					ConfigDir:  "./testdata/invalidprofileconfig",
				},
			},
			wantErr: true,
		}})
}

func TestLoadFromFile(t *testing.T) {
	tests := []ConfigTestParam{
		{
			name:         "file config",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{
				ConfigFile: &ConfigFileSettings{
					ConfigDir: "./testdata/noprofiles",
				},
			},
			want: &Config{
				Number: 55,
				Part:   &ConfigPart{Str: "filepartbase"},
			},
		},
		{
			name:         "file config",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{
				ConfigFile: &ConfigFileSettings{
					ConfigDir:      "./testdata/noprofiles/prefix",
					FileNamePrefix: "config",
				},
			},
			want: &Config{
				Number: 555,
				Part:   &ConfigPart{Str: "filepartbaseprefix"},
			},
		},
		{
			name:         "file config with profile set by base file",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{
				ConfigFile: &ConfigFileSettings{
					ConfigDir: "./testdata/profiles",
				},
			},
			want: &Config{
				Number: 44,
				Part:   &ConfigPart{Str: "filepartdevpr"},
			},
		},
		{
			name:         "file config with profile set by base file and env (env should override file)",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{
				ConfigFile: &ConfigFileSettings{
					ConfigDir: "./testdata/profiles",
				},
			},
			env: map[string]string{"APP_PROFILE": "stage"},
			want: &Config{
				Number: 77,
				Part:   &ConfigPart{Str: "filepartstagepr"},
			},
		},
		{
			name:         "file config with prefix",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{
				ConfigFile: &ConfigFileSettings{
					ConfigDir:      "./testdata/profiles/prefix",
					FileNamePrefix: "config",
				},
			},
			want: &Config{
				Number: 11,
				Part:   &ConfigPart{Str: "filepartbaseprefix"},
			},
		},
		{
			name:         "file config with profile and prefix",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{
				ConfigFile: &ConfigFileSettings{
					ConfigDir:      "./testdata/profiles/prefix",
					FileNamePrefix: "config",
				},
			},
			env: map[string]string{"APP_PROFILE": "dev"},
			want: &Config{
				Number: 22,
				Part:   &ConfigPart{Str: "filepartdevprefix"},
			},
		},
	}

	runTests(t, tests)
}

func TestLoadFromFileAndEnv(t *testing.T) {
	tests := []ConfigTestParam{
		{
			name:         "env without prefix",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{
				LoadSysEnv: true,
				ConfigFile: &ConfigFileSettings{
					ConfigDir: "./testdata/noprofiles",
				},
			},
			env: map[string]string{"NUMBER": "99", "PART_STR": "part"},
			want: &Config{
				Number: 99,
				Part:   &ConfigPart{Str: "part"},
			},
		},
		{
			name:         "env with prefix",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{
				LoadSysEnv: true,
				EnvPrefix:  "app",
				ConfigFile: &ConfigFileSettings{
					ConfigDir: "./testdata/noprofiles",
				},
			},
			env: map[string]string{"NUMBER": "99", "APP_PART_STR": "part"},
			want: &Config{
				Number: 55,
				Part:   &ConfigPart{Str: "part"},
			},
		},
		{
			name:         "env without prefix and profile from file",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{
				LoadSysEnv: true,
				ConfigFile: &ConfigFileSettings{
					ConfigDir: "./testdata/profiles",
				},
			},
			env: map[string]string{"NUMBER": "99", "PART_STR": "part"},
			want: &Config{
				Number: 99,
				Part:   &ConfigPart{Str: "part"},
			},
		},
		{
			name:         "env with prefix and profile from file",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{
				LoadSysEnv: true,
				EnvPrefix:  "app",
				ConfigFile: &ConfigFileSettings{
					ConfigDir: "./testdata/profiles",
				},
			},
			env: map[string]string{"NUMBER": "99", "APP_PART_STR": "part"},
			want: &Config{
				Number: 44,
				Part:   &ConfigPart{Str: "part"},
			},
		},
		{
			name:         "env without prefix and profile from env",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{
				LoadSysEnv: true,
				ConfigFile: &ConfigFileSettings{
					ConfigDir: "./testdata/profiles",
				},
			},
			env: map[string]string{"APP_PROFILE": "stage", "NUMBER": "99", "PART_STR": "part"},
			want: &Config{
				Number: 99,
				Part:   &ConfigPart{Str: "part"},
			},
		},
		{
			name:         "env with prefix and profile from env",
			resultConfig: &Config{},
			loadSettings: &LoadSettings{
				LoadSysEnv: true,
				EnvPrefix:  "myapp",
				ConfigFile: &ConfigFileSettings{
					ConfigDir: "./testdata/profiles",
				},
			},
			env: map[string]string{"MYAPP_APP_PROFILE": "stage", "NUMBER": "99", "MYAPP_PART_STR": "part"},
			want: &Config{
				Number: 77,
				Part:   &ConfigPart{Str: "part"},
			},
		},
	}

	runTests(t, tests)
}

func runTests(t *testing.T, tests []ConfigTestParam) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range test.env {
				if err := os.Setenv(k, v); err != nil {
					t.Errorf("cant set env var: [key: %s, val: %s]", k, v)
				}
			}
			fmt.Println()
			config := test.resultConfig
			err := Load(config, test.loadSettings)
			if test.wantErr {
				if err == nil {
					t.Error("want error")
				}
				return
			}
			if !reflect.DeepEqual(config, test.want) {
				t.Errorf("got = %+v, want %+v", config, test.want)
			}
		})
	}
}
