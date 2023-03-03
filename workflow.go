package menuscreen

import (
	"context"
	"log"

	"github.com/sshelll/sinfra/util"
)

type Workflow interface {
	// ID returns the id of the workflow. This id should be unique.
	ID() string

	// Ctx returns the context of the workflow. See RunWorkflow for more details.
	// Do not call this method in your code.
	Ctx() context.Context

	// SetCtx sets the context of the workflow. See RunWorkflow for more details.
	// Do not call this method in your code.
	SetCtx(ctx context.Context)

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

	// GetSelected returns the selected item of the workflow.
	GetSelected() (idx int, line string, ok bool)
}

// SimpleWorkflow is a simple implementation of Workflow.
type SimpleWorkflow struct {
	// workflow info
	id  string
	ctx context.Context

	// menu info
	title string
	items []string

	// next workflow info
	callback    func(int, string)
	nextMap     map[int]Workflow
	nextDefault Workflow
	nextGlobal  Workflow
}

type selectResult struct {
	selected bool
	idx      int
	line     string
}

func NewSimpleWorkflow(title string, items []string) *SimpleWorkflow {
	return &SimpleWorkflow{
		id:      util.UUID(),
		title:   title,
		items:   items,
		nextMap: make(map[int]Workflow),
	}
}

func (w *SimpleWorkflow) ID() string {
	return w.id
}

func (w *SimpleWorkflow) Ctx() context.Context {
	return w.ctx
}

func (w *SimpleWorkflow) SetCtx(ctx context.Context) {
	w.ctx = ctx
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
	if w.nextGlobal != nil {
		return w.nextGlobal
	}
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

func (w *SimpleWorkflow) GetSelected() (idx int, line string, ok bool) {
	v := w.Ctx().Value(w.ID())
	if v == nil {
		return 0, "", false
	}
	selectResult := v.(*selectResult)
	return selectResult.idx, selectResult.line, selectResult.selected
}

func RunWorkflow(w Workflow) {

	var (
		idx  int
		line string
		ok   bool
	)

	ctx := context.Background()

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

		ctx = context.WithValue(ctx, w.ID(), &selectResult{
			selected: ok,
			idx:      idx,
			line:     line,
		})
		w.SetCtx(ctx)

		screen.Fini()

		if !ok {
			w = w.NextDefault()
		} else {
			w = w.Next(idx)
		}

	}

}
