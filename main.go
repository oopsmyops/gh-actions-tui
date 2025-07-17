package main

import (
	"flag"
	"fmt"
	"os"

	"gh-actions-tui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	repo := flag.String("r", "", "Repository to view workflows from (e.g., 'owner/repo')")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("gh-actions-tui %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built: %s\n", date)
		return
	}

	if *repo == "" {
		fmt.Println("Error: Repository flag (-r) is required")
		fmt.Println("Usage: gh-actions-tui -r owner/repo")
		fmt.Println("Example: gh-actions-tui -r microsoft/vscode")
		os.Exit(1)
	}

	p := tea.NewProgram(tui.InitialModel(*repo))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}
