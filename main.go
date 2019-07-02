package main

import (
	"flag"
	"github.com/mlogclub/mlog-wxbot/config"
	"github.com/mlogclub/simple"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"github.com/songtianyi/wechat-go/wxweb"

	"github.com/mlogclub/mlog-wxbot/wxbot"
)

var configFlag = flag.String("config", "./config.yaml", "配置文件路径")

func init() {
	flag.Parse()
	config.InitConfig(*configFlag)

	// 连接数据库
	simple.OpenDB(&simple.DBConfiguration{
		Dialect:        "mysql",
		MaxIdle:        5,
		MaxActive:      20,
		Url:            config.Conf.MySqlUrl,
		EnableLogModel: config.Conf.ShowSql,
		Models:         []interface{}{&wxbot.WxArticle{}},
	})
}

func main() {
	session, err := wxweb.CreateSession(nil, nil, wxweb.TERMINAL_MODE)
	if err != nil {
		return
	}

	wxbot.Register(session)

	for {
		if err := session.LoginAndServe(false); err != nil {
			logrus.Error("session exit, %s", err)
			for i := 0; i < 3; i++ {
				logrus.Info("trying re-login with cache")
				if err := session.LoginAndServe(true); err != nil {
					logrus.Error("re-login error or session down, %s", err)
				}
				time.Sleep(3 * time.Second)
			}
			if session, err = wxweb.CreateSession(nil, session.HandlerRegister, wxweb.TERMINAL_MODE); err != nil {
				logrus.Error("create new session failed, %s", err)
				break
			}
		} else {
			logrus.Info("closed by user")
			break
		}
	}
}
