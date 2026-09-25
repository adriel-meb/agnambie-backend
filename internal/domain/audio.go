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
	ChapterStart  int      `json:"chapter_start"`
	ChapterEnd    *int     `json:"chapter_end"`
	VerseStart    *string  `json:"verse_start"`
	VerseEnd      *string  `json:"verse_end"`
	Timestamp     *float64 `json:"timestamp"`
	Path          string   `json:"path"`            // CDN URL for the audio file
	Duration      float64  `json:"duration"`         // Duration in seconds
	FilesizeBytes int64    `json:"filesize_in_bytes"` // File size for download estimation
	Thumbnail     *string  `json:"thumbnail"`
}
