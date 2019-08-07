package models

import (
	//go操作数据库的模块
	_ "github.com/go-sql-driver/mysql"
	"plant/gateway/utils"
	"github.com/astaxie/beego/orm"
)

//场景属性表
type SceneData struct {
	Id             int    `json:"scenedata_id"`
	SceneCode      string `orm:"size(256)" json:"scene_code"`
	Name           string `orm:"size(256)" json:"name"`
	OutputAddition int    `json:"output_addition"`
	Permanent      int    `json:"permanent"`
	TimeLimit      int    `json:"time_limit"`
	PriceDiamonds  int    `json:"price_diamonds"`
	SmallUrl       string `orm:"size(256)" json:"small_url"`
	MaxUrl         string `orm:"size(256)" json:"max_url"`
}

//各场景地块的开启金币数
type MassifData struct {
	Id         int    `json:"Massifdata_id"`
	SceneCode  int    `json:"scene_code"`
	OnePrice   string `orm:"size(256)" json:"one_price"`
	TwoPrice   string `orm:"size(256)" json:"two_price"`
	ThreePrice string `orm:"size(256)" json:"three_price"`
	FourPrice  string `orm:"size(256)" json:"four_price"`
	FivePrice  string `orm:"size(256)" json:"five_price"`
	SixPrice   string `orm:"size(256)" json:"six_price"`
	SevenPrice string `orm:"size(256)" json:"seven_price"`
	EightPrice string `orm:"size(256)" json:"eight_price"`
	NinePrice  string `orm:"size(256)" json:"nine_price"`
}

//农作物属性表
type Crops struct {
	Id                            int    `json:"id"`
	Name                          string `orm:"size(256)" json:"name"`
	CropCode                      string `orm:"size(256)" json:"crop_code"`
	DefaultUnlock                 int    `json:"default_unlock"`
	Subordinate                   int    `json:"subordinate"`
	CommonlyYield                 string `orm:"size(256)" json:"commonly_yield"`
	GreenYield                    string `orm:"size(256)" json:"green_yield"`
	BlueYield                     string `orm:"size(256)" json:"blue_yield"`
	VioletYield                   string `orm:"size(256)" json:"violet_yield"`
	OrangeYield                   string `orm:"size(256)" json:"orange_yield"`
	BasicProductionTime           int    `json:"basic_production_time"`
	LimitProductionTime           string `json:"limit_production_time"`
	RegistrationRequiredPromotion int    `json:"registration_required_promotion"`
	StageLiftingSpeed             int    `json:"stage_lifting_speed"`
	UpgradeBaseGold               string `orm:"size(256)" json:"upgrade_base_gold"`
	ExhibitionUrl                 string `orm:"size(256)" json:"exhibition_url"`
	DynamicUrl                    string `orm:"size(256)" json:"dynamic_url"`
}

//全部收益数据
type AllProceeds struct {
	Id               int    `json:"id"`
	Code             int    `json:"code"`
	ProceedsMultiple int    `json:"proceeds_multiple"`
	Consume          string `orm:"size(256)" json:"consume"`
	ArtsUrl          string `orm:"size(256)" json:"arts_url"`
}

//金币强化表
type GoldEnhancementData struct {
	Id                             int    `json:"id"`
	MassifCode                     int    `json:"massif_code"`
	Level                          int    `json:"level"`
	EarningsMultiplier             int    `json:"earnings_multiplier"`
	IntensiveConsumption           string `orm:"size(256)" json:"intensive_consumption"`
	AutomatedProductionConsumption string `orm:"size(256)" json:"automated_production_consumption"`
	ArtsUrl                        string `orm:"size(256)" json:"arts_url"`
}

//现金强化表
type MoneyEnhancementData struct {
	Id                   int    `json:"id"`
	Level                int    `json:"level"`                                 //等级
	EarningsMultiplier   int    `json:"earnings_multiplier"`                   //收益倍数
	IntensiveConsumption string `orm:"size(256)" json:"intensive_consumption"` //强化消耗
	ArtsUrl              string `orm:"size(256)" json:"arts_url"`              //美术资源
}

