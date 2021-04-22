package config

import "gopkg.in/ini.v1"

type Config struct {
	FfmpegBinPath  string `ini:"FfmpegBinPath"`
	FfprobeBinPath string `ini:"FfprobeBinPath"`
	Format         string `ini:"Format"`
	Resolution     []int
	Bitrate        []string
}

func Load(path string) (*Config, error) {
	conf := Config{}
	if err := ini.MapTo(&conf, path); err != nil {
		return nil, err
	}

	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}
	conf.Resolution = cfg.Section("").Key("Resolution").Ints(",")
	conf.Bitrate = cfg.Section("").Key("Bitrate").Strings(",")

	return &conf, nil
}
