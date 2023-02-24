package main

import "github.com/sshelll/menuscreen"

func main() {

	root := menuscreen.NewSimpleWorkflow("root", []string{"1st line", "2nd line", "3rd line"})

	r1Node := menuscreen.NewSimpleWorkflow("r-1 node", []string{"1st line", "2nd line", "3rd line"})
	r1Node.SetCallback(func(idx int, ln string) {
		println("r1Node callback", "idx:", idx, "line:", ln)
	})

	r2Node := menuscreen.NewSimpleWorkflow("r-2 node", []string{"1st line", "2nd line", "3rd line"})
	r2Node.SetCallback(func(idx int, ln string) {
		println("r2Node callback", "idx:", idx, "line:", ln)
	})

	r3Node := menuscreen.NewSimpleWorkflow("r-3 node", []string{"1st line", "2nd line", "3rd line"})
	r3Node.SetCallback(func(idx int, ln string) {
		println("r3Node callback", "idx:", idx, "line:", ln)
	})

	root.SetNext(0, r1Node)
	root.SetNext(1, r2Node)
	root.SetNext(2, r3Node)
	menuscreen.RunWorkflow(root)

}
