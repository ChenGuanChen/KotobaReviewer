package main

import (
	"github.com/ChenGuanChen/kotobaReviewer/gui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("Kotoba Trainer")
	w.Resize(fyne.NewSize(600, 560))

	studyContent := container.NewVBox()
	statsContent := container.NewVBox()
	manageContent := container.NewVBox()

	tabs := container.NewAppTabs(
		container.NewTabItem("Study", studyContent),
		container.NewTabItem("Manage", manageContent),
		container.NewTabItem("Stats", statsContent),
	)

	gui.ShowMenu(w, studyContent)
	gui.ShowManage(manageContent)
	gui.ShowStats(statsContent)

	w.SetContent(tabs)
	w.ShowAndRun()
}
