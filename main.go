package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	runtime.LockOSThread()
	rand.Seed(time.Now().UnixNano())

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	cfg, err := newConfig()
	if err != nil {
		log.Fatal(err)
	}

	ch8 := newChip8(cfg)

	if err := ch8.loadROM(cfg.ROMPath); err != nil {
		log.Fatal(err)
	}

	ch8.run()
}
