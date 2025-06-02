package emulator

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/cedrick-ah/chip8-go/debugger"
)

var debug = false

var sprites = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

type CPU struct {
	Vx     [16]uint8   // cpu registers V0-VF
	Key    *[16]uint8   // input key
	Stack  [16]uint16  // program counter stack
	Oc uint16 // current opcode
	Pc uint16 // program counter
	Sp uint16 // stack pointer
	Iv uint16 // index register
	DelayTimer uint8 // delay timer
	SoundTimer uint8 // sound timer
}

var memory [4096]uint8

func (c *CPU) LoadSprites() {
	copy(memory[:len(sprites)], sprites)
}

func (c *CPU) ExecuteInstruction(r *Renderer) {
	c.Oc = (uint16(memory[c.Pc]) << 8) | uint16(memory[c.Pc+1])

	switch c.Oc & 0xF000 {
	case 0x0000:
		switch c.Oc & 0x000F {
		case 0x0000: // 0x00E0 Clears screen
			r.Clear()
			c.Pc = c.Pc + 2
		case 0x000E: // 0x00EE Returns from a subroutine
			c.Sp = c.Sp - 1
			c.Pc = c.Stack[c.Sp]
			c.Pc = c.Pc + 2
		default:
			fmt.Printf("Invalid opcode %X\n", c.Oc)
		}
	case 0x1000: // 0x1NNN Jump to address NNN
		c.Pc = c.Oc & 0x0FFF
	case 0x2000: // 0x2NNN Calls subroutine at NNN
		c.Stack[c.Sp] = c.Pc // store current program counter
		c.Sp = c.Sp + 1      // increment stack pointer
		c.Pc = c.Oc & 0x0FFF // jump to address NNN
	case 0x3000: // 0x3XNN Skips the next instruction if VX equals NN
		if uint16(c.Vx[(c.Oc&0x0F00)>>8]) == c.Oc&0x00FF {
			c.Pc = c.Pc + 4
		} else {
			c.Pc = c.Pc + 2
		}
	case 0x4000: // 0x4XNN Skips the next instruction if VX doesn't equal NN
		if uint16(c.Vx[(c.Oc&0x0F00)>>8]) != c.Oc&0x00FF {
			c.Pc = c.Pc + 4
		} else {
			c.Pc = c.Pc + 2
		}
	case 0x5000: // 0x5XY0 Skips the next instruction if VX equals VY
		if c.Vx[(c.Oc&0x0F00)>>8] == c.Vx[(c.Oc&0x00F0)>>4] {
			c.Pc = c.Pc + 4
		} else {
			c.Pc = c.Pc + 2
		}
	case 0x6000: // 0x6XNN Sets VX to NN
		c.Vx[(c.Oc&0x0F00)>>8] = uint8(c.Oc & 0x00FF)
		c.Pc = c.Pc + 2
	case 0x7000: // 0x7XNN Adds NN to VX
		c.Vx[(c.Oc&0x0F00)>>8] = c.Vx[(c.Oc&0x0F00)>>8] + uint8(c.Oc&0x00FF)
		c.Pc = c.Pc + 2
	case 0x8000:
		switch c.Oc & 0x000F {
		case 0x0000: // 0x8XY0 Sets VX to the value of VY
			c.Vx[(c.Oc&0x0F00)>>8] = c.Vx[(c.Oc&0x00F0)>>4]
			c.Pc = c.Pc + 2
		case 0x0001: // 0x8XY1 Sets VX to VX or VY
			c.Vx[(c.Oc&0x0F00)>>8] = c.Vx[(c.Oc&0x0F00)>>8] | c.Vx[(c.Oc&0x00F0)>>4]
			c.Pc = c.Pc + 2
		case 0x0002: // 0x8XY2 Sets VX to VX and VY
			c.Vx[(c.Oc&0x0F00)>>8] = c.Vx[(c.Oc&0x0F00)>>8] & c.Vx[(c.Oc&0x00F0)>>4]
			c.Pc = c.Pc + 2
		case 0x0003: // 0x8XY3 Sets VX to VX xor VY
			c.Vx[(c.Oc&0x0F00)>>8] = c.Vx[(c.Oc&0x0F00)>>8] ^ c.Vx[(c.Oc&0x00F0)>>4]
			c.Pc = c.Pc + 2
		case 0x0004: // 0x8XY4 Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't
			if c.Vx[(c.Oc&0x00F0)>>4] > 0xFF-c.Vx[(c.Oc&0x0F00)>>8] {
				c.Vx[0xF] = 1
			} else {
				c.Vx[0xF] = 0
			}
			c.Vx[(c.Oc&0x0F00)>>8] = c.Vx[(c.Oc&0x0F00)>>8] + c.Vx[(c.Oc&0x00F0)>>4]
			c.Pc = c.Pc + 2
		case 0x0005: // 0x8XY5 VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			if c.Vx[(c.Oc&0x00F0)>>4] > c.Vx[(c.Oc&0x0F00)>>8] {
				c.Vx[0xF] = 0
			} else {
				c.Vx[0xF] = 1
			}
			c.Vx[(c.Oc&0x0F00)>>8] = c.Vx[(c.Oc&0x0F00)>>8] - c.Vx[(c.Oc&0x00F0)>>4]
			c.Pc = c.Pc + 2
		case 0x0006: // 0x8XY6 Shifts VY right by one and stores the result to VX (VY remains unchanged). VF is set to the value of the least significant bit of VY before the shift
			c.Vx[0xF] = c.Vx[(c.Oc&0x0F00)>>8] & 0x1
			c.Vx[(c.Oc&0x0F00)>>8] = c.Vx[(c.Oc&0x0F00)>>8] >> 1
			c.Pc = c.Pc + 2
		case 0x0007: // 0x8XY7 Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			if c.Vx[(c.Oc&0x0F00)>>8] > c.Vx[(c.Oc&0x00F0)>>4] {
				c.Vx[0xF] = 0
			} else {
				c.Vx[0xF] = 1
			}
			c.Vx[(c.Oc&0x0F00)>>8] = c.Vx[(c.Oc&0x00F0)>>4] - c.Vx[(c.Oc&0x0F00)>>8]
			c.Pc = c.Pc + 2
		case 0x000E: // 0x8XYE Shifts VY left by one and copies the result to VX. VF is set to the value of the most significant bit of VY before the shift
			c.Vx[0xF] = c.Vx[(c.Oc&0x0F00)>>8] >> 7
			c.Vx[(c.Oc&0x0F00)>>8] = c.Vx[(c.Oc&0x0F00)>>8] << 1
			c.Pc = c.Pc + 2
		default:
			fmt.Printf("Invalid opcode %X\n", c.Oc)
		}
	case 0x9000: // 9XY0 Skips the next instruction if VX doesn't equal VY
		if c.Vx[(c.Oc&0x0F00)>>8] != c.Vx[(c.Oc&0x00F0)>>4] {
			c.Pc = c.Pc + 4
		} else {
			c.Pc = c.Pc + 2
		}
	case 0xA000: // 0xANNN Sets I to the address NNN
		c.Iv = c.Oc & 0x0FFF
		c.Pc = c.Pc + 2
	case 0xB000: // 0xBNNN Jumps to the address NNN plus V0
		c.Pc = (c.Oc & 0x0FFF) + uint16(c.Vx[0x0])
	case 0xC000: // 0xCXNN Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN
		c.Vx[(c.Oc&0x0F00)>>8] = uint8(rand.Intn(256)) & uint8(c.Oc&0x00FF)
		c.Pc = c.Pc + 2
	case 0xD000: // 0xDXYN Draws a sprite at coordinate (VX, VY)
		x := c.Vx[(c.Oc&0x0F00)>>8]
		y := c.Vx[(c.Oc&0x00F0)>>4]
		h := c.Oc & 0x000F
		c.Vx[0xF] = 0
		var j uint16 = 0
		var i uint16 = 0
		for j = 0; j < h; j++ {
			pixel := memory[c.Iv+j]
			for i = 0; i < 8; i++ {
				if (pixel & (0x80 >> i)) != 0 {
					if (r.SetPixel(x + uint8(i), y + uint8(j))) {
						c.Vx[0xF] = 1
					}
				}
			}
		}
		c.Pc = c.Pc + 2
	case 0xE000:
		switch c.Oc & 0x00FF {
		case 0x009E: // 0xEX9E Skips the next instruction if the key stored in VX is pressed
			if c.Key[c.Vx[(c.Oc&0x0F00)>>8]] == 1 {
				c.Pc = c.Pc + 4
			} else {
				c.Pc = c.Pc + 2
			}
		case 0x00A1: // 0xEXA1 Skips the next instruction if the key stored in VX isn't pressed
			if c.Key[c.Vx[(c.Oc&0x0F00)>>8]] == 0 {
				c.Pc = c.Pc + 4
			} else {
				c.Pc = c.Pc + 2
			}
		default:
			fmt.Printf("Invalid opcode %X\n", c.Oc)
		}
	case 0xF000:
		switch c.Oc & 0x00FF {
		case 0x0007: // 0xFX07 Sets VX to the value of the delay timer
			c.Vx[(c.Oc&0x0F00)>>8] = c.DelayTimer
			c.Pc = c.Pc + 2
		case 0x000A: // 0xFX0A A key press is awaited, and then stored in VX
			pressed := false
			for i := 0; i < len(c.Key); i++ {
				if c.Key[i] != 0 {
					c.Vx[(c.Oc&0x0F00)>>8] = uint8(i)
					pressed = true
				}
			}
			if !pressed {
				return
			}
			c.Pc = c.Pc + 2
		case 0x0015: // 0xFX15 Sets the delay timer to VX
			c.DelayTimer = c.Vx[(c.Oc&0x0F00)>>8]
			c.Pc = c.Pc + 2
		case 0x0018: // 0xFX18 Sets the sound timer to VX
			c.SoundTimer = c.Vx[(c.Oc&0x0F00)>>8]
			c.Pc = c.Pc + 2
		case 0x001E: // 0xFX1E Adds VX to I
			if c.Iv+uint16(c.Vx[(c.Oc&0x0F00)>>8]) > 0xFFF {
				c.Vx[0xF] = 1
			} else {
				c.Vx[0xF] = 0
			}
			c.Iv = c.Iv + uint16(c.Vx[(c.Oc&0x0F00)>>8])
			c.Pc = c.Pc + 2
		case 0x0029: // 0xFX29 Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font
			c.Iv = uint16(c.Vx[(c.Oc&0x0F00)>>8]) * 0x5
			c.Pc = c.Pc + 2
		case 0x0033: // 0xFX33 Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2
			memory[c.Iv] = c.Vx[(c.Oc&0x0F00)>>8] / 100
			memory[c.Iv+1] = (c.Vx[(c.Oc&0x0F00)>>8] / 10) % 10
			memory[c.Iv+2] = (c.Vx[(c.Oc&0x0F00)>>8] % 100) / 10
			c.Pc = c.Pc + 2
		case 0x0055: // 0xFX55 Stores V0 to VX (including VX) in memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int((c.Oc&0x0F00)>>8)+1; i++ {
				memory[uint16(i)+c.Iv] = c.Vx[i]
			}
			c.Iv = ((c.Oc & 0x0F00) >> 8) + 1
			c.Pc = c.Pc + 2
		case 0x0065: // 0xFX65 Fills V0 to VX (including VX) with values from memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int((c.Oc&0x0F00)>>8)+1; i++ {
				c.Vx[i] = memory[c.Iv+uint16(i)]
			}
			c.Iv = ((c.Oc & 0x0F00) >> 8) + 1
			c.Pc = c.Pc + 2
		default:
			fmt.Printf("Invalid opcode %X\n", c.Oc)
		}
	default:
		fmt.Printf("Invalid opcode %X\n", c.Oc)
	}
}

