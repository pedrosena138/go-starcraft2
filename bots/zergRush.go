package bots

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/zerg"
	"github.com/chippydip/go-sc2ai/search"
)

type bot struct {
	*botutil.Bot

	myStartLocation    api.Point2D
	myNaturalLocation  api.Point2D
	enemyStartLocation api.Point2D

	camera api.Point2D
}

func ZergRush(info client.AgentInfo) {
	bot := bot{Bot: botutil.NewBot(info)}
	bot.LogActionErrors()

	bot.init()
	for bot.IsInGame() {
		bot.strategy()
		bot.tactics()
		if err := bot.Step(1); err != nil {
			log.Print(err)
			break
		}
	}
}

func (bot *bot) init() {
	bot.myStartLocation = bot.Self[zerg.Hatchery].First().Pos2D()
	bot.enemyStartLocation = *bot.GameInfo().GetStartRaw().GetStartLocations()[0]
	bot.camera = bot.myStartLocation

	expansions := search.CalculateBaseLocations(bot.Bot, true)
	query := make([]*api.RequestQueryPathing, len(expansions))

	for i, exp := range expansions {
		pos := exp.Location
		query[i] = &api.RequestQueryPathing{
			Start: &api.RequestQueryPathing_StartPos{
				StartPos: &bot.myStartLocation,
			},
			EndPos: &pos,
		}
	}

	response := bot.Query(api.RequestQuery{Pathing: query})
	best, minDist := -1, float32(256)
	for i, result := range response.GetPathing() {
		if result.Distance < minDist && result.Distance > 5 {
			best, minDist = i, result.Distance
		}
	}

	bot.myNaturalLocation = expansions[best].Location
}

func (bot *bot) strategy() {
	hatchesCount := bot.Self.Count(zerg.Hatchery)
	pool := bot.Self[zerg.SpawningPool].First()

	maxDrones := 14
	SPARE_SUPPLIES := hatchesCount * 3

	if pool.IsNil() {
		pos := bot.myStartLocation.Offset(bot.enemyStartLocation, 5)
		if !bot.BuildUnitAt(zerg.Drone, ability.Build_SpawningPool, pos) {
			return
		}
	}

	if bot.FoodLeft() <= SPARE_SUPPLIES && bot.Self.CountInProduction(zerg.Overlord) == 0 {
		if !bot.BuildUnit(zerg.Larva, ability.Train_Overlord) {
			return
		}
	}

	if hatchesCount > 1 {
		maxDrones = 16
	}

	droneCount := bot.Self.CountAll(zerg.Drone)
	bot.BuildUnits(zerg.Larva, ability.Train_Drone, maxDrones-droneCount)

	if pool.IsNil() || pool.BuildProgress < 1 {
		return
	}

	bot.BuildUnits(zerg.Larva, ability.Train_Zergling, 100)
	bot.BuildUnits(zerg.Hatchery, ability.Train_Queen, hatchesCount-bot.Self.CountAll(zerg.Queen))

	if hatchesCount < 2 {
		bot.BuildUnitAt(zerg.Drone, ability.Build_Hatchery, bot.myNaturalLocation)
	}
}

func (bot *bot) tactics() {}
