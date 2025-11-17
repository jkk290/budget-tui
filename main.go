package main

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

type Transaction struct {
	id          uuid.UUID
	amount      float32
	account     string
	created_at  time.Time
	updated_at  time.Time
	description string
}

var transactions = make([]Transaction, 0)

var pages = tview.NewPages()
var txInfo = tview.NewTextView()
var app = tview.NewApplication()
var newTxForm = tview.NewForm()
var txList = tview.NewList().ShowSecondaryText(false)
var flex = tview.NewFlex()
var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("(a) to add a new transaction \n(q) to quit")

func main() {
	txList.SetSelectedFunc(func(index int, name, second_name string, shortcut rune) {
		setTxInfo(&transactions[index])
	})

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(txList, 0, 1, true).
			AddItem(txInfo, 0, 4, false), 0, 6, false).
		AddItem(text, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			app.Stop()
		} else if event.Rune() == 97 {
			newTxForm.Clear(true)
			addNewTxForm()
			pages.SwitchToPage("Add Transaction")
		}
		return event
	})

	pages.AddPage("Menu", flex, true, true)
	pages.AddPage("Add Transaction", newTxForm, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func addNewTxForm() {
	transaction := Transaction{}

	newTxForm.AddInputField("Amount", "", 20, nil, func(amount string) {
		val, _ := strconv.ParseFloat(amount, 32)
		// if err != nil {
		// 	log.Printf("error parsing amount to float: %v", err)
		// }
		transaction.amount = float32(val)
	})

	newTxForm.AddInputField("Account", "", 20, nil, func(account string) {
		transaction.account = account
	})

	newTxForm.AddInputField("Transaction Date(YYY-MM-DD HH:MM)", "", 20, nil, func(txDate string) {
		formatted := txDate + ":00"
		timeStamp, _ := time.Parse(time.DateTime, formatted)
		// if err != nil {
		// 	log.Printf("error parsing transaction date: %v", err)
		// }
		transaction.created_at = timeStamp
	})

	newTxForm.AddInputField("Description", "", 40, nil, func(description string) {
		transaction.description = description
	})

	newTxForm.AddButton("Save", func() {
		transaction.id = uuid.New()
		transaction.updated_at = time.Now()
		transactions = append(transactions, transaction)
		addTxList()
		pages.SwitchToPage("Menu")
	})
}

func addTxList() {
	txList.Clear()
	for i, tx := range transactions {
		mainText := tx.description + " -- $" + strconv.FormatFloat(float64(tx.amount), 'f', 2, 32)
		txList.AddItem(mainText, "", rune(49+i), nil)
	}
}

func setTxInfo(transaction *Transaction) {
	txInfo.Clear()
	info := transaction.description + "\n" + "$" + strconv.FormatFloat(float64(transaction.amount), 'f', 2, 32) + "\n" + "Date: " + transaction.created_at.String() + "\n" + "Account: " + transaction.account
	txInfo.SetText(info)
}
