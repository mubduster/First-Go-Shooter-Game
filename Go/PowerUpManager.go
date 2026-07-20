package main

import (
	"math/rand"
	// rl "github.com/gen2brain/raylib-go/raylib"
)

// this will be the file that Manages the powerups on how they spawn nad work.

var PowerUps []PowerUp

func RandomPowerUpSpawner(SpawnPoints []spawnPoints, PowerUpsData []powerUpData) []PowerUp{
	
	PointIndex := rand.Intn(len(SpawnPoints))
	
	if !SpawnPoints[PointIndex].Occupied {
		SpawnPoints[PointIndex].Occupied = true
		PowerUps = append(PowerUps, PowerUp{Pos: SpawnPoints[PointIndex].Pos, SpawnIndex: PointIndex, Type: GetRandomPowerUp(PowerUpsData)})
	}

	return PowerUps
	
}

func GetRandomPowerUp(PowerUpsData []powerUpData) int {
	TotalWeight := 0

	for _,p := range PowerUpsData {
		TotalWeight += p.Weight
	}

	Roll := rand.Intn(TotalWeight)

	for _,p := range PowerUpsData {
		if Roll < p.Weight{
			return p.Type
		}
		Roll -= p.Weight
	}

	return PUHealth

}

func CheckPlayerPowerUpPickUp() {

}

func CheckForPowerUpDespawn() {

}
