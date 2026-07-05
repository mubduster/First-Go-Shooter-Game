package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"

)

type screen struct {
	X float32
	Y float32
}
type world struct {
	X float32
	Y float32
}
type bean struct {
	Pos rl.Vector2
	Width float32
	Height	float32
	Radius	float32
	Speed rl.Vector2
	MaxSpeed float32
	Acceleration float32
	Drag float32
	Jump float32
	CanJump bool
}
type gravity struct {
	Max float32
	Bean float32
}
type Map struct {
	Border rl.Rectangle
}
type Platform struct {
	Rect rl.Rectangle
}
type mapColl struct {
	Floor bool
	onPlatform bool
}
var MapColl mapColl
var isGrounded bool
var hasJumped bool
var oldHeight float32
var Platforms []Platform
 
func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowMaximized)
	rl.InitWindow(0, 0, "Go Shooter Game")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	Screen := screen{X: float32(rl.GetScreenWidth()), Y: float32(rl.GetScreenHeight())}
	
	World := world{X: 6000, Y:4000}

	Map := Map{Border: rl.NewRectangle(0, 0, World.X, World.Y)}

	Platforms = []Platform{
		{Rect: rl.NewRectangle(0, 3680, 1000, 20)},
	}

	Bean := bean{Pos: rl.NewVector2(World.X/2, World.Y/2), Width: 40, Height: 100, Radius: 20, Speed: rl.NewVector2(0, 0), MaxSpeed: 30, Acceleration: 10, Drag: 5, Jump: 70}

	Gravity := gravity{Max: 50, Bean: 0}

	Camera := rl.Camera2D{Offset: rl.NewVector2(Screen.X/2, Screen.Y/2), Target: Bean.Pos, Rotation: 0.0, Zoom: 0.2}

	for !rl.WindowShouldClose(){
		
		// Map Collisions ----------------------------------------------------------------------------------------------------------------------------------------------------------
		MapColl.onPlatform = false
		for _, p := range Platforms{
			feetY := Bean.Pos.Y + Bean.Height
			if feetY >= p.Rect.Y && feetY < p.Rect.Y + p.Rect.Height && Bean.Pos.X < p.Rect.X + p.Rect.Width && Bean.Pos.X + Bean.Width > p.Rect.X {
				Bean.Pos.Y = p.Rect.Y - Bean.Height 
				Bean.Speed.Y = 0.0
				MapColl.onPlatform = true
				break
			}
		}

		if Bean.Pos.Y + Bean.Height > Map.Border.Height - 50 { // Map Collision for floor
			Bean.Pos.Y = Map.Border.Height - Bean.Height - 35
			MapColl.Floor = true
			}else if Bean.Pos.Y + Bean.Height < Map.Border.Height - 200 {
			MapColl.Floor = false
		}
		//if (Bean.Pos.Y + Bean.Height > Map.Line1.Y - 50) && ((Bean.Pos.X < Map.Line1End.X) && (Bean.Pos.X + Bean.Width > Map.Line1.X)) && !(Bean.Pos.Y > Map.Line1.Y) { // Map collison for Line 1
		//	Bean.Pos.Y = Map.Line1.Y - Bean.Height - 10
		//	MapColl.Line1 = true
		//	}else {
		//	MapColl.Line1 = false
		//}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		if MapColl.Floor || MapColl.onPlatform {
			isGrounded = true
		}else {
			isGrounded = false
		}
		if rl.IsKeyReleased(rl.KeyW) || rl.IsKeyReleased(rl.KeyUp) {
			hasJumped = false
		}
			
		oldHeight = Bean.Height
		if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
			Bean.Height = 50
			Bean.MaxSpeed = 20
		}else {
			Bean.Height = 100
			Bean.MaxSpeed = 30
		}

		Bean.Pos.Y += oldHeight - Bean.Height

		// Drag --------------------------------------------------------------------------------------------------------------------------------------------------------------------
		if Bean.Speed.Y < 0 {
			Bean.Speed.Y += Bean.Drag
		}else if !(rl.IsKeyDown(rl.KeyS) || rl. IsKeyDown(rl.KeyDown)) && Bean.Speed.Y > 0 {
			Bean.Speed.Y -= Bean.Drag
		}
		if !(rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft)) && Bean.Speed.X < 0 {
			Bean.Speed.X += Bean.Drag
		}else if !(rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight)) && Bean.Speed.X > 0 {
			Bean.Speed.X -= Bean.Drag
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		// Movement ----------------------------------------------------------------------------------------------------------------------------------------------------------------
		if (rl.IsKeyPressed(rl.KeyW) || rl.IsKeyPressed(rl.KeyUp)) && -Bean.Speed.Y < Bean.MaxSpeed  && (isGrounded && !hasJumped){
			Bean.Speed.Y -= Bean.Jump
			isGrounded = false
			hasJumped = true
		
		// }else if (rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown)) && Bean.Speed.Y < Bean.MaxSpeed && !isGrounded {
		// 	Bean.Speed.Y += Bean.Acceleration
		} 
		if (rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft)) && -Bean.Speed.X < Bean.MaxSpeed {
			Bean.Speed.X -= Bean.Acceleration
		}else if (rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight)) && Bean.Speed.X < Bean.MaxSpeed {
			Bean.Speed.X += Bean.Acceleration
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		// Gravity -----------------------------------------------------------------------------------------------------------------------------------------------------------------
		if MapColl.Floor || MapColl.onPlatform {
			Gravity.Bean = 0.0
		}else {
            if Gravity.Bean < Gravity.Max {
				Gravity.Bean += 1
			}
			Bean.Pos.Y += Gravity.Bean
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------


		Bean.Pos.X += Bean.Speed.X
		Bean.Pos.Y += Bean.Speed.Y

		Camera.Target = rl.NewVector2(World.X/2, World.Y/2)

		rl.BeginDrawing()

		rl.ClearBackground(rl.GetColor(0x0034a9ff))

		
		rl.BeginMode2D(Camera)
		rl.DrawRectangle(3000, 3000, 500, 500, rl.GetColor(0xff0000ff))
		
		rl.DrawRectangleV(Bean.Pos, rl.NewVector2(Bean.Width, Bean.Height), rl.GetColor(0x00ffffff))
		Bean.Pos.X += Bean.Width/2
		Bean.Pos.Y -= 19
		rl.DrawCircleV(Bean.Pos, Bean.Radius, rl.GetColor(0x00ffffff))
		Bean.Pos.X -= Bean.Width/2
		Bean.Pos.Y += 19
		
		rl.DrawRectangleLinesEx(Map.Border, 30, rl.GetColor(0x000000ff))
		
		for _, p := range Platforms {
			rl.DrawRectangleRec(p.Rect, rl.GetColor(0x000000ff))
		}
		
		rl.EndMode2D()

		rl.DrawText(fmt.Sprintf("SpeedX: %0.1f\nSpeedY: %0.1f\nGravity Bean: %0.1f",Bean.Speed.X, Bean.Speed.Y, Gravity.Bean), 10, 10, 30, rl.GetColor(0xffffffff))
		
		rl.EndDrawing()
	}
}