package common

import "github.com/spf13/viper"

type config struct {
	*viper.Viper
}

var BaseDir string

func init() {
	BaseDir = GetConfig("common.basedir")
}

func loadConfigFromYaml(c *config) error {
	c.Viper = viper.New()
	c.SetConfigType("yaml")
	c.SetConfigName("config")
	c.AddConfigPath("../")
	c.AddConfigPath("./")
	c.AddConfigPath("../dataserver")
	if err := c.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func GetConfig(arg string) string {
	c := &config{}
	if err := loadConfigFromYaml(c); err != nil {
		return ""
	}
	return c.GetString(arg)
}

func GetConfigArray(arg string) []string {
	c := &config{}
	if err := loadConfigFromYaml(c); err != nil {
		return nil
	}
	return c.GetStringSlice(arg)
}
