package modules

import (
	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
)

func onReady(s *discordgo.Session, e *discordgo.Ready) {
	out.Infoln("websocket started")

	out.Infoln("authorized as:", session.State.User.String())
	out.Debugln("token:", s.Token)

	loadModules()

	if config.Bot.LogChannel == nil {
		return
	}

	data := map[string]interface{}{
		"Name": e.User,
	}

	SendLog("started.xml", data)
}
