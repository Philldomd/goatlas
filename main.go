package main

import (
  //"bytes"
  "context"
  "fmt"
  "time"

  "workspace/client"
  "github.com/rivo/tview"
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

func main() {
  ctx := context.Background()
  cli := client.GetClient()
  srv := client.GetService(ctx, cli)
  tmin := time.Now().Format(time.RFC3339)
  tmax := time.Now().Add(time.Duration(time.Hour * 24)).Format(time.RFC3339)
  fmt.Println(tmin, tmax)
  calendars, err := srv.CalendarList.List().Fields("items").Do()
  if err != nil {
    fmt.Println("No calendars found")
  }

  fmt.Printf("calendars: %v\n", calendars)

  for _, element := range calendars.Items {
    fmt.Printf("Name: %s      Id: %s\n" , element.Summary, element.Id)
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
  /*app := tview.NewApplication()
  dataForm := GetForm(app)
  box := tview.NewBox().SetBorder(true).SetTitle(dataForm.GetFormItemByLabel("Username").(*tview.InputField).GetText())
  if err := app.SetRoot(box, true).SetFocus(box).Run(); err != nil {
  	panic(err)
  }
  print(dataForm.GetFormItemByLabel("Password").(*tview.InputField).GetText())
  if err != nil {
  	panic(err)
  }*/
}
