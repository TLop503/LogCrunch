package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

type inputMode int

const (
	normalMode inputMode = iota
	editingMode
)

type model struct {
	choices   []string          // items in the file (IPs)
	cursor    int               // which item our cursor is pointing at
	aliases   map[string]string // map of IP addresses to their aliases
	mode      inputMode         // current input mode
	inputText string            // buffer for alias text input
}

func initialModel(choices []string) model {
	return model{
		choices: choices,
		aliases: make(map[string]string),
		mode:    normalMode,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle key presses differently based on mode
		if m.mode == normalMode {
			return m.updateNormalMode(msg)
		} else {
			return m.updateEditingMode(msg)
		}

	case []string: // File update message
		m.choices = msg
		// If we're editing and the selected IP is no longer in the file,
		// go back to normal mode
		if m.mode == editingMode && (m.cursor >= len(m.choices) || m.cursor < 0) {
			m.mode = normalMode
			m.inputText = ""
		}
	}

	return m, nil
}

func (m model) updateNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}

	case "enter", " ":
		if m.cursor < len(m.choices) {
			// Switch to editing mode
			m.mode = editingMode
			// Initialize input with existing alias if any
			if alias, exists := m.aliases[m.choices[m.cursor]]; exists {
				m.inputText = alias
			} else {
				m.inputText = ""
			}
		}
	}

	return m, nil
}

func (m model) updateEditingMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel editing
		m.mode = normalMode
		m.inputText = ""

	case "enter":
		// Save the alias
		if m.cursor < len(m.choices) {
			ip := m.choices[m.cursor]
			if m.inputText == "" {
				// Empty input means remove the alias
				delete(m.aliases, ip)
			} else {
				m.aliases[ip] = m.inputText
			}
		}
		// Return to normal mode
		m.mode = normalMode
		m.inputText = ""

	case "backspace":
		// Delete the last character
		if len(m.inputText) > 0 {
			m.inputText = m.inputText[:len(m.inputText)-1]
		}

	default:
		// Add character to input if it's a printable character
		if len(msg.String()) == 1 {
			m.inputText += msg.String()
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.mode == normalMode {
		return m.normalView()
	}
	return m.editingView()
}

func (m model) normalView() string {
	// The header
	s := "LogCrunchers active (agents):\n\n"

	// Iterate over our choices
	for i, ip := range m.choices {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Display IP with alias if it exists
		display := ip
		if alias, exists := m.aliases[ip]; exists && alias != "" {
			display = fmt.Sprintf("%s / %s", ip, alias)
		} else {
			display = fmt.Sprintf("%s / ", ip)
		}

		// Render the row
		s += fmt.Sprintf("%s [%s]\n", cursor, display)
	}

	// The footer
	s += "\nPress Enter to add/edit alias"
	s += "\nPress q to quit.\n"

	return s
}

func (m model) editingView() string {
	s := m.normalView()

	if m.cursor < len(m.choices) {
		ip := m.choices[m.cursor]
		s += fmt.Sprintf("\nEditing alias for %s\n", ip)
		s += fmt.Sprintf("Alias: %s", m.inputText)
		s += "\nPress Enter to save, Esc to cancel\n"
	}

	return s
}

func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}

func watchFile(filePath string, updates chan<- []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating file watcher:", err)
		os.Exit(1)
	}
	defer watcher.Close()

	err = watcher.Add(filePath)
	if err != nil {
		fmt.Println("Error watching file:", err)
		os.Exit(1)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				lines, err := readFile(filePath)
				if err == nil {
					updates <- lines
				}
			}
		case err := <-watcher.Errors:
			fmt.Println("File watcher error:", err)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: server_hud <file to watch>")
		os.Exit(1)
	}

	ipFile := os.Args[1]

	// Read the initial file contents
	lines, err := readFile(ipFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Create a channel for file updates
	updates := make(chan []string)

	// Start watching the file in a separate goroutine
	go watchFile(ipFile, updates)

	// Initialize the Bubble Tea program
	p := tea.NewProgram(initialModel(lines))

	// Run the program and handle file updates
	go func() {
		for updatedLines := range updates {
			p.Send(updatedLines)
		}
	}()

	if err := p.Start(); err != nil {
		fmt.Println("Error starting Bubble Tea program:", err)
		os.Exit(1)
	}
}
