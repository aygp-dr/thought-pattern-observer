package decision

// AllScenarios returns the five mock decision scenarios.
func AllScenarios() []*Tree {
	return []*Tree{
		DebuggingScenario(),
		ArchitectureScenario(),
		HiringScenario(),
		IncidentResponseScenario(),
		CodeReviewScenario(),
	}
}

// DebuggingScenario models the thought pattern for debugging a production issue.
func DebuggingScenario() *Tree {
	return &Tree{
		Name:        "Debugging",
		Description: "Systematic approach to debugging a production issue",
		RootID:      "d-root",
		Nodes: map[string]*DecisionNode{
			"d-root": {
				ID:       "d-root",
				Question: "Is the issue reproducible?",
				Options: []Option{
					{Label: "Yes, consistently", NextID: "d-repro-yes"},
					{Label: "Intermittent", NextID: "d-repro-intermittent"},
					{Label: "No, only in production", NextID: "d-repro-no"},
				},
			},
			"d-repro-yes": {
				ID:       "d-repro-yes",
				Question: "Where does the error originate?",
				Options: []Option{
					{Label: "Frontend/UI layer", NextID: "d-frontend"},
					{Label: "API/Backend layer", NextID: "d-backend"},
					{Label: "Database/Data layer", NextID: "d-data"},
				},
			},
			"d-repro-intermittent": {
				ID:       "d-repro-intermittent",
				Question: "What varies between occurrences?",
				Options: []Option{
					{Label: "Load/concurrency", NextID: "d-race"},
					{Label: "Input data", NextID: "d-input"},
					{Label: "Timing/order", NextID: "d-timing"},
				},
			},
			"d-repro-no": {
				ID:       "d-repro-no",
				Question: "What production-specific factors exist?",
				Options: []Option{
					{Label: "Scale/traffic differences", NextID: "d-scale"},
					{Label: "Environment config", NextID: "d-env"},
				},
			},
			"d-frontend": {
				ID: "d-frontend", Outcome: "Inspect browser devtools, check component state, review recent UI changes",
			},
			"d-backend": {
				ID: "d-backend", Outcome: "Add structured logging, trace request lifecycle, check error handling paths",
			},
			"d-data": {
				ID: "d-data", Outcome: "Check query plans, verify data integrity, review recent migrations",
			},
			"d-race": {
				ID: "d-race", Outcome: "Run race detector, review shared state, add mutex or channel synchronization",
			},
			"d-input": {
				ID: "d-input", Outcome: "Fuzz test inputs, check edge cases, add input validation",
			},
			"d-timing": {
				ID: "d-timing", Outcome: "Add distributed tracing, check timeout configs, review retry logic",
			},
			"d-scale": {
				ID: "d-scale", Outcome: "Load test locally, check resource limits, profile memory/CPU",
			},
			"d-env": {
				ID: "d-env", Outcome: "Diff environment configs, check secrets/feature flags, verify dependency versions",
			},
		},
	}
}

