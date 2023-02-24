package menuscreen

import "log"

type Workflow interface {
	// Title returns the title of the workflow, which will be shown in the menu screen.
	Title() string

	// Callback returns the callback of the workflow, which will be executed before choosing the next workflow.
	Callback() func(idx int, line string)

	// MenuItems returns the menu items of the workflow, which will be shown in the menu screen.
	// If the menu items is empty, the workflow will not create a menu screen.
	MenuItems() []string

	// Next returns the next workflow of the workflow, which will be executed after choosing the item.
	Next(chosenIdx int) Workflow

	// NextDefault returns the default next workflow of the workflow, which will be executed if the user does not choose any item.
	NextDefault() Workflow
}

// SimpleWorkflow is a simple implementation of Workflow.
type SimpleWorkflow struct {
	title       string
	items       []string
	callback    func(int, string)
	nextMap     map[int]Workflow
	nextDefault Workflow
	nextGlobal  Workflow
}

func NewSimpleWorkflow(title string, items []string) *SimpleWorkflow {
	return &SimpleWorkflow{
		title:   title,
		items:   items,
		nextMap: make(map[int]Workflow),
	}
}

func (w *SimpleWorkflow) Title() string {
	return w.title
}

func (w *SimpleWorkflow) Callback() func(idx int, line string) {
	return w.callback
}

func (w *SimpleWorkflow) MenuItems() []string {
	return w.items
}

func (w *SimpleWorkflow) Next(chosenIdx int) Workflow {
	if w.nextGlobal != nil {
		return w.nextGlobal
	}
	return w.nextMap[chosenIdx]
}

func (w *SimpleWorkflow) NextDefault() Workflow {
	return w.nextDefault
}

func (w *SimpleWorkflow) SetNext(chosenIdx int, next Workflow) {
	w.nextMap[chosenIdx] = next
}

func (w *SimpleWorkflow) SetNextDefault(next Workflow) {
	w.nextDefault = next
}

func (w *SimpleWorkflow) SetNextGlobal(next Workflow) {
	w.nextGlobal = next
}

func (w *SimpleWorkflow) SetCallback(callback func(idx int, line string)) {
	w.callback = callback
}

func RunWorkflow(w Workflow) {

	var (
		idx  int
		line string
		ok   bool
	)

	for {

		if w == nil {
			break
		}

		if w.Callback() != nil {
			w.Callback()(idx, line)
		}

		menuItem := w.MenuItems()
		if len(menuItem) == 0 {
			w = w.NextDefault()
			continue
		}

		screen, err := NewMenuScreen()
		if err != nil {
			log.Fatalln("init screen failed:", err)
		}

		idx, line, ok = screen.SetTitle(w.Title()).
			SetLines(menuItem...).
			Start().
			ChosenLine()

		screen.Fini()

		if !ok {
			w = w.NextDefault()
		} else {
			w = w.Next(idx)
		}

	}

}