//钻石强化表
type DiamondEnhancementData struct {
	Id                   int    `json:"id"`
	AttributeClass       int    `json:"attribute_class"`                       //属性类别
	EarningsMultiplier   int    `json:"earnings_multiplier"`                   //增加全部产出倍数
	ImmediateRevenue     int    `json:"immediate_revenue"`                     //立即收益
	IntensiveConsumption string `orm:"size(256)" json:"intensive_consumption"` //强化消耗
	ArtsUrl              string `orm:"size(256)" json:"arts_url"`              //美术资源
}

//植物碎片表
type PlantDebris struct {
	Id           int    `json:"id"`
	PlantName    string `orm:"size(256)" json:"plant_name"`    //植物名称
	PlantCode    string `orm:"size(256)" json:"plant_code"`    //植物id
	PlantLevel   string `orm:"size(256)" json:"plant_level"`   //植物品阶
	FragmentCode string `orm:"size(256)" json:"fragment_code"` //碎片id
}

//邀请新人数值表
type NewPeople struct {
	Id           int    `json:"id"`
	Number       int    `json:"number"`        //邀请新人数
	RewardCode   int    `json:"reward_code"`   //奖励id
	RewardNumber int    `json:"reward_number"` //奖励数量
	GetStatus    string `orm:"size(256)" json:"get_status"`
}

//每日签到数值表
type SignIn struct {
	Id           int    `json:"id"`
	Days         string `orm:"size(256)" json:"days"`          //天数
	RewardCode   string `orm:"size(256)" json:"reward_code"`   //奖励id
	RewardNumber string `orm:"size(256)" json:"reward_number"` //奖励数量
}

//好友助理数值表
type FriendHelp struct {
	Id            int    `json:"id"`
	PeopleNumber  string `orm:"size(256)" json:"people_number"`   //需要邀请人数
	Reward1Code   string `orm:"size(256)" json:"reward_1_code"`   //奖励1id
	Reward1Number string `orm:"size(256)" json:"reward_1_number"` //奖励1数量
	Reward2Code   string `orm:"size(256)" json:"reward_2_code"`   //奖励2id
	Reward2Number string `orm:"size(256)" json:"reward_2_number"` //奖励2数量
	Reward3Code   string `orm:"size(256)" json:"reward_3_code"`   //奖励3id
	Reward3Number string `orm:"size(256)" json:"reward_3_number"` //奖励3数量
}

//罐子可抽取数值表
type PotFragments struct {
	Id           int    `json:"id"`
	PlantName    string `orm:"size(256)" json:"plant_name"`
	PlantCode    string `orm:"size(256)" json:"plant_code"`    //植物id
	PlantLevel   string `orm:"size(256)" json:"plant_level"`   //植物品阶
	FragmentCode string `orm:"size(256)" json:"fragment_code"` //碎片id
}

//道具表
type Prop struct {
	Id        int    `json:"id"`
	PropCode  string `orm:"size(256)" json:"prop_code"` //道具id
	PropName  string `orm:"size(256)" json:"prop_name"` //道具名字
	Introduce string `orm:"size(256)" json:"introduce"` //介绍
}

//随机道具表
type RandomProp struct {
	Id           int    `json:"id"`
	PropCode     string `orm:"size(256)" json:"prop_code"`      //道具id
	PropName     string `orm:"size(256)" json:"prop_name"`      //道具名字
	Reward1      string `orm:"size(256)" json:"reward_1"`       //奖励1
	Reward1Max   string `orm:"size(256)" json:"reward_1_max"`   //奖励1最大值
	Reward1Small string `orm:"size(256)" json:"reward_1_small"` //奖励1最小值
	Reward2      string `orm:"size(256)" json:"reward_2"`       //奖励2
	Reward2Max   string `orm:"size(256)" json:"reward_2_max"`   //奖励2最大值
	Reward2Small string `orm:"size(256)" json:"reward_2_small"` //奖励2最小值
	Reward3      string `orm:"size(256)" json:"reward_3"`       //奖励2
	OutputTime   string `orm:"size(256)" json:"output_time"`    //当前金币产出量时间
}

//抽奖数值表
type LuckDraw struct {
	Id          int    `json:"id"`
	PlantCode   string `orm:"size(256)" json:"plant_code"`   //植物id
	PlantNumber string `orm:"size(256)" json:"plant_number"` //奖励数量
	Probability string `orm:"size(256)" json:"probability"`  //几率
	PropName    string `orm:"size(256)" json:"prop_name"`    //道具名字
}

