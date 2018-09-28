package agent

import (
	"github.com/k3rn3l-p4n1c/entanglement"
	"github.com/k3rn3l-p4n1c/goaway/configuration"
	"github.com/sirupsen/logrus"
)

var (
	instance *Agent
	//once     sync.Once
)

type Agent struct {
	ProjectName *entanglement.Entanglement
}

func Start() {
	conf := configuration.GetInstance()

	isMaster := conf.GetBool("master.bootstrap")
	isSlave := conf.IsSet("slave.join")

	if isMaster {
		logrus.Debug("starting master agent")
		StartMaster(conf)
	} else if isSlave {
		logrus.Debug("starting slave agent")
		StartSlave(conf)
		//panic("not implemented")
	} else {
		panic("not master nor slave")
	}
}

func GetInstance() *Agent {
	//once.Do(func() {
	//	conf := entanglement.DefaultConfig()
	//	system := entanglement.Bootstrap(conf)
	//	instance = &Agent{
	//		ProjectName: system.New("projectName"),
	//	}
	//})
	return instance

}
