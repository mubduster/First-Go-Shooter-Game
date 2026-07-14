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
func resolveMapBulletCollison (platforms []Platform, Bullet *bullet) {
// 	bulletRect := rl.NewRectangle(Bullet.Pos.X - Bullet.Radius, Bullet.Pos.Y - Bullet.Radius, Bullet.Radius * 2, Bullet.Radius * 2) 

// 	for _,b := range platforms {
// 		overlap := getSignedCollisionRec(b.Rect, bulletRect)
		
// 		if Math.Abs(float64(overlap.Width)) < Math.Abs(float64(overlap.Height)) {
// 			Bullet.Speed.X *= -1
// 			// Bullet.Pos.X = overlap.Width
// 		}else {
// 			Bullet.Speed.Y *= -1
// 			// Bullet.Pos.Y += overlap.Height
// 		}
// 	} 
	for _,p := range platforms {
		if rl.CheckCollisionCircleRec(Bullet.Pos, Bullet.Radius, p.Rect) {
			ClosestXPoint := rl.Clamp(Bullet.Pos.X, p.Rect.X, p.Rect.X + p.Rect.Width)
			ClosestYPoint := rl.Clamp(Bullet.Pos.Y, p.Rect.Y, p.Rect.Y + p.Rect.Height)

			Dx := Bullet.Pos.X - ClosestXPoint
			Dy := Bullet.Pos.Y - ClosestYPoint

			if Math.Abs(float64(Dx)) > Math.Abs(float64(Dy)) {
				Bullet.Speed.X *= -1 

				if Dx > 0 {
					Bullet.Pos.X = ClosestXPoint + Bullet.Radius
				}else {
					Bullet.Pos.X = ClosestXPoint - Bullet.Radius
				}
			}else {
				Bullet.Speed.Y *= -1

				if Dy > 0 {
					Bullet.Pos.Y = ClosestYPoint + Bullet.Radius
				}else {
					Bullet.Pos.Y = ClosestYPoint - Bullet.Radius
				}
			}
			break	

			// if (Bullet.PrevPos.X - Bullet.Radius >= p.Rect.X + p.Rect.Width && Bullet.Pos.X - Bullet.Radius <= p.Rect.X + p.Rect.Width) || (Bullet.PrevPos.X + Bullet.Radius <= p.Rect.X && Bullet.Pos.X >= p.Rect.X) {
			// 	Bullet.Speed.X *= -1
			// }
			// if (Bullet.PrevPos.Y - Bullet.Radius >= p.Rect.Y + p.Rect.Height && Bullet.Pos.Y - Bullet.Radius <= p.Rect.Y + p.Rect.Height) || (Bullet.PrevPos.Y + Bullet.Radius <= p.Rect.Y && Bullet.Pos.Y + Bullet.Radius > p.Rect.Y) {
			// 	Bullet.Speed.Y *= -1
			// }
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------