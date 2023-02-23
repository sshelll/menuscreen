package menuscreen

import "log"

type Workflow interface {
	// Title returns the title of the workflow, which will be shown in the menu screen.
	Title() string

	// Callback returns the callback of the workflow, which will be executed before choosing the next workflow.
	Callback() func()

	// MenuItems returns the menu items of the workflow, which will be shown in the menu screen.
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
	callback    func()
	next        map[int]Workflow
	nextDefault Workflow
}

func NewSimpleWorkflow(title string, items []string) *SimpleWorkflow {
	return &SimpleWorkflow{
		title: title,
		items: items,
		next:  make(map[int]Workflow),
	}
}

func (w *SimpleWorkflow) Title() string {
	return w.title
}

func (w *SimpleWorkflow) Callback() func() {
	return w.callback
}

func (w *SimpleWorkflow) MenuItems() []string {
	return w.items
}

func (w *SimpleWorkflow) Next(chosenIdx int) Workflow {
	return w.next[chosenIdx]
}

func (w *SimpleWorkflow) NextDefault() Workflow {
	return w.nextDefault
}

func (w *SimpleWorkflow) SetNext(chosenIdx int, next Workflow) {
	w.next[chosenIdx] = next
}

func (w *SimpleWorkflow) SetNextDefault(next Workflow) {
	w.nextDefault = next
}

func (w *SimpleWorkflow) SetCallback(callback func()) {
	w.callback = callback
}

func RunWorkflow(w Workflow) {

	for {

		if w == nil {
			break
		}

		if w.Callback() != nil {
			w.Callback()()
		}

		screen, err := NewMenuScreen()
		if err != nil {
			log.Fatalln("init screen failed:", err)
		}

		idx, _, ok := screen.SetTitle(w.Title()).
			SetLines(w.MenuItems()...).
			Start().
			ChosenLine()

		if !ok {
			w = w.NextDefault()
		} else {
			w = w.Next(idx)
		}

	}

}
