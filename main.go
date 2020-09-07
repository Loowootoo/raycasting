package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"raycasting/input"
	"raycasting/tex"

	"raycasting/vec3"
)

var worldMap = [24][24]int{
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 7, 7, 7, 7, 7, 7, 7, 7},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 4, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 7, 7, 0, 7, 7, 7, 7, 7},
	{4, 0, 5, 0, 0, 0, 0, 5, 0, 5, 0, 5, 0, 5, 0, 5, 7, 0, 0, 0, 7, 7, 7, 1},
	{4, 0, 6, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 0, 0, 0, 0, 8},
	{4, 0, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 7, 7, 1},
	{4, 0, 8, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 0, 0, 0, 0, 8},
	{4, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 0, 7, 7, 7, 1},
	{4, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 0, 5, 5, 5, 5, 7, 7, 7, 7, 7, 7, 7, 1},
	{6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
	{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
	{6, 6, 6, 6, 6, 6, 0, 6, 6, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
	{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 6, 0, 6, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 2, 0, 0, 5, 0, 0, 2, 0, 0, 0, 2},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2},
	{4, 0, 6, 0, 6, 0, 0, 0, 0, 4, 6, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 2},
	{4, 0, 0, 5, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2},
	{4, 0, 6, 0, 6, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 5, 0, 0, 2, 0, 0, 0, 2},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2},
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 1, 1, 1, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3},
}

const (
	winWidth  int = 512
	winHeight int = 256
	texSize   int = 64
)

type Frame struct {
	LineHigh  int
	DrawStart int
	DrawEnd   int
	Buf       []byte
	Texture   [9]*tex.Texture
}

func NewFrame() *Frame {
	fr := new(Frame)
	buf := make([]byte, winWidth*winHeight*4)
	fr.Buf = buf
	fr.Texture[0] = tex.LoadFromFile("assets/BGTile3.png")
	fr.Texture[1] = tex.LoadFromFile("assets/BGTile4.png")
	fr.Texture[2] = tex.LoadFromFile("assets/BGTile5.png")
	fr.Texture[3] = tex.LoadFromFile("assets/BGTile6.png")
	fr.Texture[4] = tex.LoadFromFile("assets/Box.png")
	fr.Texture[5] = tex.LoadFromFile("assets/DoorLocked.png")
	fr.Texture[6] = tex.LoadFromFile("assets/DoorOpen.png")
	fr.Texture[7] = tex.LoadFromFile("assets/Acid2.png")
	fr.Texture[8] = tex.LoadFromFile("assets/Tile5.png")
	return fr
}

type RayCasting struct {
	Pos          vec3.Vec2
	Dir          vec3.Vec2
	CameraPlan   vec3.Vec2
	RayDir       vec3.Vec2
	SideDist     vec3.Vec2
	DeltaDist    vec3.Vec2
	PrepWallDist float64
	CameraX      float64
	MapX, MapY   int
	StepX, StepY int
	Hit          int
	Side         int
	ScrFrame     *Frame
}

func NewRacCasting() *RayCasting {
	rc := new(RayCasting)
	rc.Pos.X = 10
	rc.Pos.Y = 10
	rc.Dir.X = -1
	rc.Dir.Y = 0
	rc.CameraPlan.X = 0
	rc.CameraPlan.Y = 0.66
	rc.ScrFrame = NewFrame()
	return rc
}
func (rc *RayCasting) ClearBuf() {
	for i := range rc.ScrFrame.Buf {
		rc.ScrFrame.Buf[i] = 0
	}
}

func (rc *RayCasting) GetInput() {
	//speed modifiers
	frameTime := 1 / ebiten.CurrentTPS()
	moveSpeed := frameTime * 6.0 //the constant value is in squares/second
	rotSpeed := frameTime * 3.0  //the constant value is in radians/second
	//move forward if no wall in front of you
	if input.IsWKey() {
		if worldMap[int(rc.Pos.X+rc.Dir.X*moveSpeed)][int(rc.Pos.Y)] == 0 {
			rc.Pos.X += rc.Dir.X * moveSpeed
		}
		if worldMap[int(rc.Pos.X)][int(rc.Pos.Y+rc.Dir.Y*moveSpeed)] == 0 {
			rc.Pos.Y += rc.Dir.Y * moveSpeed
		}
	}
	//move backwards if no wall behind you
	if input.IsSKey() {
		if worldMap[int(rc.Pos.X-rc.Dir.X*moveSpeed)][int(rc.Pos.Y)] == 0 {
			rc.Pos.X -= rc.Dir.X * moveSpeed
		}
		if worldMap[int(rc.Pos.X)][int(rc.Pos.Y-rc.Dir.Y*moveSpeed)] == 0 {
			rc.Pos.Y -= rc.Dir.Y * moveSpeed
		}
	}
	//rotate to the right
	if input.IsDKey() {
		//both camera direction and camera plane must be rotated
		oldDirX := rc.Dir.X
		rc.Dir.X = rc.Dir.X*math.Cos(-rotSpeed) - rc.Dir.Y*math.Sin(-rotSpeed)
		rc.Dir.Y = oldDirX*math.Sin(-rotSpeed) + rc.Dir.Y*math.Cos(-rotSpeed)
		oldPlaneX := rc.CameraPlan.X
		rc.CameraPlan.X = rc.CameraPlan.X*math.Cos(-rotSpeed) - rc.CameraPlan.Y*math.Sin(-rotSpeed)
		rc.CameraPlan.Y = oldPlaneX*math.Sin(-rotSpeed) + rc.CameraPlan.Y*math.Cos(-rotSpeed)
	}
	//rotate to the left
	if input.IsAKey() {
		//both camera direction and camera plane must be rotated
		oldDirX := rc.Dir.X
		rc.Dir.X = rc.Dir.X*math.Cos(rotSpeed) - rc.Dir.Y*math.Sin(rotSpeed)
		rc.Dir.Y = oldDirX*math.Sin(rotSpeed) + rc.Dir.Y*math.Cos(rotSpeed)
		oldPlaneX := rc.CameraPlan.X
		rc.CameraPlan.X = rc.CameraPlan.X*math.Cos(rotSpeed) - rc.CameraPlan.Y*math.Sin(rotSpeed)
		rc.CameraPlan.Y = oldPlaneX*math.Sin(rotSpeed) + rc.CameraPlan.Y*math.Cos(rotSpeed)
	}
}

