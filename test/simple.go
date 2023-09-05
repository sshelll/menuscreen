package main

import (
	"fmt"

	"github.com/sshelll/menuscreen"
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
		SetLine(5, "第五行a").
		SetLine(6, "第六行ba").
		SetLine(7, "7TH LINE").
		Start()
	idx, ln, ok := menu.ChosenLine()
	if !ok {
		fmt.Println("you did not chose any items.")
		return
	}
	fmt.Printf("you've chosen %d line, content is: %s\n", idx, ln)
}
