package main

import(
	"math"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func Swetp_AABB(Moving, Static rl.Rectangle, Vel rl.Vector2) (collisionTime float32, normal rl.Vector2) {
	var xInvEntry, yInvEntry, xInvExit, yInvExit, xEntry, yEntry, xExit, yExit float32

	if Vel.X > 0 {
		xInvEntry = Static.X - (Moving.X + Moving.Width)
		xInvExit = (Static.X + Static.Width) - Moving.X
	}else {
		xInvEntry = (Static.X + Static.Width) - Moving.X
		xInvExit = Static.X - (Moving.X + Moving.Width)
	}
	if Vel.Y > 0 {
		yInvEntry = Static.Y - (Moving.Y + Moving.Height)
		yInvExit = (Static.Y + Static.Height) - Moving.Y
	}else {
		yInvEntry = (Static.Y + Static.Width) - Moving.Y
		yInvExit = Static.Y - (Moving.Y + Moving.Height)
	}

	if Vel.X == 0 {
		xEntry = float32(math.Inf(-1))
		xExit = float32(math.Inf(1))	
	}else {
		xEntry	= xInvEntry / Vel.X
		xExit = xInvExit / Vel.X
	}
	if Vel.Y == 0 {
		yEntry = float32(math.Inf(-1))
		yExit = float32(math.Inf(1))
	}else {
		yEntry = yInvEntry / Vel.Y
		yExit = yInvExit / Vel.Y
	}

	EntryTime := float32(math.Max(float64(xEntry), float64(yEntry)))
	ExitTime := float32(math.Min(float64(xExit), float64(yExit)))

	if EntryTime > ExitTime || xEntry < 0.0 && yEntry < 0.0 || xEntry > 1.0 || yEntry > 1.0 {
		return 1.0, rl.NewVector2(0,0)
	}

	if xEntry > yEntry {
		if xInvEntry < 0 {
			normal = rl.NewVector2(1,0)
		}else {
			normal = rl.NewVector2(-1,0)
		}
	}else {
		if yInvEntry < 0 {
			normal = rl.NewVector2(0,1)
		}else {
			normal = rl.NewVector2(0,-1)
		}
	}

	return EntryTime, normal
}