func(c *CPU) Cycle(r *Renderer, s *Speaker) {
		c.ExecuteInstruction(r);
		c.updateTimers();
		r.Render()
		
		if c.SoundTimer > 0 {
			if c.SoundTimer == 1 {
				s.Play();
			}
		}
}

func (c *CPU) updateTimers() {
	if (c.DelayTimer > 0) {
		c.DelayTimer -= 1;
	}

	if (c.SoundTimer > 0) {
		c.SoundTimer -= 1;
	}
}

func (c *CPU) LoadProgram(fileName string) error {
	file, fileErr := os.OpenFile("roms/"+fileName, os.O_RDONLY, 0777)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	fStat, fStatErr := file.Stat()
	if fStatErr != nil {
		return fStatErr
	}
	if int64(len(memory)-512) < fStat.Size() {
		return fmt.Errorf("program size bigger than memory")
	}

	buffer := make([]byte, fStat.Size())
	if _, readErr := file.Read(buffer); readErr != nil {
		return readErr
	}

	for i := 0; i < len(buffer); i++ {
		memory[i+512] = buffer[i]
	}

	if debug {
		file, err := os.Create("log.txt")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return err
		}
		defer file.Close()
		debugger.Print_mem(file, memory[:], 0, 4096)
	}
	
	return nil
}
