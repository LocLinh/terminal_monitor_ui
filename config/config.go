package config

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

// App config struct
type AppConfig struct {
	Kafka          KafkaConfig     `yaml:"Kafka"`
	SqlServer      SqlServerConfig `yaml:"SqlServer"`
	ConveyorDevice string          `yaml:"ConveyorDevice"`
}

type KafkaConfig struct {
	Addrs           []string    `yaml:"Addrs"`
	Group           string      `yaml:"Group"`
	MaxMessageBytes int         `yaml:"MaxMessageBytes"`
	Compress        bool        `yaml:"Compress"`
	Newest          bool        `yaml:"Newest"`
	Version         string      `yaml:"Version"`
	Topics          []string    `yaml:"Topics"`
	TopicNames      KafkaTopics `yaml:"TopicNames"`
	Partition       int         `yaml:"Partition"`
	GroupId         string
}
type ConsumerTopics struct {
	HasakiNowOrder string `yaml:"HasakiNowOrder"`
	InsideStore    string `yaml:"InsideStore"`
}
type SqlServerConfig struct {
	ServerName string `yaml:"ServerName"`
	User       string `yaml:"User"`
	Password   string `yaml:"Password"`
	Port       string `yaml:"Port"`
	Database   string `yaml:"Database"`
}

type KafkaTopics struct {
	NowOrder            string `yaml:"NowOrder"`
	InsideStore         string `yaml:"InsideStore"`
	NowPickupLocation   string `yaml:"NowPickupLocation"`
	ShippingOrder       string `yaml:"ShippingOrder"`
	MappingTrackingCode string `yaml:"MappingTrackingCode"`
	User                string `yaml:"User"`
	ConveyorOrder       string `yaml:"ConveyorOrder"`
	ProductTester       string `yaml:"ProductTester"`
}

func (k KafkaTopics) ToList() []string {
	var ret []string
	v := reflect.ValueOf(k)

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.IsZero() {
			ret = append(ret, fmt.Sprintf("%v", f.Interface()))
		}
	}
	return ret
}

// Load config file from given path
func LoadConfig() (*AppConfig, error) {
	appConfig := &AppConfig{}
	file, err := os.Open("./config/config.yml")
	if err != nil {
		log.Default().Fatalf("Error opening config.yml: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&appConfig); err != nil {
		return nil, err
	}

	return appConfig, nil
}
