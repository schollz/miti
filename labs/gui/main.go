package main

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

func main() {
	app := app.New()

	w := app.NewWindow("Hello")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Start", func() {
			dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
				fmt.Println(file, err)

			}, w)
		}),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
