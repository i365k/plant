package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-web"
	"github.com/julienschmidt/httprouter"
	"plant/gateway/handler"
)

func main() {
	// create new web service
	//创建一个新的web服务
	service := web.NewService(
		web.Name("go.micro.web.gateway"),
		web.Version("latest"),
		web.Address(":8888"),
	)

	// initialise service
	//服务初始化
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	//使用路由中间件映射页面
	rou := httprouter.New()

	//映射静态页面
	//rou.NotFound = http.FileServer(http.Dir("html"))

	//ping健康检查
	rou.HEAD("/ping",handler.Ping)

	//登陆服务
	rou.POST("/api/plant/Login", handler.Login)
	rou.POST("/api/plant/login", handler.Login)

	//初始化服务
	rou.POST("/api/plant/postconfig",handler.Config)
	rou.POST("/api/plant/config",handler.Config)

	//时间戳服务
	rou.POST("/api/plant/nowtime",handler.NowTime)

	//获取公告服务
	rou.POST("/api/plant/TestServer", handler.Notice)
	rou.POST("/api/plant/notice", handler.Notice)

	//同步数据服务
	rou.POST("/api/plant/posttest", handler.Data)
	rou.POST("/api/plant/data", handler.Data)

	//获取邀请新人列表
	rou.POST("/api/plant/pullnew",handler.PullNew)

	//钻石助力服务
	rou.POST("/api/plant/diamond",handler.Diamond)

	//获取钻石助力列表服务
	rou.POST("/api/plant/diamlest",handler.Diamlest)

	//体力助力服务
	rou.POST("/api/plant/strength",handler.Strength)

	//获取体力助力数据服务
	rou.POST("/api/plant/strelest",handler.Strelest)

	//灯助力服务
	rou.POST("/api/plant/lamp",handler.Lamp)

	//获取灯助力数据服务
	rou.POST("/api/plant/lamplest",handler.Lamplest)

	//宝图服务
	rou.POST("/api/plant/treasuremap",handler.TreasureMap)
	rou.POST("/api/plant/treasuremapa",handler.TreasureMap)

	//获取灯助力数据服务
	rou.POST("/api/plant/achievement",handler.Achievement)

	//##########################################################

	//登陆服务
	rou.POST("/api/lamp/login", handler.LoginLamp)
	rou.POST("/api/lamp/Login", handler.LoginLamp)

	//初始化服务
	rou.POST("/api/lamp/config",handler.ConfigLamp)

	//获取邀请新人列表
	rou.POST("/api/lamp/pullnew",handler.PullNewLamp)

	//金币变化服务
	rou.POST("/api/lamp/gold",handler.GoldLamp)

	//服务器排名服务
	rou.POST("/api/lamp/ranking",handler.RankingLamp)

	//关卡服务
	rou.POST("/api/lamp/checkpoint",handler.CheckpointLamp)

	//设置灯服务
	rou.POST("/api/lamp/setlamp",handler.SetLamp)
	rou.POST("/get",handler.SetLamp)

	//获取灯设置服务
	rou.POST("/api/lamp/getlamp",handler.GetLamp)
	rou.POST("/set",handler.GetLamp)

	//购买灯的服务
	rou.POST("/api/lamp/buylamp",handler.BuyLamp)
	// register html handler
	service.Handle("/", rou)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}