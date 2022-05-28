package core

import (
	"fmt"
	"io"
	"os"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v2"

	"github.com/miniscruff/changie/shared"
)

const configEnvVar = "CHANGIE_CONFIG_PATH"
const CreateFileMode os.FileMode = 0644
const CreateDirMode os.FileMode = 0755
const timeFormat string = "20060102-150405"

var ConfigPaths []string = []string{
	".changie.yaml",
	".changie.yml",
}

// GetVersions will return, in semver sorted order, all released versions
type GetVersions func(shared.ReadDirer, Config) ([]*semver.Version, error)

type KindConfig struct {
	Label             string   `yaml:"label"`
	Format            string   `yaml:"format,omitempty"`
	ChangeFormat      string   `yaml:"changeFormat,omitempty"`
	SkipGlobalChoices bool     `yaml:"skipGlobalChoices,omitempty"`
	SkipBody          bool     `yaml:"skipBody,omitempty"`
	AdditionalChoices []Custom `yaml:"additionalChoices,omitempty"`
}

func (kc KindConfig) String() string {
	return kc.Label
}

type BodyConfig struct {
	MinLength *int64 `yaml:"minLength,omitempty"`
	MaxLength *int64 `yaml:"minLength,omitempty"`
}

func (b BodyConfig) CreatePrompt(stdinReader io.ReadCloser) Prompt {
	p, _ := Custom{
		Label:     "Body",
		Type:      CustomString,
		MinLength: b.MinLength,
		MaxLength: b.MaxLength,
	}.CreatePrompt(stdinReader)

	return p
}

type NewLineConfig struct {
	BeforeVersion    int `yaml:"beforeVersion,omitempty"`
	AfterVersion     int `yaml:"afterVersion,omitempty"`
	BeforeComponent  int `yaml:"beforeComponent,omitempty"`
	AfterComponent   int `yaml:"afterComponent,omitempty"`
	BeforeHeader     int `yaml:"beforeHeader,omitempty"`
	AfterHeader      int `yaml:"afterHeader,omitempty"`
	BeforeFooter     int `yaml:"beforeFooter,omitempty"`
	AfterFooter      int `yaml:"afterFooter,omitempty"`
	BeforeHeaderFile int `yaml:"beforeHeaderFile,omitempty"`
	AfterHeaderFile  int `yaml:"afterHeaderFile,omitempty"`
	BeforeFooterFile int `yaml:"beforeFooterFile,omitempty"`
	AfterFooterFile  int `yaml:"afterFooterFile,omitempty"`
	BeforeKindHeader int `yaml:"beforeKindHeader,omitempty"`
	AfterKindHeader  int `yaml:"afterKindHeader,omitempty"`
	BeforeChanges    int `yaml:"beforeChanges,omitempty"`
	AfterChanges     int `yaml:"afterChanges,omitempty"`
}

// Config handles configuration for a changie project
type Config struct {
	ChangesDir        string `yaml:"changesDir"`
	UnreleasedDir     string `yaml:"unreleasedDir"`
	HeaderPath        string `yaml:"headerPath"`
	VersionHeaderPath string `yaml:"versionHeaderPath"`
	VersionFooterPath string `yaml:"versionFooterPath"`
	ChangelogPath     string `yaml:"changelogPath"`
	VersionExt        string `yaml:"versionExt"`
	// formats
	FragmentFileFormat string `yaml:"fragmentFileFormat,omitempty"`
	VersionFormat      string `yaml:"versionFormat"`
	ComponentFormat    string `yaml:"componentFormat,omitempty"`
	KindFormat         string `yaml:"kindFormat,omitempty"`
	ChangeFormat       string `yaml:"changeFormat"`
	HeaderFormat       string `yaml:"headerFormat"`
	FooterFormat       string `yaml:"footerFormat"`
	// custom
	Body          BodyConfig    `yaml:"body,omitempty"`
	Components    []string      `yaml:"components,omitempty"`
	Kinds         []KindConfig  `yaml:"kinds,omitempty"`
	CustomChoices []Custom      `yaml:"custom,omitempty"`
	Replacements  []Replacement `yaml:"replacements,omitempty"`
	NewLines      NewLineConfig `yaml:"newlines,omitempty"`
}

func (c Config) KindHeader(label string) string {
	for _, kindConfig := range c.Kinds {
		if kindConfig.Format != "" && kindConfig.Label == label {
			return kindConfig.Format
		}
	}

	return c.KindFormat
}

func (c Config) ChangeFormatForKind(label string) string {
	for _, kindConfig := range c.Kinds {
		if kindConfig.ChangeFormat != "" && kindConfig.Label == label {
			return kindConfig.ChangeFormat
		}
	}

	return c.ChangeFormat
}

// Save will save the config as a yaml file to the default path
func (c Config) Save(wf shared.WriteFiler) error {
	bs, _ := yaml.Marshal(&c)
	return wf(ConfigPaths[0], bs, CreateFileMode)
}

// LoadConfig will load the config from the default path
func LoadConfig(rf shared.ReadFiler) (Config, error) {
	var (
		c   Config
		bs  []byte
		err error
	)

	customPath := os.Getenv(configEnvVar)
	if customPath != "" {
		bs, err = rf(customPath)
	} else {
		for _, path := range ConfigPaths {
			bs, err = rf(path)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(bs, &c)
	if err != nil {
		return c, err
	}

	// load backward incompatible configs
	if c.FragmentFileFormat == "" {
		if len(c.Components) > 0 {
			c.FragmentFileFormat = "{{.Component}}-"
		}

		if len(c.Kinds) > 0 {
			c.FragmentFileFormat += "{{.Kind}}-"
		}

		c.FragmentFileFormat += fmt.Sprintf("{{.Time.Format \"%v\"}}", timeFormat)
	}

	return c, nil
}
