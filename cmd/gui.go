package cmd

import (
	"fmt"
	"io"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

// Gui entrypoint for the GUI application
func Gui(out io.Writer) int {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")
	myWindow.SetContent(widget.NewLabel("Hello"))

	myWindow.Show()
	myApp.Run()

	return postGuiTidy()
}

func postGuiTidy() int {
	fmt.Println("Exited")

	return 0
}
