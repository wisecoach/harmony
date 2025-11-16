package main

import (
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type Validator struct {
	Index        int
	Address      string // account HMY address
	EthAddr      string // account ETH address
	BLSPublicKey string // account public BLS key
	ShardID      uint32 // shardID of the account
	EcdsaKeyPath string
	BLSKeyPATH   string
}

func writeYaml(path string, data any) {
	b, _ := yaml.Marshal(data)
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer file.Close()

	// 将字符串写入文件
	_, _ = file.Write(b)
}

type ShardConfig struct {
	name      string // test name
	shard     uint32 // shard num
	validator int    // validator num per shard
	ssc       int    // ssc member num per shard
	delay     uint64
}

func buildConfig(config ShardConfig) {
	dev_path := fmt.Sprintf("test/configs/%s/shard=%d_validator=%d_ssc=%d_delay=%d", "dev", config.shard, config.validator, config.ssc, config.delay)
	local_path := fmt.Sprintf("test/configs/%s/shard=%d_validator=%d_ssc=%d_delay=%d", "local", config.shard, config.validator, config.ssc, config.delay)
	os.MkdirAll(dev_path, 0755)
	os.MkdirAll(local_path, 0755)

	shard2validators := make(map[uint32][]*Validator)
	file, _ := os.Open(".hmy/validators.json")
	defer file.Close()
	json.NewDecoder(file).Decode(&shard2validators)

	local_launch_config_lines := make([]string, 0)
	dev_launch_config_lines := make([]string, 0)
	for i := uint32(0); i < config.shard; i++ {
		for j := 0; j < config.validator; j++ {
			v := shard2validators[i][j]
			local_launch_config_lines = append(local_launch_config_lines, fmt.Sprintf("%s %s %s %d %s %d",
				v.Address, v.EthAddr, v.BLSKeyPATH, v.ShardID, "127.0.0.1", 9000+40*int(i)+j*2))
			dev_launch_config_lines = append(dev_launch_config_lines, fmt.Sprintf("%s %s %s %d %s %d",
				v.Address, v.EthAddr, v.BLSKeyPATH, v.ShardID, fmt.Sprintf("10.7.95.%d", 200+i), 9000+40*int(i)+j*2))
		}
	}
	os.WriteFile(local_path+"/"+"launch_config_local.txt", []byte(strings.Join(local_launch_config_lines, "\n")), 0644)
	os.WriteFile(dev_path+"/"+"launch_config_dev.txt", []byte(strings.Join(dev_launch_config_lines, "\n")), 0644)

	type GenesisConfig struct {
		GenesisAccountsDir string `json:"genesis_accounts_dir" yaml:"genesis_accounts_dir"`
	}
	local_genesis := &GenesisConfig{
		GenesisAccountsDir: ".hmy/expr_accounts",
	}
	dev_genesis := &GenesisConfig{
		GenesisAccountsDir: ".hmy/expr_accounts",
	}
	writeYaml(local_path+"/"+"genesis_config_local.yaml", local_genesis)
	writeYaml(dev_path+"/"+"genesis_config_dev.yaml", dev_genesis)

	localTomlConfig, _ := toml.LoadFile("test/configs/local/default_config_local.toml")
	devTomlConfig, _ := toml.LoadFile("test/configs/dev/default_config_dev.toml")
	localGeneral := localTomlConfig.Get("General").(*toml.Tree)
	devGeneral := devTomlConfig.Get("General").(*toml.Tree)
	localGeneral.Set("GenesisConfigFile", local_path+"/"+"genesis_config_local.yaml")
	devGeneral.Set("GenesisConfigFile", dev_path+"/"+"genesis_config_dev.yaml")
	localToml, _ := os.Create(local_path + "/" + "default_config_local.toml")
	devToml, _ := os.Create(dev_path + "/" + "default_config_dev.toml")
	localTomlConfig.Set("General", localGeneral)
	localTomlConfig.WriteTo(localToml)
	devTomlConfig.Set("General", devGeneral)
	devTomlConfig.WriteTo(devToml)
}

func main() {
	configs := []ShardConfig{
		{
			name:      "基准测试",
			shard:     4,
			validator: 4,
			ssc:       1,
			delay:     5,
		},
		{
			name:      "基准测试，时延10",
			shard:     4,
			validator: 4,
			ssc:       1,
			delay:     10,
		},
		{
			name:      "基准测试，时延20",
			shard:     4,
			validator: 4,
			ssc:       1,
			delay:     20,
		},
		{
			name:      "不同分片数_2",
			shard:     2,
			validator: 4,
			ssc:       1,
			delay:     5,
		},
		{
			name:      "不同分片数_8",
			shard:     8,
			validator: 4,
			ssc:       1,
			delay:     5,
		},
		{
			name:      "不同分片数_16",
			shard:     16,
			validator: 4,
			ssc:       1,
			delay:     5,
		},
		{
			name:      "不同分片数_32",
			shard:     32,
			validator: 4,
			ssc:       1,
			delay:     5,
		},
		{
			name:      "安全性测试",
			shard:     2,
			validator: 10,
			ssc:       4,
			delay:     5,
		},
	}
	for _, config := range configs {
		buildConfig(config)
	}
}
