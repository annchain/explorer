package main

import (
	"github.com/annchain/explorer/job"
	"github.com/annchain/explorer/repository"
	_ "github.com/annchain/explorer/routers"
	"github.com/annchain/explorer/rpc"
	"github.com/astaxie/beego"
)

func main() {
	rpc.Init()
	repository.Init()
	go job.SyncTimingTask()
	beego.Run()
}
