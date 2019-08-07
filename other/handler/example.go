package handler

import (
	"context"
	example "plant/other/proto/example"
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/gomodule/redigo/redis"
	_ "github.com/garyburd/redigo/redis"
	"github.com/garyburd/redigo/redis"
	"plant/gateway/utils"
	"plant/gateway/models"
	"github.com/astaxie/beego/logs"
	"github.com/medivhzhan/weapp"
	"time"
	"github.com/astaxie/beego"
	"strconv"
)

//配置缓存参数
var redis_conf = map[string]string{
	"key": utils.G_server_name,
	//127.0.0.1:6379
	"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
	"dbNum":    utils.G_redis_dbnum,
	"password": utils.G_redis_password,
}

//将map进行转化成为json
var redis_conf_js, _ = json.Marshal(redis_conf)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PullNew(ctx context.Context, req *example.PullNewRequest, rsp *example.PullNewResponse) error {

	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/other/other.log"}`)

	logs.Info("PullNew  新人列表服务  /api/plant/PostData")

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
	var Invitations []models.Invitation
	o.QueryTable("invitation").Filter("user_id__hxopid", opidtemp).All(&Invitations)

	logs.Info("查询数据成功")

	for _, v := range Invitations {
		var templist example.New
		templist.Url = utils.Encrypt(v.Invitation_Avatar_url)
		templist.Name = utils.Encrypt(v.Invitation_Name)
		rsp.Newlist = append(rsp.Newlist, &templist)
	}

	rsp.Newlist = rsp.Newlist
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	return nil

}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Strength(ctx context.Context, req *example.StrengthRequest, rsp *example.StrengthResponse) error {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/other/other.log"}`)
	logs.Info("体力助力服务 url：/api/plant/strength")

	//判断是否本人点击
	if req.Skey == req.Shareskey {
		logs.Info("本人点击助力失败")
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//使用分享人skey从缓存中取出分享人的hxopid
	sharehxopid := bm.Get(req.Shareskey)
	if sharehxopid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//将取出的数据转换为string
	sharehxopidtemp, _ := redis.String(sharehxopid, nil)
	logs.Info("分享人poid为", sharehxopidtemp)

	//使用用户skey从缓存中取出hxopid
	hxopid := bm.Get(req.Skey)
	if hxopid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//将取出的数据转换为string
	hxopidtemp, _ := redis.String(hxopid, nil)
	logs.Info("用户poid为", hxopidtemp)

	o := orm.NewOrm()
	//查询分享人信息
	var shareuser models.User
	err = o.QueryTable("user").Filter("hxopid", sharehxopidtemp).One(&shareuser)
	if err != nil {
		logs.Info("查询分享人数据失败")
	}
	//查询助力人信息
	var user models.User
	err = o.QueryTable("user").Filter("hxopid", hxopidtemp).One(&user)
	if err != nil {
		logs.Info("查询用户数据失败")
	}

	//判断分享条件
	var tempa []models.Strength
	num, err := o.QueryTable("strength").Filter("user_id__hxopid", sharehxopidtemp).All(&tempa)
	if err != nil {
		logs.Info("判断分享条件时查询数据库失败")
	}
	if num >= 5 {
		logs.Info("助力人超过五人,点击助力失败")
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	for _, v := range tempa {
		if v.Invitation_OpenId == hxopidtemp {
			logs.Info("此人已经存在,点击助力失败")
			rsp.Errno = utils.RECODE_PARAMERR
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
	}

	//插入分享人助力信息
	Strengtha := models.Strength{
		Invitation_OpenId:     user.Hxopid,
		Invitation_Name:       user.Name,
		Invitation_Avatar_url: user.Avatar_url,
		Invitation_time:       req.Sharetime,
		GetStatus:             "true",
		User:                  &shareuser,
	}
	o.Insert(&Strengtha)
	logs.Info("分享人助力信息插入成功")
	//插入助力人助力信息
	Strengthb := models.Strength{
		Invitation_OpenId:     shareuser.Hxopid,
		Invitation_Name:       shareuser.Name,
		Invitation_Avatar_url: shareuser.Avatar_url,
		Invitation_time:       req.Sharetime,
		GetStatus:             "true",
		User:                  &user,
	}
	o.Insert(&Strengthb)
	logs.Info("助力人助力信息插入成功")
	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Strelest(ctx context.Context, req *example.StrelestRequest, rsp *example.StrelestResponse) error {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/other/other.log"}`)
	logs.Info("获取体力助力数据服务 url：/api/plant/strelest")

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("redis连接成功")

	//从缓存中取出opid
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

	if len(req.Url) > 10 {
		var TInvitations []models.Strength
		//按照opid查出所有助力列表
		o.QueryTable("strength").Filter("user_id__hxopid", opidtemp).All(&TInvitations)
		//遍历查询结果
		for _, v := range TInvitations {
			var TTInvitations models.Strength
			//找到助力人url的数据
			if v.Invitation_Avatar_url == req.Url {
				//利用找到数据的id查询出数据
				o.QueryTable("strength").Filter("id", v.Id).One(&TTInvitations)
				if err == orm.ErrNoRows {
					//返回请求参数错误
					rsp.Errno = utils.RECODE_DBERR
					rsp.Errmsg = utils.RecodeText(rsp.Errno)
					return nil
				}
				//修改字段数据后更新数据库
				TTInvitations.GetStatus = "false"
				num, err := o.Update(&TTInvitations)
				if err != nil {
					logs.Info("更新领取状态失败", err)
				}
				logs.Info(num, "更新领取状态成功")
			}
		}
	}

	//原来版本会出现不能同步重复数据的get状态
	//	logs.Info(req.Url)
	//	var Invitations models.Strength
	//	err :=o.QueryTable("strength").Filter("invitation_avatar_url",req.Url).One(&Invitations)
	//	if err ==orm.ErrNoRows{
	//		//返回请求参数错误
	//		rsp.Errno = utils.RECODE_DBERR
	//		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//		return nil
	//	}
	//	logs.Info(Invitations)
	//	Invitations.GetStatus = "false"
	//	num, err := o.Update(&Invitations)
	//	if err != nil {
	//		logs.Info("更新领取状态失败", err)
	//	}
	//	logs.Info(num,"更新领取状态成功")
	//}

	var Invitations []models.Strength
	o.QueryTable("strength").Filter("user_id__hxopid", opidtemp).All(&Invitations)

	logs.Info("查询数据成功")

	for _, v := range Invitations {
		var templist example.New
		templist.Url = utils.Encrypt(v.Invitation_Avatar_url)
		templist.Name = utils.Encrypt(v.Invitation_Name)
		templist.GetStatus = utils.Encrypt(v.GetStatus)
		rsp.Newlist = append(rsp.Newlist, &templist)
	}

	rsp.Newlist = rsp.Newlist
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Diamond(ctx context.Context, req *example.DiamondRequest, rsp *example.DiamondResponse) error {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/other/other.log"}`)
	logs.Info("钻石卡包服务 url：/api/plant/nowtime")

	//判断是否本人点击
	if req.Skey == req.Shareskey {
		logs.Info("本人点击助力失败")
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//使用分享人skey从缓存中取出分享人的hxopid
	sharehxopid := bm.Get(req.Shareskey)
	if sharehxopid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//将取出的数据转换为string
	sharehxopidtemp, _ := redis.String(sharehxopid, nil)
	logs.Info("分享人poid为", sharehxopidtemp)

	//使用用户skey从缓存中取出hxopid
	hxopid := bm.Get(req.Skey)
	if hxopid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//将取出的数据转换为string
	hxopidtemp, _ := redis.String(hxopid, nil)
	logs.Info("用户poid为", hxopidtemp)

	o := orm.NewOrm()
	//查询分享人信息
	var shareuser models.User
	err = o.QueryTable("user").Filter("hxopid", sharehxopidtemp).One(&shareuser)
	if err != nil {
		logs.Info("查询分享人数据失败")
	}
	//查询助力人信息
	var user models.User
	err = o.QueryTable("user").Filter("hxopid", hxopidtemp).One(&user)
	if err != nil {
		logs.Info("查询用户数据失败")
	}

	//判断分享条件
	var tempa []models.Diamondcard
	num, err := o.QueryTable("diamondcard").Filter("user_id__hxopid", sharehxopidtemp).All(&tempa)
	if err != nil {
		logs.Info("判断分享条件时查询数据库失败")
	}
	if num >= 5 {
		logs.Info("助力人超过五人,点击助力失败")
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	for _, v := range tempa {
		if v.Invitation_OpenId == hxopidtemp {
			logs.Info("此人已经存在,点击助力失败")
			rsp.Errno = utils.RECODE_PARAMERR
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
	}

	//插入分享人助力信息
	Strengtha := models.Diamondcard{
		Invitation_OpenId:     user.Hxopid,
		Invitation_Name:       user.Name,
		Invitation_Avatar_url: user.Avatar_url,
		Invitation_time:       req.Sharetime,
		GetStatus:             "true",
		User:                  &shareuser,
	}
	o.Insert(&Strengtha)
	logs.Info("分享人助力信息插入成功")
	//插入助力人助力信息
	Strengthb := models.Diamondcard{
		Invitation_OpenId:     shareuser.Hxopid,
		Invitation_Name:       shareuser.Name,
		Invitation_Avatar_url: shareuser.Avatar_url,
		Invitation_time:       req.Sharetime,
		GetStatus:             "true",
		User:                  &user,
	}
	o.Insert(&Strengthb)
	logs.Info("助力人助力信息插入成功")

	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Diamlest(ctx context.Context, req *example.DiamlestRequest, rsp *example.DiamlestResponse) error {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/other/other.log"}`)
	logs.Info("获取钻石助力列表服务 url：/api/plant/diamlest")

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

	if len(req.Url) > 10 {
		var TInvitations []models.Diamondcard
		//按照opid查出所有助力列表
		o.QueryTable("diamondcard").Filter("user_id__hxopid", opidtemp).All(&TInvitations)
		//遍历查询结果
		for _, v := range TInvitations {
			var TTInvitations models.Diamondcard
			//找到助力人url的数据
			if v.Invitation_Avatar_url == req.Url {
				//利用找到数据的id查询出数据
				o.QueryTable("diamondcard").Filter("id", v.Id).One(&TTInvitations)
				if err == orm.ErrNoRows {
					//返回请求参数错误
					rsp.Errno = utils.RECODE_DBERR
					rsp.Errmsg = utils.RecodeText(rsp.Errno)
					return nil
				}
				//修改字段数据后更新数据库
				TTInvitations.GetStatus = "false"
				num, err := o.Update(&TTInvitations)
				if err != nil {
					logs.Info("更新领取状态失败", err)
				}
				logs.Info(num, "更新领取状态成功")
			}
		}
	}

	var Invitations []models.Diamondcard
	o.QueryTable("diamondcard").Filter("user_id__hxopid", opidtemp).All(&Invitations)

	logs.Info("查询数据成功")

	for _, v := range Invitations {
		var templist example.New
		templist.Url = utils.Encrypt(v.Invitation_Avatar_url)
		templist.Name = utils.Encrypt(v.Invitation_Name)
		templist.GetStatus = utils.Encrypt(v.GetStatus)
		rsp.Newlist = append(rsp.Newlist, &templist)
	}

	rsp.Newlist = rsp.Newlist
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Lamp(ctx context.Context, req *example.LampRequest, rsp *example.LampResponse) error {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/other/other.log"}`)
	logs.Info("灯助力服务 url：/api/plant/lamp")

	//判断是否本人点击
	if req.Skey == req.Shareskey {
		logs.Info("本人点击助力失败")
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//使用分享人skey从缓存中取出分享人的hxopid
	sharehxopid := bm.Get(req.Shareskey)
	if sharehxopid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("获取分享人hxopid成功")

	//将取出的数据转换为string
	sharehxopidtemp, _ := redis.String(sharehxopid, nil)

	logs.Info("分享人poid为", sharehxopidtemp)

	//使用用户skey从缓存中取出hxopid
	hxopid := bm.Get(req.Skey)
	if hxopid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("获取用户hxopid成功")

	//将取出的数据转换为string
	hxopidtemp, _ := redis.String(hxopid, nil)

	logs.Info("用户poid为", hxopidtemp)

	o := orm.NewOrm()
	//判断是否新用户
	var shareuser models.User
	err = o.QueryTable("user").Filter("hxopid", sharehxopidtemp).One(&shareuser)
	if err != nil {
		logs.Info("查询分享人数据失败")
	}
	var user models.User
	err = o.QueryTable("user").Filter("hxopid", hxopidtemp).One(&user)
	if err != nil {
		logs.Info("查询用户数据失败")
	}

	//判断分享条件
	var tempa []models.Lampshare
	num, err := o.QueryTable("lampshare").Filter("user_id__hxopid", sharehxopidtemp).All(&tempa)
	if err != nil {
		logs.Info("判断分享条件时查询数据库失败")
	}
	if num >= 5 {
		logs.Info("助力人超过五人,点击助力失败")
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	for _, v := range tempa {
		if v.Invitation_OpenId == hxopidtemp {
			logs.Info("此人已经存在,点击助力失败")
			rsp.Errno = utils.RECODE_PARAMERR
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
	}

	Lamp := models.Lampshare{
		Invitation_OpenId:     user.Hxopid,
		Invitation_Name:       user.Name,
		Invitation_Avatar_url: user.Avatar_url,
		Invitation_time:       req.Sharetime,
		GetStatus:             "true",
		User:                  &shareuser,
	}
	o.Insert(&Lamp)
	logs.Info("助力人信息插入成功")

	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Lamplest(ctx context.Context, req *example.LamplestRequest, rsp *example.LamplestResponse) error {

	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/other/other.log"}`)
	logs.Info("获取灯助力数据服务 url：/api/plant/lamplest")

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

	//创建切片准备接收查询数据
	var Invitations []models.Lampshare
	switch req.Type {
	case "0":
		//将所有数据返回
		//查询数据库并将结果保存在切片当中
		o.QueryTable("lampshare").Filter("user_id__hxopid", opidtemp).All(&Invitations)
		logs.Info("查询数据成功")
		//遍历查询结果赋值给要返回要求结构的切片
		for _, v := range Invitations {
			var templist example.New
			templist.Url = utils.Encrypt(v.Invitation_Avatar_url)
			templist.Name = utils.Encrypt(v.Invitation_Name)
			templist.GetStatus = utils.Encrypt(v.GetStatus)
			rsp.Newlist = append(rsp.Newlist, &templist)
			logs.Info(v.GetStatus)
		}
	case "1":
		//使用poid查出所有助力信息
		o.QueryTable("lampshare").Filter("user_id__hxopid", opidtemp).All(&Invitations)
		logs.Info("修改领取数据成功")
		//遍历所有助力信息
		for _, v := range Invitations {
			//定义一个数据模型
			var temp models.Lampshare
			//使用遍历出的name查出对应的数据改变其getstatus字段值
			o.QueryTable("lampshare").Filter("invitation_name", v.Invitation_Name).One(&temp)
			temp.GetStatus = "false"
			//更新到表中
			o.Update(&temp)
		}
		//使用opid查出所有助力信息返回给前端
		o.QueryTable("lampshare").Filter("user_id__hxopid", opidtemp).All(&Invitations)
		for _, v := range Invitations {
			var templist example.New
			templist.Url = utils.Encrypt(v.Invitation_Avatar_url)
			templist.Name = utils.Encrypt(v.Invitation_Name)
			templist.GetStatus = utils.Encrypt(v.GetStatus)
			rsp.Newlist = append(rsp.Newlist, &templist)
		}
		//type为2删除邀请数据
	case "2":
		//使用opid查出并删除当前opid下的所有助力信息
		o.QueryTable("lampshare").Filter("user_id__hxopid", opidtemp).Delete()
		logs.Info("删除邀请数据成功")
	}
	rsp.Newlist = rsp.Newlist
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) TreasureMap(ctx context.Context, req *example.TreasureMapRequest, rsp *example.TreasureMapResponse) error {

	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/other/other.log"}`)
	logs.Info("宝图服务 url：/api/plant/treasuremap")
	logs.Info("请求code为", req.Code)
	logs.Info("请求skey为", req.Skey)
	logs.Info("请求偏移量", req.Iv)
	logs.Info("请求加密数据", req.Data)
	logs.Info("领取奖励url", req.Url)
	logs.Info("分享人skey", req.Shareskey)

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

	//取出分享人opid
	sharesopid := bm.Get(req.Shareskey)
	if opid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("获取opid成功")
	//将取出的数据转换为string
	sharesopidtemp, _ := redis.String(sharesopid, nil)

	o := orm.NewOrm()

	//分享人数据
	var sharesuser models.User
	o.QueryTable("user").Filter("hxopid", sharesopidtemp).One(&sharesuser)

	//助力人数据
	var user models.User
	o.QueryTable("user").Filter("hxopid", opidtemp).One(&user)

	switch req.Code {
	//0  助力人请求
	case "17e565a7e2346f19852cc308cb6ed2be":

		logs.Info("助力人请求", user.Hxopid)

		if user.Hxopid == sharesuser.Hxopid {
			rsp.Errno = utils.RECODE_PARAMERR
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
		mapkey := "mapdata" + sharesuser.Hxopid
		beego.Info(mapkey)
		openGidkey := "openGid" + sharesuser.Hxopid
		beego.Info(openGidkey)
		userkey := "user" + sharesuser.Hxopid
		beego.Info(userkey)

		//获取群号
		data := req.Data
		iv := req.Iv
		ssk := user.Session_key

		openGid, err := weapp.DecryptShareInfo(ssk, data, iv)
		if err != nil {
			logs.Info("请求微信失败", err)
		}
		logs.Info(openGid)

		//从redis中取出群数据

		opidtemp_s := bm.Get(openGidkey)
		if opidtemp_s == nil {
			logs.Info("redis取出数据为空")
		}

		//解码
		tempGid := make([]string, 0)
		err = json.Unmarshal(opidtemp_s.([]byte), &tempGid)
		if err != nil {
			logs.Info("解码失败")
		}

		//判断群号是否存在
		index := -1
		for i := 0; i < len(tempGid); i++ {
			if openGid == tempGid[i] {
				index = i
			}
		}

		//如果不存在

		if index == -1 {
			tempGid = append(tempGid, openGid)
			index := len(tempGid) - 1
			openGidslicejson, _ := json.Marshal(&tempGid)

			err = bm.Put(openGidkey, openGidslicejson, time.Second*86400)
			if err != nil {
				logs.Info("存入缓存失败", err)
			}
			logs.Info("群id存入缓存成功")
			//从redis中取出map数据
			//1获取缓存数据
			areas_info_value := bm.Get(mapkey)
			if areas_info_value == nil {
				logs.Info("redis取出数据为空")
			}
			//用来存放解码的json
			tempmap := []example.TreasureMap{}
			//解码
			err = json.Unmarshal(areas_info_value.([]byte), &tempmap)
			if err != nil {

				logs.Info("解码失败")
			}
			logs.Info("获取map数据成功")
			logs.Info(index)

			tempmap[index].State = "017d4757e0d0d9e91edec7ae166d0a49"
			tempmap[index].Url = user.Avatar_url

			Maplistjson, _ := json.Marshal(&tempmap)

			err = bm.Put(mapkey, Maplistjson, time.Second*86400)
			if err != nil {
				logs.Info("存入缓存失败", err)
			}
			logs.Info("宝图数据存入缓存成功")

			usertemp_s := bm.Get(userkey)
			if usertemp_s == nil {
				logs.Info("redis取出数据为空")
			}

			//解码
			usertemp := make([]string, 0)
			err = json.Unmarshal(usertemp_s.([]byte), &usertemp)
			if err != nil {
				logs.Info("解码失败")
			}

			usertemp = append(usertemp, user.Hxopid)

			logs.Info(usertemp)
			usertempjson, _ := json.Marshal(&usertemp)

			err = bm.Put(userkey, usertempjson, time.Second*86400)
			if err != nil {
				logs.Info("存入缓存失败", err)
			}
			logs.Info("宝图数据存入缓存成功")

			rsp.Goldtime = "65e32c010d2544cec4a6b36a1592961e"
			rsp.Name = utils.Encrypt(sharesuser.Name)
			rsp.Errno = utils.RECODE_OK
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
			//如果存在
		} else {

			//从redis中取出群数据

			usertemp_s := bm.Get(userkey)
			if usertemp_s == nil {
				logs.Info("redis取出数据为空")
			}

			//解码
			usertemp := make([]string, 0)
			err = json.Unmarshal(usertemp_s.([]byte), &usertemp)
			if err != nil {
				logs.Info("解码失败")
			}

			userindex := -1

			//判断是否领取过金币
			for i := 0; i < len(usertemp); i++ {
				if user.Hxopid == usertemp[i] {
					userindex = i
				}
			}

			if userindex == -1 {
				usertemp = append(usertemp, user.Hxopid)

				logs.Info(usertemp)
				usertempjson, _ := json.Marshal(&usertemp)

				err = bm.Put(userkey, usertempjson, time.Second*86400)
				if err != nil {
					logs.Info("存入缓存失败", err)
				}
				logs.Info("宝图数据存入缓存成功")
				logs.Info("放入数据成功", usertempjson)
				rsp.Goldtime = utils.Encrypt("600")
				rsp.Name = utils.Encrypt("-1")
				rsp.Errno = utils.RECODE_OK
				rsp.Errmsg = utils.RecodeText(rsp.Errno)
				return nil
			}
			rsp.Errno = utils.RECODE_PARAMERR
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
		//1 分享人请求
	case "017d4757e0d0d9e91edec7ae166d0a49":
		mapkey := "mapdata" + user.Hxopid
		beego.Info(mapkey)
		openGidkey := "openGid" + user.Hxopid
		beego.Info(openGidkey)
		userkey := "user" + user.Hxopid
		beego.Info(userkey)

		//1获取缓存数据
		areas_info_value := bm.Get(mapkey)
		if areas_info_value == nil {
			//####################################################################
			//新用户
			beego.Info("redis取出数据为空")
			var confmaps []models.TreasureMap
			o.QueryTable("treasure_map").All(&confmaps)

			var confsboxs []models.BoxSmall
			o.QueryTable("box_small").All(&confsboxs)

			var confbboxs []models.BoxBig
			o.QueryTable("box_big").All(&confbboxs)

			for _, v := range confmaps {
				var maplist example.TreasureMap
				if v.Type == "0" {
					for _, v1 := range confsboxs {
						var box example.Box
						box.Type = utils.Encrypt(v1.Code)
						box.Number = utils.Encrypt(v1.Number)
						box.Name = utils.Encrypt(v1.Name)
						maplist.Reward = append(maplist.Reward, &box)
					}
					maplist.State = utils.Encrypt(v.State)
					maplist.Type = utils.Encrypt(v.Type)
					maplist.Url = v.Url
					rsp.Maplist = append(rsp.Maplist, &maplist)
				} else if v.Type == "1" {
					for _, v1 := range confbboxs {
						var box example.Box
						box.Type = utils.Encrypt(v1.Code)
						box.Number = utils.Encrypt(v1.Number)
						box.Name = utils.Encrypt(v1.Name)
						maplist.Reward = append(maplist.Reward, &box)
					}
					maplist.State = utils.Encrypt(v.State)
					maplist.Type = utils.Encrypt(v.Type)
					maplist.Url = v.Url
					rsp.Maplist = append(rsp.Maplist, &maplist)
				}
			}

			Maplistjson, _ := json.Marshal(&rsp.Maplist)

			err = bm.Put(mapkey, Maplistjson, time.Second*86400)
			if err != nil {
				logs.Info("存入缓存失败", err)
			}
			logs.Info("宝图数据存入缓存成功")

			var openGidslice = make([]string, 0)
			openGidslicejson, _ := json.Marshal(&openGidslice)

			err = bm.Put(openGidkey, openGidslicejson, time.Second*86400)
			if err != nil {
				logs.Info("存入缓存失败", err)
			}
			logs.Info("群id存入缓存成功")

			var userkeyslice = make([]string, 0)
			userkeyslicejson, _ := json.Marshal(&userkeyslice)

			err = bm.Put(userkey, userkeyslicejson, time.Second*86400)
			if err != nil {
				logs.Info("存入缓存失败", err)
			}
			logs.Info("助力用户存入缓存成功")

			rsp.Errno = utils.RECODE_OK
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}

		//用来存放解码的json
		tempmap := []example.TreasureMap{}
		//解码
		err = json.Unmarshal(areas_info_value.([]byte), &tempmap)
		if err != nil {

			beego.Info("解码失败")
		}

		aaa := -1
		//有url的情况
		if len(req.Url) > 10 {
			for i := 0; i < len(tempmap); i++ {
				if req.Url == tempmap[i].Url {
					aaa = i
					break
				}
			}

			if aaa != -1 {
				tempmap[aaa].State = "b10c953cfabad6db861eb9c8a90ec058"
				Maplistjson, _ := json.Marshal(&tempmap)

				err = bm.Put(mapkey, Maplistjson, time.Second*86400)
				if err != nil {
					logs.Info("存入缓存失败", err)
				}
				logs.Info("宝图数据存入缓存成功")
			}
		}

		for _, v := range tempmap {
			var temp example.TreasureMap
			temp.Reward = v.Reward
			temp.Url = v.Url
			temp.State = v.State
			temp.Type = v.Type
			rsp.Maplist = append(rsp.Maplist, &temp)
		}
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
		//2删除
	case "b10c953cfabad6db861eb9c8a90ec058":
		mapkey := "mapdata" + user.Hxopid
		logs.Info(mapkey)
		openGidkey := "openGid" + user.Hxopid
		logs.Info(openGidkey)
		userkey := "user" + user.Hxopid
		logs.Info(userkey)

		bm.Delete(mapkey)
		bm.Delete(openGidkey)
		bm.Delete(userkey)

		logs.Info("清空数据成功")
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil

	default:
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Achievement(ctx context.Context, req *example.AchievementRequest, rsp *example.AchievementResponse) error {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/other/other.log"}`)
	logs.Info("灯助力服务 url：/api/plant/lamp")

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_js))
	if err != nil {
		logs.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//使用用户skey从缓存中取出hxopid
	hxopid := bm.Get(req.Skey)
	if hxopid == nil {
		logs.Info("获取opid失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logs.Info("获取用户hxopid成功")

	//将取出的数据转换为string
	hxopidtemp, _ := redis.String(hxopid, nil)

	logs.Info("用户poid为", hxopidtemp)

	o := orm.NewOrm()
	//判断是否新用户
	//var shareuser models.User
	//err = o.QueryTable("user").Filter("hxopid", sharehxopidtemp).One(&shareuser)
	//if err != nil {
	//	logs.Info("查询分享人数据失败")
	//}
	//var user models.User
	//err = o.QueryTable("user").Filter("hxopid", hxopidtemp).One(&user)
	//if err != nil {
	//	logs.Info("查询用户数据失败")
	//}

	temp:=utils.Decrypt(req.Code)
	temp1 , _:=strconv.Atoi(temp)
	Code:=strconv.Itoa(temp1+1)

	//判断分享条件
	var temps models.Achievement
	err = o.QueryTable("achievement").Filter("title_code",Code ).One(&temps)
	if err != nil {
		logs.Info("判断分享条件时查询数据库失败")
	}

	var tempreward example.Reward

	tempreward.Code = utils.Encrypt(temps.Reward4Code)
	tempreward.Number = utils.Encrypt(temps.Reward4Number)
	rsp.Reward = append(rsp.Reward, &tempreward)

	var tempreward1 example.Reward
	tempreward1.Code = utils.Encrypt(temps.Reward1Code)
	tempreward1.Number = utils.Encrypt(temps.Reward1Number)
	rsp.Reward = append(rsp.Reward, &tempreward1)

	var tempreward2 example.Reward
	tempreward2.Code = utils.Encrypt(temps.Reward2Code)
	tempreward2.Number = utils.Encrypt(temps.Reward2Number)
	rsp.Reward = append(rsp.Reward, &tempreward2)

	var tempreward3 example.Reward
	tempreward3.Code = utils.Encrypt(temps.Reward3Code)
	tempreward3.Number = utils.Encrypt(temps.Reward3Number)
	rsp.Reward = append(rsp.Reward, &tempreward3)
	//nnumber, _ := strconv.Atoi(Code)

	rsp.NextNumber = utils.Encrypt(temps.TitleCode)
	rsp.Name = utils.Encrypt(temps.TitlName)
	rsp.ConditionA = utils.Encrypt(temps.CumulativeMoney)
	rsp.ConditionB = utils.Encrypt(temps.HarvestTimes)
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	return nil
}
