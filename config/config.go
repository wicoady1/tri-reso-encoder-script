package config

import "gopkg.in/ini.v1"

type Config struct {
	FfmpegBinPath  string `ini:"FfmpegBinPath"`
	FfprobeBinPath string `ini:"FfprobeBinPath"`
}

func Load(path string) (*Config, error) {
	conf := Config{}
	if err := ini.MapTo(&conf, path); err != nil {
		return nil, err
	}

	return &conf, nil
}
