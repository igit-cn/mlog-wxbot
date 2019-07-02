package wxbot

import (
	"github.com/sirupsen/logrus"
	"github.com/songtianyi/wechat-go/wxweb"
)

// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *wxweb.Session) {
	// 将插件注册到session
	// 第一个参数: 指定消息类型, 所有该类型的消息都会被转发到此插件
	// 第二个参数: 指定消息处理函数, 消息会进入此函数
	// 第三个参数: 自定义插件名，不能重名，switcher插件会用到此名称
	if err := session.HandlerRegister.Add(wxweb.MSG_LINK, wxweb.Handler(collector), "collector"); err != nil {
		logrus.Error(err)
	}
	if err := session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(print), "printText"); err != nil {
		logrus.Error(err)
	}

	// 开启插件
	if err := session.HandlerRegister.EnableByName("collector"); err != nil {
		logrus.Error(err)
	}
	if err := session.HandlerRegister.EnableByName("printText"); err != nil {
		logrus.Error(err)
	}
}

// 消息处理函数
func collector(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	go func() {
		wxArticle := collect(msg.Url)
		if wxArticle == nil {
			return
		}
		wxArticle = save(wxArticle)
		publish(wxArticle)
	}()
}

func print(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	logrus.Info("收到消息，发送人：" + msg.FromUserName + ", 内容：" + msg.Content)
}
