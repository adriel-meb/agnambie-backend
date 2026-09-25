// File: internal/domain/gabon_test.go
// Purpose: Unit tests for the Gabon allowlist logic.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package domain

import (
	"testing"
)

func TestIsAllowedLanguage(t *testing.T) {
	tests := []struct {
		code string
		want bool
	}{
		{"FAN", true},
		{"fan", true}, // Case insensitive
		{"FRA", true},
		{"ENG", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			if got := IsAllowedLanguage(tt.code); got != tt.want {
				t.Errorf("IsAllowedLanguage(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestIsAllowedFileset(t *testing.T) {
	tests := []struct {
		filesetID string
		want      bool
	}{
		{"FANBSGN2DA", true},
		{"FANBSGN2DA-opus16", true}, // Data-saver variant
		{"FRNTLSO2DA-opus16", true},
		{"ENGESVN1DA", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.filesetID, func(t *testing.T) {
			if got := IsAllowedFileset(tt.filesetID); got != tt.want {
				t.Errorf("IsAllowedFileset(%q) = %v, want %v", tt.filesetID, got, tt.want)
			}
		})
	}
}

func TestIsAllowedBibleID(t *testing.T) {
	tests := []struct {
		bibleID string
		want    bool
	}{
		{"FANBSG", true}, // 6 chars from FANBSGN2DA
		{"FRNTLS", true},
		{"ENGESV", false},
		{"FAN", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.bibleID, func(t *testing.T) {
			if got := IsAllowedBibleID(tt.bibleID); got != tt.want {
				t.Errorf("IsAllowedBibleID(%q) = %v, want %v", tt.bibleID, got, tt.want)
			}
		})
	}
}

func TestFindLanguage(t *testing.T) {
	lang, ok := FindLanguage("fan")
	if !ok || lang == nil {
		t.Fatal("FindLanguage('fan') expected to find Language")
	}
	if lang.Code != "FAN" {
		t.Errorf("Expected language code FAN, got %s", lang.Code)
	}

	_, ok = FindLanguage("ENG")
	if ok {
		t.Fatal("FindLanguage('ENG') expected to not find Language")
	}
}

func TestAllFilesetIDs(t *testing.T) {
	ids := AllFilesetIDs()
	if len(ids) == 0 {
		t.Fatal("AllFilesetIDs() returned empty list")
	}
	
	// Ensure no opus16 in the AllFilesetIDs output
	for _, id := range ids {
		if id == "" {
			t.Error("AllFilesetIDs() returned empty string element")
		}
		if len(id) > 6 && id[len(id)-7:] == "-opus16" {
			t.Errorf("AllFilesetIDs() returned data-saver variant %s", id)
		}
	}
}