//抽罐子价格表
type PotPrice struct {
	Id             int    `json:"id"`
	PotCode        string `orm:"size(256)" json:"pot_code"`        //罐子id
	PotName        string `orm:"size(256)" json:"pot_name"`        //罐子名字
	CurrencyCode   string `orm:"size(256)" json:"currency_code"`   //货币id
	CurrencyNumber string `orm:"size(256)" json:"currency_number"` //花费数量
	FreeTime       string `orm:"size(256)" json:"free_time"`       //免费机会时间（h）
	LuckDrawCode   string `orm:"size(256)" json:"luck_draw_code"`  //抽奖类型
}

//成就数值表
type Achievement struct {
	Id                int    `json:"id"`
	UserLevel         int    `json:"user_level"`                         //用户等级
	TitleCode         string `orm:"size(256)" json:"title_code"`         //称号奖励id
	TitlName          string `orm:"size(256)" json:"titl_name"`          //称号名字
	CumulativeDiamond string `orm:"size(256)" json:"cumulative_diamond"` //累计拥有钻石数
	CumulativeGold    string `orm:"size(256)" json:"consume_gold"`       //累计拥有金币数
	CumulativeMoney   string `orm:"size(256)" json:"cumulative_money"`   //累计拥有现金数
	HarvestTimes      string `orm:"size(256)" json:"harvest_times"`      //累计首歌次数
	Reward1Code       string `orm:"size(256)" json:"reward_1_code"`      //奖励1id
	Reward1Number     string `orm:"size(256)" json:"reward_1_number"`    //奖励1数量
	Reward2Code       string `orm:"size(256)" json:"reward_2_code"`      //奖励2id
	Reward2Number     string `orm:"size(256)" json:"reward_2_number"`    //奖励2数量
	Reward3Code       string `orm:"size(256)" json:"reward_3_code"`      //奖励3id
	Reward3Number     string `orm:"size(256)" json:"reward_3_number"`    //奖励4数量
	Reward4Code       string `orm:"size(256)" json:"reward_4_code"`      //奖励3id
	Reward4Number     string `orm:"size(256)" json:"reward_4_number"`    //奖励4数量
	Url               string `orm:"size(256)" json:"url"`                //美术资源
	Promote           string `orm:"size(256)" json:"promote"`            //产出倍数永久提升
	D1                string `orm:"size(256)" json:"d_1"`                //地块1
	D2                string `orm:"size(256)" json:"d_2"`                //地块2
	D3                string `orm:"size(256)" json:"d_3"`                //地块3
	D4                string `orm:"size(256)" json:"d_4"`                //地块4
	D5                string `orm:"size(256)" json:"d_5"`                //地块5
	D6                string `orm:"size(256)" json:"d_6"`                //地块6
	D7                string `orm:"size(256)" json:"d_7"`                //地块7
	D8                string `orm:"size(256)" json:"d_8"`                //地块8
	D9                string `orm:"size(256)" json:"d_9"`                //地块9
}

//收割后产出倍数区间表
type Section struct {
	Id        int
	SmallGold string `orm:"size(256)" json:"small_gold"` //最小现金
	MaxGold   string `orm:"size(256)" json:"max_gold"`   //最大现金
	Multiple  string `orm:"size(256)" json:"multiple"`   //产出倍数
}

//金币自动生存
type AutoGold struct {
	Id          int    `json:"id"`
	SceneCode   int    `json:"scene_code"`                   //地块id
	Auto        int    `json:"auto"`                         //自动生存初始状态
	AutoConsume string `orm:"size(256)" json:"auto_consume"` //自动生存金币消耗
	Url         string `orm:"size(256)" json:"url"`          //美术资源
}

//收割皮肤表
type Harvest struct {
	Id         int    `json:"id"`
	CarSkin    string `orm:"size(256)" json:"car_skin"`    //汽车皮肤
	CarName    string `orm:"size(256)" json:"car_name"`    //汽车名词
	CarProduce string `orm:"size(256)" json:"car_produce"` //汽车产出
	Permanent  string `orm:"size(256)" json:"permanent"`   //是否永久
	CarTime    string `orm:"size(256)" json:"car_time"`    //汽车时限
	Price      string `orm:"size(256)" json:"price"`       //出售价格
	Url        string `orm:"size(256)" json:"url"`         //美术资源
}

