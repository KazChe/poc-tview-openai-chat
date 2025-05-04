package main

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"poc-tview-openai-chat/chat" // Update to your actual module path
)

func main() {
	app := tview.NewApplication()

	// Chat display area
	chatView := tview.NewTextView()
	chatView.SetDynamicColors(true).
		SetWrap(true).
		SetTitle("Chat").
		SetBorder(true)

	// User input area
	input := tview.NewInputField()
	input.SetLabel("> ").
		SetFieldWidth(0).
		SetBorder(true).
		SetTitle("Input")

	// Menu options below input
	menu := tview.NewList().
		AddItem("Clear Chat", "Remove all messages", 'c', nil).
		AddItem("Help", "Show help menu", 'h', nil).
		AddItem("Exit", "Quit the app", 'q', func() {
			app.Stop()
		})

	menu.SetBorder(true).SetTitle("Menu")

	// Chat client
	openaiClient := chat.NewOpenAIClient()

	// Input handler
	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			userMsg := input.GetText()
			chatView.Write([]byte(fmt.Sprintf("[yellow]You: %s\n", userMsg)))
			input.SetText("")

			go func() {
				reply, err := openaiClient.GetResponse(context.Background(), userMsg)
				if err != nil {
					reply = fmt.Sprintf("Error: %v", err)
				}
				app.QueueUpdateDraw(func() {
					chatView.Write([]byte(fmt.Sprintf("[green]AI: %s\n", reply)))
				})
			}()
		}
	})

	// Right pane for input (can add more widgets here later)
	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(input, 3, 0, true).
		AddItem(menu, 0, 1, false)

	// Split screen: Chat on left, Input on right
	root := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(chatView, 0, 2, false).
		AddItem(rightPanel, 0, 1, true)

	if err := app.SetRoot(root, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
