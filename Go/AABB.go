package main

import(
	"math"
	rl "github.com/gen2brain/raylib-go/raylib"
)

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

var maxArea float32

func resolveMapCollision(platforms []Platform, actor *bean) {
	actorAABB := rl.NewRectangle(actor.Pos.X, actor.Pos.Y, actor.Width, actor.Height)
	AABB := actorAABB
	
	for i := 0; i < 10; i++ {
		var mostOverlap rl.Rectangle
		maxArea = -1

		for _, p := range platforms {
			if p.OneWay && ( actor.Speed.Y < 0 || actor.isCrouched){
				continue
			}

			overlap := getSignedCollisionRec(p.Rect, AABB)
			area := float32(math.Abs(float64(overlap.Width * overlap.Height)))
			if area > maxArea{
				maxArea = area
				mostOverlap = overlap
			}

		}
		if maxArea <= 0 {
			break
		}
		if float32(math.Abs(float64(mostOverlap.Width))) < float32(math.Abs(float64(mostOverlap.Height))) {
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
			actor.isGrounded = false
		}
		actor.Pos.X = AABB.X
		actor.Pos.Y = AABB.Y
	}
}