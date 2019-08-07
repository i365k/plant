package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"plant/other/handler"
	example "plant/other/proto/example"
	"github.com/micro/go-grpc"
	"github.com/robfig/cron"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

func main() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/other/other.log"}`)


	//定时任务
	c:=cron.New()
	c.AddFunc("0 0 0 * * ?", func() {
		o:=orm.NewOrm()
		num,err:=o.QueryTable("strength").Filter("id__gte",0).Delete()
		if err !=nil{
			logs.Info("体力卡包清理失败",err)
		}
		logs.Info("体力卡包列表清理成功",num)
		num,err=o.QueryTable("diamondcard").Filter("id__gte",0).Delete()
		if err !=nil{
			logs.Info("钻石卡包清理失败",err)
		}
		logs.Info("钻石卡包列表清理成功",num)
	})
	c.Start()

	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.other"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
