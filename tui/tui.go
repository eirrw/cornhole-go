// Package tui handles the setup and functionality of the main management interface for the application. I'm not a
// functional programmer but closures are fun.
package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"sort"
	"strconv"
	"strings"
	"virunus.com/cornhole/config"
	"virunus.com/cornhole/database"
)

const (
	pageEventList = "events"
	pageEventDetails = "details"
	pageModal = "modal"

	cmdCreateNew = 'c'
	cmdEdit = 'e'
	cmdDelete = 'x'
	cmdGenBracket = 'g'
)

func Get(config *config.Config) *tview.Application {
	var (
		err error
		eventId int
	)

	app := tview.NewApplication()
	pages := tview.NewPages()

	commandRow := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(nil)

	paneRow := tview.NewFlex()
	paneRow.SetBorder(false)
	paneRow.AddItem(pages, 0, 1, true)

	display := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(paneRow, 0, 1, true).
		AddItem(commandRow, 1, 1, false)

	// events
	events := tview.NewList().ShowSecondaryText(false)
	events.SetBorder(true).SetTitle("events")
	eventInfo := tview.NewTextView()
	eventInfo.SetBorder(true).SetTitle("event info")

	// teams
	teams := tview.NewTable()
	teams.SetBorder(true).SetTitle("teams")
	teams.SetSeparator(tview.Borders.Vertical).SetSelectable(true, false)

	// games
	games := tview.NewTable()
	games.SetBorder(true).SetTitle("games")
	games.SetSeparator(tview.Borders.Vertical).SetSelectable(true, true)

	//flex pages
	eventsFlex := tview.NewFlex().
		AddItem(events, 0, 1, true).
		AddItem(eventInfo, 0, 2, false)

	eventDetailsFlex := tview.NewFlex().
		AddItem(teams, 0, 1, true).
		AddItem(games, 0, 1, false)

	pages.AddPage(pageEventList, eventsFlex, true, true)
	pages.AddPage(pageEventDetails, eventDetailsFlex, true, false)

	// default page
	app.SetRoot(display, true)

	// load data
	eventList, err := database.GetEvents(config)
	if err != nil {
		log.Fatal(err)
	}

	updateEventInfo := func(e *database.Event) {
		infoItems := make([]string, 3)
		infoItems[0] = e.Name
		infoItems[1] = e.Date
		infoItems[2] = getEventStyles()[e.Style]

		eventInfo.Clear()
		eventInfo.SetText(strings.Join(infoItems, "\n"))
	}

	insertTeamRow := func(team *database.Team, row int) {
		teams.SetCell(row, 0,
			tview.NewTableCell(strconv.Itoa(team.TeamId)).
				SetAlign(tview.AlignCenter))
		teams.SetCell(row, 1,
			tview.NewTableCell(team.PlayerOne).
				SetAlign(tview.AlignCenter))
		teams.SetCell(row, 2,
			tview.NewTableCell(team.PlayerTwo).
				SetAlign(tview.AlignCenter))
	}

	for _, event := range eventList {
		events.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			eventId = index
			e := eventList[index]

			teamList, err := e.GetTeams(config)
			if err != nil {
				log.Fatal(err)
			}

			// remove old entries
			teams.Clear()

			// set headers
			teams.SetCell(0, 0, getHeaderCell("team id").SetExpansion(1)).
				SetCell(0, 1, getHeaderCell("player one").SetExpansion(3)).
				SetCell(0,2,getHeaderCell("player two").SetExpansion(3))
			teams.SetFixed(1, 0)

			// load teams
			for idx, team := range teamList {
				insertTeamRow(team, idx+1)
			}

			setCmdInfo(getTeamCommands(), commandRow)
			pages.SwitchToPage(pageEventDetails)
		})
		events.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			updateEventInfo(eventList[index])
		})
		events.AddItem(event.Name, strconv.Itoa(event.EventId), 0, nil)
	}

	editEvent := func(idx int) {
		var form *tview.Form
		event := eventList[idx]

		cancelFunc := func() {
			form.Clear(true)
			pages.RemovePage(pageModal)
			app.SetFocus(events)
		}

		saveFunc := func() {
			if _, err := event.Save(config); err != nil {
				log.Panic(err)
			}

			eventList[idx] = event
			events.SetItemText(idx, event.Name, "")
			events.SetCurrentItem(idx)
			updateEventInfo(event)

			cancelFunc()
		}

		form = tview.NewForm().
			AddInputField("Name", event.Name, 0, nil, func(text string) {
			event.Name = text
		}).
			AddInputField("Date", event.Date, 0, nil, func(text string) {
			event.Date = text
		}).
			AddDropDown("Style", getEventStyles(), event.Style, func(_ string, optionIndex int) {
			event.Style = optionIndex
		}).
			AddButton("Save", saveFunc).
			AddButton("Cancel", cancelFunc).
			SetCancelFunc(cancelFunc).
			SetButtonsAlign(tview.AlignCenter)
		form.SetTitle("Edit Event").SetBorder(true)

		pages.AddPage(pageModal, customModal(form, 40, 15), true, true)
	}

	createEvent := func() {
		var form *tview.Form
		event := new(database.Event)

		cancelFunc := func() {
			form.Clear(true)
			pages.RemovePage(pageModal)
			app.SetFocus(events)
		}

		saveFunc := func() {
			if _, err := event.Save(config); err != nil {
				log.Panic(err)
			}

			eventList = append(eventList, event)
			events.AddItem(event.Name, "", 0, nil)
			events.SetCurrentItem(-1)
			updateEventInfo(event)

			cancelFunc()
		}

		form = tview.NewForm().
			AddInputField("Name", event.Name, 0, nil, func(text string) {
				event.Name = text
			}).
			AddInputField("Date", event.Date, 0, nil, func(text string) {
				event.Date = text
			}).
			AddDropDown("Style", getEventStyles(), event.Style, func(_ string, optionIndex int) {
				event.Style = optionIndex
			}).
			AddButton("Save", saveFunc).
			AddButton("Cancel", cancelFunc).
			SetCancelFunc(cancelFunc).
			SetButtonsAlign(tview.AlignCenter)
		form.SetTitle("Edit Event").SetBorder(true)

		pages.AddPage(pageModal, customModal(form, 40, 15), true, true)
	}

	deleteEvent := func(idx int) {
		event := eventList[idx]

		if err := event.Delete(config); err != nil {
			log.Panic(err)
		}

		newList := make([]*database.Event, 0)
		newList = append(newList, eventList[:idx]...)
		eventList = append(newList, eventList[idx+1:]...)

		events.RemoveItem(idx)
	}

	// loadTeam is a function that loads team information from the given row in the teams table.
	loadTeam := func(row int) *database.Team {
		teamId, _ := strconv.Atoi(teams.GetCell(row, 0).Text)
		team := database.Team{
			TeamId: teamId,
			EventId: eventList[eventId].EventId,
			PlayerOne: teams.GetCell(row, 1).Text,
			PlayerTwo: teams.GetCell(row, 2).Text,
		}

		return &team
	}

	editTeam := func(row int) {
		var form *tview.Form

		team := loadTeam(row)

		cancelFunc := func() {
			form.Clear(true)
			pages.RemovePage(pageModal)
			app.SetFocus(teams)
		}

		saveFunc := func() {
			if _, err := team.Save(config); err != nil {
				log.Panic(err)
			}

			teams.GetCell(row, 1).SetText(team.PlayerOne)
			teams.GetCell(row, 2).SetText(team.PlayerTwo)

			cancelFunc()
		}

		form = tview.NewForm().
			AddInputField("Player One", team.PlayerOne, 0, nil, func(text string) {
			team.PlayerOne = text
		}).
			AddInputField("Player Two", team.PlayerTwo, 0, nil, func(text string) {
			team.PlayerTwo = text
		}).
			AddButton("Save", saveFunc).
			AddButton("Cancel", cancelFunc).
			SetCancelFunc(cancelFunc).
			SetButtonsAlign(tview.AlignCenter)
		form.SetTitle("Edit Team").SetBorder(true)

		modal := customModal(form, 40, 10)
		pages.AddPage(pageModal, modal, true, true)
	}

	createTeam := func() {
		var (
			form *tview.Form
		)
		team := new(database.Team)

		cancelFunc := func() {
			form.Clear(true)
			pages.RemovePage(pageModal)
			app.SetFocus(teams)
		}

		saveFunc := func() {
			team.EventId = eventList[eventId].EventId
			team, err = team.Save(config)
			if err != nil {
				log.Panic(err)
			}

			idx := teams.GetRowCount()
			insertTeamRow(team, idx)

			cancelFunc()
		}

		form = tview.NewForm().
			AddInputField("Player One", "", 0, nil, func(text string) {
				team.PlayerOne = text
			}).
			AddInputField("Player Two", "", 0, nil, func(text string) {
				team.PlayerTwo = text
			}).
			AddButton("Save", saveFunc).
			AddButton("Cancel", cancelFunc).
			SetCancelFunc(cancelFunc).
			SetButtonsAlign(tview.AlignCenter)
		form.SetTitle("Edit Team").SetBorder(true)

		modal := customModal(form, 40, 10)
		pages.AddPage(pageModal, modal, true, true)
	}

	deleteTeam := func(row int) {
		team := loadTeam(row)

		// TODO: confirm deletion
		if err := team.Delete(config); err != nil {
			log.Panic(err)
		}

		teams.RemoveRow(row)
	}

	// key handlers
	events.SetDoneFunc(func() {
		app.Stop()
	})
	events.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		eventIdx := events.GetCurrentItem()
		switch event.Rune() {
		case cmdCreateNew:
			createEvent()
			return nil
		case cmdEdit:
			editEvent(eventIdx)
			return nil
		case cmdDelete:
			deleteEvent(eventIdx)
			return nil
		default:
			return event
		}
	})

	teams.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEsc {
			setCmdInfo(getEventCommands(), commandRow)
			pages.SwitchToPage(pageEventList)
		}
	})
	teams.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := teams.GetSelection()
		switch event.Rune() {
		case cmdEdit:
			editTeam(row)
			return nil
		case cmdCreateNew:
			createTeam()
			return nil
		case cmdDelete:
			deleteTeam(row)
			return nil
		case cmdGenBracket:
			// generate games
			return nil
		default:
			return event
		}
	})

	setCmdInfo(getEventCommands(), commandRow)

	//app.EnableMouse(true)
	return app
}

