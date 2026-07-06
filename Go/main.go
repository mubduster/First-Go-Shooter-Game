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
	OneWay bool
}
type mapColl struct {
	Floor bool
	onPlatform bool
}
var MapColl mapColl
var isGrounded bool
var hasJumped bool
var isCrouched bool = false
var oldHeight float32
var Platforms []Platform
var colorOneWay rl.Color = rl.GetColor(0x444444ff)
var colorSolid rl.Color = rl.GetColor(0x000000ff)

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowMaximized)
	rl.InitWindow(0, 0, "Go Shooter Game")
	defer rl.CloseWindow()
	
	rl.SetTargetFPS(60)

	Screen := screen{X: float32(rl.GetScreenWidth()), Y: float32(rl.GetScreenHeight())}
	
	World := world{X: 6000, Y:4000}

	Map := Map{Border: rl.NewRectangle(0, 0, World.X, World.Y)}

	Platforms = []Platform{
		{Rect: rl.NewRectangle(30, 3650, 500, 20), OneWay: true},
		{Rect: rl.NewRectangle(501, 3650, 1500, 20), OneWay: false},
		{Rect: rl.NewRectangle(30, 3370, 2000, 20), OneWay: false},
	}
	
	Bean := bean{Pos: rl.NewVector2(World.X/2, World.Y/2), Width: 40, Height: 100, Radius: 20, Speed: rl.NewVector2(0, 0), MaxSpeed: 1000, Acceleration: 100, Drag: 100, Jump: 10500}
	
	Gravity := gravity{Max: 1000, Bean: 0}
	
	Camera := rl.Camera2D{Offset: rl.NewVector2(Screen.X/2, Screen.Y/2), Target: Bean.Pos, Rotation: 0.0, Zoom: 0.2}
	
	for !rl.WindowShouldClose(){

		dT := rl.GetFrameTime()
		
		oldHeight = Bean.Height
		if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
			Bean.Height = 50
			Bean.MaxSpeed = 500
			isCrouched = true
		}else {
			Bean.Height = 100
			Bean.MaxSpeed = 1000
			isCrouched = false
		}
	
		Bean.Pos.Y += oldHeight - Bean.Height
		
		
		// Drag --------------------------------------------------------------------------------------------------------------------------------------------------------------------
		if Bean.Speed.Y < 0 {
			Bean.Speed.Y += Bean.Drag * dT
		}
		if !(rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft)) && Bean.Speed.X < 0{
			Bean.Speed.X += Bean.Drag * dT
			if Bean.Speed.X > 0 {
				Bean.Speed.X = 0.0
			}
		}else if !(rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight)) && Bean.Speed.X > 0{
			Bean.Speed.X -= Bean.Drag * dT
			if Bean.Speed.X < 0 {
				Bean.Speed.X = 0.0
			}
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		// Movement ----------------------------------------------------------------------------------------------------------------------------------------------------------------
		if (rl.IsKeyPressed(rl.KeyW) || rl.IsKeyPressed(rl.KeyUp)) && -Bean.Speed.Y < Bean.MaxSpeed  && (isGrounded && !hasJumped){
			Bean.Speed.Y -= Bean.Jump * dT
			isGrounded = false
			hasJumped = true
		} 
		if rl.IsKeyDown(rl.KeyA) && rl.IsKeyDown(rl.KeyD) {
			Bean.Speed.X = 0.0
		}else if (rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight)) && Bean.Speed.X < Bean.MaxSpeed {
			Bean.Speed.X += Bean.Acceleration * dT
		}else if (rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft)) && -Bean.Speed.X < Bean.MaxSpeed {
			Bean.Speed.X -= Bean.Acceleration * dT
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		if MapColl.Floor || MapColl.onPlatform {
			isGrounded = true
		}else {
			isGrounded = false
		}
		if rl.IsKeyReleased(rl.KeyW) || rl.IsKeyReleased(rl.KeyUp) {
			hasJumped = false
		}
		
		
		// Gravity -----------------------------------------------------------------------------------------------------------------------------------------------------------------
		if MapColl.Floor || MapColl.onPlatform {
			Gravity.Bean = 0
		}else {
			if Gravity.Bean < Gravity.Max && Bean.Speed.Y < Gravity.Max{
				Gravity.Bean = 250
				Bean.Speed.Y += Gravity.Bean * dT
			}else{
				Bean.Speed.Y = Gravity.Bean * dT
			}
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
			
			
		Bean.Pos.X += Bean.Speed.X * dT
		Bean.Pos.Y += Bean.Speed.Y * dT

		// Map Collisions ----------------------------------------------------------------------------------------------------------------------------------------------------------
			MapColl.onPlatform = false
			for _, p := range Platforms{
				feetY := Bean.Pos.Y + Bean.Height
				headY := Bean.Pos.Y + (Bean.Radius*2)
			
				if p.OneWay && !isCrouched && feetY >= p.Rect.Y && feetY < p.Rect.Y + p.Rect.Height && Bean.Pos.X < p.Rect.X + p.Rect.Width && Bean.Pos.X + Bean.Width > p.Rect.X { // top collision for OneWay platform
					Bean.Pos.Y = p.Rect.Y - Bean.Height 
					Bean.Speed.Y = 0.0
					MapColl.onPlatform = true
					break
				}
			
				if !p.OneWay{
					if feetY >= p.Rect.Y && feetY < p.Rect.Y + p.Rect.Height && Bean.Pos.X < p.Rect.X + p.Rect.Width && Bean.Pos.X + Bean.Width > p.Rect.X { // top collision
						Bean.Pos.Y = p.Rect.Y - Bean.Height  
						Bean.Speed.Y = 0.0 
						MapColl.onPlatform = true 
						break 
					}
					if Bean.Speed.Y < 0  && headY <= p.Rect.Y  && headY > p.Rect.Y - p.Rect.Height - 10  &&  Bean.Pos.X < p.Rect.X + p.Rect.Width && Bean.Pos.X + Bean.Width > p.Rect.X { //bottom collision 
						Bean.Pos.Y = p.Rect.Y + (Bean.Radius * 2)
						Bean.Speed.Y = -1.0
						break	 
					}
					if rl.CheckCollisionPointRec(rl.NewVector2(p.Rect.X + p.Rect.Width + 1, p.Rect.Y + (p.Rect.Height/2)), rl.NewRectangle(Bean.Pos.X, Bean.Pos.Y, Bean.Width, Bean.Height)) {// right side collision 
						Bean.Pos.X = p.Rect.X + p.Rect.Width + 1  
						Bean.Speed.X = 0.0 
					}
					if rl.CheckCollisionPointRec(rl.NewVector2(p.Rect.X - 10, p.Rect.Y + (p.Rect.Height/2)), rl.NewRectangle(Bean.Pos.X, Bean.Pos.Y, Bean.Width, Bean.Height)){ //left side collision 
						Bean.Pos.X =  p.Rect.X - Bean.Width
						Bean.Speed.X = 0.0
					}
				}

			}

			if Bean.Pos.Y + Bean.Height > Map.Border.Height - 50 { // Map Collision for floor
				Bean.Pos.Y = Map.Border.Height - Bean.Height - 35
				Bean.Speed.Y = 0.0
				MapColl.Floor = true
				}else if Bean.Pos.Y + Bean.Height < Map.Border.Height - 200 {
				MapColl.Floor = false
			}
			if Bean.Pos.Y + (Bean.Radius * 2) < Map.Border.Y + 50 {
				Bean.Pos.Y = (Bean.Radius*2) + Bean.Height + 35
			}
			if Bean.Pos.X  < Map.Border.X + 50 {
				Bean.Speed.X = 0.0
				Bean.Pos.X = Map.Border.X + Bean.Width + 10
			}else if Bean.Pos.X + Bean.Width > Map.Border.Width - 50 {
				Bean.Speed.X = 0.0
				Bean.Pos.X = Map.Border.Width - Bean.Width - 40
			}

			//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------	



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
			if p.OneWay == true{
			rl.DrawRectangleRec(p.Rect, colorOneWay)
			}else {
				rl.DrawRectangleRec(p.Rect, colorSolid)
			}
		}
		
		rl.EndMode2D()

		rl.DrawText(fmt.Sprintf("SpeedX: %0.1f\nSpeedY: %0.1f\nGravity Bean: %0.1f",Bean.Speed.X, Bean.Speed.Y, Gravity.Bean), 10, 10, 30, rl.GetColor(0xffffffff))
		
		rl.EndDrawing()
	}
}