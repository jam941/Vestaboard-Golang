package main

import (
	"fmt"
	"os"
	"github.com/jam941/vestaboard-go/vestaboard"
)

func main() {
	token := os.Getenv("VBOARD_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "VBOARD_TOKEN environment variable not set")
		os.Exit(1)
	}

	// Change to vestaboard.NewNote(token) for a 3×15 board.
	client := vestaboard.New(token)

	// Read the current board state.
	msg, err := client.GetMessage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetMessage: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Current message ID: %s\n", msg.ID)
	fmt.Printf("Layout (%d rows):\n", len(msg.Layout))
	for _, row := range msg.Layout {
		fmt.Println(row)
	}

	// Send plain text.
	result, err := client.SendText("HELLO WORLD", false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SendText: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SendText OK — id: %s\n", result.ID)



	
	// Compose a centered multi-component layout and send it.
	sendResult, err := client.ComposeAndSend(vestaboard.ComposeRequest{
		Components: []vestaboard.Component{
			{
				Template: "VESTABOARD",
				Style:    &vestaboard.ComponentStyle{Justify: "center", Align: "center"},
			},
		},
	}, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ComposeAndSend: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("ComposeAndSend OK — id: %s\n", sendResult.ID)

	// Get and update transition.
	ti, err := client.GetTransition()
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetTransition: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Current transition: %s (%s)\n", ti.Transition, ti.TransitionSpeed)

	ti, err = client.SetTransition(vestaboard.TransitionWave, vestaboard.SpeedGentle)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SetTransition: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Transition set to: %s (%s)\n", ti.Transition, ti.TransitionSpeed)
}
