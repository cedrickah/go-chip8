package emulator

import (
	"github.com/veandco/go-sdl2/sdl"
)

var KeysPressed = [16]uint8{}

func WatchKeyDown(event interface{}) {
		if e, ok := event.(*sdl.KeyboardEvent); ok {
			key := e.Keysym.Sym
			if e.Type == sdl.KEYDOWN {
				switch key {
				case sdl.K_1:
					KeysPressed[0x1] = 1
				case sdl.K_2:
					KeysPressed[0x2] = 1
				case sdl.K_3:
					KeysPressed[0x3] = 1
				case sdl.K_4:
					KeysPressed[0xc] = 1
				case sdl.K_q:
					KeysPressed[0x4] = 1
				case sdl.K_w:
					KeysPressed[0x5] = 1
				case sdl.K_e:
					KeysPressed[0x6] = 1
				case sdl.K_r:
					KeysPressed[0xD] = 1
				case sdl.K_a:
					KeysPressed[0x7] = 1
				case sdl.K_s:
					KeysPressed[0x8] = 1
				case sdl.K_d:
					KeysPressed[0x9] = 1
				case sdl.K_f:
					KeysPressed[0xE] = 1
				case sdl.K_z:
					KeysPressed[0xA] = 1
				case sdl.K_x:
					KeysPressed[0x0] = 1
				case sdl.K_c:
					KeysPressed[0xB] = 1
				case sdl.K_v:
					KeysPressed[0xF] = 1
				}
			} 
		}
}

func WatchKeyUp(event interface{}) {
		if e, ok := event.(*sdl.KeyboardEvent); ok {
			key := e.Keysym.Sym
			if e.Type == sdl.KEYUP {
				switch key {
				case sdl.K_1:
					KeysPressed[0x1] = 0
				case sdl.K_2:
					KeysPressed[0x2] = 0
				case sdl.K_3:
					KeysPressed[0x3] = 0
				case sdl.K_4:
					KeysPressed[0xc] = 0
				case sdl.K_q:
					KeysPressed[0x4] = 0
				case sdl.K_w:
					KeysPressed[0x5] = 0
				case sdl.K_e:
					KeysPressed[0x6] = 0
				case sdl.K_r:
					KeysPressed[0xD] = 0
				case sdl.K_a:
					KeysPressed[0x7] = 0
				case sdl.K_s:
					KeysPressed[0x8] = 0
				case sdl.K_d:
					KeysPressed[0x9] = 0
				case sdl.K_f:
					KeysPressed[0xE] = 0
				case sdl.K_z:
					KeysPressed[0xA] = 0
				case sdl.K_x:
					KeysPressed[0x0] = 0
				case sdl.K_c:
					KeysPressed[0xB] = 0
				case sdl.K_v:
					KeysPressed[0xF] = 0
				}
			} 
		}
}
