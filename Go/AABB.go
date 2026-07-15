package main

import(
	Math "math"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Collision Helper ----------------------------------------------------------------------------------------------------------
func getSignedCollisionRec(rect1, rect2 rl.Rectangle) rl.Rectangle {
	r := rl.GetCollisionRec(rect1, rect2)
	if rect2.X < rect1.X {
		r.Width = -r.Width
	}
	if rect2.Y < rect1.Y {
		r.Height = -r.Height
	}
	return r
}
//----------------------------------------------------------------------------------------------------------------------------

// AABB Collision Handler ----------------------------------------------------------------------------------------------------
var maxArea float32

func resolveMapCollision(platforms []Platform, actor *bean) {
	actorAABB := rl.NewRectangle(actor.Pos.X, actor.Pos.Y, actor.Width, actor.Height)
	AABB := actorAABB
	playerBottom := actor.Pos.Y + actor.Height
	
	for i := 0; i < 10; i++ {
		var mostOverlap rl.Rectangle
		maxArea = -1

		for Index, p := range platforms {
			if p.OneWay && playerBottom > p.Rect.Y && actor.Speed.Y < 0 {
				continue
			}
			if Index == actor.ignoredPlatformIndex {
				continue
			}

			overlap := getSignedCollisionRec(p.Rect, AABB)
			area := float32(Math.Abs(float64(overlap.Width * overlap.Height)))
			if area > maxArea{
				maxArea = area
				mostOverlap = overlap
			}

		}
		if maxArea <= 0 {
			break
		}
		if float32(Math.Abs(float64(mostOverlap.Width))) < float32(Math.Abs(float64(mostOverlap.Height))) {
			AABB.X += mostOverlap.Width
		}else {
			AABB.Y += mostOverlap.Height
		}
	}

	if AABB.X != actorAABB.X || AABB.Y != actorAABB.Y {
		if AABB.Y < actorAABB.Y {
			actor.Speed.Y = 0
			actor.isGrounded = true
			actor.restingOnPlatform = true
		}else if AABB.Y > actorAABB.Y {
			actor.Speed.Y = 0
		}
		if AABB.X != actorAABB.X {
			actor.Speed.X = 0
		}
		actor.Pos.X = AABB.X
		actor.Pos.Y = AABB.Y
	}
}
//-------------------------------------------------------------------------------------------------------------------------------

// Map Bullet Collision Handler -------------------------------------------------------------------------------------------------
func resolveMapBulletCollison (Platforms []Platform, Bullet *bullet) {
	bulletRect := rl.NewRectangle(Bullet.Pos.X - Bullet.Radius, Bullet.Pos.Y - Bullet.Radius, Bullet.Radius * 2, Bullet.Radius * 2) 

	var mostOverlap rl.Rectangle
	maxArea = -1

	for _,p := range Platforms {
		if p.OneWay {
			continue
		}

		overlap := getSignedCollisionRec(p.Rect, bulletRect)
		area := float32(Math.Abs(float64(overlap.Width * overlap.Height)))

		if area > maxArea {
			maxArea = area
			mostOverlap = overlap
		}
		
	}

	if maxArea <= 0 {
		return
	}

	if float32(Math.Abs(float64(mostOverlap.Width))) < float32(Math.Abs(float64(mostOverlap.Height))) {
		bulletRect.X += mostOverlap.Width
		Bullet.Speed.X *= -1
	}else {
		bulletRect.Y += mostOverlap.Height
		Bullet.Speed.Y *= -1
	}

	Bullet.Pos.X = bulletRect.X + Bullet.Radius
	Bullet.Pos.Y = bulletRect.Y + Bullet.Radius
}

//-------------------------------------------------------------------------------------------------------------------------------

// Stop Gun from penetrating walls ----------------------------------------------------------------------------------------------
func CheckGunWallCollision (Gun *gun, Platforms []Platform) {
	
}
//-------------------------------------------------------------------------------------------------------------------------------