// ArchitectureScenario models the thought pattern for architecture decisions.
func ArchitectureScenario() *Tree {
	return &Tree{
		Name:        "Architecture",
		Description: "Evaluating architecture decisions for a new system component",
		RootID:      "a-root",
		Nodes: map[string]*DecisionNode{
			"a-root": {
				ID:       "a-root",
				Question: "What is the primary constraint?",
				Options: []Option{
					{Label: "Latency (real-time)", NextID: "a-latency"},
					{Label: "Throughput (batch)", NextID: "a-throughput"},
					{Label: "Consistency (correctness)", NextID: "a-consistency"},
				},
			},
			"a-latency": {
				ID:       "a-latency",
				Question: "What is the acceptable P99 latency?",
				Options: []Option{
					{Label: "< 10ms", NextID: "a-ultra-low"},
					{Label: "< 100ms", NextID: "a-low"},
					{Label: "< 1s", NextID: "a-moderate"},
				},
			},
			"a-throughput": {
				ID:       "a-throughput",
				Question: "What is the data volume?",
				Options: []Option{
					{Label: "< 1GB/day", NextID: "a-small-batch"},
					{Label: "1-100GB/day", NextID: "a-medium-batch"},
					{Label: "> 100GB/day", NextID: "a-large-batch"},
				},
			},
			"a-consistency": {
				ID:       "a-consistency",
				Question: "What level of consistency is required?",
				Options: []Option{
					{Label: "Strong (linearizable)", NextID: "a-strong"},
					{Label: "Eventual (convergent)", NextID: "a-eventual"},
				},
			},
			"a-ultra-low": {
				ID: "a-ultra-low", Outcome: "In-memory cache, pre-computed results, edge deployment",
			},
			"a-low": {
				ID: "a-low", Outcome: "Read replicas, CDN, async write-behind cache",
			},
			"a-moderate": {
				ID: "a-moderate", Outcome: "Standard request-response, connection pooling, indexed queries",
			},
			"a-small-batch": {
				ID: "a-small-batch", Outcome: "Cron job, simple ETL script, single-node processing",
			},
			"a-medium-batch": {
				ID: "a-medium-batch", Outcome: "Message queue with workers, partitioned processing",
			},
			"a-large-batch": {
				ID: "a-large-batch", Outcome: "Distributed processing (Spark/Flink), columnar storage, data lake",
			},
			"a-strong": {
				ID: "a-strong", Outcome: "Single-leader DB, consensus protocol (Raft/Paxos), serializable transactions",
			},
			"a-eventual": {
				ID: "a-eventual", Outcome: "CRDTs, event sourcing, multi-leader replication",
			},
		},
	}
}

// HiringScenario models the thought pattern for evaluating a candidate.
func HiringScenario() *Tree {
	return &Tree{
		Name:        "Hiring",
		Description: "Structured evaluation of an engineering candidate",
		RootID:      "h-root",
		Nodes: map[string]*DecisionNode{
			"h-root": {
				ID:       "h-root",
				Question: "Does the candidate meet the technical bar?",
				Options: []Option{
					{Label: "Strong yes", NextID: "h-tech-yes"},
					{Label: "Borderline", NextID: "h-tech-maybe"},
					{Label: "No", NextID: "h-tech-no"},
				},
			},
			"h-tech-yes": {
				ID:       "h-tech-yes",
				Question: "How is their system design thinking?",
				Options: []Option{
					{Label: "Considers trade-offs well", NextID: "h-design-good"},
					{Label: "Needs guidance on trade-offs", NextID: "h-design-ok"},
				},
			},
			"h-tech-maybe": {
				ID:       "h-tech-maybe",
				Question: "Is the gap trainable within 3 months?",
				Options: []Option{
					{Label: "Yes, with mentoring", NextID: "h-trainable"},
					{Label: "Unlikely", NextID: "h-not-trainable"},
				},
			},
			"h-tech-no": {
				ID: "h-tech-no", Outcome: "No hire: technical skills below required level",
			},
			"h-design-good": {
				ID:       "h-design-good",
				Question: "Culture and collaboration fit?",
				Options: []Option{
					{Label: "Strong collaborator", NextID: "h-hire-strong"},
					{Label: "Prefers solo work", NextID: "h-hire-conditional"},
				},
			},
			"h-design-ok": {
				ID: "h-design-ok", Outcome: "Hire at mid-level: strong coding, needs design mentoring",
			},
			"h-trainable": {
				ID: "h-trainable", Outcome: "Hire with structured onboarding plan and 90-day check-in",
			},
			"h-not-trainable": {
				ID: "h-not-trainable", Outcome: "No hire: gap too large for current team capacity",
			},
			"h-hire-strong": {
				ID: "h-hire-strong", Outcome: "Strong hire: recommend for senior role",
			},
			"h-hire-conditional": {
				ID: "h-hire-conditional", Outcome: "Hire with team fit discussion: strong technically, clarify collaboration expectations",
			},
		},
	}
}

