package file

import (
	"log"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Service struct {
	Name        string
	DockerImage string `yaml:"docker-image"`
}

type StackFile struct {
	Services []Service
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Run() {
	t := StackFile{}

	data, err := ioutil.ReadFile("./file/dummy.yml")
	check(err)
	fmt.Print(string(data))

	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", t)

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", m)

	d, err = yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}
