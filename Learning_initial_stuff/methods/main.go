package main

import (
	"github.com/andlabs/ui"
)

func main() {
	err := ui.Main(func() {
		window := ui.NewWindow("Method expression demo", 100, 100, false)
		box := ui.NewVerticalBox()
		myCheckbox := ui.NewCheckbox("I become disabled when checked!")
		// Normally this function takes a func(*Checkbox), so
		// we transform the method into a function with a method expression
		myCheckbox.OnToggled((*ui.Checkbox).Disable)
		box.Append(myCheckbox, false)
		window.SetChild(box)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
