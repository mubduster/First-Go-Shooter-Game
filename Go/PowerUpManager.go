package main

import (
	"math/rand"
	rl "github.com/gen2brain/raylib-go/raylib"

)

// this will be the file that Manages the powerups on how they spawn nad work.

var PowerUps []PowerUp
var Spawned bool
var i int

func RandomPowerUpSpawner(SpawnPoints []spawnPoints, PowerUpsData []powerUpData) []PowerUp { // Spawns a randomly selected Power Up at a randomly selected Spawn Point.
	Spawned = false

	for !Spawned {
		i++

		PointIndex := rand.Intn(len(SpawnPoints))

		if !SpawnPoints[PointIndex].Occupied {
			Spawned = true
			SpawnPoints[PointIndex].Occupied = true
			PowerUps = append(PowerUps, PowerUp{Pos: SpawnPoints[PointIndex].Pos, SpawnIndex: PointIndex, Type: GetRandomPowerUp(PowerUpsData), LifeSpan: 25})
		}

		if i >= len(SpawnPoints) { // Breaks loop if all Spawn Points are occupied so loop doesn become infinite.
			break
		}

	}

	return PowerUps

}

func GetRandomPowerUp(PowerUpsData []powerUpData) int { // Selects a random Power Up.
	TotalWeight := 0

	for _, p := range PowerUpsData {
		TotalWeight += p.Weight
	}

	Roll := rand.Intn(TotalWeight)

	for _, p := range PowerUpsData {
		if Roll < p.Weight {
			return p.Type
		}
		Roll -= p.Weight
	}

	return PUHealth

}

func CheckPlayerPowerUpPickUp(actor *bean) { // Checks if a Power Up is touched by a player and consumes it .
	actorRect := rl.NewRectangle(actor.Pos.X, actor.Pos.Y - (actor.Radius*2), actor.Width, actor.Height + (actor.Radius*2))
	for i := 0; i < len(PowerUps); {
		PowerUpRect := rl.NewRectangle(PowerUps[i].Pos.X, PowerUps[i].Pos.Y, 92, 92)

		if rl.CheckCollisionRecs(actorRect, PowerUpRect) {

		}
	}
}

func CheckForPowerUpDespawn(dT float32, SpawnPoints []spawnPoints) { // Desapwns Power Up after a certain time has passed set by RandomPowerUpSpawner() and a player hasn't consumed it.
	for i := 0; i < len(PowerUps); {
		PowerUps[i].LifeSpan -= dT

		if PowerUps[i].LifeSpan <= 0.0 {
			PowerUps[i] = PowerUps[len(PowerUps)-1]
			PowerUps = PowerUps[:len(PowerUps)-1]
			SpawnPoints[PowerUps[i].SpawnIndex].Occupied = false 
		}else {
			i++
		}
	}
}
