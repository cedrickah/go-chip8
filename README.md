# go-chip8

Implementing CHIP-8 in Golang to learn more about emulation and system programming.

## Installing

- Install [`sdl2`](https://www.libsdl.org/) on your machine

- Install go-sdl2

```
go get -u github.com/veandco/go-sdl2/sdl
```

## Running

Use one of the ROMs from the `roms` directory. They end with `ch8`.

```
go run main.go <ROM>
```

## Key Bindings

```
Chip8 keypad         Keyboard mapping
1 | 2 | 3 | C        1 | 2 | 3 | 4
4 | 5 | 6 | D   =>   Q | W | E | R
7 | 8 | 9 | E   =>   A | S | D | F
A | 0 | B | F        Z | X | C | V
```

## Sources

- [How to write an emulator chip-8 interpreter](http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/)
- [Cowgod's Chip-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)
- [Chip-8 opcode table](https://en.wikipedia.org/wiki/CHIP-8)
- [go-chip8] (https://github.com/skatiyar/go-chip8)
