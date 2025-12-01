package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	baseURL := "http://localhost:8080/api/v1"

	client := newClient(baseURL)

	p := tea.NewProgram(initialModel(client))
	if _, err := p.Run(); err != nil {
		fmt.Printf("An error occurred: %v", err)
		os.Exit(1)
	}
}
