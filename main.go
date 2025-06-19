package main

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/runner"
	"github.com/pedrosena138/go-starcraft2/bots"
)

func main() {
	agent := client.AgentFunc(bots.ZergRush)
	runner.SetComputer(
		api.Race_Random,
		api.Difficulty_Easy,
		api.AIBuild_RandomBuild,
	)
	runner.RunAgent(client.NewParticipant(
		api.Race_Zerg,
		agent,
		"Zerg Rush",
	))
}
