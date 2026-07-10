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
	MaxSpeedY float32
	Acceleration float32
	Drag float32
	Jump float32
	CanJump bool
	isGrounded bool
	hasJumped bool
	isCrouched bool
	ignoredPlatformIndex int
	IgnoredCooldown float32
	CurrentPlatformIndex int
	restingOnPlatform bool
}
type gravity struct {
	Max float32
	Bean float32
	Force float32
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

// var isGrounded bool
// var hasJumped bool
// var isCrouched bool = false
var oldHeight float32

var Platforms []Platform
var colorOneWay rl.Color = rl.GetColor(0x444444ff)
var colorSolid rl.Color = rl.GetColor(0x000000ff)

var earliestTime float32 = 1.0
var collisionNormal rl.Vector2 
var hitAny bool
// var ignoredPlatformIndex int = -1
// var ignoreCooldown float32 
var hitPlatformIndex int = -1
var remainingTime float32 
var restingOnPlatform bool
const epsilon float32 = 2.5

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowMaximized)
	rl.InitWindow(0, 0, "Go Shooter Game")
	defer rl.CloseWindow()
	
	rl.SetTargetFPS(60)

	Screen := screen{X: float32(rl.GetScreenWidth()), Y: float32(rl.GetScreenHeight())}
	
	World := world{X: 6000, Y:4000}

	Map := Map{Border: rl.NewRectangle(0, 0, World.X, World.Y)}

	Platforms = []Platform{
		{Rect: rl.NewRectangle(30, 3650, 500, 56), OneWay: true},
		{Rect: rl.NewRectangle(501, 3650, 1500, 56), OneWay: false},
		{Rect: rl.NewRectangle(30, 3370, 1200, 56), OneWay: false},
		{Rect: rl.NewRectangle(1230, 3370, 350, 56), OneWay: true},
		{Rect: rl.NewRectangle(1581, 3370, 350, 56), OneWay: false},
		{Rect: rl.NewRectangle(4000, 3650, 600, 56), OneWay: false},
		{Rect: rl.NewRectangle(4601, 3650, 500, 56), OneWay: true},
		{Rect: rl.NewRectangle(5102, 3650, 898, 56), OneWay: false},
		{Rect: rl.NewRectangle(5470, 3370, 500, 56), OneWay: true},
		{Rect: rl.NewRectangle(3600, 3370, 1870, 56), OneWay: false},	
		{Rect: rl.NewRectangle(2600, 3700, 600, 180), OneWay: false}, // rectangle
		{Rect: rl.NewRectangle(1400, 3070, 1200, 56), OneWay:  false},
		{Rect: rl.NewRectangle(30, 3070, 800, 56), OneWay:  false},
		{Rect: rl.NewRectangle(831, 3070, 570, 56), OneWay: true},
		{Rect: rl.NewRectangle(2600, 3070, 600, 56), OneWay: true}, // middle platform
		{Rect: rl.NewRectangle(3201, 3070, 1500, 56), OneWay: false},
		{Rect: rl.NewRectangle(4702, 3070, 600, 56), OneWay: true},
		{Rect: rl.NewRectangle(5302, 3070, 698, 56), OneWay: false},

		{Rect: rl.NewRectangle(530, 2790, 4940, 56), OneWay: false},
		{Rect: rl.NewRectangle(30, 2790, 500, 56), OneWay: true},
		{Rect: rl.NewRectangle(5470, 2790, 500, 56), OneWay: true},

		{Rect: rl.NewRectangle(30, 2510, 1000, 56), OneWay: false},
		{Rect: rl.NewRectangle(4970, 2510, 1000, 56), OneWay: false},
		{Rect: rl.NewRectangle(1800, 2510, 600, 56), OneWay: false},
		{Rect: rl.NewRectangle(3550, 2510, 600, 56), OneWay: false},

		{Rect: rl.NewRectangle(2950, 2235, 95, 300), OneWay: false},
		{Rect: rl.NewRectangle(1800, 2235, 100, 240), OneWay: false},
		{Rect: rl.NewRectangle(4050, 2235, 100, 240), OneWay: false},

		// {Rect: rl.NewRectangle(30, 2230, 500, 56), OneWay: true},
		{Rect: rl.NewRectangle(530, 2230, 4940, 56), OneWay: false},
		// {Rect: rl.NewRectangle(5471, 2230, 500, 56), OneWay: true},

		{Rect: rl.NewRectangle(30, 1950, 900, 56), OneWay: false},
		{Rect: rl.NewRectangle(5070, 1950, 900, 56), OneWay: false},
		{Rect: rl.NewRectangle(1400, 1950, 900, 56), OneWay: false},
		{Rect: rl.NewRectangle(3700, 1950, 900, 56), OneWay: false},
		{Rect: rl.NewRectangle(2700, 1950, 600, 56), OneWay: false},
		{Rect: rl.NewRectangle(2300, 1950, 400, 56), OneWay: true},
		{Rect: rl.NewRectangle(3300, 1950, 400, 56), OneWay: true},

		{Rect: rl.NewRectangle(530, 1670, 1100, 56), OneWay: false},
		{Rect: rl.NewRectangle(1630, 1600, 50, 125), OneWay: false},
		{Rect: rl.NewRectangle(4370, 1670, 1100, 56), OneWay: false},
		{Rect: rl.NewRectangle(4320, 1600, 50, 125), OneWay: false},

		{Rect: rl.NewRectangle(1500, 1390, 1200, 56), OneWay: false},
		{Rect: rl.NewRectangle(3300, 1390, 1200, 56), OneWay: false},
		{Rect: rl.NewRectangle(2700, 1390, 600, 56), OneWay: true},

		{Rect: rl.NewRectangle(30, 1390, 700, 56), OneWay:	false},
		{Rect: rl.NewRectangle(5300, 1390, 700, 56), OneWay: false},

		{Rect: rl.NewRectangle(500, 1110, 600, 56), OneWay: true},
		{Rect: rl.NewRectangle(1101, 1110, 3800, 56), OneWay: false},
		{Rect: rl.NewRectangle(4900, 1110,600, 56), OneWay: true},
	}
	
	Bean := bean{Pos: rl.NewVector2(World.X/2, World.Y/2), Width: 40, Height: 100, Radius: 20, Speed: rl.NewVector2(0, 0), MaxSpeed: 2000, Acceleration: 500, Drag: 460, Jump: 3000, CurrentPlatformIndex: -1, ignoredPlatformIndex: -1, restingOnPlatform: false}
	
	Gravity := gravity{Max: 1500, Bean: 0, Force: 60}
	
	Camera := rl.Camera2D{Offset: rl.NewVector2(Screen.X/2, Screen.Y/2), Target: Bean.Pos, Rotation: 0.0, Zoom: 0.2}
	
	for !rl.WindowShouldClose(){

		dT := rl.GetFrameTime()  // Delta Time for allowing the game to run at any FPS and framerate

		// OneWay Collision Helper -------------------------------------------------------------------------------------------------------------------------------------------------
		if Bean.IgnoredCooldown > 0 {  // OneWay playform collision ignore timer countdown
			Bean.IgnoredCooldown -= dT
			if Bean.IgnoredCooldown <= 0 {
				Bean.ignoredPlatformIndex = -1  // resets OneWay collision
			}
		}

		if Bean.ignoredPlatformIndex != -1 {
			p := Platforms[Bean.ignoredPlatformIndex]
			if Bean.Pos.Y > p.Rect.Y + p.Rect.Height && Bean.IgnoredCooldown <= 0 {
				Bean.ignoredPlatformIndex = -1
			}
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

		// Crouching Handler -------------------------------------------------------------------------------------------------------------------------------------------------------
		oldHeight = Bean.Height
		if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
			Bean.isCrouched = true
		}else {
			Bean.isCrouched = false
		}
		if Bean.isCrouched {
			Bean.Height = 50
			Bean.MaxSpeed = 500
			Bean.Jump = 1500
			Bean.isCrouched = true
		}else {
			Bean.Height = 100
			Bean.MaxSpeed = 1000
			Bean.Jump = 3000
			Bean.isCrouched = false
		}
	
		Bean.Pos.Y += oldHeight - Bean.Height
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		// Drag --------------------------------------------------------------------------------------------------------------------------------------------------------------------
		if !(rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft)) && Bean.Speed.X < 0{
			Bean.Speed.X += Bean.Drag 
			if Bean.Speed.X > 0 {
				Bean.Speed.X = 0.0
			}
		}else if !(rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight)) && Bean.Speed.X > 0{
			Bean.Speed.X -= Bean.Drag 
			if Bean.Speed.X < 0 {
				Bean.Speed.X = 0.0
			}
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		// Movement ----------------------------------------------------------------------------------------------------------------------------------------------------------------
		if (rl.IsKeyPressed(rl.KeyW) || rl.IsKeyPressed(rl.KeyUp)) && -Bean.Speed.Y < Bean.MaxSpeed  && (Bean.isGrounded && !Bean.hasJumped){
			Bean.Speed.Y -= Bean.Jump
			Bean.isGrounded = false
			Bean.hasJumped = true
		} 
		if rl.IsKeyDown(rl.KeyA) && rl.IsKeyDown(rl.KeyD) {
			Bean.Speed.X = 0.0
		}else if (rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight)) && Bean.Speed.X < Bean.MaxSpeed {
			Bean.Speed.X += Bean.Acceleration 
		}else if (rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft)) && -Bean.Speed.X < Bean.MaxSpeed {
			Bean.Speed.X -= Bean.Acceleration 
		}

		if Bean.Speed.X > Bean.MaxSpeed  && rl.IsKeyDown(rl.KeyD) {
			Bean.Speed.X = Bean.MaxSpeed
		}
		if Bean.Speed.X < -Bean.MaxSpeed && rl.IsKeyDown(rl.KeyA) {
			Bean.Speed.X = -Bean.MaxSpeed
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

		
		// Map Border Collisions ---------------------------------------------------------------------------------------------------------------------------------------------------
		if Bean.Speed.Y >= 0 && Bean.Pos.Y + Bean.Height > Map.Border.Height - 50 { // Map Collision for floor
			Bean.Pos.Y = Map.Border.Height - Bean.Height - 35
			Bean.Speed.Y = 0.0
			MapColl.Floor = true
		}else if Bean.Pos.Y + Bean.Height < Map.Border.Height + 200 {
			MapColl.Floor = false
		}
		if Bean.Pos.Y + (Bean.Radius * 2) < Map.Border.Y + 50 {  // Map Collision for Roof
			Bean.Pos.Y = (Bean.Radius*2) + Bean.Height + 35
		}
		if Bean.Speed.X <= 0 && Bean.Pos.X  < Map.Border.X + 60 {  // Map Collision for Left Wall
			Bean.Speed.X = 0.0
			Bean.Pos.X = Map.Border.X + Bean.Width + 10
		}else if Bean.Speed.X >= 0 && Bean.Pos.X + Bean.Width > Map.Border.Width - 60 {  // Map Collision for Right Wall
			Bean.Speed.X = 0.0
			Bean.Pos.X = Map.Border.Width - Bean.Width - 40
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------	
		
		// Add Speed to Position of the Bean aka movement --------------------------------------------------------------------------------------------------------------------------
		Bean.Pos.X += Bean.Speed.X * dT
		Bean.Pos.Y += Bean.Speed.Y * dT
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		// Gravity -----------------------------------------------------------------------------------------------------------------------------------------------------------------
		if MapColl.Floor || Bean.isGrounded {
			Gravity.Bean = 0
		}else {
			if Gravity.Bean < Gravity.Max && Bean.Speed.Y < Gravity.Max{
				Gravity.Bean += Gravity.Force
				Bean.Speed.Y += Gravity.Bean 
			}else{
				Bean.Speed.Y = Gravity.Bean 
			}
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

		// Platform Collisions -----------------------------------------------------------------------------------------------------------------------------------------------------
		if Bean.isCrouched && Bean.restingOnPlatform && Bean.CurrentPlatformIndex != -1 && Platforms[Bean.CurrentPlatformIndex].OneWay {
			Bean.ignoredPlatformIndex = Bean.CurrentPlatformIndex
			Bean.IgnoredCooldown = 0.25
			Bean.isGrounded	 = false
			Bean.restingOnPlatform = false
		}

		Bean.CurrentPlatformIndex = -1
		Bean.restingOnPlatform = false
		
		for Index,p := range Platforms {
			if Index == Bean.ignoredPlatformIndex {
				continue
			}

			beanBottom := Bean.Pos.Y + Bean.Height

			if beanBottom >= p.Rect.Y - epsilon && beanBottom <= p.Rect.Y + epsilon && Bean.Pos.X + Bean.Width > p.Rect.X && Bean.Pos.X < p.Rect.X + p.Rect.Width {
				Bean.restingOnPlatform = true
				Bean.CurrentPlatformIndex = Index
				break
			}
		}

		resolveMapCollision(Platforms, &Bean)  // AABB Collision function
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

		// Collision Check for Gravity and Jump Check ------------------------------------------------------------------------------------------------------------------------------
		if MapColl.Floor || Bean.restingOnPlatform{
			Bean.isGrounded = true
		}else {
			Bean.isGrounded = false
		}
		if rl.IsKeyReleased(rl.KeyW) || rl.IsKeyReleased(rl.KeyUp) {
			Bean.hasJumped = false
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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

		rl.DrawText(fmt.Sprintf("SpeedX: %0.1f\nSpeedY: %0.1f\nGravity Bean: %0.1f\nGrounded: %v\n Crouched: %v",Bean.Speed.X, Bean.Speed.Y, Gravity.Bean, Bean.isGrounded, Bean.isCrouched), 10, 10, 30, rl.GetColor(0xffffffff))
		
		rl.EndDrawing()
	}
}