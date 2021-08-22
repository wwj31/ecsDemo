package main

import (
	"ecsDemo/graphUI"
	"ecsDemo/world"
	"github.com/wwj31/dogactor/actor"
	"github.com/wwj31/dogactor/actor/cmd"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)
var (
	ETCD_ADDR   = "127.0.0.1:2379"
	ETCD_PREFIX = "demo/"
)

func main()  {
	log.Init(log.TAG_DEBUG_I, nil, "./_log", "demo", 1)
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sys, err := actor.NewSystem(
		actor.Addr("127.0.0.1:8888"),
		actor.WithCMD(cmd.New()),
		//cluster.WithRemote(ETCD_ADDR, ETCD_PREFIX),
	)
	expect.Nil(err)

	center := actor.New("CenterWorld", &world.CenterWorld{})
	sys.Regist(center)
	time.Sleep(2*time.Second)
	sys.Regist(actor.New("GraphUI", &graphUI.GraphUI{}, actor.SetMailBoxSize(10000)))
	graphUI.Init(exit)

	<-exit
	sys.Stop()
	log.Stop()
	return
}