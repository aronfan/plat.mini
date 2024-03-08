package test

import "fmt"

type Pet struct {
	PttID int    `redis:"pttid"`
	Level int    `redis:"level"`
	Name  string `redis:"name"`
}

func (pet *Pet) Init() {
	fmt.Println("init")
}

func (pet *Pet) LevelUp() int {
	pet.Level += 1
	return pet.Level
}
