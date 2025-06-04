package bots

import (
	"log"

	"github.com/chippydip/go-sc2ai/botutil"
	"github.com/chippydip/go-sc2ai/client"
)

type bot struct {
	*botutil.Bot
}

func RunAgent(info client.AgentInfo) {
	bot := bot{Bot: botutil.NewBot(info)}
	bot.LogActionErrors()

	bot.init()
	for bot.IsInGame() {
		bot.doSmt()
		log.Print(bot.GameInfo().GetStartRaw())
		if err := bot.Step(1); err != nil {
			log.Print(err)
			break
		}
	}
}
func (bot *bot) init() {
}

func (bot *bot) doSmt() {}
