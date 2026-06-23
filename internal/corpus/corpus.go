// Package corpus provides the read-only language data the adaptive engine
// needs: a frequency order over letters, a frequency-ranked list of ngrams,
// and a context -> next-character transition graph for pseudo-word generation.
//
// Per ADR 0013 this is static reference data, derived from a source corpus by
// a committed generator (cmd/corpusgen) and embedded in the binary with
// go:embed. It is identical for every user and every deployment, so it is NOT
// a database-backed bounded context: there is no handler.go, no repository.go,
// and no SQL. The package serves its data in-process, which is what lets the
// engine run both behind the API and in the offline / anonymous client
// (ADR 0014).
package corpus

// Candidate is one possible next character in the transition graph, paired
// with the base frequency of the ngram it would form. Consumed by the
// generator in internal/adaptive.
type Candidate struct {
	Char rune
	Freq float64
}

// Provider serves the embedded corpus. It satisfies the Corpus interface that
// internal/adaptive defines and consumes: the consumer owns the interface,
// this package owns the data.
//
// TODO(phase-3): back this with the go:embed'd generated artifact, e.g.
//
//	//go:embed data/corpus.json
//	var corpusData []byte
//
// and implement the methods the engine requires:
//   - StartingKeys() int
//   - KeyOrder() []rune                 // frequency order, for unlocking
//   - NgramsByFrequency() []string      // frequency-ranked, defines tiers
//   - Transitions(context string) []Candidate
type Provider struct{}
