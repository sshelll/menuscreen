package main

import "github.com/sshelll/menuscreen"

func main() {
	root := menuscreen.NewSimpleWorkflow("root", []string{"1st line", "2nd line", "3rd line"})

	r1Node := menuscreen.NewSimpleWorkflow("r-1 node", []string{"1st line", "2nd line", "3rd line"})
	r1Node.SetCallback(func() {
		println("r1Node callback")
	})

	root.SetNext(0, r1Node)
	root.SetNext(1, menuscreen.NewSimpleWorkflow("r-2 node", []string{"1st line", "2nd line", "3rd line"}))
	root.SetNext(2, menuscreen.NewSimpleWorkflow("r-3 node", []string{"1st line", "2nd line", "3rd line"}))
	menuscreen.RunWorkflow(root)
}