func (rc *RayCasting) Gen() {
	for x := 0; x < winWidth; x++ {
		// 計算位置和方向
		rc.CameraX = 2*float64(x)/float64(winWidth) - 1
		rc.RayDir.X = rc.Dir.X + rc.CameraPlan.X*rc.CameraX
		rc.RayDir.Y = rc.Dir.Y + rc.CameraPlan.Y*rc.CameraX
		rc.MapX = int(rc.Pos.X)
		rc.MapY = int(rc.Pos.Y)

		rc.DeltaDist.X = math.Abs(1 / rc.RayDir.X)
		rc.DeltaDist.Y = math.Abs(1 / rc.RayDir.Y)

		if rc.RayDir.X < 0 {
			rc.StepX = -1
			rc.SideDist.X = (rc.Pos.X - float64(rc.MapX)) * rc.DeltaDist.X
		} else {
			rc.StepX = 1
			rc.SideDist.X = (float64(rc.MapX) + 1.0 - rc.Pos.X) * rc.DeltaDist.X
		}
		if rc.RayDir.Y < 0 {
			rc.StepY = -1
			rc.SideDist.Y = (rc.Pos.Y - float64(rc.MapY)) * rc.DeltaDist.Y
		} else {
			rc.StepY = 1
			rc.SideDist.Y = (float64(rc.MapY) + 1.0 - rc.Pos.Y) * rc.DeltaDist.Y
		}
		rc.Hit = 0
		for rc.Hit == 0 {
			// jump to next map square, OR in x-direction, OR in y-direction
			if rc.SideDist.X < rc.SideDist.Y {
				rc.SideDist.X += rc.DeltaDist.X
				rc.MapX += rc.StepX
				rc.Side = 0
			} else {
				rc.SideDist.Y += rc.DeltaDist.Y
				rc.MapY += rc.StepY
				rc.Side = 1
			}
			if worldMap[rc.MapX][rc.MapY] > 0 {
				rc.Hit = 1
			}
		}
		var wallX float64
		//Calculate distance projected on camera direction (Euclidean distance will give fisheye effect!)
		if rc.Side == 0 {
			rc.PrepWallDist = (float64(rc.MapX) - rc.Pos.X + (1-float64(rc.StepX))/2) / rc.RayDir.X
			wallX = rc.Pos.Y + rc.PrepWallDist*rc.RayDir.Y
		} else {
			rc.PrepWallDist = (float64(rc.MapY) - rc.Pos.Y + (1-float64(rc.StepY))/2) / rc.RayDir.Y
			wallX = rc.Pos.X + rc.PrepWallDist*rc.RayDir.X
		}
		/*
			if x == winWidth/2 {
				wallDistance := rc.PrepWallDist
			}
		*/
		wallX -= math.Floor(wallX)
		texX := int(wallX * float64(texSize))
		//Calculate height of line to draw on screen
		rc.ScrFrame.LineHigh = int(float64(winHeight) / rc.PrepWallDist)
		if rc.ScrFrame.LineHigh < 1 {
			rc.ScrFrame.LineHigh = 1
		}
		//calculate lowest and highest pixel to fill in current stripe
		rc.ScrFrame.DrawStart = -rc.ScrFrame.LineHigh/2 + winHeight/2

		if rc.ScrFrame.DrawStart < 0 {
			rc.ScrFrame.DrawStart = 0
		}

		rc.ScrFrame.DrawEnd = rc.ScrFrame.LineHigh/2 + winHeight/2

		if rc.ScrFrame.DrawEnd >= winHeight {
			rc.ScrFrame.DrawEnd = winHeight - 1
		}
		if rc.Side == 0 && rc.RayDir.X > 0 {
			texX = texSize - texX - 1
		}
		if rc.Side == 1 && rc.RayDir.Y < 0 {
			texX = texSize - texX - 1
		}

		texNum := worldMap[rc.MapX][rc.MapY] - 1
		ww := winWidth * 4
		st := (rc.ScrFrame.DrawStart*winWidth + x) * 4
		for y := rc.ScrFrame.DrawStart; y < rc.ScrFrame.DrawEnd; y++ {
			//d := y*256 - winHeight*128 + rc.ScrFrame.LineHigh*128 //256 and 128 factors to avoid floats
			d := (y << 8) - (winHeight << 7) + (rc.ScrFrame.LineHigh << 7)
			// TODO: avoid the division to speed this up
			texY := ((d * 64) / rc.ScrFrame.LineHigh) / 256
			index := (64*texY + texX) * 4
			rc.ScrFrame.Buf[st] = rc.ScrFrame.Texture[texNum].Pixels[index]
			rc.ScrFrame.Buf[st+1] = rc.ScrFrame.Texture[texNum].Pixels[index+1]
			rc.ScrFrame.Buf[st+2] = rc.ScrFrame.Texture[texNum].Pixels[index+2]
			//			if rc.Side == 1 {
			//				rc.ScrFrame.Buf[st+3] = rc.ScrFrame.Texture[texNum].Pixels[index+3] - 127
			//			} else {
			cr := rc.ScrFrame.Texture[texNum].Pixels[index+3]
			cr = byte(float64(cr) / rc.PrepWallDist)
			rc.ScrFrame.Buf[st+3] = cr
			//			}
			st += ww
		}
		/////////////////////////////////////////////////////////////
		var floorWall vec3.Vec2

		if rc.Side == 0 && rc.RayDir.X > 0 {
			floorWall.X = float64(rc.MapX)
			floorWall.Y = float64(rc.MapY) + wallX
		} else if rc.Side == 0 && rc.RayDir.X < 0 {
			floorWall.X = float64(rc.MapX) + 1.0
			floorWall.Y = float64(rc.MapY) + wallX
		} else if rc.Side == 1 && rc.RayDir.Y > 0 {
			floorWall.X = float64(rc.MapX) + wallX
			floorWall.Y = float64(rc.MapY)
		} else {
			floorWall.X = float64(rc.MapX) + wallX
			floorWall.Y = float64(rc.MapY) + 1.0
		}

		distWall, distPlayer := rc.PrepWallDist, 0.0

		ww = winWidth * 4
		st = ((rc.ScrFrame.DrawEnd+1)*winWidth + x) * 4
		st1 := ((winHeight-rc.ScrFrame.DrawEnd+1)*winWidth + x) * 4
		var currentFloor vec3.Vec2
		for y := rc.ScrFrame.DrawEnd + 1; y < winHeight; y++ {
			currentDist := float64(winHeight) / (2.0*float64(y) - float64(winHeight))
			weight := (currentDist - distPlayer) / (distWall - distPlayer)
			currentFloor.X = weight*floorWall.X + (1.0-weight)*rc.Pos.X
			currentFloor.Y = weight*floorWall.Y + (1.0-weight)*rc.Pos.Y
			fx := int(currentFloor.X*float64(texSize)) % texSize
			fy := int(currentFloor.Y*float64(texSize)) % texSize
			idx := (fy*64 + fx) * 4
			rc.ScrFrame.Buf[st] = rc.ScrFrame.Texture[7].Pixels[idx]
			rc.ScrFrame.Buf[st+1] = rc.ScrFrame.Texture[7].Pixels[idx+1]
			rc.ScrFrame.Buf[st+2] = rc.ScrFrame.Texture[7].Pixels[idx+2]
			rc.ScrFrame.Buf[st+3] = rc.ScrFrame.Texture[7].Pixels[idx+3] - 127
			rc.ScrFrame.Buf[st1] = rc.ScrFrame.Texture[8].Pixels[idx]
			rc.ScrFrame.Buf[st1+1] = rc.ScrFrame.Texture[8].Pixels[idx+1]
			rc.ScrFrame.Buf[st1+2] = rc.ScrFrame.Texture[8].Pixels[idx+2]
			rc.ScrFrame.Buf[st1+3] = rc.ScrFrame.Texture[8].Pixels[idx+3]
			st += ww
			st1 -= ww
		}
	}
}

func (rc *RayCasting) Run() {
	rc.ClearBuf()
	rc.Gen()
	rc.GetInput()
}

func update(screen *ebiten.Image) error {
	rayCasting.Run()
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.ReplacePixels(rayCasting.ScrFrame.Buf)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	return nil
}

var rayCasting *RayCasting

func main() {
	rayCasting = NewRacCasting()
	err := ebiten.Run(update, winWidth, winHeight, 1, "Raycasting")
	if err != nil {
		panic(err)
	}
}
