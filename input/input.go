package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 5
		interval = 1
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func IsWKey() bool {
	return repeatingKeyPressed(ebiten.KeyW)
}
func IsSKey() bool {
	return repeatingKeyPressed(ebiten.KeyS)
}
func IsAKey() bool {
	return repeatingKeyPressed(ebiten.KeyA)
}
func IsDKey() bool {
	return repeatingKeyPressed(ebiten.KeyD)
}
func IsSpaceKey() bool {
	return repeatingKeyPressed(ebiten.KeySpace)
}
