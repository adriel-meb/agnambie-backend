// File: internal/domain/bible.go
// Purpose: Bible, book, and fileset summary types matching Bible Brain responses.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package domain

// BibleSummary represents a Bible translation from a listing response.
// Fields match the Bible Brain GET /bibles response shape.
type BibleSummary struct {
	ID              string                     `json:"id"`
	Abbr            string                     `json:"abbr"`
	Name            string                     `json:"name"`
	VName           string                     `json:"vname"`
	Language        string                     `json:"language"`
	LanguageID      int                        `json:"language_id"`
	LanguageAutonym string                     `json:"language_autonym"`
	ISO             string                     `json:"iso"`
	Date            int                        `json:"date"`
	Filesets        map[string][]FilesetEntry  `json:"filesets"`
}

// FilesetEntry is a fileset within a Bible listing response.
type FilesetEntry struct {
	ID          string `json:"id"`
	SetTypeCode string `json:"set_type_code"`
	SetSizeCode string `json:"set_size_code"`
}

// Book represents a book entry from the Bible Brain GET /bibles/{id}/book response.
type Book struct {
	BookID         string `json:"book_id"`
	BookIDUsfx     string `json:"book_id_usfx"`
	BookIDOsis     string `json:"book_id_osis"`
	Name           string `json:"name"`
	Testament      string `json:"testament"`
	TestamentOrder int    `json:"testament_order"`
	BookOrder      int    `json:"book_order"`
	BookGroup      string `json:"book_group"`
	Chapters       []int  `json:"chapters"`
}