// customModal creates a flex element to float a generic primitive in the centre of the screen with the given width and
// height.
func customModal(p tview.Primitive, width, height int) *tview.Grid {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func setCmdInfo(commands map[rune]string, view *tview.TextView) {
	view.Clear()
	var keys []rune
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, cmd := range keys {
		fmt.Fprintf(view, "[yellow]%c)[white] %s\t", cmd, commands[cmd])
	}
}

func getEventCommands() map[rune]string {
	commands := map[rune]string{
		cmdCreateNew: "create new",
		cmdEdit:      "edit",
		cmdDelete:    "delete",
	}

	return commands
}

func getTeamCommands() map[rune]string {
	commands := map[rune]string{
		cmdCreateNew: "add new",
		cmdEdit:      "edit",
		cmdDelete:    "delete",
		cmdGenBracket: "generate games",
	}

	return commands
}

func getEventStyles() []string {
	styles := make([]string, 2)
	styles[0] = "single elimination"
	styles[1] = "double elimination"

	return styles
}

func getHeaderCell(headerText string) *tview.TableCell {
	cell := tview.NewTableCell(headerText)
	cell.SetAlign(tview.AlignCenter)
	cell.SetSelectable(false)
	cell.SetTextColor(tcell.ColorYellow)

	return cell
}

