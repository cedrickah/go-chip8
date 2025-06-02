package main

import (
	"os"

	"github.com/cedrick-ah/chip8-go/emulator"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if len(os.Args) < 2 {
		panic("Please provide a program to run")
	}

	fileName := os.Args[1]

	cpu :=  &emulator.CPU{
		Vx: [16]uint8{},
		Key:         &emulator.KeysPressed,
		Stack:     [16]uint16{},
		Pc: 0x200,
		Iv: 0,
		DelayTimer: 0,
		SoundTimer: 0,
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	speaker, speakerErr := emulator.NewSpeaker()
	if speakerErr != nil {
		panic(speakerErr)
	}
	defer speaker.Close()

	cpu.LoadSprites()
	if loadErr := cpu.LoadProgram(fileName); loadErr != nil {
		panic(loadErr)
	}

	window, windowErr := sdl.CreateWindow("Chip8-emulator", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
	640, 320, sdl.WINDOW_SHOWN)
	if windowErr != nil {
		panic(windowErr)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	renderer := emulator.NewRenderer(window, surface)
	
	running := true
	for running {
        cpu.Cycle(renderer, speaker)

        for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
            switch event.(type) {
            case *sdl.QuitEvent:
                println("Quit")
                running = false
            }
			emulator.WatchKeyDown(event)
			emulator.WatchKeyUp(event)
        }
		sdl.Delay(1000 / 60)
    }
}
