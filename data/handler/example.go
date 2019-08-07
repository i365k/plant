package handler

import (
	example "plant/data/proto/example"
	"plant/gateway/utils"
	"plant/gateway/models_lamp"
	"context"
	"encoding/json"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/gomodule/redigo/redis"
	"time"
	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego/logs"
	"github.com/medivhzhan/weapp"
	"github.com/goEncrypt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
)

type Example struct{}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//同步服务
func (e *Example) Data(ctx context.Context, req *example.Request, rsp *example.Response) error {

	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/data/data.log"}`)
	logs.Info("同步数据服务 url : /api/plant/data")
	//将接收到的数据解码到map中
	var request = make(map[string]interface{})
	json.Unmarshal(req.Data, &request)
	//从map中取出skey
	redis_key := request["skey"]
	if redis_key == nil {
		logs.Info("请求参数为空")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//配置缓存参数
	redis_conf := map[string]string{
		"key": utils.G_server_name,
		//127.0.0.1:6379
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": utils.G_redis_password,
	}
	logs.Info(redis_conf)
	//将map进行转化成为json
	redis_conf_js, _ := json.Marshal(redis_conf)
	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("redis连接成功")
	//将同步来的data存入redis
	//从缓存中取出哈希过的opid
	opid := bm.Get(redis_key.(string))
	if opid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("获取opid成功")

	//将取出的数据转换为string
	opidtemp, _ := redis.String(opid, nil)

	logs.Info(opidtemp)
	//使用哈希opid为key把用户data放入缓存
	err = bm.Put(opidtemp, req.Data, time.Second*36000000)
	if err != nil {
		logs.Info("存入redis失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("存入缓存成功")

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	return nil
}

//用户登陆
func (e *Example) LoginLamp(ctx context.Context, req *example.LoginLampRequest, rsp *example.LoginLampResponse) error {

	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("LoginLamp 登陆服务 /api/lamp/login")
	//将接收到的数据解码
	var request = make(map[string]interface{})
	json.Unmarshal(req.Data, &request)
	//取出userinfu数据
	tempuserinfo := request["userInfo"].(map[string]interface{})
	nickName := tempuserinfo["nickName"]
	logs.Info(nickName)
	gender := tempuserinfo["gender"]
	logs.Info(gender)
	language := tempuserinfo["language"]
	logs.Info(language)
	city := tempuserinfo["city"]
	logs.Info(city)
	province := tempuserinfo["province"]
	logs.Info(province)
	country := tempuserinfo["country"]
	logs.Info(country)
	avatarUrl := tempuserinfo["avatarUrl"]
	logs.Info(avatarUrl)

	//请求微信接口获取session_key和openid
	appID := "wxb1314e33cdda4c6c"
	secret := "21fd3db3a6a2ed3ee4b90faa99bb3407"
	code := request["code"].(string)
	//调用微信接口获取数据
	res, err := weapp.Login(appID, secret, code)
	if err != nil {
		logs.Info(err)
		logs.Info("code验证失败")
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	ssk := res.SessionKey
	//调用sha256方法将session_key哈希得到skey
	skey := goEncrypt.GetStringHash256(ssk)
	//将从微信获取到的opid哈希处理
	hxopid := GetMd5String(res.OpenID)

	//配置缓存参数
	redis_conf := map[string]string{
		"key": utils.G_server_name_lamp,
		//127.0.0.1:6379
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum_lamp,
		"password": utils.G_redis_password,
	}
	logs.Info(redis_conf)
	//将map进行转化成为json
	redis_conf_js, _ := json.Marshal(redis_conf)
	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	o := orm.NewOrm()
	//判断是否新用户
	var usertemp models_lamp.User
	err = o.QueryTable("user").Filter("hxopid", hxopid).One(&usertemp)
	if err == orm.ErrNoRows {
		//判断为新用户
		logs.Info("新用户登陆")
		//插入用户数据
		usertemp.OpenId = res.OpenID
		usertemp.Hxopid = hxopid
		usertemp.Name = nickName.(string)
		//usertemp.Gender = gender.(string)
		usertemp.Session_key = res.SessionKey
		usertemp.Skey = skey
		usertemp.Avatar_url = avatarUrl.(string)
		usertemp.Landing_time = strconv.Itoa(int(time.Now().Unix()))
		usertemp.MaxLevel = "1"
		usertemp.CurrentLevel = "1"
		usertemp.CurrentLampCode = "1001"
		usertemp.UserMoney = "0"
		num, err := o.Insert(&usertemp)
		if err != nil {
			logs.Info("插入新用户失败", err)
		}
		logs.Info("插入新用户成功", num)

		var lampconfigs []models_lamp.Lampconfig
		o.QueryTable("lampconfig").All(&lampconfigs)

		var user models_lamp.User
		o.QueryTable("user").Filter("hxopid", hxopid).One(&user)

		var templamp models_lamp.Lamp
		for _, v := range lampconfigs {
			templamp.Id = 0
			templamp.LampId = v.LampId
			templamp.LampPrice = v.LampPrice
			templamp.LampUrl = v.LampUrl
			templamp.IsHave = v.IsHave
			templamp.User = &user
			num, err := o.Insert(&templamp)
			if err != nil {
				logs.Info("插入数据失败", err)
			}
			logs.Info("插入灯初始数据成功", num)
		}

		//将skey对应的openid写入缓存
		bm.Put(skey, hxopid, time.Second*36000)
		if err != nil {
			logs.Info("存入redis失败", err)
			rsp.Errno = utils.RECODE_DBERR
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
		logs.Info("opid插入redis成功")

		//判断是否为受邀请用户
		shareskey := request["skey"].(string)
		if len(shareskey) > 15 {
			logs.Info(shareskey)
			//从缓存中取出哈希过的opid
			opid := bm.Get(shareskey)
			if opid == nil {
				logs.Info("超过分享时间,判定为普通用户登陆")
				//rsp.Errno = utils.RECODE_PARAMERR
				//rsp.Errmsg = utils.RecodeText(rsp.Errno)
				//return nil
			}
			logs.Info("获取分享人信息成功,判定为邀请用户登陆")
			//将取出的数据转换为string
			opidtemp, _ := redis.String(opid, nil)

			var user models_lamp.User
			o.QueryTable("user").Filter("hxopid", opidtemp).One(&user)
			invitations := models_lamp.Invitation{
				Invitation_OpenId:     hxopid,
				Invitation_Name:       nickName.(string),
				Invitation_Avatar_url: avatarUrl.(string),
				GetStatus:             "true",
				User:                  &user,
			}
			o.Insert(&invitations)
			logs.Info("插入邀请新人表成功")
		}
		//返回数据
		rsp.Skey = skey
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//判断为老用户
	logs.Info("老用户登陆", usertemp.Name)
	//更新用户数据
	usertemp.Skey = skey
	usertemp.Name = nickName.(string)
	//usertemp.Gender = gender.(string)
	usertemp.Avatar_url = avatarUrl.(string)
	usertemp.Landing_time = strconv.Itoa(int(time.Now().Unix()))
	num, err := o.Update(&usertemp)
	if err != nil {
		logs.Info("老用户更新数据失败", err)
	}
	logs.Info(num)

	//将skey对应的openid写入缓存
	bm.Put(skey, hxopid, time.Second*36000)
	if err != nil {
		logs.Info("存入redis失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("opid插入redis成功")

	//返回数据
	rsp.Skey = skey
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	return nil

	/*
		//使用session_key解密用户信息
		// 解密手机号码
		// @ssk 通过 Login 向微信服务端请求得到的 session_key
		// @data 小程序通过 api 得到的加密数据(encryptedData)
		data := req.EncryptedData
		// @iv 小程序通过 api 得到的初始向量(iv)
		iv := req.Iv
		phone , _ := weapp.DecryptPhoneNumber(ssk, data, iv)

		fmt.Printf("手机数据: %#v", phone)


		// 解密转发信息的加密数据
		// @ssk 通过 Login 向微信服务端请求得到的 session_key
		// @data 小程序通过 api 得到的加密数据(encryptedData)
		// @iv 小程序通过 api 得到的初始向量(iv)

		// @gid 小程序唯一群号
		openGid , _ := weapp.DecryptShareInfo(ssk, data, iv )
		fmt.Println(openGid)


		// 解密用户信息
		// @rawData 不包括敏感信息的原始数据字符串, 用于计算签名。
		rawData:=req.Rawdata
		// @encryptedData 包括敏感数据在内的完整用户信息的加密数据
		// @signature 使用 sha1( rawData + session_key ) 得到字符串, 用于校验用户信息
		signature:=req.Signature
		// @iv 加密算法的初始向量
		// @ssk 微信 session_key
		ui, _ := weapp.DecryptUserInfo(rawData, data, signature, iv, ssk )

		fmt.Printf("用户数据: %#v", ui)

	*/
}

//获取邀请新人列表
func (e *Example) PullNewLamp(ctx context.Context, req *example.PullNewLampRequest, rsp *example.PullNewLampResponse) error {

	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("PullNewLamp  新人列表服务  /api/lamp/pullnew")

	//配置缓存参数
	var redis_conf = map[string]string{
		"key": utils.G_server_name_lamp,
		//127.0.0.1:6379
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum_lamp,
		"password": utils.G_redis_password,
	}
	//将map进行转化成为json
	var redis_conf_js, _ = json.Marshal(redis_conf)

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("redis连接成功")
	//将同步来的data存入redis
	//从缓存中取出哈希过的opid
	opid := bm.Get(req.Skey)
	if opid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("获取opid成功")
	//将取出的数据转换为string
	opidtemp, _ := redis.String(opid, nil)

	o := orm.NewOrm()

	var user models_lamp.User
	o.QueryTable("user").Filter("hxopid",opidtemp).One(&user)
	if len(req.Url) > 10 {
		var TInvitations []models_lamp.Invitation
		o.QueryTable("invitation").Filter("user_id__hxopid", opidtemp).All(&TInvitations)
		for _, v := range TInvitations {
			var TTInvitations models_lamp.Invitation
			if v.Invitation_Avatar_url == req.Url {
				o.QueryTable("invitation").Filter("id", v.Id).One(&TTInvitations)
				if err == orm.ErrNoRows {
					rsp.Errno = utils.RECODE_DBERR
					rsp.Errmsg = utils.RecodeText(rsp.Errno)
					return nil
				}
				TTInvitations.GetStatus = "false"
				num, err := o.Update(&TTInvitations)
				if err != nil {
					logs.Info("更新领取状态失败", err)
				}
				usqq,_:=strconv.Atoi(user.UserMoney)
				usqq+=100
				user.UserMoney = strconv.Itoa(usqq)
				num, err = o.Update(&user)
				if err != nil {
					logs.Info("更新用户金币失败", err)
				}
				logs.Info(num, "更新用户金币成功")
			}
		}
	}

	var Invitations []models_lamp.Invitation
	o.QueryTable("invitation").Filter("user_id__hxopid", opidtemp).All(&Invitations)

	logs.Info("查询数据成功")

	for _, v := range Invitations {
		var templist example.New
		templist.Url = utils.Encrypt(v.Invitation_Avatar_url)
		templist.Name = utils.Encrypt(v.Invitation_Name)
		templist.GetStatus = utils.Encrypt(v.GetStatus)
		rsp.Newlist = append(rsp.Newlist, &templist)
	}
	rsp.UserMoney= user.UserMoney
	rsp.Newlist = rsp.Newlist
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	return nil

}

//初始化服务
func (e *Example) ConfigLamp(ctx context.Context, req *example.ConfigLampRequest, rsp *example.ConfigLampResponse) error {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("ConfigLamp  初始化服务 url：/api/lamp/config")

	//配置缓存参数
	redis_conf := map[string]string{
		"key": utils.G_server_name_lamp,
		//127.0.0.1:6379
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum_lamp,
		"password": utils.G_redis_password,
	}
	logs.Info(redis_conf)
	//将map进行转化成为json
	redis_conf_js, _ := json.Marshal(redis_conf)
	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//从缓存中获取opid
	hxopid := bm.Get(req.Skey)
	logs.Info(hxopid)
	if hxopid == nil {
		logs.Info("获取opid失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	var data example.Data
	//将缓存中取出的opid转成string
	temp, _ := redis.String(hxopid, nil)
	logs.Info(temp)

	o := orm.NewOrm()
	//准备用户数据
	var User models_lamp.User
	err = o.QueryTable("user").Filter("hxopid", temp).One(&User)
	if err != nil {
		logs.Info("查询用户数据失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//准备所有配置灯数据
	var lampconfigs []models_lamp.Lampconfig
	o.QueryTable("lampconfig").All(&lampconfigs)
	//准备用户所有的灯
	var lamps []models_lamp.Lamp
	_, err = o.QueryTable("lamp").Filter("user_id__hxopid", temp).All(&lamps)
	if err != nil {
		logs.Info("查询用户数据失败", err)
	}
	logs.Info("查询用户数据成功")
	//补全灯列表
	var templamp models_lamp.Lamp
	for i := 0; i < len(lampconfigs); i++ {
		index := -1
		aaa := lampconfigs[i]
		for _, v := range lamps {
			if v.LampId == aaa.LampId {
				index = i
				break
			}
		}
		if index == -1 {
			templamp.Id = 0
			templamp.LampId = aaa.LampId
			templamp.LampPrice = aaa.LampPrice
			templamp.LampUrl = aaa.LampUrl
			templamp.IsHave = aaa.IsHave
			templamp.User = &User
			num, err := o.Insert(&templamp)
			if err != nil {
				logs.Info("插入新灯失败", err)
			}
			logs.Info("插入新灯成功", num)
		}
	}
	//更新灯列表
	for _, v1 := range lampconfigs {
		for _, v := range lamps {
			if v.LampId == v1.LampId {
				var temp models_lamp.Lamp
				err = o.QueryTable("lamp").Filter("id", v.Id).One(&temp)
				if err == orm.ErrNoRows {
					logs.Info("没有这个灯的数据")
					rsp.Errno = utils.RECODE_PARAMERR
					rsp.Errmsg = utils.RecodeText(rsp.Errno)
					return nil
				}
				temp.LampPrice = v1.LampPrice
				temp.LampUrl = v1.LampUrl
				o.Update(&temp)
			}
		}
	}

	data.CurrentLevel = utils.Encrypt(User.CurrentLevel)
	data.MaxLevel = utils.Encrypt(User.MaxLevel)
	data.CurrentLampCode = utils.Encrypt(User.CurrentLampCode)
	data.UserMoney = utils.Encrypt(User.UserMoney)

	beego.Info("使用中的灯", data.CurrentLampCode)
	beego.Info("金币数", data.UserMoney)
	beego.Info("最高关卡数", data.MaxLevel)
	beego.Info("当前关卡数", data.CurrentLevel)
	//返回数据给web端
	rsp.Datalist = &data
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	return nil
}

//金币变化接口
func (e *Example) GoldLamp(ctx context.Context, req *example.GoldLampRequest, rsp *example.GoldLampResponse) error {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("GoldLamp  金币变化服务  /api/lamp/gold")

	//配置缓存参数
	var redis_conf = map[string]string{
		"key": utils.G_server_name_lamp,
		//127.0.0.1:6379
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum_lamp,
		"password": utils.G_redis_password,
	}
	//将map进行转化成为json
	var redis_conf_js, _ = json.Marshal(redis_conf)

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("redis连接成功")
	//从缓存中取出哈希过的opid
	opid := bm.Get(req.Skey)
	if opid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("获取opid成功")
	//将取出的数据转换为string
	opidtemp, _ := redis.String(opid, nil)

	o := orm.NewOrm()
	var user models_lamp.User
	err = o.QueryTable("user").Filter("hxopid", opidtemp).One(&user)
	if err != nil {
		logs.Info("查询用户数据失败")
	}
	logs.Info("查询数据成功")

	usergold, _ := strconv.Atoi(user.UserMoney)
	num, _ := strconv.Atoi(req.Number)

	//判断参数
	if num < 0 {
		logs.Info("金币数不合法")
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	switch req.Type {
	case "0":
		user.UserMoney = strconv.Itoa(usergold + num)
		o.Update(&user)
		logs.Info("增加金币成功", num)
	case "1":
		if usergold < num {
			logs.Info("金币数不合法")
			rsp.Errno = utils.RECODE_PARAMERR
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
		user.UserMoney = strconv.Itoa(usergold - num)
		o.Update(&user)
		logs.Info("减少金币成功", num)
	default:
		logs.Info("参数不合法")
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	return nil

}

//服务器排名服务
func (e *Example) RankingLamp(ctx context.Context, req *example.RankingLampRequest, rsp *example.RankingLampResponse) error {

	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)

	logs.Info("RankingLamp  服务器排名服务  /api/lamp/ranking")

	//配置缓存参数
	var redis_conf = map[string]string{
		"key": utils.G_server_name_lamp,
		//127.0.0.1:6379
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum_lamp,
		"password": utils.G_redis_password,
	}
	//将map进行转化成为json
	var redis_conf_js, _ = json.Marshal(redis_conf)

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("redis连接成功")

	//从缓存中取出哈希过的opid
	opid := bm.Get(req.Skey)
	if opid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("获取opid成功")
	//将取出的数据转换为string
	opidtemp, _ := redis.String(opid, nil)

	o := orm.NewOrm()
	var users []models_lamp.User
	o.QueryTable("user").Filter("hxopid", opidtemp).All(&users)

	logs.Info("查询数据成功")

	for _, v := range users {
		var templist example.New
		templist.Url = utils.Encrypt(v.Avatar_url)
		templist.Name = utils.Encrypt(v.Name)
		rsp.Newlist = append(rsp.Newlist, &templist)
	}

	rsp.Newlist = rsp.Newlist
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	return nil
}

//关卡服务
func (e *Example) CheckpointLamp(ctx context.Context, req *example.CheckpointLampRequest, rsp *example.CheckpointLampResponse) error {

	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)

	logs.Info("CheckpointLamp  关卡服务  /api/lamp/checkpoint")

	if len(req.Code) < 1 {
		logs.Info("关卡数不合法")
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//beego.Info("原始数据",req.Code)
	//
	//Key := []byte("1234567887654321")
	//yuan := []byte(req.Code)
	//jiadata:=goEncrypt.AesCBC_Encrypt(yuan,Key)
	//cipher := hex.EncodeToString(jiadata)
	//beego.Info("加密后数据",cipher)
	//
	////jiadata:=[]byte("b10c953cfabad6db861eb9c8a90ec058")
	//jiedata:=goEncrypt.AesCBC_Decrypt(jiadata,Key)
	//cipher = hex.EncodeToString(jiedata)
	//beego.Info("节密后数据",cipher)
	//
	//
	//

	num, _ := strconv.Atoi(req.Code)

	//判断参数
	if num < 0 {
		logs.Info("关卡数不合法")
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//配置缓存参数
	var redis_conf = map[string]string{
		"key": utils.G_server_name_lamp,
		//127.0.0.1:6379
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum_lamp,
		"password": utils.G_redis_password,
	}
	//将map进行转化成为json
	var redis_conf_js, _ = json.Marshal(redis_conf)

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("redis连接成功")
	//将同步来的data存入redis
	//从缓存中取出哈希过的opid
	opid := bm.Get(req.Skey)
	if opid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("获取opid成功")
	//将取出的数据转换为string
	opidtemp, _ := redis.String(opid, nil)

	temptime := time.Now().Unix()
	o := orm.NewOrm()
	var User models_lamp.User
	o.QueryTable("user").Filter("hxopid", opidtemp).One(&User)
	beego.Info("请求code为", req.Code)
	User.MaxLevel = req.Code
	User.CurrentLevel = req.Code
	User.LevelTime = strconv.Itoa(int(temptime))
	o.Update(&User)
	logs.Info("设置关卡成功")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	return nil

}

//设置灯服务
func (e *Example) SetLamp(ctx context.Context, req *example.SetLampRequest, rsp *example.SetLampResponse) error {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("SetLamp  设置灯服务  /api/lamp/setlamp")

	//配置缓存参数
	var redis_conf = map[string]string{
		"key": utils.G_server_name_lamp,
		//127.0.0.1:6379
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum_lamp,
		"password": utils.G_redis_password,
	}
	//将map进行转化成为json
	var redis_conf_js, _ = json.Marshal(redis_conf)

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("redis连接成功")
	//从缓存中取出哈希过的opid
	opid := bm.Get(req.Skey)
	if opid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("获取opid成功")
	//将取出的数据转换为string
	opidtemp, _ := redis.String(opid, nil)

	o := orm.NewOrm()
	var user models_lamp.User
	err = o.QueryTable("user").Filter("hxopid", opidtemp).One(&user)
	if err != nil {
		logs.Info("查询用户数据失败")
	}
	logs.Info("查询数据成功")

	user.CurrentLampCode = req.Code
	o.Update(&user)
	logs.Info("设置用户灯信息成功")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	return nil
}

//获取灯列表服务
func (e *Example) GetLamp(ctx context.Context, req *example.GetLampRequest, rsp *example.GetLampResponse) error {

	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)

	logs.Info("GetLamp  获取灯设置服务  /api/lamp/getlamp")

	//配置缓存参数
	var redis_conf = map[string]string{
		"key": utils.G_server_name_lamp,
		//127.0.0.1:6379
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum_lamp,
		"password": utils.G_redis_password,
	}
	//将map进行转化成为json
	var redis_conf_js, _ = json.Marshal(redis_conf)

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("redis连接成功")

	//从缓存中取出哈希过的opid
	opid := bm.Get(req.Skey)
	if opid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("获取opid成功")
	//将取出的数据转换为string
	opidtemp, _ := redis.String(opid, nil)

	o := orm.NewOrm()
	var users []models_lamp.Lamp
	o.QueryTable("lamp").Filter("user_id__hxopid", opidtemp).OrderBy("-is_have", "id").All(&users)

	logs.Info("查询数据成功")

	for _, v := range users {
		var templist example.Lamp
		templist.LampId = utils.Encrypt(v.LampId)
		templist.LampPrice = utils.Encrypt(v.LampPrice)
		templist.LampUrl = utils.Encrypt(v.LampUrl)
		templist.IsHave = utils.Encrypt(v.IsHave)
		rsp.Lamplist = append(rsp.Lamplist, &templist)
	}

	//for _, v := range users {
	//	var templist example.Lamp
	//	templist.LampId = v.LampId
	//	templist.LampPrice = v.LampPrice
	//	templist.LampUrl = v.LampUrl
	//	templist.IsHave = v.IsHave
	//	rsp.Lamplist = append(rsp.Lamplist, &templist)
	//}

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	return nil

}

//购买灯服务
func (e *Example) BuyLamp(ctx context.Context, req *example.BuyLampRequest, rsp *example.BuyLampResponse) error {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/lamp/lamp.log"}`)
	logs.Info("SetLamp  设置灯服务  /api/lamp/setlamp")

	//配置缓存参数
	var redis_conf = map[string]string{
		"key": utils.G_server_name_lamp,
		//127.0.0.1:6379
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum_lamp,
		"password": utils.G_redis_password,
	}
	//将map进行转化成为json
	var redis_conf_js, _ = json.Marshal(redis_conf)

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("redis连接成功")
	//从缓存中取出哈希过的opid
	opid := bm.Get(req.Skey)
	if opid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("获取opid成功")
	//将取出的数据转换为string
	opidtemp, _ := redis.String(opid, nil)

	o := orm.NewOrm()
	//扣除用户金币,准备用户数据
	var user models_lamp.User
	o.QueryTable("user").Filter("hxopid", opidtemp).One(&user)
	//准备灯数据
	var lamps []models_lamp.Lamp
	//找出要要修改的灯
	_, err = o.QueryTable("lamp").Filter("user_id__hxopid", opidtemp).All(&lamps)
	if err != nil {
		logs.Info("查询用户数据失败")
	}
	logs.Info("查询用户数据成功")

	for _, v := range lamps {
		if v.LampId == req.Code {
			var temp models_lamp.Lamp
			err = o.QueryTable("lamp").Filter("id", v.Id).One(&temp)
			if err != nil {
				logs.Info("查询数据失败")
				rsp.Errno = utils.RECODE_PARAMERR
				rsp.Errmsg = utils.RecodeText(rsp.Errno)
				return nil
			}
			userqq, _ := strconv.Atoi(user.UserMoney)
			lampqq, _ := strconv.Atoi(temp.LampPrice)
			if userqq < lampqq {
				logs.Info("金币不足购买失败")
				rsp.UserMoney = user.UserMoney
				rsp.Errno = utils.RECODE_BUYERR
				rsp.Errmsg = utils.RecodeText(rsp.Errno)
				return nil
			} else {
				buemoney := userqq - lampqq
				user.UserMoney = strconv.Itoa(buemoney)
				temp.IsHave = "true"
				o.Update(&user)
				o.Update(&temp)
			}
		}
	}

	logs.Info("购买灯成功")
	rsp.UserMoney = user.UserMoney
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	return nil
}
