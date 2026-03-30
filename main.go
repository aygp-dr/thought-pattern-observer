package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aygp-dr/thought-pattern-observer/decision"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	activeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	outcomeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	pathStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	freqStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	questionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
	borderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
)

type viewState int

const (
	viewDashboard viewState = iota
	viewTree
)

type model struct {
	scenarios []*decision.Tree
	trackers  []*decision.PathTracker
	cursor    int
	view      viewState
	active    int // index of active scenario in tree view
}

func initialModel() model {
	scenarios := decision.AllScenarios()
	trackers := make([]*decision.PathTracker, len(scenarios))
	for i, s := range scenarios {
		trackers[i] = decision.NewPathTracker(s.RootID)
	}
	return model{
		scenarios: scenarios,
		trackers:  trackers,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.view == viewTree {
				m.view = viewDashboard
				m.cursor = m.active
				return m, nil
			}
			return m, tea.Quit
		case "j", "down":
			m.cursor++
			m.clampCursor()
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if m.view == viewDashboard {
				m.active = m.cursor
				m.view = viewTree
				m.cursor = 0
				return m, nil
			}
			// In tree view: select an option
			tree := m.scenarios[m.active]
			tracker := m.trackers[m.active]
			node := tree.Node(tracker.Current())
			if node != nil && !node.IsLeaf() && m.cursor < len(node.Options) {
				tracker.Advance(node.Options[m.cursor].NextID)
				m.cursor = 0
			}
		case "backspace", "esc":
			if m.view == viewTree {
				tracker := m.trackers[m.active]
				if !tracker.Back() {
					m.view = viewDashboard
					m.cursor = m.active
				} else {
					m.cursor = 0
				}
				return m, nil
			}
		case "r":
			if m.view == viewTree {
				tree := m.scenarios[m.active]
				m.trackers[m.active].Reset(tree.RootID)
				m.cursor = 0
				return m, nil
			}
		case "?":
			// Help is shown inline, no separate view needed
		}
	}
	return m, nil
}

func (m *model) clampCursor() {
	max := 0
	if m.view == viewDashboard {
		max = len(m.scenarios) - 1
	} else {
		tree := m.scenarios[m.active]
		tracker := m.trackers[m.active]
		node := tree.Node(tracker.Current())
		if node != nil && !node.IsLeaf() {
			max = len(node.Options) - 1
		}
	}
	if m.cursor > max {
		m.cursor = max
	}
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("ThoughtPatternObserver"))
	b.WriteString(dimStyle.Render(" — Decision Tree Visualization"))
	b.WriteString("\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 60)))
	b.WriteString("\n\n")

	switch m.view {
	case viewDashboard:
		b.WriteString(m.viewDashboard())
	case viewTree:
		b.WriteString(m.viewTree())
	}

	b.WriteString("\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 60)))
	b.WriteString("\n")
	b.WriteString(m.viewHelp())
	return b.String()
}

func (m model) viewDashboard() string {
	var b strings.Builder
	b.WriteString(activeStyle.Render("Decision Scenarios"))
	b.WriteString("\n\n")

	for i, s := range m.scenarios {
		cursor := "  "
		if i == m.cursor {
			cursor = activeStyle.Render("> ")
		}
		name := s.Name
		if i == m.cursor {
			name = activeStyle.Render(name)
		}
		desc := dimStyle.Render(s.Description)
		depth := fmt.Sprintf("depth:%d", s.Depth(s.RootID))

		// Show visit stats
		tracker := m.trackers[i]
		totalVisits := 0
		for _, count := range tracker.Frequency {
			totalVisits += count
		}
		stats := freqStyle.Render(fmt.Sprintf("[%s, visits:%d]", depth, totalVisits))

		b.WriteString(fmt.Sprintf("%s%-20s %s %s\n", cursor, name, desc, stats))
	}
	return b.String()
}

func (m model) viewTree() string {
	var b strings.Builder
	tree := m.scenarios[m.active]
	tracker := m.trackers[m.active]

	// Header
	b.WriteString(activeStyle.Render(tree.Name))
	b.WriteString(dimStyle.Render(fmt.Sprintf(" — %s", tree.Description)))
	b.WriteString("\n\n")

	// Show current path breadcrumb
	b.WriteString(pathStyle.Render("Path: "))
	for i, nodeID := range tracker.Path {
		node := tree.Node(nodeID)
		if node == nil {
			continue
		}
		if i > 0 {
			b.WriteString(dimStyle.Render(" → "))
		}
		label := node.Question
		if node.IsLeaf() {
			label = "✦"
		} else if len(label) > 30 {
			label = label[:27] + "..."
		}
		if nodeID == tracker.Current() {
			b.WriteString(activeStyle.Render(label))
		} else {
			b.WriteString(pathStyle.Render(label))
		}
	}
	b.WriteString("\n\n")

	// Current node
	node := tree.Node(tracker.Current())
	if node == nil {
		b.WriteString("Error: node not found\n")
		return b.String()
	}

	visits := tracker.VisitCount(node.ID)
	visitLabel := freqStyle.Render(fmt.Sprintf("[visited %dx]", visits))

	if node.IsLeaf() {
		b.WriteString(outcomeStyle.Render("Outcome:"))
		b.WriteString(fmt.Sprintf(" %s\n", visitLabel))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("  %s\n", node.Outcome))
		b.WriteString("\n")
		b.WriteString(dimStyle.Render("  Press backspace to go back, r to restart"))
		b.WriteString("\n")
	} else {
		b.WriteString(questionStyle.Render(node.Question))
		b.WriteString(fmt.Sprintf(" %s\n\n", visitLabel))

		for i, opt := range node.Options {
			cursor := "  "
			if i == m.cursor {
				cursor = activeStyle.Render("> ")
			}
			label := opt.Label
			if i == m.cursor {
				label = activeStyle.Render(label)
			}

			// Show frequency hint for next node
			nextVisits := tracker.VisitCount(opt.NextID)
			hint := ""
			if nextVisits > 0 {
				hint = freqStyle.Render(fmt.Sprintf(" (%dx)", nextVisits))
			}

			b.WriteString(fmt.Sprintf("%s%s%s\n", cursor, label, hint))
		}
	}

	// Frequency summary
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("── Frequency ──"))
	b.WriteString("\n")
	for id, count := range tracker.Frequency {
		node := tree.Node(id)
		if node == nil {
			continue
		}
		label := node.Question
		if node.IsLeaf() {
			if len(node.Outcome) > 40 {
				label = node.Outcome[:37] + "..."
			} else {
				label = node.Outcome
			}
		} else if len(label) > 40 {
			label = label[:37] + "..."
		}
		onPath := " "
		if tracker.IsOnPath(id) {
			onPath = pathStyle.Render("●")
		}
		b.WriteString(fmt.Sprintf("%s %s %s\n", onPath, freqStyle.Render(fmt.Sprintf("%2d", count)), dimStyle.Render(label)))
	}

	return b.String()
}

func (m model) viewHelp() string {
	if m.view == viewDashboard {
		return helpStyle.Render("j/k: navigate  enter: select scenario  q: quit")
	}
	return helpStyle.Render("j/k: navigate  enter: choose option  backspace/esc: back  r: restart  q: dashboard")
}

func main() {
	jsonFlag := flag.Bool("json", false, "Output scenarios as JSON and exit")
	flag.Parse()

	if *jsonFlag {
		scenarios := decision.AllScenarios()
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(scenarios); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
