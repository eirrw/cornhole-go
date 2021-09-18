package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
	"strings"
	"virunus.com/cornhole/config"
	"virunus.com/cornhole/database"
)

var eventList []database.Event

func New(config *config.Config) *tview.Application {
	app := tview.NewApplication()

	// events
	events := tview.NewList().ShowSecondaryText(false)
	events.SetBorder(true).SetTitle("events")
	eventInfo := tview.NewTextView()
	eventInfo.SetBorder(true).SetTitle("event info")

	// teams
	teams := tview.NewList().ShowSecondaryText(false)
	teams.SetBorder(true).SetTitle("teams")

	//flex pages
	eventsFlex := tview.NewFlex().
		AddItem(events, 0, 1, true).
		AddItem(eventInfo, 0, 2, false)

	teamsFlex := tview.NewFlex().
		AddItem(teams, 0, 1, true)

	pages := tview.NewPages().
		AddPage("events", eventsFlex, true, true).
		AddPage("teams", teamsFlex, true, false)

	// default page
	app.SetRoot(pages, true)

	// load data
	var err error
	eventList, err = database.GetEvents(config)
	if err != nil {
		log.Fatal(err)
	}

	for _, event := range eventList {
		events.AddItem(event.Name, strconv.Itoa(event.EventId), 0, nil)
		events.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			e := eventList[index]

			teamList, err := e.GetTeams(config)
			if err != nil {
				log.Fatal(err)
			}

			teams.Clear()
			for _, team := range teamList {
				teams.AddItem(team.PlayerOne + "/" + team.PlayerTwo, "", 0, nil)
			}

			pages.SwitchToPage("teams")
		})
		events.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			e := eventList[index]

			infoItems := make([]string, 3)
			infoItems[0] = e.Name
			infoItems[1] = e.Date
			infoItems[2] = strconv.Itoa(e.Style)

			eventInfo.Clear()
			eventInfo.SetText(strings.Join(infoItems, "\n"))
		})
	}

	// key handlers
	events.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
			return nil
		}
		return event
	})
	teams.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			pages.SwitchToPage("events")
			return nil
		}
		return event
	})


	return app
}

