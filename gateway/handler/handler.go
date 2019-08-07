package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	LOGIN "plant/login/proto/example"
	DATA "plant/data/proto/example"
	OTHER "plant/other/proto/example"
	"github.com/micro/go-grpc"
	"github.com/julienschmidt/httprouter"
	"plant/gateway/utils"
	"io/ioutil"
	"github.com/goEncrypt"
	"encoding/hex"
	"strconv"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
)

//ping健康检查
func Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//返回数据给前端
	response := map[string]interface{}{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//用户登陆
func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("登陆服务  url : /api/plant/Login")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logs.Info("转换错误")
	}

	//创建服务句柄
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := LOGIN.NewExampleService("go.micro.srv.login", server.Client())
	rsp, err := exampleClient.Login(context.TODO(), &LOGIN.LoginRequest{
		Data: body,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"skey":   rsp.Skey,
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//设置放回json格式
	w.Header().Set("Content-Type", "application/json")

	//将返回数据转为json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//初始化服务
func Config(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("初始化服务 url：/api/plant/config")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建连接
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := LOGIN.NewExampleService("go.micro.srv.login", service.Client())
	rsp, err := exampleClient.Config(context.TODO(), &LOGIN.ConfigRequest{
		Skey: request["skey"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var resp = make(map[string]interface{})
	json.Unmarshal(rsp.Datalist, &resp)

	aaa := resp["data"]

	// we want to augment the response
	response := map[string]interface{}{
		"configlist": rsp.Configlist,
		"datalist":   aaa,
		"errno":      rsp.Errno,
		"errmsg":     rsp.Errmsg,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//时间戳服务
func NowTime(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("时间戳服务 url：/api/plant/nowtime")

	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//对skey进行判断
	tempskey := request["skey"]
	if len(tempskey.(string)) < 10 {

		response := map[string]interface{}{
			"errno":  utils.RECODE_PARAMERR,
			"errmsg": utils.RecodeText(utils.RECODE_PARAMERR),
		}

		w.Header().Set("Content-Type", "application/json")

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 501)
			return
		}
		return
	}

	//准备时间戳
	nowtime := time.Now().Unix()
	//准备加密
	Key := []byte("1234567887654321")
	ciphernowtime := goEncrypt.AesCBC_Encrypt([]byte(strconv.Itoa(int(nowtime))), Key)
	encodenowtime := hex.EncodeToString(ciphernowtime)

	// we want to augment the response
	response := map[string]interface{}{
		"nowtime": encodenowtime,
		"ref":     time.Now().UnixNano(),
		"errno":   utils.RECODE_OK,
		"errmsg":  utils.RecodeText(utils.RECODE_OK),
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//获取公告服务
func Notice(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("获取公告服务 url：/api/plant/TestServer")

	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//判断skey是否合法
	tempskey := request["skey"]
	if len(tempskey.(string)) < 10 {
		logs.Info("请求参数不合法")
		response := map[string]interface{}{
			"errno":  utils.RECODE_PARAMERR,
			"errmsg": utils.RecodeText(utils.RECODE_PARAMERR),
		}

		w.Header().Set("Content-Type", "application/json")

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 501)
			return
		}
		return
	}

	//准备公告数据
	type Notice struct {
		Code   string
		Notice string
	}
	tempnot := Notice{Code: "10002", Notice: "亲爱的玩家，欢迎您来到花园大亨"}

	//返回数据给前端
	response := map[string]interface{}{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
		"data":   tempnot,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//同步数据服务
func Data(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("同步数据服务 url : /api/plant/PostTest")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logs.Info("转换错误")
	}
	//创建连接
	service := grpc.NewService()
	service.Init()

	exampleClient := DATA.NewExampleService("go.micro.srv.data", service.Client())
	rsp, err := exampleClient.Data(context.TODO(), &DATA.Request{
		Data: body,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//邀请新人列表
func PullNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("邀请新人列表服务 url：/api/plant/pullnew")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for k, v := range request {
		logs.Info(k, v)
	}

	//创建服务句柄
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := OTHER.NewExampleService("go.micro.srv.other", server.Client())
	rsp, err := exampleClient.PullNew(context.TODO(), &OTHER.PullNewRequest{
		Skey: request["skey"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"list":   rsp.Newlist,
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//设置放回json格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//体力助力服务
func Strength(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("体力助力服务 url：/api/plant/strength")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建连接
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := OTHER.NewExampleService("go.micro.srv.other", service.Client())
	rsp, err := exampleClient.Strength(context.TODO(), &OTHER.StrengthRequest{
		Skey:      request["skey"].(string),
		Shareskey: request["shareskey"].(string),
		Sharetime: request["sharetime"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//获取体力助力数据服务
func Strelest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("获取体力助力数据服务 url：/api/plant/strelest")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建连接
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := OTHER.NewExampleService("go.micro.srv.other", service.Client())
	rsp, err := exampleClient.Strelest(context.TODO(), &OTHER.StrelestRequest{
		Skey: request["skey"].(string),
		Url:  request["url"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"list":   rsp.Newlist,
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//钻石卡包服务
func Diamond(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("钻石卡包服务 url：/api/plant/nowtime")

	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建连接
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := OTHER.NewExampleService("go.micro.srv.other", service.Client())
	rsp, err := exampleClient.Diamond(context.TODO(), &OTHER.DiamondRequest{
		Skey:      request["skey"].(string),
		Shareskey: request["shareskey"].(string),
		Sharetime: request["sharetime"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//获取钻石助力列表服务
func Diamlest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("获取钻石助力列表服务 url：/api/plant/diamlest")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建连接
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := OTHER.NewExampleService("go.micro.srv.other", service.Client())
	rsp, err := exampleClient.Diamlest(context.TODO(), &OTHER.DiamlestRequest{
		Skey: request["skey"].(string),
		Url:  request["url"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"list":   rsp.Newlist,
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//灯助力服务
func Lamp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("灯助力服务 url：/api/plant/lamp")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建连接
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := OTHER.NewExampleService("go.micro.srv.other", service.Client())
	rsp, err := exampleClient.Lamp(context.TODO(), &OTHER.LampRequest{
		Skey:      request["skey"].(string),
		Shareskey: request["shareskey"].(string),
		Sharetime: request["sharetime"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//获取灯助力数据服务
func Lamplest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("获取灯助力数据服务 url：/api/plant/lamplest")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建连接
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := OTHER.NewExampleService("go.micro.srv.other", service.Client())
	rsp, err := exampleClient.Lamplest(context.TODO(), &OTHER.LamplestRequest{
		Skey: request["skey"].(string),
		Type: request["type"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"list":   rsp.Newlist,
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//宝图服务
func TreasureMap(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("宝图服务 url：/api/plant/treasuremap")
	beego.Info("宝图服务 url：/api/plant/treasuremap")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if _, ok := request["skey"].(string); !ok {
	   logs.Info("请求参数格式错误")
	}

	//创建连接
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := OTHER.NewExampleService("go.micro.srv.other", service.Client())
	rsp, err := exampleClient.TreasureMap(context.TODO(), &OTHER.TreasureMapRequest{
		Skey:      request["skey"].(string),
		Iv:        request["iv"].(string),
		Data:      request["data"].(string),
		Code:      request["code"].(string),
		Shareskey: request["shareskey"].(string),
		Url: request["url"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"list":   rsp.Maplist,
		"goldTime":rsp.Goldtime,
		"name":rsp.Name,
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//获取灯助力数据服务
func Achievement(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/gateway/gateway.log"}`)
	logs.Info("获取灯助力数据服务 url：/api/plant/lamplest")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建连接
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := OTHER.NewExampleService("go.micro.srv.other", service.Client())
	rsp, err := exampleClient.Achievement(context.TODO(), &OTHER.AchievementRequest{
		Skey: request["skey"].(string),
		Code: request["code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"reward":   rsp.Reward,
		"nextnumber":   rsp.NextNumber,
		"name":   rsp.Name,
		"conditionA":   rsp.ConditionA,
		"conditionB":   rsp.ConditionB,
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//##################################################################################

//用户登陆
func LoginLamp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("登陆服务  url : /api/plant/Login")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logs.Info("转换错误")
	}

	//创建服务句柄
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := DATA.NewExampleService("go.micro.srv.data", server.Client())
	rsp, err := exampleClient.LoginLamp(context.TODO(), &DATA.LoginLampRequest{
		Data: body,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"skey":   rsp.Skey,
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//设置放回json格式
	w.Header().Set("Content-Type", "application/json")

	//将返回数据转为json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//邀请新人列表
func PullNewLamp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("邀请新人列表服务 url：/api/plant/pullnew")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建服务句柄
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := DATA.NewExampleService("go.micro.srv.data", server.Client())
	rsp, err := exampleClient.PullNewLamp(context.TODO(), &DATA.PullNewLampRequest{
		Skey: request["skey"].(string),
		Url:  request["url"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"usermoney": rsp.UserMoney,
		"list":      rsp.Newlist,
		"errno":     rsp.Errno,
		"errmsg":    rsp.Errmsg,
	}

	//设置放回json格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//初始化服务
func ConfigLamp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("初始化服务 url：/api/plant/config")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建连接
	service := grpc.NewService()
	service.Init()

	// call the backend service
	exampleClient := DATA.NewExampleService("go.micro.srv.data", service.Client())
	rsp, err := exampleClient.ConfigLamp(context.TODO(), &DATA.ConfigLampRequest{
		Skey: request["skey"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//var resp = make(map[string]interface{})
	//json.Unmarshal(rsp.Datalist, &resp)
	//
	//aaa := resp["data"]

	// we want to augment the response
	response := map[string]interface{}{
		"datalist": rsp.Datalist,
		"errno":    rsp.Errno,
		"errmsg":   rsp.Errmsg,
	}

	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//金币变化接口
func GoldLamp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("金币变化接口 url：/api/lamp/gold")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建服务句柄
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := DATA.NewExampleService("go.micro.srv.data", server.Client())
	rsp, err := exampleClient.GoldLamp(context.TODO(), &DATA.GoldLampRequest{
		Skey:   request["skey"].(string),
		Type:   request["type"].(string),
		Number: request["number"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"list":   rsp.Newlist,
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//设置返回json格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//服务器排名服务
func RankingLamp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("邀请新人列表服务 url：/api/plant/pullnew")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for k, v := range request {
		logs.Info(k, v)
	}

	//创建服务句柄
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := DATA.NewExampleService("go.micro.srv.data", server.Client())
	rsp, err := exampleClient.RankingLamp(context.TODO(), &DATA.RankingLampRequest{
		Skey: request["skey"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"list":   rsp.Newlist,
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//设置放回json格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//关卡服务
func CheckpointLamp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("邀请新人列表服务 url：/api/plant/pullnew")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for k, v := range request {
		logs.Info(k, v)
	}

	//创建服务句柄
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := DATA.NewExampleService("go.micro.srv.data", server.Client())
	rsp, err := exampleClient.CheckpointLamp(context.TODO(), &DATA.CheckpointLampRequest{
		Skey: request["skey"].(string),
		Code: request["code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//设置放回json格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//设置灯服务
func SetLamp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("邀请新人列表服务 url：/api/plant/pullnew")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建服务句柄
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := DATA.NewExampleService("go.micro.srv.data", server.Client())
	rsp, err := exampleClient.SetLamp(context.TODO(), &DATA.SetLampRequest{
		Skey: request["skey"].(string),
		Code: request["code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//设置放回json格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//获取灯列表服务
func GetLamp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("邀请新人列表服务 url：/api/plant/pullnew")

	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建服务句柄
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := DATA.NewExampleService("go.micro.srv.data", server.Client())
	rsp, err := exampleClient.GetLamp(context.TODO(), &DATA.GetLampRequest{
		Skey: request["skey"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"lamp":   rsp.Lamplist,
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//设置放回json格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}

//购买灯服务
func BuyLamp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("设置灯服务   url：/api/lamp/buylamp")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建服务句柄
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := DATA.NewExampleService("go.micro.srv.data", server.Client())
	rsp, err := exampleClient.BuyLamp(context.TODO(), &DATA.BuyLampRequest{
		Skey: request["skey"].(string),
		Code: request["code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"usermoney": rsp.UserMoney,
		"errno":     rsp.Errno,
		"errmsg":    rsp.Errmsg,
	}

	//设置放回json格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}
}
