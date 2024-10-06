package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var Configs = DefaultConfig()

// Config 定义配置结构体
type Config struct {
	ServerAddress string `json:"server_address"`
	StorageRoot   string `json:"storage_root"`
	Rabbitmq      struct {
		RabbitmqServer string `json:"rabbitmq_server"`
	} `json:"rabbitmq"`
}

// DefaultConfig 返回默认的配置
func DefaultConfig() Config {
	return Config{
		ServerAddress: ":8080",
		StorageRoot:   "./storage",
		Rabbitmq: struct {
			RabbitmqServer string `json:"rabbitmq_server"`
		}{
			RabbitmqServer: "amqp://guest:guest@127.0.0.1:5672/",
		},
	}
}

// LoadConfig 尝试加载配置文件，如果文件不存在，则创建默认配置
func LoadConfig(filename string) (Config, error) {
	var config Config

	// 检查配置文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// 如果文件不存在，创建默认配置文件
		config = DefaultConfig()
		err := SaveConfig(filename, config)
		if err != nil {
			return config, fmt.Errorf("failed to create default config: %v", err)
		}
		log.Printf("Default config created at %s\n", filename)
		return config, nil
	}

	// 加载配置文件
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %v", err)
	}

	// 解析配置文件内容
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse config file: %v", err)
	}

	log.Printf("Config loaded from %s\n", filename)
	return config, nil
}

// SaveConfig 将配置保存到文件
func SaveConfig(filename string, config Config) error {
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %v", err)
	}

	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	return nil
}

func Init() {
	// 配置文件路径
	configFile := "config.json"

	// 加载或创建配置文件
	var err error
	Configs, err = LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	//// 更新配置示例
	//Configs.ServerAddress = "9099"
	//err = SaveConfig(configFile, Configs)
	//if err != nil {
	//	log.Fatalf("Error saving config: %v", err)
	//}
}
