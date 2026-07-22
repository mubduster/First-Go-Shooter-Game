package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// this will be the file that Manages the powerups on how they spawn nad work.

// var PowerUps []PowerUp
var Spawned bool
var i int

func RandomPowerUpSpawner(SpawnedPowerUps []PowerUp, SpawnPoints []spawnPoints, PowerUpsData []powerUpData) []PowerUp { // Spawns a randomly selected Power Up at a randomly selected Spawn Point.
	Spawned = false

	for !Spawned {
		i++

		PointIndex := rand.Intn(len(SpawnPoints))

		if !SpawnPoints[PointIndex].Occupied {
			Spawned = true
			SpawnPoints[PointIndex].Occupied = true
			SpawnedPowerUps = append(SpawnedPowerUps, PowerUp{Pos: SpawnPoints[PointIndex].Pos, SpawnIndex: PointIndex, Type: GetRandomPowerUp(PowerUpsData), LifeSpan: 25})
		}

		if i >= len(SpawnPoints) { // Breaks loop if all Spawn Points are occupied so loop doesn become infinite.
			break
		}

	}

	return SpawnedPowerUps

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

// Remove power up, update player.
var power int = -1 // -1 mean no powers

func CheckPlayerPowerUpPickUp(SpawnedPowerUps []PowerUp, actor *bean, SpawnPoints []spawnPoints) ([]PowerUp, bool, int) { // Checks if a Power Up is touched by a player and consumes it .
	actorRect := rl.NewRectangle(actor.Pos.X, actor.Pos.Y-(actor.Radius*2), actor.Width, actor.Height+(actor.Radius*2))
	for i := 0; i < len(SpawnedPowerUps); {
		PowerUpRect := rl.NewRectangle(SpawnedPowerUps[i].Pos.X, SpawnedPowerUps[i].Pos.Y, 92, 92)

		if rl.CheckCollisionRecs(actorRect, PowerUpRect) {
			SpawnPoints[SpawnedPowerUps[i].SpawnIndex].Occupied = false
			power = SpawnedPowerUps[i].Type
			actor.PowersNumber++

			SpawnedPowerUps[i] = SpawnedPowerUps[len(SpawnedPowerUps)-1]
			SpawnedPowerUps = SpawnedPowerUps[:len(SpawnedPowerUps)-1]

			return SpawnedPowerUps, true, power

		} else {
			i++
		}
	}

	return SpawnedPowerUps, false, -1 // -1 means no Power Ups
}

func CheckForPowerUpDespawn(SpawnedPowerUps []PowerUp, dT float32, SpawnPoints []spawnPoints) []PowerUp { // Desapwns Power Up after a certain time has passed set by RandomPowerUpSpawner() and a player hasn't consumed it.
	for i := 0; i < len(SpawnedPowerUps); {
		SpawnedPowerUps[i].LifeSpan -= dT

		if SpawnedPowerUps[i].LifeSpan <= 0.0 {
			SpawnPoints[SpawnedPowerUps[i].SpawnIndex].Occupied = false
			SpawnedPowerUps[i] = SpawnedPowerUps[len(SpawnedPowerUps)-1]
			SpawnedPowerUps = SpawnedPowerUps[:len(SpawnedPowerUps)-1]
		} else {
			i++
		}
	}
	return SpawnedPowerUps
}
