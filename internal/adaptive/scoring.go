package adaptive

// Scoring is a separate file because the math is interesting enough to
// deserve its own test file (scoring_test.go) and grows independently of
// lesson selection.
//
// TODO(phase-3):
//   - Score type with combined accuracy + speed components
//   - Accuracy weight should dominate (suggested 0.7 / 0.3 split)
//   - Decay model: a score that hasn't been exercised in a while should
//     drift downward so the engine revisits it
//   - Threshold constants for unlocking new keys and advancing ngram tiers
