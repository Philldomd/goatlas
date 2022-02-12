package main

import (
  "github.com/rivo/tview"
)

func GetForm() *tview.Form {
  app := tview .NewApplication()
  form := tview.NewForm().
    AddInputField("Username", "", 20, nil, nil).
    AddPasswordField("Password", "", 20, '*', nil).
    AddButton("Submit", func() {
      app.Stop()
    })
  form.SetBorder(true).SetTitle("Submit username/passwod").SetTitleAlign(tview.AlignLeft)
  if err := app.SetRoot(form, false).SetFocus(form).Run(); err != nil {
    panic(err)
  }
  return form
}

func main() {
  dataForm := GetForm()
  print(dataForm.GetFormItemByLabel("Username"))
}
