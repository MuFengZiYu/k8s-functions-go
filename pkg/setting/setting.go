package setting

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Set struct {
	Redis Redis
}

type Redis struct {
	Host      string `yaml:"host"`
	Password  string `yaml:"password"`
	Timeout   int    `yaml:"timeout"`
	MaxActive int    `yaml:"max_active"`
	MaxIdle   int    `yaml:"max_idle"`
	Db        int
}

var Setting = Set{}

func InitSetting() {
	file, err := ioutil.ReadFile("./conf/config.yml")
	if err != nil {
		log.Fatal("fail to read file:", err)
	}

	err = yaml.Unmarshal(file, &Setting)
	if err != nil {
		log.Fatal("fail to yaml unmarshal:", err)
	}

}
