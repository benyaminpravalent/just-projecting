package config

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultType  = "json"
	configURLEnv = "CONFIG_URL"
)

var config *viper.Viper

//Load load configuration from config url, by default will load environment variable
func Load(def map[string]interface{}, urlStr string) error {

	// first lets load .env file
	econf := viper.New()
	config = viper.New()

	for k, v := range def {
		econf.BindEnv(k)
		config.SetDefault(k, v)
		if econf.IsSet(k) {
			config.Set(k, econf.Get(k))
		}
	}

	if urlStr == "" {
		urlStr = os.Getenv(configURLEnv)
	}

	if urlStr == "" {
		return nil
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	aconf := viper.New()

	switch u.Scheme {
	case "file":
		path := strings.TrimPrefix(urlStr, "file://")
		name := filepath.Base(path)
		t := strings.TrimPrefix(filepath.Ext(path), ".")
		path = filepath.Dir(path)
		aconf.SetConfigName(name)
		aconf.SetConfigType(t)
		aconf.AddConfigPath(path)
		if err := aconf.ReadInConfig(); err != nil {
			return err
		}
	case "consul":
		host := u.Host
		key := strings.TrimPrefix(u.Path, "/")
		aconf.AddRemoteProvider("consul", host, key)
		aconf.SetConfigType(defaultType) // Need to explicitly set this to json
		if err := aconf.ReadRemoteConfig(); err != nil {
			return err
		}
	case "etcd":
		host := u.Host
		path := u.Path
		aconf.AddRemoteProvider("etcd", host, path)
		aconf.SetConfigType(defaultType)
		if err := aconf.ReadRemoteConfig(); err != nil {
			return err
		}
	default:
		return errors.New("Unsupported config scheme")
	}

	for k := range def {
		if aconf.IsSet(k) {
			config.Set(k, aconf.Get(k))
		}
	}

	return nil
}

//Get get interface{}
func Get(k string) interface{} {
	return config.Get(k)
}

//GetString get string
func GetString(k string) string {
	return config.GetString(k)
}

//GetBool get bool
func GetBool(k string) bool {
	return config.GetBool(k)
}

//GetInt get int
func GetInt(k string) int {
	return config.GetInt(k)
}

//GetFloat64 get float64
func GetFloat64(k string) float64 {
	return config.GetFloat64(k)
}

//GetStringSlice get []string
func GetStringSlice(k string) []string {
	return config.GetStringSlice(k)
}

//GetStringMapString get map[string]string
func GetStringMapString(k string) map[string]string {
	return config.GetStringMapString(k)
}

// GetStringMap get map[string]interface{}
func GetStringMap(k string) map[string]interface{} {
	return config.GetStringMap(k)
}
