package handler

import (
	example "plant/login/proto/example"
	"context"
	"plant/gateway/utils"
	"plant/gateway/models"

	_ "github.com/goEncrypt"
	"github.com/goEncrypt"
	"github.com/astaxie/beego/orm"
	"github.com/medivhzhan/weapp"
	"encoding/json"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/gomodule/redigo/redis"
	_ "github.com/garyburd/redigo/redis"
	"github.com/garyburd/redigo/redis"
	"time"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"github.com/astaxie/beego/logs"
)

type Example struct{}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Login(ctx context.Context, req *example.LoginRequest, rsp *example.LoginResponse) error {

	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/login/login.log"}`)
	logs.Info("Login 登陆服务 /api/plant/login")
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
	appID := "wx625e64a74fff566c"
	secret := "39bbb7f1b96bbaf514b37b5bfe6d7869"

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

	//查询数据库openid判断是否新用户
	//老用户更形数据库信息
	//新用户将全部信息openid,session_key,skey,用户信息入库
	o := orm.NewOrm()
	//判断是否新用户
	var usertemp models.User
	err = o.QueryTable("user").Filter("hxopid", hxopid).One(&usertemp)
	if err == orm.ErrNoRows {
		//判断为新用户
		logs.Info("新用户登陆")
		//插入用户数据
		usertemp.OpenId = res.OpenID
		usertemp.Hxopid = hxopid
		usertemp.Name = nickName.(string)
		usertemp.Gender = strconv.Itoa(int(gender.(float64)))
		usertemp.Session_key = res.SessionKey
		usertemp.Skey = skey
		usertemp.Avatar_url = avatarUrl.(string)
		usertemp.Landing_time = strconv.Itoa(int(time.Now().Unix()))
		num, err := o.Insert(&usertemp)
		if err != nil {
			logs.Info("插入新用户失败", err)
		}
		logs.Info("插入新用户成功", num)

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

			var user models.User
			o.QueryTable("user").Filter("hxopid", opidtemp).One(&user)
			invitations := models.Invitation{
				Invitation_OpenId:     hxopid,
				Invitation_Name:       nickName.(string),
				Invitation_Avatar_url: avatarUrl.(string),
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
	usertemp.Gender = strconv.Itoa(int(gender.(float64)))
	usertemp.Avatar_url = avatarUrl.(string)
	usertemp.Session_key = res.SessionKey
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

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) Config(ctx context.Context, req *example.ConfigRequest, rsp *example.ConfigResponse) error {
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/login/login.log"}`)
	logs.Info("初始化服务 url：/api/plant/config")
	//组织config数据
	o := orm.NewOrm()
	//定义一个切片用于接收查询到的植物数据(格式使用数据库定义格式）

	var crops []models.Crops
	var scenedata []models.SceneData
	var MassifData []models.MassifData
	var AllProceeds []models.AllProceeds
	var GoldEnhancementData []models.GoldEnhancementData
	var MoneyEnhancementData []models.MoneyEnhancementData
	var DiamondEnhancementData []models.DiamondEnhancementData
	var PlantDebris []models.PlantDebris
	var NewPeople []models.NewPeople
	var SignIn []models.SignIn
	var FriendHelp []models.FriendHelp
	var PotFragments []models.PotFragments
	var Prop []models.Prop
	var RandomProp []models.RandomProp
	var LuckDraw []models.LuckDraw
	var PotPrice []models.PotPrice
	var Achievement []models.Achievement
	var Section []models.Section
	var AutoGold []models.AutoGold
	var Harvest []models.Harvest
	var Lamp []models.Lamp

	//创建查询条件
	cropsqs := o.QueryTable("crops")
	//查询全部地区
	_, err := cropsqs.All(&crops)
	if err != nil {
	}
	scenedataqs := o.QueryTable("scene_data")
	//查询全部地区
	_, err = scenedataqs.All(&scenedata)
	if err != nil {
	}
	massifDataqs := o.QueryTable("massif_data")
	//查询全部地区
	_, err = massifDataqs.All(&MassifData)
	if err != nil {
	}
	AllProceedsqs := o.QueryTable("all_proceeds")
	//查询全部地区
	_, err = AllProceedsqs.All(&AllProceeds)
	if err != nil {
	}
	GoldEnhancementDataqs := o.QueryTable("gold_enhancement_data")
	//查询全部地区
	_, err = GoldEnhancementDataqs.All(&GoldEnhancementData)
	if err != nil {
	}
	MoneyEnhancementDataqs := o.QueryTable("money_enhancement_data")
	//查询全部地区
	_, err = MoneyEnhancementDataqs.All(&MoneyEnhancementData)
	if err != nil {
	}
	DiamondEnhancementDataqs := o.QueryTable("diamond_enhancement_data")
	//查询全部地区
	_, err = DiamondEnhancementDataqs.All(&DiamondEnhancementData)
	if err != nil {
	}

	//创建查询条件句柄植物碎片表
	PlantDebrisqs := o.QueryTable("plant_debris")
	//查询全部地区
	_, err = PlantDebrisqs.All(&PlantDebris)
	if err != nil {
	}
	//创建查询条件句柄植物碎片表
	NewPeopleqs := o.QueryTable("new_people")
	//查询全部地区
	_, err = NewPeopleqs.All(&NewPeople)
	if err != nil {
	}
	//创建查询条件句柄植物碎片表
	SignInqs := o.QueryTable("sign_in")
	//查询全部地区
	_, err = SignInqs.All(&SignIn)
	if err != nil {
	}
	//创建查询条件句柄植物碎片表
	FriendHelpqs := o.QueryTable("friend_help")
	//查询全部地区
	_, err = FriendHelpqs.All(&FriendHelp)
	if err != nil {
	}
	//创建查询条件句柄植物碎片表
	PotFragmentsqs := o.QueryTable("pot_fragments")
	//查询全部地区
	_, err = PotFragmentsqs.All(&PotFragments)
	if err != nil {
	}
	//创建查询条件句柄植物碎片表
	Propqs := o.QueryTable("prop")
	//查询全部地区
	_, err = Propqs.All(&Prop)
	if err != nil {
	}
	//创建查询条件句柄植物碎片表
	RandomPropqs := o.QueryTable("random_prop")
	//查询全部地区
	_, err = RandomPropqs.All(&RandomProp)
	if err != nil {
	}
	//创建查询条件句柄植物碎片表
	LuckDrawqs := o.QueryTable("luck_draw")
	//查询全部地区
	_, err = LuckDrawqs.All(&LuckDraw)
	if err != nil {
	}
	//创建查询条件句柄植物碎片表
	PotPriceqs := o.QueryTable("pot_price")
	//查询全部地区
	_, err = PotPriceqs.All(&PotPrice)
	if err != nil {
	}
	//创建查询条件句柄植物碎片表
	Achievementqs := o.QueryTable("achievement")
	//查询全部地区
	_, err = Achievementqs.All(&Achievement)
	if err != nil {
	}
	//创建查询条件句柄植物碎片表
	Sectionqs := o.QueryTable("section")
	//查询全部地区
	_, err = Sectionqs.All(&Section)
	if err != nil {
	}
	//创建查询条件句柄植物碎片表
	AutoGoldqs := o.QueryTable("auto_gold")
	//查询全部地区
	_, err = AutoGoldqs.All(&AutoGold)
	if err != nil {
	}
	//收割皮肤表
	Harvestqs := o.QueryTable("harvest")
	//查询全部地区
	_, err = Harvestqs.All(&Harvest)
	if err != nil {
	}
	//收割皮肤表
	Lampqs := o.QueryTable("lamp")
	//查询全部地区
	_, err = Lampqs.All(&Lamp)
	if err != nil {
	}

	//定义一个中间件用于发现proto中定义格式的植物切片,接收返回给web的数据
	var conftemp example.Config
	//遍历查出的crops切片

	for _, v := range crops {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempcrops := example.Crops{
			Id:                            utils.Encrypt(strconv.Itoa(v.Id)),
			Name:                          utils.Encrypt(v.Name),
			CropCode:                      utils.Encrypt(v.CropCode),
			DefaultUnlock:                 utils.Encrypt(strconv.Itoa(v.DefaultUnlock)),
			Subordinate:                   utils.Encrypt(strconv.Itoa(v.Subordinate)),
			CommonlyYield:                 utils.Encrypt(v.CommonlyYield),
			GreenYield:                    utils.Encrypt(v.GreenYield),
			BlueYield:                     utils.Encrypt(v.BlueYield),
			VioletYield:                   utils.Encrypt(v.VioletYield),
			OrangeYield:                   utils.Encrypt(v.OrangeYield),
			BasicProductionTime:           utils.Encrypt(strconv.Itoa(v.BasicProductionTime)),
			LimitProductionTime:           utils.Encrypt(v.LimitProductionTime),
			RegistrationRequiredPromotion: utils.Encrypt(strconv.Itoa(v.RegistrationRequiredPromotion)),
			StageLiftingSpeed:             utils.Encrypt(strconv.Itoa(v.StageLiftingSpeed)),
			UpgradeBaseGold:               utils.Encrypt(v.UpgradeBaseGold),
			ExhibitionUrl:                 utils.Encrypt(v.ExhibitionUrl),
			DynamicUrl:                    utils.Encrypt(v.DynamicUrl),
		}
		//把植物信息追加成切片
		conftemp.Cropslist = append(conftemp.Cropslist, &tempcrops)
	}
	//遍历查出的scenedata切片
	for _, v := range scenedata {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempscenedata := example.SceneData{
			ID:             utils.Encrypt(strconv.Itoa(v.Id)),
			SceneCode:      utils.Encrypt(v.SceneCode),
			Name:           utils.Encrypt(v.Name),
			OutputAddition: utils.Encrypt(strconv.Itoa(v.OutputAddition)),
			Permanent:      utils.Encrypt(strconv.Itoa(v.Permanent)),
			TimeLimit:      utils.Encrypt(strconv.Itoa(v.TimeLimit)),
			PriceDiamonds:  utils.Encrypt(strconv.Itoa(v.PriceDiamonds)),
			SmallUrl:       utils.Encrypt(v.SmallUrl),
			MaxUrl:         utils.Encrypt(v.MaxUrl),
		}
		//把植物信息追加成切片
		conftemp.SceneDatalist = append(conftemp.SceneDatalist, &tempscenedata)
	}
	//遍历查出的MassifData切片
	for _, v := range MassifData {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempMassifData := example.MassifData{
			Id:         utils.Encrypt(strconv.Itoa(v.Id)),
			SceneCode:  utils.Encrypt(strconv.Itoa(v.SceneCode)),
			OnePrice:   utils.Encrypt(v.OnePrice),
			TwoPrice:   utils.Encrypt(v.TwoPrice),
			ThreePrice: utils.Encrypt(v.ThreePrice),
			FourPrice:  utils.Encrypt(v.FourPrice),
			FivePrice:  utils.Encrypt(v.FivePrice),
			SixPrice:   utils.Encrypt(v.SixPrice),
			SevenPrice: utils.Encrypt(v.SevenPrice),
			EightPrice: utils.Encrypt(v.EightPrice),
			NinePrice:  utils.Encrypt(v.NinePrice),
		}
		//把植物信息追加成切片
		conftemp.MassifDatalist = append(conftemp.MassifDatalist, &tempMassifData)
	}
	//遍历查出的scenedata切片
	for _, v := range AllProceeds {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempAllProceeds := example.AllProceeds{
			Id:               utils.Encrypt(strconv.Itoa(v.Id)),
			Code:             utils.Encrypt(strconv.Itoa(v.Code)),
			ProceedsMultiple: utils.Encrypt(strconv.Itoa(v.ProceedsMultiple)),
			Consume:          utils.Encrypt(v.Consume),
			ArtsUrl:          utils.Encrypt(v.ArtsUrl),
		}
		//把植物信息追加成切片
		conftemp.AllProceedslist = append(conftemp.AllProceedslist, &tempAllProceeds)
	}
	//遍历查出的scenedata切片
	for _, v := range GoldEnhancementData {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempGoldEnhancementData := example.GoldEnhancementData{
			Id:                             utils.Encrypt(strconv.Itoa(v.Id)),
			MassifCode:                     utils.Encrypt(strconv.Itoa(v.MassifCode)),
			Level:                          utils.Encrypt(strconv.Itoa(v.Level)),
			EarningsMultiplier:             utils.Encrypt(strconv.Itoa(v.EarningsMultiplier)),
			IntensiveConsumption:           utils.Encrypt(v.IntensiveConsumption),
			AutomatedProductionConsumption: utils.Encrypt(v.AutomatedProductionConsumption),
			ArtsUrl:                        utils.Encrypt(v.ArtsUrl),
		}
		//把植物信息追加成切片
		conftemp.GoldEnhancementDatalist = append(conftemp.GoldEnhancementDatalist, &tempGoldEnhancementData)
	}
	//遍历查出的scenedata切片
	for _, v := range MoneyEnhancementData {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempMoneyEnhancementData := example.MoneyEnhancementData{
			Id:                   utils.Encrypt(strconv.Itoa(v.Id)),
			Level:                utils.Encrypt(strconv.Itoa(v.Level)),
			EarningsMultiplier:   utils.Encrypt(strconv.Itoa(v.EarningsMultiplier)),
			IntensiveConsumption: utils.Encrypt(v.IntensiveConsumption),
			ArtsUrl:              utils.Encrypt(v.ArtsUrl),
		}
		//把植物信息追加成切片
		conftemp.MoneyEnhancementDatalist = append(conftemp.MoneyEnhancementDatalist, &tempMoneyEnhancementData)
	}
	//遍历查出的scenedata切片
	for _, v := range DiamondEnhancementData {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempDiamondEnhancementData := example.DiamondEnhancementData{
			Id:                   utils.Encrypt(strconv.Itoa(v.Id)),
			AttributeClass:       utils.Encrypt(strconv.Itoa(v.AttributeClass)),
			EarningsMultiplier:   utils.Encrypt(strconv.Itoa(v.EarningsMultiplier)),
			ImmediateRevenue:     utils.Encrypt(strconv.Itoa(v.ImmediateRevenue)),
			IntensiveConsumption: utils.Encrypt(v.IntensiveConsumption),
			ArtsUrl:              utils.Encrypt(v.ArtsUrl),
		}
		//把植物信息追加成切片
		conftemp.DiamondEnhancementDatalist = append(conftemp.DiamondEnhancementDatalist, &tempDiamondEnhancementData)
	}
	//遍历查出的scenedata切片
	for _, v := range PlantDebris {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempPlantDebris := example.PlantDebris{
			Id:           utils.Encrypt(strconv.Itoa(v.Id)),
			PlantName:    utils.Encrypt(v.PlantName),
			PlantCode:    utils.Encrypt(v.PlantCode),
			PlantLevel:   utils.Encrypt(v.PlantLevel),
			FragmentCode: utils.Encrypt(v.FragmentCode),
		}
		//把植物信息追加成切片
		conftemp.PlantDebrislist = append(conftemp.PlantDebrislist, &tempPlantDebris)
	}
	//遍历查出的scenedata切片
	for _, v := range NewPeople {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempNewPeople := example.NewPeople{
			Id:           utils.Encrypt(strconv.Itoa(v.Id)),
			Number:       utils.Encrypt(strconv.Itoa(v.Number)),
			RewardCode:   utils.Encrypt(strconv.Itoa(v.RewardCode)),
			RewardNumber: utils.Encrypt(strconv.Itoa(v.RewardNumber)),
		}
		//把植物信息追加成切片
		conftemp.NewPeoplelist = append(conftemp.NewPeoplelist, &tempNewPeople)
	}
	//遍历查出的scenedata切片
	for _, v := range SignIn {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempSignIn := example.SignIn{
			Id:           utils.Encrypt(strconv.Itoa(v.Id)),
			Days:         utils.Encrypt(v.Days),
			RewardCode:   utils.Encrypt(v.RewardCode),
			RewardNumber: utils.Encrypt(v.RewardNumber),
		}
		//把植物信息追加成切片
		conftemp.SignInlist = append(conftemp.SignInlist, &tempSignIn)
	}
	//遍历查出的scenedata切片
	for _, v := range FriendHelp {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempFriendHelp := example.FriendHelp{
			Id:            utils.Encrypt(strconv.Itoa(v.Id)),
			PeopleNumber:  utils.Encrypt(v.PeopleNumber),
			Reward1Code:   utils.Encrypt(v.Reward1Code),
			Reward1Number: utils.Encrypt(v.Reward1Code),
			Reward2Code:   utils.Encrypt(v.Reward2Code),
			Reward2Number: utils.Encrypt(v.Reward2Number),
			Reward3Code:   utils.Encrypt(v.Reward3Code),
			Reward3Number: utils.Encrypt(v.Reward3Number),
		}
		//把植物信息追加成切片
		conftemp.FriendHelplist = append(conftemp.FriendHelplist, &tempFriendHelp)
	}
	//遍历查出的scenedata切片
	for _, v := range PotFragments {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempPotFragments := example.PotFragments{
			Id:           utils.Encrypt(strconv.Itoa(v.Id)),
			PlantName:    utils.Encrypt(v.PlantName),
			PlantCode:    utils.Encrypt(v.PlantCode),
			PlantLevel:   utils.Encrypt(v.PlantLevel),
			FragmentCode: utils.Encrypt(v.FragmentCode),
		}
		//把植物信息追加成切片
		conftemp.PotFragmentslist = append(conftemp.PotFragmentslist, &tempPotFragments)
	}
	//遍历查出的scenedata切片
	for _, v := range Prop {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempProp := example.Prop{
			Id:        utils.Encrypt(strconv.Itoa(v.Id)),
			PropCode:  utils.Encrypt(v.PropCode),
			PropName:  utils.Encrypt(v.PropName),
			Introduce: utils.Encrypt(v.Introduce),
		}
		//把植物信息追加成切片
		conftemp.Proplist = append(conftemp.Proplist, &tempProp)
	}
	//遍历查出的scenedata切片
	for _, v := range RandomProp {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempRandomProp := example.RandomProp{
			Id:           utils.Encrypt(strconv.Itoa(v.Id)),
			PropCode:     utils.Encrypt(v.PropCode),
			PropName:     utils.Encrypt(v.PropName),
			Reward1:      utils.Encrypt(v.Reward1),
			Reward1Max:   utils.Encrypt(v.Reward1Max),
			Reward1Small: utils.Encrypt(v.Reward1Small),
			Reward2:      utils.Encrypt(v.Reward2),
			Reward2Max:   utils.Encrypt(v.Reward2Max),
			Reward2Small: utils.Encrypt(v.Reward2Small),
			Reward3:      utils.Encrypt(v.Reward3),
			OutputTime:   utils.Encrypt(v.OutputTime),
		}
		//把植物信息追加成切片
		conftemp.RandomProplist = append(conftemp.RandomProplist, &tempRandomProp)
	}
	//遍历查出的scenedata切片
	for _, v := range LuckDraw {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempLuckDraw := example.LuckDraw{
			Id:          utils.Encrypt(strconv.Itoa(v.Id)),
			PlantCode:   utils.Encrypt(v.PlantCode),
			PlantNumber: utils.Encrypt(v.PlantNumber),
			Probability: utils.Encrypt(v.Probability),
			PropName:    utils.Encrypt(v.PropName),
		}
		//把植物信息追加成切片
		conftemp.LuckDrawlist = append(conftemp.LuckDrawlist, &tempLuckDraw)
	}
	//遍历查出的scenedata切片
	for _, v := range PotPrice {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempPotPrice := example.PotPrice{
			Id:            utils.Encrypt(strconv.Itoa(v.Id)),
			PotCode:       utils.Encrypt(v.PotCode),
			PotName:       utils.Encrypt(v.PotName),
			CurrencyCode:  utils.Encrypt(v.CurrencyCode),
			CurrencyNumbe: utils.Encrypt(v.CurrencyNumber),
			FreeTime:      utils.Encrypt(v.FreeTime),
			LuckDrawCode:  utils.Encrypt(v.LuckDrawCode),
		}
		//把植物信息追加成切片
		conftemp.PotPricelist = append(conftemp.PotPricelist, &tempPotPrice)
	}
	//遍历查出的scenedata切片
	for _, v := range Achievement {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempAchievement := example.Achievement{
			Id:             utils.Encrypt(strconv.Itoa(v.Id)),
			UserLevel:      utils.Encrypt(strconv.Itoa(v.UserLevel)),
			TitleCode:      utils.Encrypt(v.TitleCode),
			TitlName:       utils.Encrypt(v.TitlName),
			CumulativeDiam: utils.Encrypt(v.CumulativeDiamond),
			CumulativeGold: utils.Encrypt(v.CumulativeGold),
			CumulativeMone: utils.Encrypt(v.CumulativeMoney),
			HarvestTimes:   utils.Encrypt(v.HarvestTimes),
			Reward1Code:    utils.Encrypt(v.Reward1Code),
			Reward1Number:  utils.Encrypt(v.Reward1Number),
			Reward2Code:    utils.Encrypt(v.Reward2Code),
			Reward2Number:  utils.Encrypt(v.Reward2Number),
			Reward3Code:    utils.Encrypt(v.Reward3Code),
			Reward3Number:  utils.Encrypt(v.Reward2Number),
			Url:            utils.Encrypt(v.Url),
			Promote:        utils.Encrypt(v.Promote),
			D1:             utils.Encrypt(v.D1),
			D2:             utils.Encrypt(v.D2),
			D3:             utils.Encrypt(v.D3),
			D4:             utils.Encrypt(v.D4),
			D5:             utils.Encrypt(v.D5),
			D6:             utils.Encrypt(v.D6),
			D7:             utils.Encrypt(v.D7),
			D8:             utils.Encrypt(v.D8),
			D9:             utils.Encrypt(v.D9),
		}
		//把植物信息追加成切片
		conftemp.Achievementlist = append(conftemp.Achievementlist, &tempAchievement)
	}
	//遍历查出的scenedata切片
	for _, v := range Section {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempSection := example.Section{
			Id:        utils.Encrypt(strconv.Itoa(v.Id)),
			SmallGold: utils.Encrypt(v.SmallGold),
			MaxGold:   utils.Encrypt(v.MaxGold),
			Multiple:  utils.Encrypt(v.Multiple),
		}
		//把植物信息追加成切片
		conftemp.Sectionlist = append(conftemp.Sectionlist, &tempSection)
	}
	//遍历查出的scenedata切片
	for _, v := range AutoGold {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempAutoGold := example.AutoGold{
			Id:          utils.Encrypt(strconv.Itoa(v.Id)),
			SceneCode:   utils.Encrypt(strconv.Itoa(v.SceneCode)),
			Auto:        utils.Encrypt(strconv.Itoa(v.Auto)),
			AutoConsume: utils.Encrypt(v.AutoConsume),
			Url:         utils.Encrypt(v.Url),
		}
		//把植物信息追加成切片
		conftemp.AutoGoldlist = append(conftemp.AutoGoldlist, &tempAutoGold)
	}

	for _, v := range Harvest {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempHarvest := example.HarvestSkin{
			Id:         utils.Encrypt(strconv.Itoa(v.Id)),
			CarSkin:    utils.Encrypt(v.CarSkin),
			CarName:    utils.Encrypt(v.CarName),
			CarProduce: utils.Encrypt(v.CarProduce),
			Permanent:  utils.Encrypt(v.Permanent),
			CarTime:    utils.Encrypt(v.CarTime),
			Price:      utils.Encrypt(v.Price),
			Url:        utils.Encrypt(v.Url),
		}
		//把植物信息追加成切片
		conftemp.Harvestlist = append(conftemp.Harvestlist, &tempHarvest)
	}

	for _, v := range Lamp {
		//初始化临时crops,类型为proto定义类型,并用遍历出的数据给其赋值
		tempLamp := example.Lamp{
			Id:           utils.Encrypt(strconv.Itoa(v.Id)),
			LampSkin:     utils.Encrypt(v.LampSkin),
			LampName:     utils.Encrypt(v.LampName),
			LampAddition: utils.Encrypt(v.LampAddition),
			Permanent:    utils.Encrypt(v.Permanent),
			LampTime:     utils.Encrypt(v.LampTime),
			Price:        utils.Encrypt(v.Price),
			Url:          utils.Encrypt(v.Url),
		}
		//把植物信息追加成切片
		conftemp.Lamplist = append(conftemp.Lamplist, &tempLamp)
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

	//从缓存中获取opid
	hxopid := bm.Get(req.Skey)
	logs.Info(hxopid)
	if hxopid == nil {
		logs.Info("获取opid失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//将缓存中取出的opid转成string
	temp, _ := redis.String(hxopid, nil)
	logs.Info(temp)
	//判断新用户
	tempuser := bm.Get(temp)
	if tempuser == nil {
		// 没有找到记录
		logs.Info("这是新用户")
		//从缓存中获取opid对应的data数据
		user_data := bm.Get("test")
		if user_data == nil {
			logs.Info("获取用户数据失败")
			rsp.Errno = utils.RECODE_DBERR
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
		rsp.Datalist = user_data.([]byte)
		rsp.Configlist = &conftemp
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	logs.Info("这个老用户")
	//从缓存中获取opid对应的data数据
	user_data := bm.Get(temp)
	if user_data == nil {
		logs.Info("获取用户数据失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	logs.Info("获取数据成功")
	//返回数据给web端
	rsp.Datalist = user_data.([]byte)
	rsp.Configlist = &conftemp
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	return nil
}
