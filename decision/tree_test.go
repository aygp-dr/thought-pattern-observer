package decision

import (
	"testing"
)

func TestDecisionNodeIsLeaf(t *testing.T) {
	leaf := &DecisionNode{ID: "leaf", Outcome: "done"}
	if !leaf.IsLeaf() {
		t.Error("expected leaf node")
	}

	branch := &DecisionNode{ID: "branch", Question: "what?", Options: []Option{{Label: "a", NextID: "b"}}}
	if branch.IsLeaf() {
		t.Error("expected non-leaf node")
	}
}

func TestTreeRoot(t *testing.T) {
	tree := &Tree{
		RootID: "root",
		Nodes: map[string]*DecisionNode{
			"root": {ID: "root", Question: "start?"},
		},
	}
	root := tree.Root()
	if root == nil || root.ID != "root" {
		t.Error("Root() should return the root node")
	}
}

func TestTreeNode(t *testing.T) {
	tree := &Tree{
		RootID: "root",
		Nodes: map[string]*DecisionNode{
			"root": {ID: "root"},
			"child": {ID: "child"},
		},
	}
	if tree.Node("child") == nil {
		t.Error("Node() should find existing node")
	}
	if tree.Node("nonexistent") != nil {
		t.Error("Node() should return nil for missing node")
	}
}

func TestTreeDepth(t *testing.T) {
	tree := &Tree{
		RootID: "root",
		Nodes: map[string]*DecisionNode{
			"root": {ID: "root", Question: "q1", Options: []Option{
				{Label: "a", NextID: "mid"},
				{Label: "b", NextID: "leaf1"},
			}},
			"mid": {ID: "mid", Question: "q2", Options: []Option{
				{Label: "c", NextID: "leaf2"},
			}},
			"leaf1": {ID: "leaf1", Outcome: "end1"},
			"leaf2": {ID: "leaf2", Outcome: "end2"},
		},
	}
	depth := tree.Depth("root")
	if depth != 2 {
		t.Errorf("expected depth 2, got %d", depth)
	}
	if tree.Depth("mid") != 1 {
		t.Errorf("expected depth 1 for mid, got %d", tree.Depth("mid"))
	}
	if tree.Depth("leaf1") != 0 {
		t.Errorf("expected depth 0 for leaf, got %d", tree.Depth("leaf1"))
	}
}

func TestPathTrackerNavigation(t *testing.T) {
	pt := NewPathTracker("root")

	if pt.Current() != "root" {
		t.Errorf("expected current=root, got %s", pt.Current())
	}
	if pt.VisitCount("root") != 1 {
		t.Errorf("expected root visit count 1, got %d", pt.VisitCount("root"))
	}

	pt.Advance("child1")
	if pt.Current() != "child1" {
		t.Errorf("expected current=child1, got %s", pt.Current())
	}
	if !pt.IsOnPath("root") {
		t.Error("root should be on path")
	}
	if !pt.IsOnPath("child1") {
		t.Error("child1 should be on path")
	}
	if pt.IsOnPath("other") {
		t.Error("other should not be on path")
	}

	pt.Advance("leaf")
	if len(pt.Path) != 3 {
		t.Errorf("expected path length 3, got %d", len(pt.Path))
	}

	ok := pt.Back()
	if !ok {
		t.Error("Back() should succeed")
	}
	if pt.Current() != "child1" {
		t.Errorf("expected current=child1 after back, got %s", pt.Current())
	}
}

func TestPathTrackerBackAtRoot(t *testing.T) {
	pt := NewPathTracker("root")
	if pt.Back() {
		t.Error("Back() at root should return false")
	}
}

func TestPathTrackerReset(t *testing.T) {
	pt := NewPathTracker("root")
	pt.Advance("a")
	pt.Advance("b")
	pt.Reset("root")

	if pt.Current() != "root" {
		t.Errorf("expected current=root after reset, got %s", pt.Current())
	}
	if len(pt.Path) != 1 {
		t.Errorf("expected path length 1 after reset, got %d", len(pt.Path))
	}
	// Frequency should persist across resets
	if pt.VisitCount("root") != 2 {
		t.Errorf("expected root visit count 2, got %d", pt.VisitCount("root"))
	}
	if pt.VisitCount("a") != 1 {
		t.Errorf("expected 'a' visit count 1, got %d", pt.VisitCount("a"))
	}
}

func TestPathTrackerEmptyCurrent(t *testing.T) {
	pt := &PathTracker{Path: nil, Frequency: make(map[string]int)}
	if pt.Current() != "" {
		t.Errorf("expected empty current for nil path, got %s", pt.Current())
	}
}

func TestAllScenariosCount(t *testing.T) {
	scenarios := AllScenarios()
	if len(scenarios) != 5 {
		t.Errorf("expected 5 scenarios, got %d", len(scenarios))
	}
}

func TestScenarioIntegrity(t *testing.T) {
	for _, tree := range AllScenarios() {
		t.Run(tree.Name, func(t *testing.T) {
			if tree.Name == "" {
				t.Error("scenario must have a name")
			}
			if tree.Description == "" {
				t.Error("scenario must have a description")
			}
			if tree.Root() == nil {
				t.Fatal("root node must exist")
			}

			// Verify all option NextIDs point to existing nodes
			for id, node := range tree.Nodes {
				if node.ID != id {
					t.Errorf("node map key %q doesn't match node ID %q", id, node.ID)
				}
				for _, opt := range node.Options {
					if tree.Node(opt.NextID) == nil {
						t.Errorf("node %q option %q points to nonexistent node %q", id, opt.Label, opt.NextID)
					}
				}
				// A node should be either a branch (has options) or a leaf (has outcome), not both
				if node.IsLeaf() && len(node.Options) > 0 {
					t.Errorf("node %q has both outcome and options", id)
				}
				if !node.IsLeaf() && node.Question == "" {
					t.Errorf("branch node %q must have a question", id)
				}
			}

			// Every scenario should have at least one leaf
			hasLeaf := false
			for _, node := range tree.Nodes {
				if node.IsLeaf() {
					hasLeaf = true
					break
				}
			}
			if !hasLeaf {
				t.Error("scenario must have at least one leaf node")
			}
		})
	}
}

func TestScenarioDepths(t *testing.T) {
	for _, tree := range AllScenarios() {
		depth := tree.Depth(tree.RootID)
		if depth < 1 {
			t.Errorf("scenario %q should have depth >= 1, got %d", tree.Name, depth)
		}
		if depth > 10 {
			t.Errorf("scenario %q has suspiciously deep tree: %d", tree.Name, depth)
		}
	}
}

func TestDebuggingScenarioWalkthrough(t *testing.T) {
	tree := DebuggingScenario()
	pt := NewPathTracker(tree.RootID)

	// Walk: reproducible -> yes -> backend
	root := tree.Node(pt.Current())
	if root.Question != "Is the issue reproducible?" {
		t.Errorf("unexpected root question: %s", root.Question)
	}

	pt.Advance("d-repro-yes")
	node := tree.Node(pt.Current())
	if node.Question != "Where does the error originate?" {
		t.Errorf("unexpected question: %s", node.Question)
	}

	pt.Advance("d-backend")
	leaf := tree.Node(pt.Current())
	if !leaf.IsLeaf() {
		t.Error("d-backend should be a leaf")
	}
	if leaf.Outcome == "" {
		t.Error("leaf should have an outcome")
	}
}
