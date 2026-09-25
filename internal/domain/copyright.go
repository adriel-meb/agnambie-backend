// File: internal/domain/copyright.go
// Purpose: Copyright types matching the Bible Brain copyright response.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package domain

// CopyrightEntry represents one fileset's copyright from
// GET /bibles/{bible_id}/copyright.
type CopyrightEntry struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Size      string         `json:"size"`
	Copyright CopyrightInfo  `json:"copyright"`
}

// CopyrightInfo holds the copyright details for a fileset.
type CopyrightInfo struct {
	CopyrightDate        string `json:"copyright_date"`
	Copyright            string `json:"copyright"`
	CopyrightDescription string `json:"copyright_description"`
	OpenAccess           int    `json:"open_access"`
}