//灯皮肤表
type Lamp struct {
	Id           int    `json:"id"`
	LampSkin     string `orm:"size(256)" json:"lamp_skin"`     //灯皮肤
	LampName     string `orm:"size(256)" json:"lamp_name"`     //灯名字
	LampAddition string `orm:"size(256)" json:"lamp_addition"` //灯产出加成
	Permanent    string `orm:"size(256)" json:"permanent"`     //是否永久
	LampTime     string `orm:"size(256)" json:"lamp_time"`     //灯限时
	Price        string `orm:"size(256)" json:"price"`         //出售价格
	Url          string `orm:"size(256)" json:"url"`           //美术资源
}

//用户
type User struct {
	Id                int            `json:"user_id"`                           //用户编号
	OpenId            string         `orm:"size(128)" json:"open_id"`           //从微信获取的用户Id
	UnionId           string         `orm:"size(128)" json:"union_id"`          //从微信获取的全平台唯一Id
	Hxopid            string         `orm:"size(256)" json:"hxopid"`            //哈希处理后的opid
	Name              string         `orm:"size(128)"  json:"name"`             //用户昵称
	Password_hash     string         `orm:"size(128)" json:"password"`          //用户密码加密的
	Gender            string         `orm:"size(128)" json:"gender"`            //性别
	Age               string         `orm:"size(128)" json:"age"`               //年龄
	Session_key       string         `orm:"size(128)" json:"session_key"`       //从微信获取的sessoing_key
	Skey              string         `orm:"size(128)" json:"skey"`              //自定义登陆态
	Mobile            string         `orm:"size(11)"  json:"mobile"`            //手机号
	Real_name         string         `orm:"size(32)" json:"real_name"`          //真实姓名  实名认证
	Id_card           string         `orm:"size(20)" json:"id_card"`            //身份证号  实名认证
	Avatar_url        string         `orm:"size(256)" json:"avatar_url"`        //用户头像路径       通过fastdfs进行图片存储
	Landing_time      string         `orm:"size(256)" json:"landing_time"`      //最后登陆时间
	Registration_time string         `orm:"size(256)" json:"registration_time"` //用户注册时间
	Invitations       []*Invitation  `orm:"reverse(many)" json:"invitations"`   //邀请的新人信息
	Strengths         []*Strength    `orm:"reverse(many)" json:"strength"`      //体力卡包的助力信息
	Diamondcards      []*Diamondcard `orm:"reverse(many)" json:"diamondcards"`  //钻石卡报的助力信息
	Lampshares        []*Lampshare   `orm:"reverse(many)" json:"lampshares"`    //助力灯的信息
}

//邀请新人的信息
type Invitation struct {
	Id                    int    `json:"Invitation_id"`
	User                  *User  `orm:"rel(fk)" json:"open_id"`                 //邀请人opid
	Invitation_OpenId     string `orm:"size(128)" json:"invitation_open_id"`    //从微信获取的用户Id
	Invitation_Avatar_url string `orm:"size(256)" json:"invitation_avatar_url"` //从微信获取的头像地址
	Invitation_Name       string `orm:"size(128)" json:"invitation_name"`       //用户昵称
	Invitation_time       string `orm:"size(128)" json:"invitation_time"`       //分享时间
}

//助力体力的信息
type Strength struct {
	Id                    int    `json:"Invitation_id"`
	User                  *User  `orm:"rel(fk)" json:"open_id"`                 //邀请人opid
	Invitation_OpenId     string `orm:"size(128)" json:"invitation_open_id"`    //从微信获取的用户Id
	Invitation_Avatar_url string `orm:"size(256)" json:"invitation_avatar_url"` //从微信获取的头像地址
	Invitation_Name       string `orm:"size(128)" json:"invitation_name"`       //用户昵称
	Invitation_time       string `orm:"size(128)" json:"invitation_time"`       //分享时间
	GetStatus             string `orm:"size(128)" json:"get_status"`            //领取状态
}

