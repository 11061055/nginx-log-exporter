package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	App []*AppConfig

	original string
}

type HistogramBuckets struct {
	Start float64
	Step float64
	Num int
}

type AppConfig struct {
	Name   string `yaml:"name"`
	Format string `yaml:"format"`

	SourceFiles   []string          `yaml:"source_files"`
	StaticConfig  map[string]string `yaml:"static_config"`
	RelabelConfig *RelabelConfig    `yaml:"relabel_config"`
	HistogramBuckets *HistogramBuckets    `yaml:"histogram_buckets"`
}

func (this *AppConfig) StaticLabelValues() (labels, values []string) {
	labels = make([]string, len(this.StaticConfig))
	values = make([]string, len(this.StaticConfig))

	i := 0
	for k, v := range this.StaticConfig {
		labels[i] = k
		values[i] = v
		i++
	}

	return
}

func (this *AppConfig) DynamicLabels() (labels []string) {
	return this.RelabelConfig.SourceLabels
}

func (this *AppConfig) Prepare() {
	for _, rs := range this.RelabelConfig.Replacements {
	    for _, r := range rs {
		    for _, replaceItem := range r.Repace {
			    replaceItem.prepare()
		    }
		}
	}
}

type RelabelConfig struct {
	SourceLabels []string                   `yaml:"source_labels"`
	Replacements  map[string][]*Replacement `yaml:"replacements"`
}

type Trim struct {
    Sep string `yaml:"sep"`
    Idx int `yaml:"idx"`
}

type Replacement struct {
	Trims   []*Trim        `yaml:"trims"`
	Repace []*RepaceTarget `yaml:"replaces"`
}

type RepaceTarget struct {
	Target string `yaml:"target"`
	Value  string `yaml:"value"`

	tRex *regexp.Regexp
}

func (this *RepaceTarget) Regexp() *regexp.Regexp {
	return this.tRex
}

func (this *RepaceTarget) prepare() {
	replace, err := regexp.Compile(this.Target)
	if err != nil {
		log.Panic(err)
	}

	this.tRex = replace
}

func (this *Config) Reload() error {
	original, err := load(this.original)
	if err != nil {
		return err
	}

	this = original
	return nil
}

func LoadFile(filename string) (conf *Config, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	conf, err = load(string(content))
	return
}

func load(s string) (*Config, error) {
	var (
		cfg  = &Config{}
		apps []*AppConfig
	)

	err := yaml.Unmarshal([]byte(s), &apps)
	if err != nil {
		return nil, err
	}

    jsonAppConfig , err := json.Marshal ( apps )
    fmt.Printf("\nAppConfig is: \n $s", string(jsonAppConfig))


	cfg.original = s
	cfg.App = apps

	return cfg, nil
}
