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
	Pos                  rl.Vector2
	Width                float32
	Height               float32
	OldHeight            float32
	Radius               float32
	Speed                rl.Vector2
	MaxSpeed             float32
	Acceleration         float32
	Drag                 float32
	Jump                 float32
	CanJump              bool
	isGrounded           bool
	hasJumped            bool
	isCrouched           bool
	ignoredPlatformIndex int
	IgnoredCooldown      float32
	CurrentPlatformIndex int
	restingOnPlatform    bool
	Health               float32
	Lives                int
	PowersNumber         int
}
type gun struct {
	Dir         int
	PrevDir     int
	Pos         rl.Vector2
	Width       float32
	Height      float32
	Angle       float32
	Mag         int
	Shots       int
	CanShoot    bool
	Delay       float32
	Reloading   bool
	ReloadDelay float32
	Barrel      rl.Vector2
}
type bullet struct {
	Angle   float32
	Pos     rl.Vector2
	PrevPos rl.Vector2
	Radius  float32
	Speed   rl.Vector2
	Damage  float32
	Time    float32
}
type gravity struct {
	Max   float32
	Bean  float32
	Bean2 float32
	Force float32
}
type Map struct {
	Border rl.Rectangle
}
type Platform struct {
	Rect   rl.Rectangle
	OneWay bool
}
type onlyCrouch struct {
	Rect rl.Rectangle
}
type mapColl struct {
	Floor  bool
	Floor2 bool
}

//	type power struct {
//		NFirerate float32
//		FireRate float32
//		NDamage float32
//		Damage float32
//		NReload float32
//		Reload	float32
//		NMag float32
//		Mag float32
//		NHealth float32
//		Health float32
//		Imune bool
//	}
type powerUpData struct {
	Type     int
	Weight   int
	Duration float32
	Texture  rl.Texture2D
}
type PowerUp struct {
	Type       int
	Pos        rl.Vector2
	SpawnIndex int
	LifeSpan   float32
}
type spawnPoints struct {
	Pos      rl.Vector2
	Occupied bool
}
type beanPowers struct {
	Type     int
	LifeTime float32
	Active   bool
}

var MapColl mapColl

var Platforms []Platform
var colorOneWay rl.Color = rl.GetColor(0x444444ff)
var colorSolid rl.Color = rl.GetColor(0x000000ff)
var width float32
var ScoreP1 int
var ScoreP2 int
var Timer float32
var Minutes float32
var Hour float32
var FPS int32
var SpawnPowerUp float32 = 4
var PowerType1 int
var Power1 bool
var PowerType2 int
var Power2 bool
var BeanPowers []beanPowers
var Bean2Powers []beanPowers

var Start bool
var Menu bool = true
var Controls bool
var Quit bool

var Pause bool

const epsilon float32 = 2.5

