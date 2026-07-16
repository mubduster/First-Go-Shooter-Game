package main

import (
	Math "math"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var BulletCollisionRadius float32 = 9 // set to image size / 2
var BulletSpeed float32 = 1300
var BulletDamage float32 = 20
var BulletLifeSpan float32 = 2.5

func NewBullets (g gun) bullet {
	Radians := (g.Angle / 180) * Math.Pi

	return bullet{
		Angle: g.Angle,
		Pos: g.Barrel,
		PrevPos: g.Barrel,
		Radius: BulletCollisionRadius,
		Speed: rl.NewVector2( float32( Math.Cos(float64(Radians))) * BulletSpeed, float32(Math.Sin(float64(Radians))) * BulletSpeed),
		Damage: BulletDamage,
		Alive: true,
		Time: BulletLifeSpan,
	}
}