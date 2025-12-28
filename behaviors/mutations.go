// Package behaviors - Agent 8: Mutation Generator
// Generates behavior variations and edge cases through systematic mutation
package behaviors

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// MutationType describes the type of mutation applied
type MutationType string

const (
	MutationAddNode       MutationType = "add_node"
	MutationRemoveNode    MutationType = "remove_node"
	MutationAddEdge       MutationType = "add_edge"
	MutationRemoveEdge    MutationType = "remove_edge"
	MutationModifyLatency MutationType = "modify_latency"
	MutationInvertEdge    MutationType = "invert_edge"
	MutationDuplicateNode MutationType = "duplicate_node"
	MutationConstraint    MutationType = "constraint"
)

// Mutation represents a single mutation applied to the behavior graph
type Mutation struct {
	ID            string
	Type          MutationType
	TargetNode    string
	TargetEdge    string
	Payload       interface{}
	Applied       bool
	Results       map[string]interface{}
	Timestamp     time.Time
}

// MutationGenerator systematically generates behavior variations
type MutationGenerator struct {
	mu            sync.RWMutex
	graph         *BehaviorGraph
	mutations     []*Mutation
	mutationIndex int
	seed          int64
}

// NewMutationGenerator creates a new mutation generator
func NewMutationGenerator(bg *BehaviorGraph, seed int64) *MutationGenerator {
	return &MutationGenerator{
		graph:     bg,
		mutations: make([]*Mutation, 0),
		seed:      seed,
	}
}

// GenerateMutations generates a set of mutations for the graph
func (mg *MutationGenerator) GenerateMutations(count int) ([]*Mutation, error) {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	rng := rand.New(rand.NewSource(mg.seed))
	mutations := make([]*Mutation, 0, count)

	nodeIDs := make([]string, 0)
	for nodeID := range mg.graph.Nodes {
		nodeIDs = append(nodeIDs, nodeID)
	}

	if len(nodeIDs) == 0 {
		return mutations, fmt.Errorf("graph has no nodes")
	}

	for i := 0; i < count; i++ {
		mutationType := mg.generateRandomMutation(rng, nodeIDs)
		mutations = append(mutations, mutationType)
	}

	mg.mutations = append(mg.mutations, mutations...)
	return mutations, nil
}

func (mg *MutationGenerator) generateRandomMutation(rng *rand.Rand, nodeIDs []string) *Mutation {
	mutationTypes := []MutationType{
		MutationAddNode,
		MutationRemoveNode,
		MutationAddEdge,
		MutationRemoveEdge,
		MutationModifyLatency,
		MutationInvertEdge,
		MutationDuplicateNode,
		MutationConstraint,
	}

	mutationType := mutationTypes[rng.Intn(len(mutationTypes))]
	targetNode := nodeIDs[rng.Intn(len(nodeIDs))]

	mutation := &Mutation{
		ID:         fmt.Sprintf("mut_%d_%d", mg.mutationIndex, time.Now().UnixNano()),
		Type:       mutationType,
		TargetNode: targetNode,
		Applied:    false,
		Results:    make(map[string]interface{}),
		Timestamp:  time.Now(),
	}
	mg.mutationIndex++

	switch mutationType {
	case MutationAddNode:
		mutation.Payload = &BehaviorNode{
			ID:          fmt.Sprintf("%s_mut_%d", targetNode, mg.mutationIndex),
			Name:        fmt.Sprintf("Mutated %s", targetNode),
			Description: fmt.Sprintf("Mutated version of %s", targetNode),
			Category:    "mutation",
		}

	case MutationModifyLatency:
		latencies := []time.Duration{
			time.Millisecond * 10,
			time.Millisecond * 50,
			time.Millisecond * 100,
			time.Millisecond * 500,
		}
		mutation.Payload = latencies[rng.Intn(len(latencies))]

	case MutationConstraint:
		constraints := []string{
			"timeout_exceeded",
			"rate_limit",
			"resource_exhaustion",
			"concurrent_limit",
		}
		mutation.Payload = constraints[rng.Intn(len(constraints))]
	}

	return mutation
}

