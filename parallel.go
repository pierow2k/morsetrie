package morsetrie

import (
	"runtime"
	"sync"
)

const (
	// taskBranchFactor represents the maximum branches (continue/commit)
	// generated at a single decoding step.
	taskBranchFactor = 2

	// workerTaskBufferMult determines the capacity of the task channel
	// relative to the number of workers.
	workerTaskBufferMult = 8

	// workerResultCapacity is the pre-allocated capacity for each worker's
	// result slice to minimize allocations during traversal.
	workerResultCapacity = 1024

	// workerTargetMult determines how many independent tasks the producer
	// aims to generate relative to the number of workers.
	workerTargetMult = 4
)

// decodeTask represents an independent exploration point in the trie.
type decodeTask struct {
	idx     int    // Position in the morse sequence
	nodeIdx int16  // Position in the trie
	prefix  string // Accumulated decoded characters
}

// validateMorse checks if the sequence contains only valid Morse symbols.
func validateMorse(sequence string) bool {
	for i := range sequence {
		if sequence[i] != '.' && sequence[i] != '-' {
			return false
		}
	}

	return true
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

// expandBranches generates the next set of tasks from a given position.
// It returns tasks for both continuing the current letter and starting a new one.
func (t *Trie) expandBranches(sequence string, idx int, nodeIdx int16, prefix string) []decodeTask {
	if idx >= len(sequence) {
		return nil
	}

	symbol := sequence[idx]

	var bit int

	switch symbol {
	case '.':
		bit = 0
	case '-':
		bit = 1
	default:
		return nil
	}

	childIdx := t.Nodes[nodeIdx].Child[bit]
	if childIdx == missingNode {
		return nil
	}

	nextIdx := idx + 1
	tasks := make([]decodeTask, 0, taskBranchFactor)

	// CONTINUE path: extend current letter prefix.
	tasks = append(tasks, decodeTask{idx: nextIdx, nodeIdx: childIdx, prefix: prefix})

	// BRANCH path: commit letter and restart at root.
	if val := t.Nodes[childIdx].Val; val != 0 {
		newPrefix := prefix + string(val)
		tasks = append(tasks, decodeTask{idx: nextIdx, nodeIdx: rootIdx, prefix: newPrefix})
	}

	return tasks
}

// produceTasks generates independent tasks via breadth-first expansion
// and sends them to the task channel for workers to process.
func (t *Trie) produceTasks(sequence string, tasks chan<- decodeTask, targetTasks int) {
	frontier := []decodeTask{{idx: 0, nodeIdx: rootIdx, prefix: ""}}

	for len(frontier) > 0 && len(frontier) < targetTasks {
		current := frontier[0]
		frontier = frontier[1:]

		children := t.expandBranches(sequence, current.idx, current.nodeIdx, current.prefix)
		if len(children) > 0 {
			frontier = append(frontier, children...)
		}
	}

	for _, task := range frontier {
		tasks <- task
	}
}

// mergeResults combines per-worker result slices into a single slice.
func mergeResults(workerResults [][]string) []string {
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

// FindCandidatesParallel returns all valid decodings for a Morse code
// sequence. It parallelizes exploration by distributing branches from a
// breadth-first frontier to worker goroutines.
//
// Concurrency model:
//   - Producer: BFS traversal to generate ~numWorkers independent tasks
//   - Workers: Pull tasks and recursively explore the subtree
//   - Accumulation: Per-worker slices merged at the end
func (t *Trie) FindCandidatesParallel(sequence string) []string {
	if len(sequence) == 0 || !validateMorse(sequence) {
		return nil
	}

	numWorkers := runtime.GOMAXPROCS(0)
	tasks := make(chan decodeTask, numWorkers*workerTaskBufferMult)
	workerResults := make([][]string, numWorkers)

	var waitGroup sync.WaitGroup

	// Start workers.
	for worker := range numWorkers {
		waitGroup.Go(func() {
			localResults := make([]string, 0, workerResultCapacity)
			for task := range tasks {
				t.parallelTraverse(sequence, task.idx, task.nodeIdx, &localResults, task.prefix)
			}

			workerResults[worker] = localResults
		})
	}

	// Producer: Generate tasks via breadth-first expansion.
	t.produceTasks(sequence, tasks, numWorkers*workerTargetMult)
	close(tasks)

	waitGroup.Wait()

	return mergeResults(workerResults)
}

// FindCandidatesParallel provides a package-level function for finding candidates.
func FindCandidatesParallel(morseCode string) []string {
	return StaticTrie.FindCandidatesParallel(morseCode)
}