// IncidentResponseScenario models the thought pattern during an incident.
func IncidentResponseScenario() *Tree {
	return &Tree{
		Name:        "Incident Response",
		Description: "Decision framework for responding to a production incident",
		RootID:      "i-root",
		Nodes: map[string]*DecisionNode{
			"i-root": {
				ID:       "i-root",
				Question: "What is the user impact?",
				Options: []Option{
					{Label: "Full outage", NextID: "i-full"},
					{Label: "Degraded performance", NextID: "i-degraded"},
					{Label: "Partial (subset of users)", NextID: "i-partial"},
				},
			},
			"i-full": {
				ID:       "i-full",
				Question: "Was there a recent deployment?",
				Options: []Option{
					{Label: "Yes, within last hour", NextID: "i-rollback"},
					{Label: "No recent changes", NextID: "i-infra"},
				},
			},
			"i-degraded": {
				ID:       "i-degraded",
				Question: "Which metrics are affected?",
				Options: []Option{
					{Label: "Latency spike", NextID: "i-latency"},
					{Label: "Error rate increase", NextID: "i-errors"},
					{Label: "Resource exhaustion", NextID: "i-resources"},
				},
			},
			"i-partial": {
				ID:       "i-partial",
				Question: "Is the impact correlated with a region or feature?",
				Options: []Option{
					{Label: "Region-specific", NextID: "i-region"},
					{Label: "Feature-specific", NextID: "i-feature"},
				},
			},
			"i-rollback": {
				ID: "i-rollback", Outcome: "Initiate rollback immediately, page on-call, open incident channel",
			},
			"i-infra": {
				ID: "i-infra", Outcome: "Check infrastructure: DNS, load balancer, cloud provider status page",
			},
			"i-latency": {
				ID: "i-latency", Outcome: "Check DB slow queries, upstream dependency latency, GC pauses",
			},
			"i-errors": {
				ID: "i-errors", Outcome: "Check error logs, recent config changes, downstream service health",
			},
			"i-resources": {
				ID: "i-resources", Outcome: "Check CPU/memory/disk, scale horizontally, identify resource leak",
			},
			"i-region": {
				ID: "i-region", Outcome: "Check region-specific infra, DNS routing, failover to healthy region",
			},
			"i-feature": {
				ID: "i-feature", Outcome: "Disable feature flag, check feature-specific dependencies, isolate blast radius",
			},
		},
	}
}

// CodeReviewScenario models the thought pattern for reviewing a pull request.
func CodeReviewScenario() *Tree {
	return &Tree{
		Name:        "Code Review",
		Description: "Systematic approach to reviewing a pull request",
		RootID:      "c-root",
		Nodes: map[string]*DecisionNode{
			"c-root": {
				ID:       "c-root",
				Question: "What is the size of the change?",
				Options: []Option{
					{Label: "Small (< 100 lines)", NextID: "c-small"},
					{Label: "Medium (100-500 lines)", NextID: "c-medium"},
					{Label: "Large (> 500 lines)", NextID: "c-large"},
				},
			},
			"c-small": {
				ID:       "c-small",
				Question: "Is the intent clear from the PR description?",
				Options: []Option{
					{Label: "Yes, well-documented", NextID: "c-review-code"},
					{Label: "No, needs context", NextID: "c-ask-context"},
				},
			},
			"c-medium": {
				ID:       "c-medium",
				Question: "Does the change have tests?",
				Options: []Option{
					{Label: "Good test coverage", NextID: "c-review-design"},
					{Label: "Missing tests", NextID: "c-request-tests"},
				},
			},
			"c-large": {
				ID:       "c-large",
				Question: "Can this be split into smaller PRs?",
				Options: []Option{
					{Label: "Yes, should be split", NextID: "c-split"},
					{Label: "No, atomic change", NextID: "c-review-incremental"},
				},
			},
			"c-review-code": {
				ID: "c-review-code", Outcome: "Review line-by-line: check correctness, edge cases, naming, approve if clean",
			},
			"c-ask-context": {
				ID: "c-ask-context", Outcome: "Request PR description update, ask for linked issue, defer review until context provided",
			},
			"c-review-design": {
				ID: "c-review-design", Outcome: "Focus on architecture: check abstractions, interfaces, error handling, backward compatibility",
			},
			"c-request-tests": {
				ID: "c-request-tests", Outcome: "Request tests before detailed review, suggest specific test cases to add",
			},
			"c-split": {
				ID: "c-split", Outcome: "Request PR be split: suggest logical boundaries, offer to pair on decomposition",
			},
			"c-review-incremental": {
				ID: "c-review-incremental", Outcome: "Review in passes: first architecture, then logic, then style. Schedule dedicated time block",
			},
		},
	}
}