//助力钻石的信息
type Diamondcard struct {
	Id                    int    `json:"Invitation_id"`
	User                  *User  `orm:"rel(fk)" json:"open_id"`                 //邀请人opid
	Invitation_OpenId     string `orm:"size(128)" json:"invitation_open_id"`    //从微信获取的用户Id
	Invitation_Avatar_url string `orm:"size(256)" json:"invitation_avatar_url"` //从微信获取的头像地址
	Invitation_Name       string `orm:"size(128)" json:"invitation_name"`       //用户昵称
	Invitation_time       string `orm:"size(128)" json:"invitation_time"`       //分享时间
	GetStatus             string `orm:"size(128)" json:"get_status"`            //领取状态
}

//助力灯的信息
type Lampshare struct {
	Id                    int    `json:"Invitation_id"`
	User                  *User  `orm:"rel(fk)" json:"open_id"`                 //邀请人opid
	Invitation_OpenId     string `orm:"size(128)" json:"invitation_open_id"`    //从微信获取的用户Id
	Invitation_Avatar_url string `orm:"size(256)" json:"invitation_avatar_url"` //从微信获取的头像地址
	Invitation_Name       string `orm:"size(128)" json:"invitation_name"`       //用户昵称
	Invitation_time       string `orm:"size(128)" json:"invitation_time"`       //分享时间
	GetStatus             string `orm:"size(128)" json:"get_status"`            //领取状态
}

type TreasureMap struct {
	Id       int         `json:"id"`
	State    string      `orm:"size(128)" json:"state"`          //当前位置奖励领取状态,状态 0 ：旗子  1 未打开的宝箱  2  已打开宝箱
	Type     string      `orm:"size(128)" json:"type"`           //当前位置奖励宝箱类型,0:小宝箱,1:大宝箱
	Url      string      `orm:"size(256)" json:"url"`            //助力人头像地址
	BoxBigs  []*BoxBig   `orm:"reverse(many)" json:"box_bigs"`   //大宝箱奖励内容
	BoXSmall []*BoxSmall `orm:"reverse(many)" json:"bo_x_small"` //小宝箱奖励内容
}

type BoxBig struct {
	Id          int          `json:"id"`
	Code        string       `orm:"size(128)" json:"lucky_code"` //奖励类型
	Name        string       `orm:"size(128)" json:"name"`       //奖励名字
	Number      string       `orm:"size(128)" json:"number"`     //奖励数量
	TreasureMap *TreasureMap `orm:"rel(fk)" json:"treasure_map"`
}

type BoxSmall struct {
	Id          int          `json:"id"`
	Code        string       `orm:"size(128)" json:"lucky_code"` //奖励类型
	Name        string       `orm:"size(128)" json:"name"`       //奖励名字
	Number      string       `orm:"size(128)" json:"number"`     //奖励数量
	TreasureMap *TreasureMap `orm:"rel(fk)" json:"treasure_map"`
}

func init() {
	//注册mysql的驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// 设置默认数据库
	orm.RegisterDataBase("default", "mysql", utils.G_mysql_user+":"+utils.G_mysql_password+"@tcp("+utils.G_mysql_addr+":"+utils.G_mysql_port+")/"+utils.G_mysql_dbname+"?charset=utf8", 30)
	//orm.RegisterDataBase("default", "mysql", "root:1@tcp("+utils.G_mysql_addr+":"+utils.G_mysql_port+")/plant?charset=utf8", 30)

	//注册model
	orm.RegisterModel(new(User), new(SceneData), new(MassifData),
		new(Crops), new(AllProceeds), new(GoldEnhancementData),
		new(MoneyEnhancementData), new(DiamondEnhancementData), new(PlantDebris),
		new(NewPeople), new(SignIn), new(FriendHelp), new(PotFragments), new(Prop),
		new(RandomProp), new(LuckDraw), new(PotPrice), new(Achievement), new(AutoGold),
		new(Section), new(Harvest), new(Lamp), new(Invitation),
		new(Strength), new(Diamondcard), new(Lampshare), new(TreasureMap), new(BoxBig),
		new(BoxSmall))

	// create table
	//第一个是别名
	//第二个是是否强制替换模块   如果表变更就将false 换成true 之后再换回来表就便更好来了
	//第三个参数是如果没有则同步或创建
	orm.RunSyncdb("default", false, true)
}
