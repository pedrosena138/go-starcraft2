package main

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/runner"
	zerg "github.com/pedrosena138/go-starcraft2/bots/zerg"
)

func main() {
	agent := client.AgentFunc(zerg.RunAgent)
	runner.SetComputer(
		api.Race_Random,
		api.Difficulty_Easy,
		api.AIBuild_RandomBuild,
	)
	runner.RunAgent(client.NewParticipant(
		api.Race_Zerg,
		agent,
		"Zerg Test",
	))
}
