// Package adaptive contains the lesson-selection and difficulty engine.
//
// This is the heart of the application's interesting domain logic. It takes
// a snapshot of a user's competency (per-key and per-ngram scores) plus the
// available corpus, and returns the next lesson tailored to push that user
// toward their next threshold.
//
// Design goals:
//
//  1. PURE FUNCTIONS. The engine touches no I/O. It accepts state, returns
//     state. This makes it trivially testable and means we can reason about
//     its behaviour in isolation. All randomness flows through an injected
//     rand.Source so tests are deterministic.
//
//  2. NO HTTP, NO DATABASE, NO LOGGING. The service layer (progress, session)
//     calls into this package; the engine never calls out.
//
//  3. DECOUPLED FROM PERSISTENCE. The engine works on engine-local types
//     (CompetencyState, Lesson, ScoreUpdate). The progress service is
//     responsible for translating between persisted models and these types.
//
// The two key responsibilities are:
//
//   - Selection: given competency state, choose which keys and ngrams the
//     next lesson should target. Weight toward weak spots; introduce new
//     content as the user crosses thresholds.
//
//   - Scoring: given a completed lesson result, compute the score delta
//     for each key and ngram involved. Accuracy-weighted, with diminishing
//     returns on items already well-mastered.
package adaptive

// CompetencyState is a snapshot of a single user's typing competency at a
// point in time. It contains per-key and per-ngram score state, plus the
// set of unlocked items.
//
// TODO(phase-3): fields TBD. Consider:
//   - UnlockedKeys map[rune]Score
//   - UnlockedNgrams map[string]Score
//   - GlobalAccuracy, GlobalWPM
//   - NgramTier int (current top-N ngrams being practiced)
type CompetencyState struct{}

// Lesson is a generated practice prompt: 10-15 english-like words built
// from the currently-unlocked keys and ngrams, weighted toward weak areas.
type Lesson struct{}

// Engine is the pure-function entry point.
//
// TODO(phase-3):
//   - NextLesson(state CompetencyState, corpus Corpus) Lesson
//   - ApplyResult(state CompetencyState, result Result) CompetencyState
//
// Corpus is provided by internal/corpus (the data) but engine takes it as
// a parameter to keep the dependency pointed downward.
type Engine struct{}
