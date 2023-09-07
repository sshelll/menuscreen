package main

import (
	"fmt"

	"github.com/sshelll/menuscreen"
)

func main() {
	testItem()
}

func testLine() {
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

func testItem() {
	menu, err := menuscreen.NewMenuScreen()
	if err != nil {
		panic(err)
	}
	defer menu.Fini()
	menu.SetTitle("TEST").
		AppendItems(
			&menuscreen.MenuItem{Content: "0th line", Item: 0},
			&menuscreen.MenuItem{Content: "1st line", Item: "1"},
			&menuscreen.MenuItem{Content: "2nd line", Item: 2},
			&menuscreen.MenuItem{Content: "3rd line", Item: "3"},
		).
		Start()
	idx, item, ok := menu.ChosenItem()
	if !ok {
		fmt.Println("you did not chose any items.")
		return
	}
	fmt.Printf("you've chosen %d line, content is: %s, item is: %v\n", idx, item.Content, item.Item)
}
