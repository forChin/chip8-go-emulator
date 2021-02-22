package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	keyMap = [16]sdl.Keycode{
		sdl.K_x, sdl.K_1, sdl.K_2, sdl.K_3,
		sdl.K_q, sdl.K_w, sdl.K_e, sdl.K_a,
		sdl.K_s, sdl.K_d, sdl.K_z, sdl.K_c,
		sdl.K_4, sdl.K_r, sdl.K_f, sdl.K_v,
	}

	fontSet = [80]byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
)

type chip8 struct {
	gameMemory     [0xfff]byte
	registers      [16]byte
	addressI       word
	gameStack      stack
	programCounter word
	keyState       [16]bool

	delayTimer       byte
	soundTimer       byte
	opcodesPerSecond time.Duration

	running bool
}

func newChip8(windowW, windowH int) *chip8 {
	ch8 := chip8{}
	ch8.addressI = 0
	ch8.programCounter = 0x200

	for i, f := range fontSet {
		ch8.gameMemory[i] = f
	}

	return &ch8
}

func (c *chip8) loadROM(romPath string) error {
	gameData, err := ioutil.ReadFile(romPath)
	if err != nil {
		return err
	}

	if len(gameData) > 0xfff-0x200 {
		return errors.New("Invalid size of ROM") // change err msg
	}

	copy(c.gameMemory[0x200:], gameData)

	return nil
}

func (c *chip8) run() {
	go c.startTimers()
	go c.handleKeys()
	go c.startRender()

	ticker := time.NewTicker(c.opcodesPerSecond * time.Millisecond)
	for range ticker.C {
		next := getNextOpcode()
		fmt.Printf("0x%x\n", next)
		executeOpcode(next)
	}
}

func (c *chip8) startTimers() {
	ticker := time.NewTicker(16667 * time.Microsecond) // 60 Hz

	for range ticker.C {
		if c.running {
			if c.delayTimer > 0 {
				c.delayTimer--
			}
			if c.soundTimer > 0 {
				c.soundTimer--
			}

			if c.soundTimer > 0 {
				fmt.Println("BEEP") // make real sound
			}
		}
	}
}

func (c *chip8) handleKeys() {
	for {
		event := sdl.PollEvent()
		if event != nil {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				c.running = false
				return
			case *sdl.KeyboardEvent:
				c.updateKey(t.Keysym.Sym, t.State == sdl.PRESSED)
			}
		}
	}
}

func (c *chip8) updateKey(keyCode sdl.Keycode, state bool) {
	for i, k := range keyMap {
		if keyCode == k {
			c.keyState[i] = state
		}
	}
}
