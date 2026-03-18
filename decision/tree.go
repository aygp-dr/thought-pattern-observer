package decision

// DecisionNode represents a single node in a decision tree.
// Leaf nodes have a non-empty Outcome and no Options.
type DecisionNode struct {
	ID       string
	Question string
	Options  []Option
	Outcome  string
}

// IsLeaf returns true if this node is a terminal node with an outcome.
func (n *DecisionNode) IsLeaf() bool {
	return n.Outcome != ""
}

// Option represents a choice that leads to another node.
type Option struct {
	Label  string
	NextID string
}

// Tree represents a complete decision tree with metadata.
type Tree struct {
	Name        string
	Description string
	RootID      string
	Nodes       map[string]*DecisionNode
}

// Root returns the root node of the tree.
func (t *Tree) Root() *DecisionNode {
	return t.Nodes[t.RootID]
}

// Node returns the node with the given ID, or nil if not found.
func (t *Tree) Node(id string) *DecisionNode {
	return t.Nodes[id]
}

// Depth returns the maximum depth of the tree from a given node.
func (t *Tree) Depth(nodeID string) int {
	node := t.Nodes[nodeID]
	if node == nil || node.IsLeaf() {
		return 0
	}
	max := 0
	for _, opt := range node.Options {
		d := t.Depth(opt.NextID)
		if d > max {
			max = d
		}
	}
	return max + 1
}

// PathTracker records navigation history through a decision tree.
type PathTracker struct {
	// Path stores the sequence of node IDs visited in the current traversal.
	Path []string
	// Frequency counts how many times each node has been visited across all traversals.
	Frequency map[string]int
}

// NewPathTracker creates a new PathTracker starting at the given root node.
func NewPathTracker(rootID string) *PathTracker {
	pt := &PathTracker{
		Path:      []string{rootID},
		Frequency: make(map[string]int),
	}
	pt.Frequency[rootID]++
	return pt
}

// Current returns the ID of the current node.
func (pt *PathTracker) Current() string {
	if len(pt.Path) == 0 {
		return ""
	}
	return pt.Path[len(pt.Path)-1]
}

// Advance moves to the next node and records the visit.
func (pt *PathTracker) Advance(nodeID string) {
	pt.Path = append(pt.Path, nodeID)
	pt.Frequency[nodeID]++
}

// Back moves to the previous node in the path. Returns false if already at root.
func (pt *PathTracker) Back() bool {
	if len(pt.Path) <= 1 {
		return false
	}
	pt.Path = pt.Path[:len(pt.Path)-1]
	return true
}

// Reset restarts traversal from the given root node, preserving frequency data.
func (pt *PathTracker) Reset(rootID string) {
	pt.Path = []string{rootID}
	pt.Frequency[rootID]++
}

// IsOnPath returns true if the given node ID is in the current path.
func (pt *PathTracker) IsOnPath(nodeID string) bool {
	for _, id := range pt.Path {
		if id == nodeID {
			return true
		}
	}
	return false
}

// VisitCount returns how many times a node has been visited.
func (pt *PathTracker) VisitCount(nodeID string) int {
	return pt.Frequency[nodeID]
}
