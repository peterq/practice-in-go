package yeb_exp

import (
	"context"
	"funny/spider_client"
	"funny/yeb_exp/util"
	"gopkg.in/mgo.v2"
	"time"
)

type Config struct {
	MongoUrl,
	UserCollectionName,
	MobileColumn,
	DbName string
	InviteParam map[string]string
	ShareLink   string
	mgo         struct {
		session *mgo.Session
		user    *mgo.Collection
		db      *mgo.Database
	}
}

var appConfig Config
var apiClient *spider_client.Client
var appCtx context.Context = context.WithValue(context.Background(), "START", time.Now())

func config() {
	appConfig = Config{
		MongoUrl:           "mongodb://root:root@127.0.0.1:27017/admin",
		UserCollectionName: "user",
		MobileColumn:       "_id",
		DbName:             "yu_e_bao",
		InviteParam: map[string]string{
			"campInfo":    "p1j%2BdzkZl03BbvY4ClDID7%2FGiTlLCEEq0EmyB5yLfH2lIGLw2ZgnyTMSXAcf56tw",
			"bizType":     "c2cShare",
			"shareId":     "2088512116812823",
			"snsScene":    "yebTrialFoundSns",
			"sign":        "%2FQ5V8dp3ZBDbouaI1ISjaxAIzEkYnJ7osU3d8OQM1wI%3D",
			"_json_token": "undefined",
		},
		ShareLink: "http://render.aa43z7.com/p/f/fd-jqqeh032/pages/outside/index.html?shareid=2088512116812823&sign=%2FQ5V8dp3ZBDbouaI1ISjaxAIzEkYnJ7osU3d8OQM1wI%3D",
		mgo: struct {
			session *mgo.Session
			user    *mgo.Collection
			db      *mgo.Database
		}{},
	}
	s, err := mgo.Dial(appConfig.MongoUrl)
	if err != nil {
		panic(err)
	}
	s.SetMode(mgo.Monotonic, true)
	appConfig.mgo.session = s
	appConfig.mgo.db = s.DB(appConfig.DbName)
	appConfig.mgo.user = appConfig.mgo.db.C(appConfig.UserCollectionName)
	apiClient = spider_client.New(5, 5, 0, true)
}

func Init() {
	config()
	mobileCh := getUser()
	util.GetUa(appCtx, appConfig.ShareLink, mobileCh)
	select {}
	//uaCh := getJsonUa()
	//doCallApi(mobileCh, uaCh)
}
