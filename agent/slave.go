package agent

import (
	"fmt"
	"github.com/k3rn3l-p4n1c/entanglement"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func StartSlave(conf *viper.Viper) error {
	selfHost := conf.GetString("node.ip")
	entangleConf := entanglement.DefaultConfig()
	logrus.Debug("trying to join ", conf.GetString("slave.join"))
	entangleConf.JoinAddr = fmt.Sprintf("%s:%d", conf.GetString("slave.join"), 12701)

	entangleConf.RaftAddr = fmt.Sprintf("%s:%d", selfHost, 12700)
	entangleConf.HttpAddr = fmt.Sprintf("%s:%d", selfHost, 12701)

	system := entanglement.Bootstrap(entangleConf)
	instance = &Agent{
		ProjectName: system.New("projectName"),
	}

	return nil
}
