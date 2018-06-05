package auth

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ActiveDomain string    `yaml:"active-domain"`
	Domains      []*Domain `yaml:"domains"`
}

func (config *Config) Domain(hostname string) *Domain {
	for _, domain := range config.Domains {
		if domain.Hostname == hostname {
			return domain
		}
	}

	domain := &Domain{
		Hostname: hostname,
	}
	config.Domains = append(config.Domains, domain)

	return domain
}

func (config *Config) SetDomain(hostname, token string) {
	config.ActiveDomain = hostname

	domain := config.Domain(hostname)
	domain.Token = token
	domain.LastUpdated = time.Now()
}

func (config *Config) Write() error {
	content, err := yaml.Marshal(config)
	if err != nil {
		return errors.Trace(err)
	}

	if err := ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".kingrc"), content, 0600); err != nil {
		return errors.Trace(err)
	}

	return nil
}

type Domain struct {
	Hostname    string    `yaml:"hostname"`
	Token       string    `yaml:"token"`
	LastUpdated time.Time `yaml:"last_updated"`
}

func (domain *Domain) IsLocal() (bool, error) {
	addresses, err := net.LookupHost(domain.Hostname)
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, address := range addresses {
		if address == "127.0.0.1" {
			return true, nil
		}
	}

	return false, nil
}

func ReadConfig() (*Config, error) {
	content, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".kingrc"))
	if err != nil {
		if os.IsNotExist(err) {
			return new(Config), nil
		}

		return nil, errors.Trace(err)
	}

	config := new(Config)
	if err := yaml.Unmarshal(content, config); err != nil {
		return nil, errors.Trace(err)
	}

	return config, nil
}
