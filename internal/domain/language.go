// Package domain defines pure data types shared across all layers of the application.
//
// File: internal/domain/language.go
// Purpose: Language and Fileset types for the Gabon audio Bible project.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package domain

// Language represents a Gabonese language and its available audio filesets.
type Language struct {
	Code       string    `json:"code"`        // ISO 639-3 code, e.g. "FAN"
	Name       string    `json:"name"`        // Display name, e.g. "Fang"
	NativeName string    `json:"native_name"` // Name in the language itself
	Filesets   []Fileset `json:"filesets"`     // Available audio filesets
}

// Fileset describes a single audio fileset available on Bible Brain.
type Fileset struct {
	ID          string `json:"id"`           // Bible Brain fileset ID, e.g. "FANBSGN2DA"
	DataSaverID string `json:"data_saver_id"` // Low-bandwidth opus16 variant, e.g. "FANBSGN2DA-opus16"
	Type        string `json:"type"`         // "audio" or "audio_drama"
	Codec       string `json:"codec"`        // "mp3"
	Description string `json:"description"`  // Human-readable source description
	Size        string `json:"size"`         // Coverage: "NT", "OT", "NTP", "C"
}