// ApplyMutation applies a single mutation to the graph
func (mg *MutationGenerator) ApplyMutation(mutation *Mutation) error {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	switch mutation.Type {
	case MutationAddNode:
		node, ok := mutation.Payload.(*BehaviorNode)
		if !ok {
			return fmt.Errorf("invalid payload for add_node mutation")
		}
		mg.graph.Nodes[node.ID] = node
		mg.graph.Edges[node.ID] = []*BehaviorEdge{}
		mutation.Results["added_node"] = node.ID

	case MutationRemoveNode:
		if node, exists := mg.graph.Nodes[mutation.TargetNode]; exists {
			delete(mg.graph.Nodes, mutation.TargetNode)
			delete(mg.graph.Edges, mutation.TargetNode)
			mutation.Results["removed_node"] = node.ID
		}

	case MutationAddEdge:
		if nodeIDs := mg.getRandomNodePair(); len(nodeIDs) == 2 {
			from, to := nodeIDs[0], nodeIDs[1]
			edge := &BehaviorEdge{
				From:          from,
				To:            to,
				Condition:     func() bool { return true },
				Weight:        1,
				Latency:       time.Millisecond * 10,
				Deterministic: true,
			}
			mg.graph.Edges[from] = append(mg.graph.Edges[from], edge)
			mutation.Results["added_edge"] = fmt.Sprintf("%s->%s", from, to)
		}

	case MutationModifyLatency:
		if latency, ok := mutation.Payload.(time.Duration); ok {
			for _, edges := range mg.graph.Edges {
				for _, edge := range edges {
					edge.Latency = latency
				}
			}
			mutation.Results["modified_latencies"] = latency.String()
		}

	case MutationConstraint:
		if constraint, ok := mutation.Payload.(string); ok {
			if node, exists := mg.graph.Nodes[mutation.TargetNode]; exists {
				node.Constraints = append(node.Constraints, constraint)
				mutation.Results["added_constraint"] = constraint
			}
		}
	}

	mutation.Applied = true
	return nil
}

// GetMutationStats returns statistics about mutations
func (mg *MutationGenerator) GetMutationStats() map[string]interface{} {
	mg.mu.RLock()
	defer mg.mu.RUnlock()

	stats := make(map[string]interface{})
	typeCount := make(map[string]int)
	appliedCount := 0

	for _, mut := range mg.mutations {
		typeCount[string(mut.Type)]++
		if mut.Applied {
			appliedCount++
		}
	}

	stats["total_mutations"] = len(mg.mutations)
	stats["applied_mutations"] = appliedCount
	stats["type_distribution"] = typeCount
	stats["graph_nodes"] = len(mg.graph.Nodes)
	stats["graph_edges"] = len(mg.graph.Edges)

	return stats
}

// RevertMutation reverts a previously applied mutation
func (mg *MutationGenerator) RevertMutation(mutation *Mutation) error {
	mg.mu.Lock()
	defer mg.mu.Unlock()

	if !mutation.Applied {
		return fmt.Errorf("mutation %s was not applied", mutation.ID)
	}

	switch mutation.Type {
	case MutationAddNode:
		if node, ok := mutation.Payload.(*BehaviorNode); ok {
			delete(mg.graph.Nodes, node.ID)
			delete(mg.graph.Edges, node.ID)
		}

	case MutationRemoveNode:
		// Would need to store original node to restore
		return fmt.Errorf("cannot revert remove_node without original data")

	case MutationConstraint:
		if constraint, ok := mutation.Payload.(string); ok {
			if node, exists := mg.graph.Nodes[mutation.TargetNode]; exists {
				newConstraints := make([]string, 0)
				for _, c := range node.Constraints {
					if c != constraint {
						newConstraints = append(newConstraints, c)
					}
				}
				node.Constraints = newConstraints
			}
		}
	}

	mutation.Applied = false
	return nil
}

// getRandomNodePair returns two different random node IDs
func (mg *MutationGenerator) getRandomNodePair() []string {
	nodeIDs := make([]string, 0)
	for id := range mg.graph.Nodes {
		nodeIDs = append(nodeIDs, id)
	}

	if len(nodeIDs) < 2 {
		return nodeIDs
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := rng.Intn(len(nodeIDs))
	j := rng.Intn(len(nodeIDs))
	for j == i {
		j = rng.Intn(len(nodeIDs))
	}

	return []string{nodeIDs[i], nodeIDs[j]}
}
