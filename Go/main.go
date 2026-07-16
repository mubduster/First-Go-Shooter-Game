package main

import (
	"fmt"
	Math "math"
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
	OldHeight float32
	Radius	float32
	Speed rl.Vector2
	MaxSpeed float32
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
	Health float32
	PowersNumber int
	PowerDuration float32
}
type gun struct {
	Dir int
	PrevDir int
	Pos rl.Vector2
	Width float32
	Height float32
	Angle float32
	Mag int
	Shots int
	CanShoot bool
	Delay float32
	Barrel rl.Vector2
}
type bullet struct{
	Angle float32
	Pos rl.Vector2
	PrevPos rl.Vector2
	Radius float32
	Speed rl.Vector2
	Damage float32
	Alive bool
	Time float32
}
type gravity struct {
	Max float32
	Bean float32
	Bean2 float32
	Force float32
}
type Map struct {
	Border rl.Rectangle
}
type Platform struct {
	Rect rl.Rectangle
	OneWay bool
}
type onlyCrouch struct {
	Rect rl.Rectangle
}
type mapColl struct {
	Floor bool
	Floor2 bool
}
type power struct {
	NFirerate float32
	FireRate float32
	NDamage float32
	Damage float32
	NReload float32
	Reload	float32
	NMag float32
	Mag float32
	NHealth float32
	Health float32
	Imune bool
}
type powerBox struct {
	Power string
	Pos rl.Vector2
	Width float32
	Height float32
	Alive bool
	Time float32
}

var MapColl mapColl

