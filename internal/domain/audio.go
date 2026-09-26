// File: internal/domain/audio.go
// Purpose: Audio chapter type matching the Bible Brain fileset content response.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package domain

// AudioChapter represents a single audio entry from
// GET /bibles/filesets/{id}/{book}/{chapter}.
type AudioChapter struct {
	BookID        string   `json:"book_id"`
	BookName      string   `json:"book_name"`
	ChapterStart  any      `json:"chapter_start"`
	ChapterEnd    any      `json:"chapter_end"`
	VerseStart    any      `json:"verse_start"`
	VerseEnd      any      `json:"verse_end"`
	Timestamp     any      `json:"timestamp"`
	Path          string   `json:"path"`            // CDN URL for the audio file
	Duration      any      `json:"duration"`         // Duration in seconds
	FilesizeBytes any      `json:"filesize_in_bytes"` // File size for download estimation
	Thumbnail     *string  `json:"thumbnail"`
}
