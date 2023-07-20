package main

import (
	"context"
	"log"
	"os"

	"gioui.org/unit"
	"github.com/sashabaranov/go-openai"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	theme  = material.NewTheme(gofont.Collection())
	button = new(widget.Clickable)
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	go func() {
		w := app.NewWindow()
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func getGPTResponse(client *openai.Client, prompt string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func run(w *app.Window) error {
	var ops op.Ops
	var labelText = "Hi, I'm Giggles, I can tell you jokes about Go developers"
	var client = openai.NewClient("open-ai-secret-key")

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			if button.Clicked() {
				go func() {
					response, err := getGPTResponse(client, "Tell me a joke about Go developers")
					if err != nil {
						log.Printf("GPT Err: %s", err)
						return
					}
					labelText = response
					w.Invalidate()
				}()
			}

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						lbl := material.Label(theme, unit.Sp(20), labelText)
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						btn := material.Button(theme, button, "Tell me a joke")
						return btn.Layout(gtx)
					})
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
