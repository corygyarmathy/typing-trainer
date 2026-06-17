package config

import "log/slog"

type Secret string

func (Secret) String() string       { return "REDACTED" }
func (Secret) LogValue() slog.Value { return slog.StringValue("REDACTED") }
func (s Secret) Reveal() string     { return string(s) } // explicit, greppable
