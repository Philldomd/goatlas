package main

import (
	//"bytes"
	"context"
	"fmt"
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

var tday = time.Now().Weekday()
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
		workHours:  0,//8,
		shownHours: 10,
	}, {
		name:        "Wednesday",
		workHours:  0,//8,
		shownHours: 10,
	}, {
		name:       "Thursday",
		workHours:  0,//8,
		shownHours: 10,
	}, {
		name:       "Friday",
		workHours:  0,//8,
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

var calendarInfo []CalendarInfo

func selectCalendarForm(c *calendar.CalendarList, p *tview.Pages) *tview.Form {
	form := tview.NewForm()
	calInfo := CalendarInfo{}
	var calendarNames []string
	for _, entry := range c.Items {
		calendarNames = append(calendarNames, entry.Summary)
	}

	form.AddDropDown("Select Calendar: ", calendarNames, 0, func(entry string, index int) {
		calInfo.name = entry
		calInfo.id = c.Items[index].Id
	})
	form.AddButton("Save", func() {
		calendarInfo = append(calendarInfo, calInfo)
		p.SwitchToPage("Calendar")
	})
	return form
}

func createCalendarDay(d Weekday) *tview.Flex{
    day := tview.NewFlex()
	day.SetTitle(d.name)
	return day
}

func createCalendar() *tview.Flex {
	flex := tview.NewFlex()
	flex.SetBorder(true).SetTitle("Calendar")
	for _, day := range weekdays {
    if day.workHours > 0 {
		  flex.SetDirection(tview.FlexRow).AddItem(createCalendarDay(day), 0, 1, true)
    }
    flex.SetDirection(tview.FlexColumn)
	}
	return flex
}

func main() {
	ctx := context.Background()
	cli := client.GetClient()
	srv := client.GetService(ctx, cli)
	var pages = tview.NewPages()
	tmin := time.Now().Add(-time.Hour * 110).Format(time.RFC3339)
	tmax := time.Now().Add(time.Hour * 110).Format(time.RFC3339)
	app := tview.NewApplication()
  app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    if event.Rune() == 27 {
	  mod := tview.NewModal().
	    SetText("Do you want to quit?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
			} else {
				app.SetRoot(pages, true)
			}
		})
	  app.SetRoot(mod,false)
	  }
    return event
  })
	fmt.Println(tmin, tmax)
	calendars, err := srv.CalendarList.List().Fields("items").Do()
	if err != nil {
		fmt.Println("No calendars found")
		panic(err)
	}

	events, err := srv.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(tmin).
		TimeMax(tmax).
		MaxResults(10).
		OrderBy("startTime").
		Do()
	if err != nil {
		print("Error fetching events")
		panic(err)
	}
	fmt.Println("Upcomming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}
	pages.AddPage("Calendar", createCalendar(), true, false)
	pages.AddPage("Selection", selectCalendarForm(calendars, pages), true, true)
  fmt.Println(tday)
  
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		fmt.Println("App could not initialize")
		panic(err)
	}
}
