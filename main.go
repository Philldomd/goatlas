package main

import (
  "./network"
  "github.com/rivo/tview"
  "bytes"
)

func GetForm() *tview.Form {
  app := tview.NewApplication()
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
  print(dataForm.GetFormItemByLabel("Username").(*tview.InputField).GetText())
  print(dataForm.GetFormItemByLabel("Password").(*tview.InputField).GetText())
  n := network.NewNetwork()
  buf := new(bytes.Buffer)
  err := n.Header.Write(buf)
  if err != nil {
    panic(err)
  }
  print(buf.String())
}
