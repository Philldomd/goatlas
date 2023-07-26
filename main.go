package main

import (
	//"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"workspace/client"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/calendar/v3"
)

func GetForm(app *tview.Application) *tview.Form {
	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Submit", func() {
			app.Stop()
		}).SetButtonsAlign(tview.AlignCenter)
	form.SetBorder(true).SetTitle("Submit username/password").SetTitleAlign(tview.AlignCenter)
	if err := app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}
	return form
}

var tday = time.Now()

type Weekday struct {
	name       string
	workHours  int
	shownHours int
}

var weekdays = []Weekday{
	{
		name:       "Monday",
		workHours:  8,
		shownHours: 10,
	}, {
		name:       "Tuesday",
		workHours:  8,
		shownHours: 10,
	}, {
		name:       "Wednesday",
		workHours:  8,
		shownHours: 10,
	}, {
		name:       "Thursday",
		workHours:  8,
		shownHours: 10,
	}, {
		name:       "Friday",
		workHours:  8,
		shownHours: 10,
	}, {
		name:       "Saturday",
		workHours:  0,
		shownHours: 10,
	}, {
		name:       "Sunday",
		workHours:  0,
		shownHours: 10,
	},
}

type CalendarInfo struct {
	name string
	id   string
}

func selectCalendarForm(srv *calendar.Service, p *tview.Pages) *tview.Form {
	updateTextView := func(label string, text map[string]string, textView *tview.TextView, form *tview.Form) tview.FormItem {
		var split []string
		for key := range text {
			split = append(split, key)
		}
		textView.SetText(strings.Join(split, ", ")).SetBorder(false)
		form.RemoveFormItem(form.GetFormItemIndex(label))
		return textView
	}
	calendarInfo := map[string]string{}
	textViewLabel := "Choosen"
	calendars, err := srv.CalendarList.List().Fields("items").Do()
	if err != nil {
		fmt.Println("No calendars found")
		panic(err)
	}
	form := tview.NewForm()
	choosenCalendars := tview.NewTextView()
	choosenCalendars.SetLabel("Choosen")

	calInfo := CalendarInfo{}
	var calendarNames []string
	for _, entry := range calendars.Items {
		calendarNames = append(calendarNames, entry.Summary)
	}

	form.AddDropDown("Select Calendar: ", calendarNames, 0, func(entry string, index int) {
		calInfo.name = entry
		calInfo.id = calendars.Items[index].Id
	})

	form.AddFormItem(choosenCalendars)
	form.AddButton("Add", func() {
		calendarInfo[calInfo.name] = calInfo.id
		form.AddFormItem(updateTextView(textViewLabel, calendarInfo, choosenCalendars, form))
	})

	form.AddButton("Remove", func() {
		delete(calendarInfo, calInfo.name)
		form.AddFormItem(updateTextView(textViewLabel, calendarInfo, choosenCalendars, form))
	})

	form.AddButton("Save", func() {
		p.AddPage("Calendar", createCalendar(srv, calendarInfo), true, false)
		p.SwitchToPage("Calendar")
	})
	return form
}

func createCalendarDay(pastWeekday *bool, d Weekday, count int) *tview.Table {
	day := tview.NewTable().SetSelectable(true, true)
	date := fmt.Sprint(tday.Add(time.Hour*24*time.Duration(-int(tday.Weekday())+count)).Day()) + "/" + fmt.Sprint(int(tday.Month()))
	day.SetTitle(d.name + " " + date).SetBorder(true).SetBorderColor(tcell.ColorBlack)
	if d.name == tday.Weekday().String() {
		day.SetBorderColor(tcell.ColorCornflowerBlue)
		*pastWeekday = true
	} else if *pastWeekday {
		day.SetBorderColor(tcell.ColorDarkGreen)
	}
	
	day.SetBorder(true)
	day.SetCell(0, 0, tview.NewTableCell("testar").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	day.SetCell(1, 0, tview.NewTableCell("testar").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	day.SetCell(2, 0, tview.NewTableCell("testar").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	day.SetCell(3, 0, tview.NewTableCell("testar").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	day.SetCell(5, 0, tview.NewTableCell("testar").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	day.SetCell(6, 0, tview.NewTableCell("testar").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	day.SetCell(8, 0, tview.NewTableCell("testar").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
    day.SetSelectedFunc(func(row int, column int) {
		day.GetCell(row, column).SetTextColor(tcell.ColorRed)
	})
	return day
}

func getEvents(srv *calendar.Service, calendarInfo map[string]string) [][]*calendar.Event {
	eventList := [][]*calendar.Event{}
	for _, val := range calendarInfo {
		events, err := srv.Events.List(val).
			ShowDeleted(false).
			SingleEvents(true).
			TimeMin(time.Now().Add(-time.Hour * 220).Format(time.RFC3339)).
			TimeMax(time.Now().Add(time.Hour * 220).Format(time.RFC3339)).
			MaxResults(10).
			OrderBy("startTime").
			Do()
		if err != nil {
			print("Error fetching events")
			panic(err)
		}
		eventList = append(eventList, events.Items)
    }
	return eventList
}

func createCalendar(srv *calendar.Service, calendarInfo map[string]string) *tview.Flex {
	flex := tview.NewFlex()
	//events := getEvents(srv, calendarInfo)
	pastWeekday := false
	calendarName := "Calendar"
	if len(calendarInfo) > 0 {
		calendarName = ""
		var choosenCalendars []string
		for key := range calendarInfo {
			choosenCalendars = append(choosenCalendars, key)
		}
		calendarName = strings.Join(choosenCalendars, ", ")
	}
	flex.SetBorder(true).SetTitle(calendarName).SetBorderColor(tcell.ColorBlack)
	count := 1
	for _, day := range weekdays {
		if day.workHours > 0 {
			flex.SetDirection(tview.FlexRow).AddItem(createCalendarDay(&pastWeekday, day, count), 0, 1, false)
			count++
		}
		flex.SetDirection(tview.FlexColumn)
	}
	return flex
}

func quitMod(app *tview.Application, p *tview.Pages) *tview.Modal {
	mod := tview.NewModal().
		SetText("Do you want to quit?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
			} else {
				app.SetRoot(p, true)
			}
		})
	return mod
}

func main() {
	ctx := context.Background()
	cli := client.GetClient()
	srv := client.GetService(ctx, cli)
	var pages = tview.NewPages()
	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 27 {
			app.SetRoot(quitMod(app, pages), false)
		}
		return event
	})
	pages.AddPage("Selection", selectCalendarForm(srv, pages), true, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		fmt.Println("App could not initialize")
		panic(err)
	}
}
