package agent

import (
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"github.com/k3rn3l-p4n1c/entanglement"
	"github.com/k3rn3l-p4n1c/goaway/scheduler"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
)

func StartMaster(conf *viper.Viper) error {
	selfHost := conf.GetString("node.ip")

	isBootstrap := conf.GetBool("master.bootstrap")
	if !isBootstrap {
		panic("NOT IMPLEMENTED YET!!! Multi-master is not supported in this version")
	}

	entangleConf := entanglement.DefaultConfig()

	entangleConf.RaftAddr = fmt.Sprintf("%s:%d", selfHost, 12700)
	entangleConf.HttpAddr = fmt.Sprintf("%s:%d", selfHost, 12701)

	system := entanglement.Bootstrap(entangleConf)
	instance = &Agent{
		Cluster: nil,
		ProjectName: system.New("projectName"),
	}

	mem := sigar.Mem{}

	cluster := &scheduler.Cluster{}
	if err := cluster.AddOrUpdateDataCenter(conf.GetString("node.datacenter")); err != nil {
		logrus.WithError(err).Fatal("unable to register self in cluster")
		return err
	}
	if err := cluster.AddOrUpdateServer(conf.GetString("node.name"),
		conf.GetString("node.datacenter"),
		mem.Total,
		net.ParseIP(selfHost)); err != nil {
		logrus.WithError(err).Fatal("unable to register self in cluster")
		return err
	}
	instance.Cluster = cluster
	GetInstance().ProjectName.Set(conf.GetString("name"))
	return nil
}
