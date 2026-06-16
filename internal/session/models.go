// Package session holds the domain types, business logic, and persistence
// for the session bounded context.
//
// Layering:
//
//	handler.go      - HTTP transport (request parsing, response shaping)
//	service.go      - business logic, orchestrates repositories
//	repository.go   - data access, returns domain types not raw rows
//	models.go       - domain types shared across the layers
package session

// TODO: domain types live here.