var Platforms []Platform
var colorOneWay rl.Color = rl.GetColor(0x444444ff)
var colorSolid rl.Color = rl.GetColor(0x000000ff) 
var width float32

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
		{Rect: rl.NewRectangle(2600, 3700, 600, 200), OneWay: false}, // rectangle
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

		{Rect: rl.NewRectangle(530, 2230, 4940, 56), OneWay: false},

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

		{Rect: rl.NewRectangle(2075, 900, 500, 150), OneWay: false},
		{Rect: rl.NewRectangle(3410, 900, 500, 150), OneWay: false},

		{Rect: rl.NewRectangle(30, 830, 400, 56), OneWay: false},
		{Rect: rl.NewRectangle(5570, 830, 400, 56), OneWay: false},

		{Rect: rl.NewRectangle(530, 570, 1000, 56), OneWay: false},
		{Rect: rl.NewRectangle(1830, 570, 1000, 56), OneWay: false},
		{Rect: rl.NewRectangle(4470, 570, 1000, 56), OneWay: false},
		{Rect: rl.NewRectangle(3170, 570, 1000, 56), OneWay: false},

		{Rect: rl.NewRectangle(760, 300, 500, 56), OneWay: false},
		{Rect: rl.NewRectangle(2075, 300, 500, 56), OneWay: false},
		{Rect: rl.NewRectangle(4710, 300, 500, 56), OneWay: false},
		{Rect: rl.NewRectangle(3410, 300, 500, 56), OneWay: false},
	}

	OnlyCrouch := []onlyCrouch{
		{Rect: rl.NewRectangle(2615, 3800, 585, 200)},
		{Rect: rl.NewRectangle(2090, 1000, 485, 100)},
		{Rect: rl.NewRectangle(3425, 1000, 485, 100)},
	}
	
	Bean := bean{Pos: rl.NewVector2(50, 3950), Width: 40, Height: 100, Radius: 20, Speed: rl.NewVector2(0, 0), MaxSpeed: 1000, Acceleration: 500, Drag: 460, Jump: 3000, CurrentPlatformIndex: -1, ignoredPlatformIndex: -1, restingOnPlatform: false, Health: 100.0}
	Bean2 := bean{Pos:rl.NewVector2(5950, 3950), Width: 40, Height: 100, Radius: 20, Speed: rl.NewVector2(0,0), MaxSpeed: 1000, Acceleration: 500, Drag: 460, Jump: 3000, CurrentPlatformIndex: -1, ignoredPlatformIndex: -1, restingOnPlatform: false, Health: 100.0}
	
	Gun := gun{Dir: 1, PrevDir: 1, Pos: rl.NewVector2(Bean.Pos.X + 25,  Bean.Pos.Y + 20), Width: 70, Height: 30, Angle: 0.0, Mag: 15, Shots: 0, CanShoot: true, Delay: 0.15}
	Gun2 := gun{Dir: -1, PrevDir: 1, Pos: rl.NewVector2(Bean2.Pos.X - 25, Bean2.Pos.Y + 20), Width:  70, Height: 30, Angle: 0.0, Mag: 15, Shots: 0, CanShoot: true, Delay: 0.15}

	Bullets := []bullet{}

	Gravity := gravity{Max: 1500, Bean: 0, Bean2: 0, Force: 60}
	
	Camera := rl.Camera2D{Offset: rl.NewVector2(Screen.X/2, Screen.Y/2), Target: rl.NewVector2(World.X/2, World.Y/2), Rotation: 0.0, Zoom: 0.2}

	TextureStand := rl.LoadTexture("./Textures/model_player.png")
	TextureCrouch := rl.LoadTexture("./Textures/model_player_crouch.png")

	TextureBullet := rl.LoadTexture("./Textures/Bullet.png")  //Bullet center 9,9 and size 18,18
	TextureGun := rl.LoadTexture("./Textures/Gun.png")
	TextureGunFlipped := rl.LoadTexture("./Textures/Gun_flipped.png")
	
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
		
		if Bean2.IgnoredCooldown > 0 {
			Bean2.IgnoredCooldown -= dT
			if Bean2.IgnoredCooldown <= 0 {
				Bean2.ignoredPlatformIndex = -1
			}
		}

		if Bean2.ignoredPlatformIndex != -1 {
			p := Platforms[Bean2.ignoredPlatformIndex]
			if Bean2.Pos.Y > p.Rect.Y + p.Rect.Height && Bean2.IgnoredCooldown <= 0 {
				Bean2.ignoredPlatformIndex = -1
			}
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

		// Crouching Handler -------------------------------------------------------------------------------------------------------------------------------------------------------
		Bean.OldHeight = Bean.Height
		Bean2.OldHeight = Bean2.Height

		if rl.IsKeyDown(rl.KeyS) {
			Bean.isCrouched = true
		}else {
			Bean.isCrouched = false
		}
		
		if rl.IsKeyDown(rl.KeyK) {
			Bean2.isCrouched = true
		}else {
			Bean2.isCrouched = false
		}

		for _,n := range OnlyCrouch {
			if rl.CheckCollisionRecs(rl.NewRectangle(Bean.Pos.X, Bean.Pos.Y, Bean.Width, Bean.Height), n.Rect) {
				Bean.isCrouched = true
				Bean.hasJumped = true
			}
			if rl.CheckCollisionRecs(rl.NewRectangle(Bean2.Pos.X, Bean2.Pos.Y, Bean2.Width, Bean2.Height), n.Rect) {
				Bean2.isCrouched = true
				Bean2.hasJumped = true
			}
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
		
		if Bean2.isCrouched {
			Bean2.Height = 50
			Bean2.MaxSpeed = 500
			Bean2.Jump = 1500
			Bean2.isCrouched = true
		}else {
			Bean2.Height = 100
			Bean2.MaxSpeed = 1000
			Bean2.Jump = 3000
			Bean2.isCrouched = false
		}
	
		Bean.Pos.Y += Bean.OldHeight - Bean.Height
		Bean2.Pos.Y += Bean2.OldHeight - Bean2.Height
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		// Drag --------------------------------------------------------------------------------------------------------------------------------------------------------------------
		if !rl.IsKeyDown(rl.KeyA) && Bean.Speed.X < 0{
			Bean.Speed.X += Bean.Drag 
			if Bean.Speed.X > 0 {
				Bean.Speed.X = 0.0
			}
		}else if !rl.IsKeyDown(rl.KeyD) && Bean.Speed.X > 0{
			Bean.Speed.X -= Bean.Drag 
			if Bean.Speed.X < 0 {
				Bean.Speed.X = 0.0
			}
		}
		
		if !rl.IsKeyDown(rl.KeyJ) && Bean2.Speed.X < 0 {
			Bean2.Speed.X += Bean2.Drag
			if Bean2.Speed.X > 0 {
				Bean2.Speed.X = 0.0
			}
		}else if !rl.IsKeyDown(rl.KeyL) && Bean2.Speed.X > 0 {
			Bean2.Speed.X-= Bean2.Drag
			if Bean2.Speed.X < 0 {
				Bean2.Speed.X = 0.0
			}
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		// Movement ----------------------------------------------------------------------------------------------------------------------------------------------------------------
		if rl.IsKeyPressed(rl.KeyW) && -Bean.Speed.Y < Bean.MaxSpeed  && (Bean.isGrounded && !Bean.hasJumped){
			Bean.Speed.Y -= Bean.Jump
			Bean.isGrounded = false
			Bean.hasJumped = true
		} 
		if rl.IsKeyDown(rl.KeyA) && rl.IsKeyDown(rl.KeyD) {
			Bean.Speed.X = 0.0
		}else if rl.IsKeyDown(rl.KeyD) && Bean.Speed.X < Bean.MaxSpeed {
			Bean.Speed.X += Bean.Acceleration 
			Gun.Dir = 1
		}else if rl.IsKeyDown(rl.KeyA) && -Bean.Speed.X < Bean.MaxSpeed {
			Bean.Speed.X -= Bean.Acceleration 
			Gun.Dir = -1
		}

		if Bean.Speed.X > Bean.MaxSpeed  && rl.IsKeyDown(rl.KeyD) {
			Bean.Speed.X = Bean.MaxSpeed
		}
		if Bean.Speed.X < -Bean.MaxSpeed && rl.IsKeyDown(rl.KeyA) {
			Bean.Speed.X = -Bean.MaxSpeed
		}

		if rl.IsKeyPressed(rl.KeyI) && -Bean2.Speed.Y < Bean2.MaxSpeed && (Bean2.isGrounded && !Bean2.hasJumped){
			Bean2.Speed.Y -= Bean2.Jump
			Bean2.isGrounded = false
			Bean2.hasJumped = true
		}
		if rl.IsKeyDown(rl.KeyJ) && rl.IsKeyDown(rl.KeyL) {
			Bean2.Speed.X = 0.0
		}else if rl.IsKeyDown(rl.KeyL) && Bean2.Speed.X < Bean2.MaxSpeed {
			Bean2.Speed.X += Bean2.Acceleration
			Gun2.Dir = 1
		}else if rl.IsKeyDown(rl.KeyJ) && -Bean2.Speed.X < Bean2.MaxSpeed {
			Bean2.Speed.X -= Bean2.Acceleration
			Gun2.Dir = -1
		}

		if Bean2.Speed.X > Bean2.MaxSpeed && rl.IsKeyDown(rl.KeyL) {
			Bean2.Speed.X = Bean2.MaxSpeed
		}
		if Bean2.Speed.X < -Bean2.MaxSpeed && rl.IsKeyDown(rl.KeyJ) {
			Bean2.Speed.X = -Bean2.MaxSpeed
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
			Bean.Speed.Y = 0.0
		}
		if Bean.Speed.X <= 0 && Bean.Pos.X  < Map.Border.X + 60 {  // Map Collision for Left Wall
			Bean.Speed.X = 0.0
			Bean.Pos.X = Map.Border.X + Bean.Width + 10
		}else if Bean.Speed.X >= 0 && Bean.Pos.X + Bean.Width > Map.Border.Width - 60 {  // Map Collision for Right Wall
			Bean.Speed.X = 0.0
			Bean.Pos.X = Map.Border.Width - Bean.Width - 40
		}

		if Bean2.Speed.Y >= 0 && Bean2.Pos.Y + Bean2.Height > Map.Border.Height - 50 {
			Bean2.Pos.Y = Map.Border.Height - Bean2.Height - 35
			Bean2.Speed.Y = 0.0
			MapColl.Floor2 = true
		}else if Bean2.Pos.Y + Bean2.Height < Map.Border.Height + 200 {
			MapColl.Floor2 = false
		}
		if Bean2.Pos.Y + (Bean2.Radius*2) < Map.Border.Y + 50 {
			Bean2.Pos.Y = (Bean2.Radius*2) + Bean2.Height + 35
			Bean2.Speed.Y = 0.0
		}
		if Bean2.Speed.X <= 0 && Bean2.Pos.X < Map.Border.X + 60 {
			Bean2.Speed.X = 0.0
			Bean2.Pos.X = Map.Border.X + Bean2.Width + 10
		}else if Bean2.Speed.X >= 0 && Bean2.Pos.X + Bean2.Width > Map.Border.Width - 60 {
			Bean2.Speed.X = 0.0
			Bean2.Pos.X = Map.Border.Width - Bean2.Width - 40
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------	
		
		// Add Speed to Position of the Bean aka movement --------------------------------------------------------------------------------------------------------------------------
		Bean.Pos.X += Bean.Speed.X * dT
		Bean.Pos.Y += Bean.Speed.Y * dT
		Bean2.Pos.X += Bean2.Speed.X * dT
		Bean2.Pos.Y += Bean2.Speed.Y * dT
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		// Gravity -----------------------------------------------------------------------------------------------------------------------------------------------------------------
		if MapColl.Floor || Bean.isGrounded {
			Gravity.Bean = 0
		}else {
			if Gravity.Bean < Gravity.Max && Bean.Speed.Y < Gravity.Max {
				Gravity.Bean += Gravity.Force
				Bean.Speed.Y += Gravity.Bean 
			}else {
				Bean.Speed.Y = Gravity.Bean 
			}
		}

		if MapColl.Floor2 || Bean2.isGrounded {
			Gravity.Bean2 = 0
		}else {
			if Gravity.Bean2 < Gravity.Max && Bean2.Speed.Y < Gravity.Max {
				Gravity.Bean2 += Gravity.Force
				Bean2.Speed.Y += Gravity.Bean2
			}else {
				Bean2.Speed.Y = Gravity.Bean2
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

		if Bean2.isCrouched && Bean2.restingOnPlatform && Bean2.CurrentPlatformIndex != -1 && Platforms[Bean2.CurrentPlatformIndex].OneWay {
			Bean2.ignoredPlatformIndex = Bean2.CurrentPlatformIndex
			Bean2.IgnoredCooldown = 0.25
			Bean2.isGrounded = false
			Bean2.restingOnPlatform = false
		}

		Bean2.CurrentPlatformIndex = -1
		Bean2.restingOnPlatform = false 
		
		for Index,p := range Platforms {
			if Index == Bean.ignoredPlatformIndex {
				continue
			}
			if Index == Bean2.ignoredPlatformIndex {
				continue
			}

			beanBottom := Bean.Pos.Y + Bean.Height

			if beanBottom >= p.Rect.Y - epsilon && beanBottom <= p.Rect.Y + epsilon && Bean.Pos.X + Bean.Width > p.Rect.X && Bean.Pos.X < p.Rect.X + p.Rect.Width {
				Bean.restingOnPlatform = true
				Bean.CurrentPlatformIndex = Index
				break
			}
		}

		for Index, p := range Platforms {
			if Index == Bean2.ignoredPlatformIndex {
				continue
			}

			bean2Bottom	:= Bean2.Pos.Y + Bean2.Height

			if bean2Bottom >= p.Rect.Y - epsilon && bean2Bottom <= p.Rect.Y + epsilon && Bean2.Pos.X + Bean2.Width > p.Rect.X && Bean2.Pos.X < p.Rect.X + p.Rect.Width {
				Bean2.restingOnPlatform = true
				Bean2.CurrentPlatformIndex = Index
			}
		}

		resolveMapCollision(Platforms, &Bean)  // AABB Collision function
		resolveMapCollision(Platforms, &Bean2)


		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

		// Collision Check for Gravity and Jump Check ------------------------------------------------------------------------------------------------------------------------------
		if MapColl.Floor || Bean.restingOnPlatform {
			Bean.isGrounded = true
		}else {
			Bean.isGrounded = false
		}
		if rl.IsKeyReleased(rl.KeyW) {
			Bean.hasJumped = false
		}

		if MapColl.Floor2 || Bean2.restingOnPlatform {
			Bean2.isGrounded = true
		}else {
			Bean2.isGrounded = false
		}
		if rl.IsKeyReleased(rl.KeyI) {
			Bean2.hasJumped = false
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

		// Gun ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		// Bean1 Gun --------------------------------------------------------------------------------------------------------------------------------------------
		if rl.IsKeyDown(rl.KeyE) && rl.IsKeyDown(rl.KeyQ) {
			
		}else if Gun.Dir == 1 {
			Gun.Pos.X = Bean.Pos.X + 25
			Gun.Pos.Y = Bean.Pos.Y + 20
				
			if Gun.Dir != Gun.PrevDir {
				Gun.Angle = 180 - Gun.Angle
				Gun.PrevDir = Gun.Dir
			}

			if rl.IsKeyDown(rl.KeyE) && Gun.Angle < 90 {
				Gun.Angle += 0.7
			}else if rl.IsKeyDown(rl.KeyQ) && Gun.Angle > -90 {
				Gun.Angle -= 0.7
			}

			Gun.Barrel.X = Gun.Pos.X + float32( Math.Cos( float64( (Gun.Angle / 180) * Math.Pi ) ) * float64(Gun.Width) )
			Gun.Barrel.Y = Gun.Pos.Y + float32( Math.Sin( float64( (Gun.Angle / 180) * Math.Pi ) ) * float64(Gun.Width) )

		}else if Gun.Dir == -1 {
			Gun.Pos.X = Bean.Pos.X + 15
			Gun.Pos.Y = Bean.Pos.Y + 20

			if Gun.Dir != Gun.PrevDir {
				Gun.Angle = 180 - Gun.Angle
				Gun.PrevDir = Gun.Dir
			}

			if rl.IsKeyDown(rl.KeyQ) && Gun.Angle < 270 {
				Gun.Angle += 0.7
			}else if rl.IsKeyDown(rl.KeyE) && Gun.Angle > 90 {
				Gun.Angle -= 0.7
			}
			
			Gun.Barrel.X = Gun.Pos.X + float32( Math.Cos( float64( (Gun.Angle / 180) * Math.Pi ) ) * float64(Gun.Width) )
			Gun.Barrel.Y = Gun.Pos.Y + float32( Math.Sin( float64( (Gun.Angle / 180) * Math.Pi ) ) * float64(Gun.Width) )
		}
		
		CheckBarrelPos(&Gun, Map, Platforms)

		if rl.IsKeyDown(rl.KeyF) && Gun.CanShoot{
			Bullets = append(Bullets, NewBullets(Gun))
		}
		//-------------------------------------------------------------------------------------------------------------------------------------------------------
		
		//Bean2 Gun ---------------------------------------------------------------------------------------------------------------------------------------------
		if rl.IsKeyDown(rl.KeyU) && rl.IsKeyDown(rl.KeyO) {
			
		}else if Gun2.Dir == 1 {
			Gun2.Pos.X = Bean2.Pos.X + 25
			Gun2.Pos.Y = Bean2.Pos.Y + 20
			
			if Gun2.Dir != Gun2.PrevDir {
				Gun2.Angle = 180 - Gun2.Angle
				Gun2.PrevDir = Gun2.Dir
			}

			if rl.IsKeyDown(rl.KeyO) && Gun2.Angle < 90 {
				Gun2.Angle += 0.7
			}else if rl.IsKeyDown(rl.KeyU) && Gun2.Angle > -90 {
				Gun2.Angle -= 0.7
			}

			Gun2.Barrel.X = Gun2.Pos.X + float32( Math.Cos( float64( (Gun2.Angle / 180) * Math.Pi) ) * float64(Gun2.Width) )
			Gun2.Barrel.Y = Gun2.Pos.Y + float32( Math.Sin( float64( (Gun2.Angle / 180) * Math.Pi) ) * float64(Gun2.Width) )
		}else if Gun2.Dir == -1 {
			Gun2.Pos.X = Bean2.Pos.X + 15
			Gun2.Pos.Y = Bean2.Pos.Y + 20
			
			if Gun2.Dir != Gun2.PrevDir {
				Gun2.Angle = 180 - Gun2.Angle
				Gun2.PrevDir = Gun2.Dir
			}
			
			if rl.IsKeyDown(rl.KeyU) && Gun2.Angle < 270 {
				Gun2.Angle += 0.7
				}else if rl.IsKeyDown(rl.KeyO) && Gun2.Angle > 90 {
					Gun2.Angle -= 0.7
				}
			}
			
			Gun2.Barrel.X = Gun2.Pos.X + float32( Math.Cos( float64( (Gun2.Angle / 180) * Math.Pi) ) * float64(Gun2.Width) )
			Gun2.Barrel.Y = Gun2.Pos.Y + float32( Math.Sin( float64( (Gun2.Angle / 180) * Math.Pi) ) * float64(Gun2.Width) )

			CheckBarrelPos(&Gun2, Map, Platforms)

			if rl.IsKeyDown(rl.KeySemicolon) {
				Bullets = append(Bullets, NewBullets(Gun2))
			}
			//------------------------------------------------------------------------------------------------------------------------------------------------------


			for Bullet := 0; Bullet < len(Bullets); {
				Bullets[Bullet].PrevPos = Bullets[Bullet].Pos
				Bullets[Bullet].Pos.X += Bullets[Bullet].Speed.X * dT
				Bullets[Bullet].Pos.Y += Bullets[Bullet].Speed.Y * dT

				if Bullets[Bullet].Pos.X < Map.Border.X + 50 {
					Bullets[Bullet].Speed.X = -Bullets[Bullet].Speed.X
					Bullets[Bullet].Pos.X += 5
					Bullets[Bullet].Time += 0.2
				}else if Bullets[Bullet].Pos.X  > Map.Border.X + Map.Border.Width - 50 {
					Bullets[Bullet].Speed.X = -Bullets[Bullet].Speed.X
					Bullets[Bullet].Pos.X -= 5
					Bullets[Bullet].Time += 0.2
				}
				if Bullets[Bullet].Pos.Y > Map.Border.Y + Map.Border.Height - 45 {
					Bullets[Bullet].Speed.Y = -Bullets[Bullet].Speed.Y
					Bullets[Bullet].Pos.Y -= 5
					Bullets[Bullet].Time += 0.2
				}else if Bullets[Bullet].Pos.Y < Map.Border.Y + 45 {
					Bullets[Bullet].Speed.Y = -Bullets[Bullet].Speed.Y
					Bullets[Bullet].Pos.Y += 5
					Bullets[Bullet].Time += 0.2
				}
				
				resolveMapBulletCollision(Platforms, &Bullets[Bullet])

				CheckBulletPlayerCollision(&Bullets[Bullet], &Bean)
				if Bullets[Bullet].Time <= 0 {
					Bullets[Bullet] = Bullets[len(Bullets)-1]
					Bullets = Bullets[:len(Bullets)-1]
					// Bullet--
					}else {
						Bullets[Bullet].Time -= dT
					Bullet++
				}
				
			}

		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

		// Health Bar --------------------------------------------------------------------------------------------------------------------------------------------------------------
		if Bean.Health < 0 {
			Bean.Health = 0
		}
		//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

		rl.BeginDrawing()

		rl.ClearBackground(rl.GetColor(0x0034a9ff))

		
		rl.BeginMode2D(Camera)
		rl.DrawRectangle(3000, 3000, 500, 500, rl.GetColor(0xff0000ff))
		
		rl.DrawRectangleLinesEx(Map.Border, 30, rl.GetColor(0x000000ff))
		
		for _, p := range Platforms {
			if p.OneWay == true{
				rl.DrawRectangleRec(p.Rect, colorOneWay)
			}else {
				rl.DrawRectangleRec(p.Rect, colorSolid)
			}
		}

		for _, b := range Bullets {
			rl.DrawTexturePro(TextureBullet, rl.NewRectangle(0, 0, 18, 18), rl.NewRectangle(b.Pos.X, b.Pos.Y, 18, 18), rl.NewVector2(9,9), b.Angle, rl.GetColor(0xffffffff))
			// rl.DrawCircleV(b.Pos, b.Radius, rl.GetColor(0xff0000ff))
		}
		
		rl.DrawRectangleV(Bean.Pos, rl.NewVector2(Bean.Width, Bean.Height), rl.GetColor(0x00ffffff))
		Bean.Pos.X += Bean.Width/2
		Bean.Pos.Y -= 19
		rl.DrawCircleV(Bean.Pos, Bean.Radius, rl.GetColor(0x00ffffff))
		Bean.Pos.X -= Bean.Width/2
		Bean.Pos.Y += 19

		rl.DrawRectangleV(Bean2.Pos, rl.NewVector2(Bean2.Width, Bean2.Height), rl.GetColor(0x00ffffff))
		Bean2.Pos.X += Bean2.Width/2
		Bean2.Pos.Y -= 19
		rl.DrawCircleV(Bean2.Pos, Bean2.Radius, rl.GetColor(0x00ffffff))
		Bean2.Pos.X -= Bean2.Width/2
		Bean2.Pos.Y += 19
		
		if !Bean.isCrouched{
			rl.DrawTextureV(TextureStand, rl.NewVector2(Bean.Pos.X, Bean.Pos.Y - (Bean.Radius*2)), rl.GetColor(0xffffffff))
		}else {
			rl.DrawTextureV(TextureCrouch, rl.NewVector2(Bean.Pos.X, Bean.Pos.Y - (Bean.Radius*2)), rl.GetColor(0xffffffff))
		}

		// gun renderer ------------------------------------------------------------------------------------------------------------------
		if Gun.Dir == 1 {
			rl.DrawTextureEx(TextureGun, rl.NewVector2(Gun.Pos.X - 15, Gun.Pos.Y - 15), Gun.Angle, 2, rl.GetColor(0xffffffff))
		}else {
			rl.DrawTextureEx(TextureGunFlipped, rl.NewVector2(Gun.Pos.X + 15, Gun.Pos.Y + 45), Gun.Angle, 2, rl.GetColor(0xffffffff))
		}
		if Gun2.Dir == 1 {
			rl.DrawTextureEx(TextureGun, rl.NewVector2(Gun2.Pos.X -15, Gun2.Pos.Y - 15), Gun2.Angle , 2, rl.GetColor(0xffffffff))
		}else {
			rl.DrawTextureEx(TextureGunFlipped, rl.NewVector2(Gun2.Pos.X +15, Gun2.Pos.Y + 45), Gun2.Angle, 2, rl.GetColor(0xffffffff))
		}
		//--------------------------------------------------------------------------------------------------------------------------------

		rl.DrawRectangleLinesEx(rl.NewRectangle(Bean.Pos.X + (Bean.Width/2) - 100, Bean.Pos.Y - (Bean.Radius * 2) - 50, 200, 30), 10, rl.GetColor(0x000000ff))
		rl.DrawRectangleV(rl.NewVector2(Bean.Pos.X + (Bean.Width/2) - 90, Bean.Pos.Y - (Bean.Radius*2) - 40), rl.NewVector2((180 * (Bean.Health/100)), 10), rl.GetColor(0x00ff00ff))

		for i := float32(0); i < 180; i++ {
			if Bean2.Pos.X + (Bean2.Width/2) + 90 + i > Bean2.Pos.X + (Bean2.Width/2) + 90 {
				width = i
			}else {
				continue
			}
		}
		rl.DrawRectangleLinesEx(rl.NewRectangle(Bean2.Pos.X + (Bean2.Width/2) - 100, Bean2.Pos.Y - (Bean2.Radius * 2) - 50, 200, 30), 10, rl.GetColor(0x000000ff))
		rl.DrawRectangleV(rl.NewVector2((Bean2.Pos.X + (Bean2.Width/2) - 90) * (Bean2.Health/100), Bean2.Pos.Y - (Bean2.Radius * 2) - 40), rl.NewVector2(width, 10), rl.GetColor(0x00ff00ff))
			
		rl.EndMode2D()

		rl.DrawText(fmt.Sprintf("SpeedX: %0.1f\nSpeedY: %0.1f\nGravity Bean: %0.1f\nGrounded: %v\nCrouched: %v\nGun Angle: %0.1f\nGun2 Angle: %0.1f",Bean.Speed.X, Bean.Speed.Y, Gravity.Bean, Bean.isGrounded, Bean.isCrouched, Gun.Angle, Gun2.Angle), 10, 10, 30, rl.GetColor(0xffffffff))
		
		rl.EndDrawing()

	}
}