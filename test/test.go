package main

import (
	"fmt"
	"github.com/SCU-SJL/menuscreen"
)

func main() {
	menu, err := menuscreen.NewMenuScreen()
	if err != nil {
		panic(err)
	}
	defer menu.Fini()
	menu.SetTitle("TEST").
		SetLine(0, "0th line").
		SetLine(1, "1st line").
		SetLine(2, "2nd line").
		SetLine(4, "4th line").
		Start()
	idx, ln, ok := menu.ChosenLine()
	if !ok {
		fmt.Println("you did not chose any items.")
		return
	}
	fmt.Printf("you've chosen %d line, content is: %s\n", idx, ln)
}