const (
	PUHealth = iota
	PUFirerate
	PUDamage
	PUReload
	PUMag
	PUImune
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowMaximized)
	rl.InitWindow(0, 0, "Go Shooter Game")
	defer rl.CloseWindow()

	rl.SetExitKey(rl.KeyNull)

	rl.SetTargetFPS(60)

	Screen := screen{X: float32(rl.GetScreenWidth()), Y: float32(rl.GetScreenHeight())}

	World := world{X: 6000, Y: 4000}

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
		{Rect: rl.NewRectangle(1400, 3070, 1200, 56), OneWay: false},
		{Rect: rl.NewRectangle(30, 3070, 800, 56), OneWay: false},
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

		{Rect: rl.NewRectangle(30, 1390, 700, 56), OneWay: false},
		{Rect: rl.NewRectangle(5300, 1390, 700, 56), OneWay: false},

		{Rect: rl.NewRectangle(500, 1110, 600, 56), OneWay: true},
		{Rect: rl.NewRectangle(1101, 1110, 3800, 56), OneWay: false},
		{Rect: rl.NewRectangle(4900, 1110, 600, 56), OneWay: true},

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

	Bean := bean{Pos: rl.NewVector2(50, 3950), Width: 40, Height: 100, Radius: 20, Speed: rl.NewVector2(0, 0), MaxSpeed: 1000, Acceleration: 500, Drag: 460, Jump: 3000, CurrentPlatformIndex: -1, ignoredPlatformIndex: -1, restingOnPlatform: false, Health: 100.0, Lives: 3}
	Bean2 := bean{Pos: rl.NewVector2(5950, 3950), Width: 40, Height: 100, Radius: 20, Speed: rl.NewVector2(0, 0), MaxSpeed: 1000, Acceleration: 500, Drag: 460, Jump: 3000, CurrentPlatformIndex: -1, ignoredPlatformIndex: -1, restingOnPlatform: false, Health: 100.0, Lives: 3}

	Gun := gun{Dir: 1, PrevDir: 1, Pos: rl.NewVector2(Bean.Pos.X+25, Bean.Pos.Y+20), Width: 70, Height: 30, Angle: 0.0, Mag: 15, Shots: 0, CanShoot: true, Delay: 0.15}
	Gun2 := gun{Dir: -1, PrevDir: 1, Pos: rl.NewVector2(Bean2.Pos.X-25, Bean2.Pos.Y+20), Width: 70, Height: 30, Angle: 0.0, Mag: 15, Shots: 0, CanShoot: true, Delay: 0.15}

	Bullets := []bullet{}

	Gravity := gravity{Max: 1500, Bean: 0, Bean2: 0, Force: 60}

	Camera := rl.Camera2D{Offset: rl.NewVector2(Screen.X/2, Screen.Y/2), Target: rl.NewVector2(World.X/2, World.Y/2), Rotation: 0.0, Zoom: 0.2}

	SpawnPoints := []spawnPoints{
		{Pos: rl.NewVector2(100, 300), Occupied: false},
		{Pos: rl.NewVector2(200, 300), Occupied: false},
		{Pos: rl.NewVector2(300, 300), Occupied: false},
		{Pos: rl.NewVector2(400, 300), Occupied: false},
		{Pos: rl.NewVector2(500, 300), Occupied: false},
		{Pos: rl.NewVector2(600, 300), Occupied: false},
	}

	PowerUps := []powerUpData{
		{Type: PUHealth, Weight: 50, Duration: 0.001, Texture: rl.LoadTexture("./Textures/Powers/PowerHealth.png")},
		{Type: PUFirerate, Weight: 40, Duration: 5, Texture: rl.LoadTexture("./Textures/Powers/PowerFirerate.png")},
		{Type: PUDamage, Weight: 30, Duration: 5, Texture: rl.LoadTexture("./Textures/Powers/PowerDamage.png")},
		// {Type: PUReload, Weight: 25, Duration: 10, Texture: rl.LoadTexture("")},
		// {Type: PUMag, Weight: 20, Duration: 10, Texture: rl.LoadTexture("")},
		// {Type: PUImune, Weight: 10, Duration: 5, Texture: rl.LoadTexture("")},
	}

	SpawnedPowerUps := []PowerUp{}

	TextureStand := rl.LoadTexture("./Textures/model_player.png")
	TextureCrouch := rl.LoadTexture("./Textures/model_player_crouch.png")

	TextureBullet := rl.LoadTexture("./Textures/Bullet.png") //Bullet center 9,9 and size 18,18
	TextureGun := rl.LoadTexture("./Textures/Gun.png")
	TextureGunFlipped := rl.LoadTexture("./Textures/Gun_flipped.png")
	TextureAmmoContainer := rl.LoadTexture("./Textures/Ammo_Container.png")
	TextureAmmo := rl.LoadTexture("./Textures/Ammo.png")
	TextureHeart := rl.LoadTexture("./Textures/Heart.png")

	Start = false

	for !rl.WindowShouldClose() {

		dT := rl.GetFrameTime() // Delta Time for allowing the game to run at any FPS and framerate
		FPS = rl.GetFPS()

		if Start {

			if rl.IsKeyPressed(rl.KeySpace) {
				Pause = !Pause
			}

			if !Pause {

				// OneWay Collision Helper -------------------------------------------------------------------------------------------------------------------------------------------------
				if Bean.IgnoredCooldown > 0 { // OneWay playform collision ignore timer countdown
					Bean.IgnoredCooldown -= dT
					if Bean.IgnoredCooldown <= 0 {
						Bean.ignoredPlatformIndex = -1 // resets OneWay collision
					}
				}

				if Bean.ignoredPlatformIndex != -1 {
					p := Platforms[Bean.ignoredPlatformIndex]
					if Bean.Pos.Y > p.Rect.Y+p.Rect.Height && Bean.IgnoredCooldown <= 0 {
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
					if Bean2.Pos.Y > p.Rect.Y+p.Rect.Height && Bean2.IgnoredCooldown <= 0 {
						Bean2.ignoredPlatformIndex = -1
					}
				}
				//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

				// Crouching Handler -------------------------------------------------------------------------------------------------------------------------------------------------------
				Bean.OldHeight = Bean.Height
				Bean2.OldHeight = Bean2.Height

				if rl.IsKeyDown(rl.KeyS) {
					Bean.isCrouched = true
				} else {
					Bean.isCrouched = false
				}

				if rl.IsKeyDown(rl.KeyK) {
					Bean2.isCrouched = true
				} else {
					Bean2.isCrouched = false
				}

				for _, n := range OnlyCrouch {
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
				} else {
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
				} else {
					Bean2.Height = 100
					Bean2.MaxSpeed = 1000
					Bean2.Jump = 3000
					Bean2.isCrouched = false
				}

				Bean.Pos.Y += Bean.OldHeight - Bean.Height
				Bean2.Pos.Y += Bean2.OldHeight - Bean2.Height
				//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

				// Drag --------------------------------------------------------------------------------------------------------------------------------------------------------------------
				if !rl.IsKeyDown(rl.KeyA) && Bean.Speed.X < 0 {
					Bean.Speed.X += Bean.Drag
					if Bean.Speed.X > 0 {
						Bean.Speed.X = 0.0
					}
				} else if !rl.IsKeyDown(rl.KeyD) && Bean.Speed.X > 0 {
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
				} else if !rl.IsKeyDown(rl.KeyL) && Bean2.Speed.X > 0 {
					Bean2.Speed.X -= Bean2.Drag
					if Bean2.Speed.X < 0 {
						Bean2.Speed.X = 0.0
					}
				}
				//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

				// Movement ----------------------------------------------------------------------------------------------------------------------------------------------------------------
				if rl.IsKeyPressed(rl.KeyW) && -Bean.Speed.Y < Bean.MaxSpeed && (Bean.isGrounded && !Bean.hasJumped) {
					Bean.Speed.Y -= Bean.Jump
					Bean.isGrounded = false
					Bean.hasJumped = true
				}
				if rl.IsKeyDown(rl.KeyA) && rl.IsKeyDown(rl.KeyD) {
					Bean.Speed.X = 0.0
				} else if rl.IsKeyDown(rl.KeyD) && Bean.Speed.X < Bean.MaxSpeed {
					Bean.Speed.X += Bean.Acceleration
					Gun.Dir = 1
				} else if rl.IsKeyDown(rl.KeyA) && -Bean.Speed.X < Bean.MaxSpeed {
					Bean.Speed.X -= Bean.Acceleration
					Gun.Dir = -1
				}

				if Bean.Speed.X > Bean.MaxSpeed && rl.IsKeyDown(rl.KeyD) {
					Bean.Speed.X = Bean.MaxSpeed
				}
				if Bean.Speed.X < -Bean.MaxSpeed && rl.IsKeyDown(rl.KeyA) {
					Bean.Speed.X = -Bean.MaxSpeed
				}

				if rl.IsKeyPressed(rl.KeyI) && -Bean2.Speed.Y < Bean2.MaxSpeed && (Bean2.isGrounded && !Bean2.hasJumped) {
					Bean2.Speed.Y -= Bean2.Jump
					Bean2.isGrounded = false
					Bean2.hasJumped = true
				}
				if rl.IsKeyDown(rl.KeyJ) && rl.IsKeyDown(rl.KeyL) {
					Bean2.Speed.X = 0.0
				} else if rl.IsKeyDown(rl.KeyL) && Bean2.Speed.X < Bean2.MaxSpeed {
					Bean2.Speed.X += Bean2.Acceleration
					Gun2.Dir = 1
				} else if rl.IsKeyDown(rl.KeyJ) && -Bean2.Speed.X < Bean2.MaxSpeed {
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
				if Bean.Speed.Y >= 0 && Bean.Pos.Y+Bean.Height > Map.Border.Height-50 { // Map Collision for floor
					Bean.Pos.Y = Map.Border.Height - Bean.Height - 35
					Bean.Speed.Y = 0.0
					MapColl.Floor = true
				} else if Bean.Pos.Y+Bean.Height < Map.Border.Height+200 {
					MapColl.Floor = false
				}
				if Bean.Pos.Y+(Bean.Radius*2) < Map.Border.Y+50 { // Map Collision for Roof
					Bean.Pos.Y = (Bean.Radius * 2) + Bean.Height + 35
					Bean.Speed.Y = 0.0
				}
				if Bean.Speed.X <= 0 && Bean.Pos.X < Map.Border.X+60 { // Map Collision for Left Wall
					Bean.Speed.X = 0.0
					Bean.Pos.X = Map.Border.X + Bean.Width + 10
				} else if Bean.Speed.X >= 0 && Bean.Pos.X+Bean.Width > Map.Border.Width-60 { // Map Collision for Right Wall
					Bean.Speed.X = 0.0
					Bean.Pos.X = Map.Border.Width - Bean.Width - 40
				}

				if Bean2.Speed.Y >= 0 && Bean2.Pos.Y+Bean2.Height > Map.Border.Height-50 {
					Bean2.Pos.Y = Map.Border.Height - Bean2.Height - 35
					Bean2.Speed.Y = 0.0
					MapColl.Floor2 = true
				} else if Bean2.Pos.Y+Bean2.Height < Map.Border.Height+200 {
					MapColl.Floor2 = false
				}
				if Bean2.Pos.Y+(Bean2.Radius*2) < Map.Border.Y+50 {
					Bean2.Pos.Y = (Bean2.Radius * 2) + Bean2.Height + 35
					Bean2.Speed.Y = 0.0
				}
				if Bean2.Speed.X <= 0 && Bean2.Pos.X < Map.Border.X+60 {
					Bean2.Speed.X = 0.0
					Bean2.Pos.X = Map.Border.X + Bean2.Width + 10
				} else if Bean2.Speed.X >= 0 && Bean2.Pos.X+Bean2.Width > Map.Border.Width-60 {
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
				} else {
					if Gravity.Bean < Gravity.Max && Bean.Speed.Y < Gravity.Max {
						Gravity.Bean += Gravity.Force
						Bean.Speed.Y += Gravity.Bean
					} else {
						Bean.Speed.Y = Gravity.Bean
					}
				}

				if MapColl.Floor2 || Bean2.isGrounded {
					Gravity.Bean2 = 0
				} else {
					if Gravity.Bean2 < Gravity.Max && Bean2.Speed.Y < Gravity.Max {
						Gravity.Bean2 += Gravity.Force
						Bean2.Speed.Y += Gravity.Bean2
					} else {
						Bean2.Speed.Y = Gravity.Bean2
					}
				}
				//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

				// Platform Collisions -----------------------------------------------------------------------------------------------------------------------------------------------------
				if Bean.isCrouched && Bean.restingOnPlatform && Bean.CurrentPlatformIndex != -1 && Platforms[Bean.CurrentPlatformIndex].OneWay {
					Bean.ignoredPlatformIndex = Bean.CurrentPlatformIndex
					Bean.IgnoredCooldown = 0.25
					Bean.isGrounded = false
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

				for Index, p := range Platforms {
					if Index == Bean.ignoredPlatformIndex {
						continue
					}
					if Index == Bean2.ignoredPlatformIndex {
						continue
					}

					beanBottom := Bean.Pos.Y + Bean.Height

					if beanBottom >= p.Rect.Y-epsilon && beanBottom <= p.Rect.Y+epsilon && Bean.Pos.X+Bean.Width > p.Rect.X && Bean.Pos.X < p.Rect.X+p.Rect.Width {
						Bean.restingOnPlatform = true
						Bean.CurrentPlatformIndex = Index
						break
					}
				}

				for Index, p := range Platforms {
					if Index == Bean2.ignoredPlatformIndex {
						continue
					}

					bean2Bottom := Bean2.Pos.Y + Bean2.Height

					if bean2Bottom >= p.Rect.Y-epsilon && bean2Bottom <= p.Rect.Y+epsilon && Bean2.Pos.X+Bean2.Width > p.Rect.X && Bean2.Pos.X < p.Rect.X+p.Rect.Width {
						Bean2.restingOnPlatform = true
						Bean2.CurrentPlatformIndex = Index
					}
				}

				resolveMapCollision(Platforms, &Bean) // AABB Collision function
				resolveMapCollision(Platforms, &Bean2)

				//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

				// Collision Check for Gravity and Jump Check ------------------------------------------------------------------------------------------------------------------------------
				if MapColl.Floor || Bean.restingOnPlatform {
					Bean.isGrounded = true
				} else {
					Bean.isGrounded = false
				}
				if rl.IsKeyReleased(rl.KeyW) {
					Bean.hasJumped = false
				}

				if MapColl.Floor2 || Bean2.restingOnPlatform {
					Bean2.isGrounded = true
				} else {
					Bean2.isGrounded = false
				}
				if rl.IsKeyReleased(rl.KeyI) {
					Bean2.hasJumped = false
				}
				//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

				// Gun ---------------------------------------------------------------------------------------------------------------------------------------------------------------------

				// Bean1 Gun --------------------------------------------------------------------------------------------------------------------------------------------
				if rl.IsKeyDown(rl.KeyE) && rl.IsKeyDown(rl.KeyQ) {

				} else if Gun.Dir == 1 {
					Gun.Pos.X = Bean.Pos.X + 25
					Gun.Pos.Y = Bean.Pos.Y + 20

					if Gun.Dir != Gun.PrevDir {
						Gun.Angle = 180 - Gun.Angle
						Gun.PrevDir = Gun.Dir
					}

					if rl.IsKeyDown(rl.KeyE) && Gun.Angle < 90 {
						Gun.Angle += 0.7
					} else if rl.IsKeyDown(rl.KeyQ) && Gun.Angle > -90 {
						Gun.Angle -= 0.7
					}

					Gun.Barrel.X = Gun.Pos.X + float32(Math.Cos(float64((Gun.Angle/180)*Math.Pi))*float64(Gun.Width))
					Gun.Barrel.Y = Gun.Pos.Y + float32(Math.Sin(float64((Gun.Angle/180)*Math.Pi))*float64(Gun.Width))

				} else if Gun.Dir == -1 {
					Gun.Pos.X = Bean.Pos.X + 15
					Gun.Pos.Y = Bean.Pos.Y + 20

					if Gun.Dir != Gun.PrevDir {
						Gun.Angle = 180 - Gun.Angle
						Gun.PrevDir = Gun.Dir
					}

					if rl.IsKeyDown(rl.KeyQ) && Gun.Angle < 270 {
						Gun.Angle += 0.7
					} else if rl.IsKeyDown(rl.KeyE) && Gun.Angle > 90 {
						Gun.Angle -= 0.7
					}

					Gun.Barrel.X = Gun.Pos.X + float32(Math.Cos(float64((Gun.Angle/180)*Math.Pi))*float64(Gun.Width))
					Gun.Barrel.Y = Gun.Pos.Y + float32(Math.Sin(float64((Gun.Angle/180)*Math.Pi))*float64(Gun.Width))
				}

				CheckBarrelPos(&Gun, Map, Platforms)

				if rl.IsKeyPressed(rl.KeyR) && !Gun.Reloading && Gun.Shots < Gun.Mag {
					Gun.Reloading = true
					Gun.ReloadDelay = 1.5
				}

				if Gun.Shots >= Gun.Mag && !Gun.Reloading {
					Gun.Reloading = true
					Gun.ReloadDelay = 1.5
				}

				if Gun.Reloading {
					Gun.ReloadDelay -= dT
					if Gun.ReloadDelay <= 0 {
						Gun.Shots = 1
						Gun.Reloading = false
					}
				} else {
					if Gun.Delay > 0 {
						Gun.Delay -= dT
					}
				}

				if Gun.Delay <= 0 {
					if Gun.Shots < Gun.Mag && !Gun.Reloading {
						Gun.CanShoot = true
					}
				}

				if rl.IsKeyDown(rl.KeyF) && Gun.CanShoot && !Gun.Reloading {
					Gun.CanShoot = false
					if Gun.Shots < Gun.Mag {
						Gun.Delay = 0.25
					}
					Gun.Shots++
					Bullets = append(Bullets, NewBullets(Gun))
				}
				//-------------------------------------------------------------------------------------------------------------------------------------------------------

				//Bean2 Gun ---------------------------------------------------------------------------------------------------------------------------------------------
				if rl.IsKeyDown(rl.KeyU) && rl.IsKeyDown(rl.KeyO) {

				} else if Gun2.Dir == 1 {
					Gun2.Pos.X = Bean2.Pos.X + 25
					Gun2.Pos.Y = Bean2.Pos.Y + 20

					if Gun2.Dir != Gun2.PrevDir {
						Gun2.Angle = 180 - Gun2.Angle
						Gun2.PrevDir = Gun2.Dir
					}

					if rl.IsKeyDown(rl.KeyO) && Gun2.Angle < 90 {
						Gun2.Angle += 0.7
					} else if rl.IsKeyDown(rl.KeyU) && Gun2.Angle > -90 {
						Gun2.Angle -= 0.7
					}

					Gun2.Barrel.X = Gun2.Pos.X + float32(Math.Cos(float64((Gun2.Angle/180)*Math.Pi))*float64(Gun2.Width))
					Gun2.Barrel.Y = Gun2.Pos.Y + float32(Math.Sin(float64((Gun2.Angle/180)*Math.Pi))*float64(Gun2.Width))
				} else if Gun2.Dir == -1 {
					Gun2.Pos.X = Bean2.Pos.X + 15
					Gun2.Pos.Y = Bean2.Pos.Y + 20

					if Gun2.Dir != Gun2.PrevDir {
						Gun2.Angle = 180 - Gun2.Angle
						Gun2.PrevDir = Gun2.Dir
					}

					if rl.IsKeyDown(rl.KeyU) && Gun2.Angle < 270 {
						Gun2.Angle += 0.7
					} else if rl.IsKeyDown(rl.KeyO) && Gun2.Angle > 90 {
						Gun2.Angle -= 0.7
					}
				}

				Gun2.Barrel.X = Gun2.Pos.X + float32(Math.Cos(float64((Gun2.Angle/180)*Math.Pi))*float64(Gun2.Width))
				Gun2.Barrel.Y = Gun2.Pos.Y + float32(Math.Sin(float64((Gun2.Angle/180)*Math.Pi))*float64(Gun2.Width))

				CheckBarrelPos(&Gun2, Map, Platforms)

				// Gun2.Reloading = false

				if rl.IsKeyPressed(rl.KeyP) && !Gun2.Reloading && Gun2.Shots < Gun2.Mag {
					Gun2.Reloading = true
					Gun2.ReloadDelay = 1.5
				}

				if Gun2.Shots >= Gun2.Mag && !Gun2.Reloading {
					Gun2.Reloading = true
					Gun2.ReloadDelay = 1.5
				}

				if Gun2.Reloading {
					Gun2.ReloadDelay -= dT
					if Gun2.ReloadDelay <= 0 {
						Gun2.Shots = 1
						Gun2.Reloading = false
					}
				} else {
					if Gun2.Delay > 0 {
						Gun2.Delay -= dT
					}
				}

				if Gun2.Delay <= 0 {
					if Gun2.Shots < Gun2.Mag && !Gun2.Reloading {
						Gun2.CanShoot = true
					}
				}

				if rl.IsKeyDown(rl.KeySemicolon) && Gun2.CanShoot && !Gun2.Reloading {
					Gun2.CanShoot = false
					if Gun2.Shots < Gun2.Mag {
						Gun2.Delay = 0.25
					}
					Gun2.Shots++
					Bullets = append(Bullets, NewBullets(Gun2))
				}
				//------------------------------------------------------------------------------------------------------------------------------------------------------

				for Bullet := 0; Bullet < len(Bullets); {
					Bullets[Bullet].PrevPos = Bullets[Bullet].Pos
					Bullets[Bullet].Pos.X += Bullets[Bullet].Speed.X * dT
					Bullets[Bullet].Pos.Y += Bullets[Bullet].Speed.Y * dT

					if Bullets[Bullet].Pos.X < Map.Border.X+50 {
						Bullets[Bullet].Speed.X = -Bullets[Bullet].Speed.X
						Bullets[Bullet].Pos.X += 5
						Bullets[Bullet].Time += 0.2
					} else if Bullets[Bullet].Pos.X > Map.Border.X+Map.Border.Width-50 {
						Bullets[Bullet].Speed.X = -Bullets[Bullet].Speed.X
						Bullets[Bullet].Pos.X -= 5
						Bullets[Bullet].Time += 0.2
					}
					if Bullets[Bullet].Pos.Y > Map.Border.Y+Map.Border.Height-45 {
						Bullets[Bullet].Speed.Y = -Bullets[Bullet].Speed.Y
						Bullets[Bullet].Pos.Y -= 5
						Bullets[Bullet].Time += 0.2
					} else if Bullets[Bullet].Pos.Y < Map.Border.Y+45 {
						Bullets[Bullet].Speed.Y = -Bullets[Bullet].Speed.Y
						Bullets[Bullet].Pos.Y += 5
						Bullets[Bullet].Time += 0.2
					}

					resolveMapBulletCollision(Platforms, &Bullets[Bullet])

					CheckBulletPlayerCollision(&Bullets[Bullet], &Bean)
					CheckBulletPlayerCollision(&Bullets[Bullet], &Bean2)

					if Bullets[Bullet].Time <= 0 {
						Bullets[Bullet] = Bullets[len(Bullets)-1]
						Bullets = Bullets[:len(Bullets)-1]
					} else {
						Bullets[Bullet].Time -= dT
						Bullet++
					}

				}

				//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

				// Power Ups Handler -------------------------------------------------------------------------------------------------------------------------------------------------------
				// replace with code for Power Ups Handler lol, this got a bit difficult so i doing menus first.
				if SpawnPowerUp < 0.0 {
					SpawnPowerUp = 0.0
				} else {
					SpawnPowerUp -= dT
				}

				if SpawnPowerUp <= 0.0 {
					SpawnedPowerUps = RandomPowerUpSpawner(SpawnedPowerUps, SpawnPoints, PowerUps)
					SpawnPowerUp = 5.0
				}


				if Bean.PowersNumber < 3 {
					SpawnedPowerUps, Power1, PowerType1 = CheckPlayerPowerUpPickUp(SpawnedPowerUps, &Bean, SpawnPoints)
					BeanPowers = append(BeanPowers, beanPowers{Type: PowerType1, LifeTime: PowerUps[PowerType1].Duration, Active: false})
				}
				if Bean2.PowersNumber < 3 {
					SpawnedPowerUps, Power2, PowerType2 = CheckPlayerPowerUpPickUp(SpawnedPowerUps, &Bean2, SpawnPoints)
					Bean2Powers = append(Bean2Powers, beanPowers{Type: PowerType2, LifeTime: PowerUps[PowerType2].Duration, Active: false})
				}

				if len(BeanPowers) != 0 {
					
					for i := 0; i < len(BeanPowers)-1; {
						
						if BeanPowers[i].LifeTime > 0 {
							
							if !BeanPowers[i].Active {
								
								switch BeanPowers[i].Type {
								case PUHealth:
									if Bean.Health+25 < 100 {
										Bean.Health += 25
									}
								case PUFirerate:
									Gun.Delay = 0.15
								case PUDamage:
									Bullets[len(Bullets)-1].Damage = 25
								}

								BeanPowers[i].Active = true
							}

							BeanPowers[i].LifeTime -= dT

						} else if BeanPowers[i].LifeTime <= 0 {
							
							switch BeanPowers[i].Type {
							case PUHealth:
							case PUFirerate:
								Gun.Delay = 0.25
							case PUDamage:
							}

							Bean.PowersNumber -= 1
							BeanPowers[i] = BeanPowers[len(BeanPowers)-1]
							BeanPowers = BeanPowers[:len(BeanPowers)-1]
						}
					}
				}

				SpawnedPowerUps = CheckForPowerUpDespawn(SpawnedPowerUps, dT, SpawnPoints)
				//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

				// Health Bar --------------------------------------------------------------------------------------------------------------------------------------------------------------
				if Bean.Health <= 0 {
					Bean.Health = 0
					Bean.Lives -= 1
				}
				if Bean2.Health <= 0 {
					Bean2.Health = 0
					Bean2.Lives -= 1
				}
				width = float32(180) * (Bean2.Health / 100)
				//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

				// Timer worker -------------------------------------------------------------------------------------------------------------------------------------------------------------
				Timer += dT
				if Timer >= 60.0 {
					Minutes += 1
					Timer = 0
				}
				if Minutes >= 60 {
					Hour += 1
					Minutes = 0
				}
				//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------
			}

		}

		// Begin Rendering ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		rl.BeginDrawing()

		if !Menu && Start {

			rl.ClearBackground(rl.GetColor(0x444444ff))

			rl.BeginMode2D(Camera)

			// Map renderer -----------------------------------------------------------------------------------------------------------------

			// Map Background ----------------------------------------------------------------------
			rl.DrawRectangle(3000, 3000, 500, 500, rl.GetColor(0xff0000ff))
			//--------------------------------------------------------------------------------------

			// Map Borders -------------------------------------------------------------------------
			rl.DrawRectanglePro(Map.Border, rl.NewVector2(0, 0), 0.0, rl.GetColor(0xAAAAAAff))
			rl.DrawRectangleLinesEx(Map.Border, 30, rl.GetColor(0x000000ff))
			//--------------------------------------------------------------------------------------

			// Platforms ---------------------------------------------------------------------------
			for _, p := range Platforms {
				if p.OneWay == true {
					rl.DrawRectangleRec(p.Rect, colorOneWay)
				} else {
					rl.DrawRectangleRec(p.Rect, colorSolid)
				}
			}
			//--------------------------------------------------------------------------------------
			//-------------------------------------------------------------------------------------------------------------------------------

			// Power Ups renderer ------------------------------------------------------------------------------------------------------------------------------------------------
			for _, p := range SpawnedPowerUps {
				rl.DrawRectangleV(p.Pos, rl.NewVector2(92, 92), rl.GetColor(0x11111100))

				switch p.Type {
				case PUHealth:
					rl.DrawTextureEx(PowerUps[PUHealth].Texture, p.Pos, 0.0, 2, rl.White)
				case PUFirerate:
					rl.DrawTextureEx(PowerUps[PUFirerate].Texture, p.Pos, 0.0, 2, rl.White)
					// rl.DrawRectangleV(p.Pos, rl.NewVector2(100, 100), rl.GetColor(0x0000ffff))
				case PUDamage:
					rl.DrawTextureEx(PowerUps[PUDamage].Texture, p.Pos, 0.0, 2, rl.White)
				}

			}
			//--------------------------------------------------------------------------------------------------------------------------------------------------------------------

			// Bullets renderer --------------------------------------------------------------------------------------------------------------------------------------------------
			for _, b := range Bullets {
				rl.DrawTexturePro(TextureBullet, rl.NewRectangle(0, 0, 18, 18), rl.NewRectangle(b.Pos.X, b.Pos.Y, 18, 18), rl.NewVector2(9, 9), b.Angle, rl.GetColor(0xffffffff))
			}
			//--------------------------------------------------------------------------------------------------------------------------------------------------------------------

			// Beans renderer ----------------------------------------------------------------------------------------------------------------
			rl.DrawRectangleV(Bean.Pos, rl.NewVector2(Bean.Width, Bean.Height), rl.GetColor(0x00ffffff))
			Bean.Pos.X += Bean.Width / 2
			Bean.Pos.Y -= 19
			rl.DrawCircleV(Bean.Pos, Bean.Radius, rl.GetColor(0x00ffffff))
			Bean.Pos.X -= Bean.Width / 2
			Bean.Pos.Y += 19

			rl.DrawRectangleV(Bean2.Pos, rl.NewVector2(Bean2.Width, Bean2.Height), rl.GetColor(0xff0000ff))
			Bean2.Pos.X += Bean2.Width / 2
			Bean2.Pos.Y -= 19
			rl.DrawCircleV(Bean2.Pos, Bean2.Radius, rl.GetColor(0xffffffff))
			Bean2.Pos.X -= Bean2.Width / 2
			Bean2.Pos.Y += 19

			if !Bean.isCrouched {
				rl.DrawTextureV(TextureStand, rl.NewVector2(Bean.Pos.X, Bean.Pos.Y-(Bean.Radius*2)), rl.GetColor(0xffffffff))
			} else {
				rl.DrawTextureV(TextureCrouch, rl.NewVector2(Bean.Pos.X, Bean.Pos.Y-(Bean.Radius*2)), rl.GetColor(0xffffffff))
			}
			//--------------------------------------------------------------------------------------------------------------------------------

			// gun renderer ------------------------------------------------------------------------------------------------------------------
			if Gun.Dir == 1 {
				rl.DrawTextureEx(TextureGun, rl.NewVector2(Gun.Pos.X-15, Gun.Pos.Y-15), Gun.Angle, 2, rl.GetColor(0xffffffff))
			} else {
				rl.DrawTextureEx(TextureGunFlipped, rl.NewVector2(Gun.Pos.X+15, Gun.Pos.Y+45), Gun.Angle, 2, rl.GetColor(0xffffffff))
			}
			if Gun2.Dir == 1 {
				rl.DrawTextureEx(TextureGun, rl.NewVector2(Gun2.Pos.X-15, Gun2.Pos.Y-15), Gun2.Angle, 2, rl.GetColor(0xffffffff))
			} else {
				rl.DrawTextureEx(TextureGunFlipped, rl.NewVector2(Gun2.Pos.X+15, Gun2.Pos.Y+45), Gun2.Angle, 2, rl.GetColor(0xffffffff))
			}
			//--------------------------------------------------------------------------------------------------------------------------------

			// Health Bar Renderer -------------------------------------------------------------------------------------------------------------------------------------------
			rl.DrawRectangleLinesEx(rl.NewRectangle(Bean.Pos.X+(Bean.Width/2)-100, Bean.Pos.Y-(Bean.Radius*2)-50, 200, 30), 10, rl.GetColor(0x000000ff))
			rl.DrawRectangleV(rl.NewVector2(Bean.Pos.X+(Bean.Width/2)-90, Bean.Pos.Y-(Bean.Radius*2)-40), rl.NewVector2((180*(Bean.Health/100)), 10), rl.GetColor(0x00ff00ff))

			rl.DrawRectangleLinesEx(rl.NewRectangle(Bean2.Pos.X+(Bean2.Width/2)-100, Bean2.Pos.Y-(Bean2.Radius*2)-50, 200, 30), 10, rl.GetColor(0x000000ff))
			rl.DrawRectangleV(rl.NewVector2(Bean2.Pos.X+(Bean2.Width/2)-90+(180-width), Bean2.Pos.Y-(Bean2.Radius*2)-40), rl.NewVector2(width, 10), rl.GetColor(0x00ff00ff))
			//----------------------------------------------------------------------------------------------------------------------------------------------------------------

			rl.EndMode2D()

			// rl.DrawText(fmt.Sprintf("SpeedX: %0.1f\nSpeedY: %0.1f\nGravity Bean: %0.1f\nGrounded: %v\nCrouched: %v\nGun Angle: %0.1f\nGun2 Angle: %0.1f",Bean.Speed.X, Bean.Speed.Y, Gravity.Bean, Bean.isGrounded, Bean.isCrouched, Gun.Angle, Gun2.Angle), 10, 10, 30, rl.GetColor(0xffffffff))

			// Info Tablet -------------------------------------------------------------------------------------------------------------------------

			// Table P1 ---------------------------------------------------------------------------------------------------------------------
			rl.DrawRectanglePro(rl.NewRectangle(-10, Screen.Y/2, 260, 500), rl.NewVector2(0, 0), 0.0, rl.Blue)
			rl.DrawRectangleLinesEx(rl.NewRectangle(-10, (Screen.Y/2)-15, 275, 500), 15, rl.Black)
			rl.DrawRectanglePro(rl.NewRectangle(-10, (Screen.Y/2)+400, 260, 20), rl.NewVector2(0, 0), 0.0, rl.GetColor(0x00000066))

			rl.DrawTextureEx(TextureAmmoContainer, rl.NewVector2(10, Screen.Y-250), 0.0, 6, rl.White)
			rl.DrawTextureEx(TextureAmmo, rl.NewVector2(100, Screen.Y-220), 0.0, 6, rl.White)
			rl.DrawText(fmt.Sprintf(": %d", Gun.Mag-Gun.Shots), 160, int32(Screen.Y)-210, 60, rl.White)

			rl.DrawTextureEx(TextureHeart, rl.NewVector2(10, Screen.Y-390), 0.0, 6, rl.White)
			rl.DrawText(fmt.Sprintf(": %d", Bean.Lives), 160, int32(Screen.Y)-360, 60, rl.White)

			// rl.DrawTextureEx(PowerTexture.Health, rl.NewVector2(10, Screen.Y/2 + 30), 0.0, 1.5, rl.White)
			// rl.DrawRectangleLinesEx(rl.NewRectangle(10,Screen.Y/2 + 100, 67, 20), 5, rl.Black)
			// rl.DrawTextureEx(PowerTexture.Health, rl.NewVector2(94, Screen.Y/2 + 30), 0.0, 1.5, rl.White)
			// rl.DrawRectangleLinesEx(rl.NewRectangle(94,Screen.Y/2 + 100, 67, 20), 5, rl.Black)
			// rl.DrawTextureEx(PowerTexture.Health, rl.NewVector2(178, Screen.Y/2 + 30), 0.0, 1.5, rl.White)
			// rl.DrawRectangleLinesEx(rl.NewRectangle(178, Screen.Y/2 + 100, 67, 20), 5, rl.Black)

			//-------------------------------------------------------------------------------------------------------------------------------

			//Table P2 ----------------------------------------------------------------------------------------------------------------------
			rl.DrawRectanglePro(rl.NewRectangle(Screen.X-260, Screen.Y/2, 260, 500), rl.NewVector2(0, 0), 0.0, rl.Red)
			rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X-265, Screen.Y/2, 280, 500), 15, rl.Black)
			rl.DrawRectanglePro(rl.NewRectangle(Screen.X-255, Screen.Y/2+400, 280, 20), rl.NewVector2(0, 0), 0.0, rl.GetColor(0x00000066))

			rl.DrawTextureEx(TextureAmmoContainer, rl.NewVector2(Screen.X-100, Screen.Y-250), 0.0, 6, rl.White)
			rl.DrawTextureEx(TextureAmmo, rl.NewVector2(Screen.X-155, Screen.Y-220), 0.0, 6, rl.White)
			rl.DrawText(fmt.Sprintf("%d :", Gun2.Mag-Gun2.Shots), int32(Screen.X)-245, int32(Screen.Y)-210, 60, rl.White)

			rl.DrawTextureEx(TextureHeart, rl.NewVector2(Screen.X-130, Screen.Y-390), 0.0, 6, rl.White)
			rl.DrawText(fmt.Sprintf("%d :", Bean2.Lives), int32(Screen.X)-230, int32(Screen.Y)-360, 60, rl.White)
			//-------------------------------------------------------------------------------------------------------------------------------

			//---------------------------------------------------------------------------------------------------------------------------------------

			// Score tab ------------------------------------------------------------------------------------------------------------------------
			rl.DrawRectanglePro(rl.NewRectangle(-10, Screen.Y-120, 740, 200), rl.NewVector2(0, 0), 0.0, rl.Blue)
			rl.DrawRectangleLinesEx(rl.NewRectangle(-10, Screen.Y-130, 750, 200), 15, rl.Black)
			rl.DrawRectanglePro(rl.NewRectangle(Screen.X-740, Screen.Y-120, 740, 200), rl.NewVector2(0, 0), 0.0, rl.Red)
			rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X-750, Screen.Y-130, 790, 200), 15, rl.Black)

			rl.DrawText(fmt.Sprintf("Score Player1: %d", ScoreP1), 30, int32(Screen.Y-100), 80, rl.GetColor(0xffffffff))
			rl.DrawText(fmt.Sprintf("%d :Score Player2", ScoreP2), int32(Screen.X-710), int32(Screen.Y-100), 80, rl.GetColor(0xffffffff))
			//-----------------------------------------------------------------------------------------------------------------------------------

			// Timer and FPS --------------------------------------------------------------------------------------------------------------------
			rl.DrawText(fmt.Sprintf("%0.0f:%0.0f:%0.01f", Hour, Minutes, Timer), int32(Screen.X/2)-44, 30, 40, rl.GetColor(0xffffffff))
			rl.DrawText(fmt.Sprintf("FPS: %v", FPS), 30, 30, 30, rl.White)
			//-----------------------------------------------------------------------------------------------------------------------------------

			if Pause {
				rl.DrawRectanglePro(rl.NewRectangle(0, 0, Screen.X, Screen.Y), rl.NewVector2(0, 0), 0.0, rl.GetColor(0x44444488))
				rl.DrawText(fmt.Sprintf("FPS: %v", FPS), 30, 30, 30, rl.White)

				rl.DrawText("PAUSED", int32(Screen.X/2)-170, int32(Screen.Y/2)-61, 80, rl.GetColor(0x000000ff))
				rl.DrawText("PAUSED", int32(Screen.X/2)-165, int32(Screen.Y/2)-65, 80, rl.Red)

				rl.DrawText("Press Space to Continue", int32(Screen.X/2)-339, int32(Screen.Y/2)+13, 50, rl.Black)
				rl.DrawText("Press Space to Continue", int32(Screen.X/2)-335, int32(Screen.Y/2)+10, 50, rl.Red)

				if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X-120, 20, 100, 100)) { //Quit button in the Pause Menu
					if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
						Quit = true
						Start = false
						Controls = false
						Menu = false
					} else if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
						rl.DrawRectanglePro(rl.NewRectangle(Screen.X-110, 30, 80, 80), rl.NewVector2(0, 0), 0.0, rl.GetColor(0xaa1111ff))
						rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X-115, 30, 90, 85), 5, rl.GetColor(0x000000ff))
						rl.DrawText("X", int32(Screen.X)-97, 32, 90, rl.Black)
					} else {
						rl.DrawRectanglePro(rl.NewRectangle(Screen.X-130, 10, 120, 120), rl.NewVector2(0, 0), 0.0, rl.GetColor(0xff1111ff))
						rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X-135, 10, 130, 125), 5, rl.GetColor(0xffffffff))
						rl.DrawText("X", int32(Screen.X)-105, 18, 120, rl.White)
					}

				} else {
					rl.DrawRectanglePro(rl.NewRectangle(Screen.X-120, 20, 100, 100), rl.NewVector2(0, 0), 0.0, rl.GetColor(0xff1111ff))
					rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X-125, 20, 110, 105), 5, rl.GetColor(0xffffffff))
					rl.DrawText("X", int32(Screen.X)-99, 27, 100, rl.White)
				}
			}

		} else if Menu { // Opens menu on start-up
			rl.ClearBackground(rl.GetColor(0x333333ff))
			rl.DrawText("Bouncing Betty 2 Player", int32(Screen.X/2)-610, int32(Screen.Y/2)-290, 100, rl.GetColor(0x000000ff))
			rl.DrawText("Bouncing Betty 2 Player", int32(Screen.X/2)-600, int32(Screen.Y/2)-300, 100, rl.GetColor(0xff2222ff))

			if !rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X/2-160, Screen.Y/2, 300, 100)) { // didnt realize until later that there is a cleaner way to do this and i dont wanna touch this now so imma just leace it like this since cleaning this wont even matter as much cause its just rendering stuff. Start Button in Main Menu.
				rl.DrawRectangle(int32(Screen.X/2)-160, int32(Screen.Y/2), 300, 100, rl.GetColor(0x00ff44ff))
				rl.DrawText("Start", int32(Screen.X/2)-122, int32(Screen.Y/2)+13, 80, rl.GetColor(0x222222ff))

			} else if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X/2-160, Screen.Y/2, 300, 100)) && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				rl.DrawRectangle(int32(Screen.X/2)-150, int32(Screen.Y/2)+10, 290, 90, rl.GetColor(0x000000ff))
				rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X/2-160, Screen.Y/2, 300, 100), 10, rl.GetColor(0xffffffff))
				rl.DrawText("Start", int32(Screen.X/2)-90, int32(Screen.Y/2)+23, 60, rl.GetColor(0xffffffff))

			} else if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X/2-160, Screen.Y/2, 300, 100)) {
				rl.DrawRectangle(int32(Screen.X/2)-170, int32(Screen.Y/2)-10, 325, 130, rl.GetColor(0x00bb44ff))
				rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X/2-180, Screen.Y/2-20, 345, 140), 10, rl.GetColor(0x000000ff))
				rl.DrawText("Start", int32(Screen.X/2)-145, int32(Screen.Y/2)+5, 100, rl.GetColor(0xffffffff))
			}

			if !rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X/2-160, Screen.Y/2+150, 300, 100)) { // Controls Button in Main Menu
				rl.DrawRectangle(int32(Screen.X/2)-160, int32(Screen.Y/2)+150, 300, 100, rl.GetColor(0x1144ffff))
				rl.DrawText("Controls", int32(Screen.X/2)-155, int32(Screen.Y/2)+168, 67, rl.GetColor(0x111111ff))

			} else if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X/2-160, Screen.Y/2+150, 300, 100)) && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				rl.DrawRectangle(int32(Screen.X/2)-150, int32(Screen.Y/2)+160, 290, 90, rl.GetColor(0x000000ff))
				rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X/2-160, Screen.Y/2+150, 300, 100), 10, rl.GetColor(0xffffffff))
				rl.DrawText("Controls", int32(Screen.X/2)-122, int32(Screen.Y/2)+175, 50, rl.GetColor(0xffffffff))

			} else if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X/2-160, Screen.Y/2+150, 300, 100)) {
				rl.DrawRectangle(int32(Screen.X/2)-170, int32(Screen.Y/2)+140, 320, 120, rl.GetColor(0x1144aaff))
				rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X/2-180, Screen.Y/2+130, 340, 140), 10, rl.GetColor(0x000000ff))
				rl.DrawText("Controls", int32(Screen.X/2)-165, int32(Screen.Y/2)+168, 70, rl.GetColor(0xffffffff))
			}

			if !rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X/2-160, Screen.Y/2+300, 300, 100)) { // Quit Button in Main Menu
				rl.DrawRectangle(int32(Screen.X/2)-160, int32(Screen.Y/2)+300, 300, 100, rl.GetColor(0xff1111ff))
				rl.DrawText("Quit", int32(Screen.X/2)-90, int32(Screen.Y/2)+313, 80, rl.GetColor(0x000000ff))

			} else if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X/2-160, Screen.Y/2+300, 300, 100)) && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				rl.DrawRectangle(int32(Screen.X/2)-150, int32(Screen.Y/2)+310, 290, 90, rl.GetColor(0x000000ff))
				rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X/2-160, Screen.Y/2+300, 300, 100), 10, rl.GetColor(0xffffffff))
				rl.DrawText("Quit", int32(Screen.X/2)-68, int32(Screen.Y/2)+323, 60, rl.GetColor(0xffffffff))
			} else if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X/2-160, Screen.Y/2+300, 300, 100)) {
				rl.DrawRectangle(int32(Screen.X/2)-170, int32(Screen.Y/2)+290, 320, 120, rl.GetColor(0xaa1111ff))
				rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X/2-180, Screen.Y/2+280, 340, 140), 10, rl.GetColor(0x000000ff))
				rl.DrawText("Quit", int32(Screen.X/2)-105, int32(Screen.Y/2)+305, 100, rl.GetColor(0xffffffff))
			}

			if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X/2-160, Screen.Y/2, 300, 100)) && rl.IsMouseButtonReleased(rl.MouseButtonLeft) { // Triggers if Start Button is pressed in the Main Menu
				Menu = false
				Start = true
				Controls = false
				Quit = false
			} else if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X/2-160, Screen.Y/2+150, 300, 100)) && rl.IsMouseButtonReleased(rl.MouseButtonLeft) { // Triggers if the Controls button is pressed in the Main menu
				Controls = true
				Menu = false
				Start = false
				Quit = false
			} else if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X/2-160, Screen.Y/2+300, 300, 100)) && rl.IsMouseButtonReleased(rl.MouseButtonLeft) { // Triggers if the Quit button is pressed in the Main Menu
				Quit = true
				Menu = false
				Start = false
				Controls = false
			}

		}

		if Controls { // Opens The Controls Menu and hides the Main Menu
			rl.ClearBackground(rl.Gray)
			rl.DrawRectangle(100, 100, 650, int32(Screen.Y)-200, rl.White)
			rl.DrawRectangleLinesEx(rl.NewRectangle(90, 90, 670, Screen.Y-180), 10, rl.Black)

			rl.DrawRectangle(int32(Screen.X)-1100, 100, 1000, int32(Screen.Y)-300, rl.White)
			rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X-1110, 90, 1010, Screen.Y-290), 10, rl.Black)

			if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.NewRectangle(Screen.X-620, Screen.Y-150, 520, 100)) { // Controls Menu Back Button Logic and Function
				if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
					Menu = true
					Start = false
					Controls = false
					Quit = false
				} else if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
					rl.DrawRectangle(int32(Screen.X)-620, int32(Screen.Y)-150, 520, 100, rl.Black)
					rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X-630, Screen.Y-160, 540, 120), 10, rl.White)
					rl.DrawText("Return to Menu", int32(Screen.X)-600, int32(Screen.Y)-130, 60, rl.White)
				} else {
					rl.DrawRectangle(int32(Screen.X)-630, int32(Screen.Y)-160, 540, 120, rl.GetColor(0xff0000ff))
					rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X-640, Screen.Y-170, 560, 140), 10, rl.Black)
					rl.DrawText("Return to Menu", int32(Screen.X)-620, int32(Screen.Y)-130, 66, rl.Black)
				}
			} else {
				rl.DrawRectangle(int32(Screen.X)-620, int32(Screen.Y)-150, 520, 100, rl.GetColor(0xff0000ff))
				rl.DrawRectangleLinesEx(rl.NewRectangle(Screen.X-630, Screen.Y-160, 540, 120), 10, rl.Black)
				rl.DrawText("Return to Menu", int32(Screen.X)-600, int32(Screen.Y)-130, 60, rl.Black)
			}
		}

		if Quit { //Triggers to close the window when true
			rl.CloseWindow()
		}

		rl.EndDrawing()

		if Bean.Health <= 0 {
			ScoreP2 += 1
			Bean = bean{Pos: rl.NewVector2(50, 3950), Width: 40, Height: 100, Radius: 20, Speed: rl.NewVector2(0, 0), MaxSpeed: 1000, Acceleration: 500, Drag: 460, Jump: 3000, CurrentPlatformIndex: -1, ignoredPlatformIndex: -1, restingOnPlatform: false, Health: 100.0, Lives: Bean.Lives}
			Bean2 = bean{Pos: rl.NewVector2(5950, 3950), Width: 40, Height: 100, Radius: 20, Speed: rl.NewVector2(0, 0), MaxSpeed: 1000, Acceleration: 500, Drag: 460, Jump: 3000, CurrentPlatformIndex: -1, ignoredPlatformIndex: -1, restingOnPlatform: false, Health: 100.0, Lives: Bean2.Lives}
			Gun = gun{Dir: 1, PrevDir: 1, Pos: rl.NewVector2(Bean.Pos.X+25, Bean.Pos.Y+20), Width: 70, Height: 30, Angle: 0.0, Mag: 15, Shots: 1, CanShoot: true, Delay: 0.15}
			Gun2 = gun{Dir: -1, PrevDir: 1, Pos: rl.NewVector2(Bean2.Pos.X-25, Bean2.Pos.Y+20), Width: 70, Height: 30, Angle: 0.0, Mag: 15, Shots: 1, CanShoot: true, Delay: 0.15}
			Bullets = []bullet{}
			Timer = 0.0
			Minutes = 0.0
			Hour = 0.0
		}
		if Bean2.Health <= 0 {
			ScoreP1 += 1
			Bean = bean{Pos: rl.NewVector2(50, 3950), Width: 40, Height: 100, Radius: 20, Speed: rl.NewVector2(0, 0), MaxSpeed: 1000, Acceleration: 500, Drag: 460, Jump: 3000, CurrentPlatformIndex: -1, ignoredPlatformIndex: -1, restingOnPlatform: false, Health: 100.0, Lives: Bean.Lives}
			Bean2 = bean{Pos: rl.NewVector2(5950, 3950), Width: 40, Height: 100, Radius: 20, Speed: rl.NewVector2(0, 0), MaxSpeed: 1000, Acceleration: 500, Drag: 460, Jump: 3000, CurrentPlatformIndex: -1, ignoredPlatformIndex: -1, restingOnPlatform: false, Health: 100.0, Lives: Bean2.Lives}
			Gun = gun{Dir: 1, PrevDir: 1, Pos: rl.NewVector2(Bean.Pos.X+25, Bean.Pos.Y+20), Width: 70, Height: 30, Angle: 0.0, Mag: 15, Shots: 1, CanShoot: true, Delay: 0.15}
			Gun2 = gun{Dir: -1, PrevDir: 1, Pos: rl.NewVector2(Bean2.Pos.X-25, Bean2.Pos.Y+20), Width: 70, Height: 30, Angle: 0.0, Mag: 15, Shots: 1, CanShoot: true, Delay: 0.15}
			Bullets = []bullet{}
			Timer = 0.0
			Minutes = 0.0
			Hour = 0.0
		}
	}
}
