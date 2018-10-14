package file

import (
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

func Read(filename string) (StackFile, error) {
	t := StackFile{}

	data, err := ioutil.ReadFile(filename)
	check(err)
	fmt.Print(string(data))

	err = yaml.Unmarshal([]byte(data), &t)

	return t, err
}
