package models_lamp

import (
	//go操作数据库的模块
	_ "github.com/go-sql-driver/mysql"
	"plant/gateway/utils"
	"github.com/astaxie/beego/orm"
)

//用户
type User struct {
	Id                int           `json:"user_id"`                           //用户编号
	OpenId            string        `orm:"size(128)" json:"open_id"`           //从微信获取的用户Id
	UnionId           string        `orm:"size(128)" json:"union_id"`          //从微信获取的全平台唯一Id
	Hxopid            string        `orm:"size(256)" json:"hxopid"`            //哈希处理后的opid
	Name              string        `orm:"size(128)"  json:"name"`             //用户昵称
	Password_hash     string        `orm:"size(128)" json:"password"`          //用户密码加密的
	Gender            string        `orm:"size(128)" json:"gender"`            //性别
	Age               string        `orm:"size(128)" json:"age"`               //年龄
	Session_key       string        `orm:"size(128)" json:"session_key"`       //从微信获取的sessoing_key
	Skey              string        `orm:"size(128)" json:"skey"`              //自定义登陆态
	Mobile            string        `orm:"size(11)"  json:"mobile"`            //手机号
	Real_name         string        `orm:"size(32)" json:"real_name"`          //真实姓名  实名认证
	Id_card           string        `orm:"size(20)" json:"id_card"`            //身份证号  实名认证
	Avatar_url        string        `orm:"size(256)" json:"avatar_url"`        //用户头像路径       通过fastdfs进行图片存储
	Landing_time      string        `orm:"size(256)" json:"landing_time"`      //最后登陆时间
	Registration_time string        `orm:"size(256)" json:"registration_time"` //用户注册时间
	Checkpoint        string        `orm:"size(128)" json:"checkpoint"`        //当前关卡
	MaxLevel          string        `orm:"size(256)" json:"max_level"`         //最高关卡数
	CurrentLevel      string        `orm:"size(256)" json:"current_level"`     //当前关卡
	CurrentLampCode   string        `orm:"size(256)" json:"current_lamp_code"` //当前使用的灯的code
	UserMoney         string        `orm:"size(256)" json:"user_money"`        //当前金币数
	LevelTime         string        `orm:"size(256)" json:"level_time"`        //过关时间
	Invitations       []*Invitation `orm:"reverse(many)" json:"invitations"`   //邀请的新人信息
	Lamps             []*Lamp       `orm:"reverse(many)" json:"lamps"`         //拥有的灯信息
}

//邀请新人的信息
type Invitation struct {
	Id                    int    `json:"Invitation_id"`
	User                  *User  `orm:"rel(fk)" json:"open_id"`                 //邀请人opid
	Invitation_OpenId     string `orm:"size(128)" json:"invitation_open_id"`    //从微信获取的用户Id
	Invitation_Avatar_url string `orm:"size(256)" json:"invitation_avatar_url"` //从微信获取的头像地址
	Invitation_Name       string `orm:"size(128)" json:"invitation_name"`       //用户昵称
	Invitation_time       string `orm:"size(128)" json:"invitation_time"`       //分享时间
	GetStatus     string `orm:"size(128)" json:"get_status"`    //奖励领取
}

//灯皮肤表
type Lamp struct {
	Id        int    `json:"id"`
	User      *User  `orm:"rel(fk)" json:"open_id"`      //所属用户信息
	LampId    string `orm:"size(128)" json:"lamp_id"`    //灯的code
	LampPrice string `orm:"size(128)" json:"lamp_price"` //灯的价格
	LampUrl   string `orm:"size:(128)" json:"lamp_url"`  //美术资源
	IsHave    string `orm:"size:(128)" json:"is_have"`   //用户是否拥有
}

//灯皮肤表
type Lampconfig struct {
	Id        int    `json:"id"`
	//User      *User  `orm:"rel(fk)" json:"open_id"`      //所属用户信息
	LampId    string `orm:"size(128)" json:"lamp_id"`    //灯的code
	LampPrice string `orm:"size(128)" json:"lamp_price"` //灯的价格
	LampUrl   string `orm:"size:(128)" json:"lamp_url"`  //美术资源
	IsHave    string `orm:"size:(128)" json:"is_have"`   //用户是否拥有
}


func init() {
	//注册mysql的驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// 设置默认数据库
	orm.RegisterDataBase("default", "mysql", utils.G_mysql_user_lamp+":"+utils.G_mysql_password_lamp+"@tcp("+utils.G_mysql_addr+":"+utils.G_mysql_port+")/"+utils.G_mysql_dbname_lamp+"?charset=utf8", 30)
	//orm.RegisterDataBase("default", "mysql", "root:1@tcp("+utils.G_mysql_addr+":"+utils.G_mysql_port+")/plant?charset=utf8", 30)

	//注册model
	orm.RegisterModel(new(User), new(Lamp), new(Invitation),new(Lampconfig))

	// create table
	//第一个是别名
	//第二个是是否强制替换模块   如果表变更就将false 换成true 之后再换回来表就便更好来了
	//第三个参数是如果没有则同步或创建
	orm.RunSyncdb("default", false, true)
}
