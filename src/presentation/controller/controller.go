package controller

import (
	"fmt"
	"os"
	"os/signal"
	"rogue/data"
	"rogue/domain"
	"rogue/presentation/input"
	"rogue/presentation/render"
	"syscall"

	"github.com/gbin/goncurses"
)

func GameCycle() {
	mainScr, err := goncurses.Init()
	if err != nil {
		fmt.Println("goncurses init:", err)
		return
	}
	defer goncurses.End()
	if !checkMainScreenSize(mainScr) {
		return
	}

	mainScr.Timeout(1000)

	err = goncurses.StartColor()
	if checkError(err) {
		return
	}

	g, err := domain.NewGameInfo(data.NewData())
	if checkError(err) {
		return
	}

	fieldScr, logScr, statScr, backpackScr, miniMapScr, err := spawnWindows(g)
	if checkError(err) {
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	endCycle := false

	for !g.Quit && !endCycle {
		render.Render(g, mainScr, fieldScr, logScr, statScr, backpackScr, miniMapScr)
		var action domain.Action
		if !g.SkipInput {
			action = input.Input(mainScr)
		}
		if g.WaitNameInput {
			g.SetPlayerName(input.GetPlayerName(mainScr))
		}
		err = g.Update(action)
		if checkError(err) {
			return
		}

		select {
		case <-sigChan:
			endCycle = true
		default:
		}
	}
}

func checkMainScreenSize(mainScr *goncurses.Window) (correct bool) {
	maxY, maxX := mainScr.MaxYX()
	correct = true
	if maxY < 24 || maxX < 80 {
		fmt.Println("terminal window must be at least 24 x 80")
		correct = false
	}
	return
}

func spawnWindows(g *domain.GameInfo) (fieldScr, logScr, statScr, backpackScr, miniMapScr *goncurses.Window, err error) {
	fieldScr, err = goncurses.NewWindow(g.Field.Height, g.Field.Width, 1, 0)
	if err != nil {
		return
	}
	logScr, err = goncurses.NewWindow(1, g.Field.Width, 0, 0)
	if err != nil {
		return
	}
	statScr, err = goncurses.NewWindow(1, g.Field.Width, g.Field.Height+1, 0)
	if err != nil {
		return
	}
	const (
		backpackHeight = 14
		backpackWidth  = 42
		backpackY      = 1
		backpackX      = 38
	)
	backpackScr, err = goncurses.NewWindow(backpackHeight, backpackWidth, backpackY, backpackX)
	if err != nil {
		return
	}

	const (
		miniMapHeight = 13
		miniMapWidth  = 37
		miniMapY      = 1
		miniMapX      = 0
	)
	miniMapScr, err = goncurses.NewWindow(miniMapHeight, miniMapWidth, miniMapY, miniMapX)
	return
}

func checkError(err error) (hasError bool) {
	if err != nil {
		goncurses.End()
		fmt.Println(err.Error())
		hasError = true
	}
	return
}
