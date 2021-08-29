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
	loadSettings *LoadSettings
	env          map[string]string
	want         *Config
	wantErr      bool
}

func TestLoadFromFile(t *testing.T) {
	tests := []ConfigTestParam{
		{
			name: "file config",
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
			name: "file config",
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
			name: "file config with profile set by base file",
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
			name: "file config with profile set by base file and env (env should override file)",
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
			name: "file config with prefix",
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
			name: "file config with profile and prefix",
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
			name: "env without prefix",
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
			name: "env with prefix",
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
			name: "env without prefix and profile from file",
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
			name: "env with prefix and profile from file",
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
			name: "env without prefix and profile from env",
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
			name: "env with prefix and profile from env",
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
			config := &Config{}
			fmt.Println()
			err := Load(config, test.loadSettings)
			if test.wantErr {
				if err == nil {
					t.Error("want error")
				}
			}
			if !reflect.DeepEqual(config, test.want) {
				t.Errorf("got = %+v, want %+v", config, test.want)
			}
		})
	}
}
