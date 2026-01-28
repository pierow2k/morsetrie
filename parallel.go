package morsetrie

import (
	"runtime"
	"sync"
)

// decodeTask represents an independent exploration point in the trie.
type decodeTask struct {
	idx     int    // Position in the morse sequence
	nodeIdx int16  // Position in the trie
	prefix  string // Accumulated decoded characters
}

// FindCandidates returns all valid decodings for a Morse code sequence.
// It parallelizes exploration by distributing branches from a breadth-first
// frontier to worker goroutines.
//
// Concurrency model:
//   - Producer: BFS traversal to generate ~numWorkers independent tasks
//   - Workers: Pull tasks and recursively explore the subtree
//   - Accumulation: Per-worker slices merged at the end
func (t *Trie) FindCandidatesParallel(sequence string) []string {
	if len(sequence) == 0 {
		return nil
	}

	// Validate input: only '.' and '-' are allowed.
	for i := range sequence {
		if sequence[i] != '.' && sequence[i] != '-' {
			return nil
		}
	}

	numWorkers := runtime.GOMAXPROCS(0)

	// Channel capacity: enough to keep workers busy without over-allocating.
	tasks := make(chan decodeTask, numWorkers*8)

	// Per-worker results eliminate mutex contention.
	workerResults := make([][]string, numWorkers)

	var wg sync.WaitGroup

	// Start workers.
	for w := range numWorkers {
		wg.Go(func() {
			// Pre-allocate a reasonable capacity for each worker's results.
			localResults := make([]string, 0, 1024)

			for task := range tasks {
				t.parallelTraverse(sequence, task.idx, task.nodeIdx, &localResults, task.prefix)
			}

			workerResults[w] = localResults
		})
	}

	// Producer: Generate tasks via breadth-first expansion.
	// We expand the tree until we have enough independent tasks to saturate workers.
	frontier := []decodeTask{{idx: 0, nodeIdx: rootIdx, prefix: ""}}
	targetTasks := numWorkers * 4

	for len(frontier) > 0 && len(frontier) < targetTasks {
		current := frontier[0]
		frontier = frontier[1:]

		// Expand the current node into its children.
		children := t.expandBranches(sequence, current.idx, current.nodeIdx, current.prefix)

		if len(children) > 0 {
			frontier = append(frontier, children...)
		}
	}

	// Distribute the frontier to workers.
	for _, task := range frontier {
		tasks <- task
	}
	close(tasks)

	wg.Wait()

	// Merge results.
	total := 0
	for _, r := range workerResults {
		total += len(r)
	}

	merged := make([]string, 0, total)
	for _, r := range workerResults {
		merged = append(merged, r...)
	}

	return merged
}

// expandBranches generates the next set of tasks from a given position.
// It returns tasks for both continuing the current letter and starting a new one.
func (t *Trie) expandBranches(sequence string, idx int, nodeIdx int16, prefix string) []decodeTask {
	if idx >= len(sequence) {
		return nil
	}

	symbol := sequence[idx]
	var bit int
	if symbol == '.' {
		bit = 0
	} else if symbol == '-' {
		bit = 1
	} else {
		return nil
	}

	childIdx := t.Nodes[nodeIdx].Child[bit]
	if childIdx == missingNode {
		return nil
	}

	nextIdx := idx + 1
	tasks := make([]decodeTask, 0, 2)

	// CONTINUE path: extend current letter prefix.
	tasks = append(tasks, decodeTask{idx: nextIdx, nodeIdx: childIdx, prefix: prefix})

	// BRANCH path: commit letter and restart at root.
	if val := t.Nodes[childIdx].Val; val != 0 {
		newPrefix := prefix + string(val)
		tasks = append(tasks, decodeTask{idx: nextIdx, nodeIdx: rootIdx, prefix: newPrefix})
	}

	return tasks
}

// parallelTraverse recursively explores the trie from a given starting point.
// It is executed by workers to process independent subtrees.
func (t *Trie) parallelTraverse(sequence string, idx int, nodeIdx int16, results *[]string, prefix string) {
	// Base case: end of sequence.
	if idx == len(sequence) {
		if nodeIdx == rootIdx && prefix != "" {
			*results = append(*results, prefix)
		}
		return
	}

	symbol := sequence[idx]
	var bit int

	switch symbol {
	case '.':
		bit = 0
	case '-':
		bit = 1
	default:
		return
	}

	childIdx := t.Nodes[nodeIdx].Child[bit]
	if childIdx == missingNode {
		return
	}

	nextIdx := idx + 1

	// STRATEGY 1: CONTINUE
	// Extend the current prefix without committing to a letter boundary.
	t.parallelTraverse(sequence, nextIdx, childIdx, results, prefix)

	// STRATEGY 2: BRANCH
	// If the current path forms a valid letter, commit it and restart at root.
	if val := t.Nodes[childIdx].Val; val != 0 {
		newPrefix := prefix + string(val)
		t.parallelTraverse(sequence, nextIdx, rootIdx, results, newPrefix)
	}
}

// FindCandidates provides a package-level function for finding candidates.
func FindCandidatesParallel(morseCode string) []string {
	return StaticTrie.FindCandidatesParallel(morseCode)
}
