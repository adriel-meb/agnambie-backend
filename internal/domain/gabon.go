// File: internal/domain/gabon.go
// Purpose: Gabon language/fileset allowlist data and lookup helpers.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package domain

import "strings"

// GabonLanguages is the complete, authoritative list of Gabonese languages
// and their audio filesets available on Bible Brain. This is the single source
// of truth used to enforce the Gabon-only scope server-side.
var GabonLanguages = []Language{
	{
		Code: "MYE", Name: "Myène", NativeName: "Myène",
		Filesets: []Fileset{
			{ID: "MYEBSGN2DA", DataSaverID: "MYEBSGN2DA-opus16", Type: "audio_drama", Codec: "mp3", Description: "1999 Bible Society of Gabon", Size: "NT"},
			{ID: "MYEBSGN1DA", DataSaverID: "MYEBSGN1DA-opus16", Type: "audio", Codec: "mp3", Description: "1999 Bible Society of Gabon", Size: "NT"},
			{ID: "MYEDPIO1DA", DataSaverID: "MYEDPIO1DA-opus16", Type: "audio", Codec: "mp3", Description: "Davar Partners International", Size: "OT"},
		},
	},
	{
		Code: "FAN", Name: "Fang", NativeName: "Fang",
		Filesets: []Fileset{
			{ID: "FANBSGN2DA", DataSaverID: "FANBSGN2DA-opus16", Type: "audio_drama", Codec: "mp3", Description: "Fang - Bible Society of Gabon", Size: "NT"},
		},
	},
	{
		Code: "PUU", Name: "Punu", NativeName: "Punu",
		Filesets: []Fileset{
			{ID: "PUUBSGN2DA", DataSaverID: "PUUBSGN2DA-opus16", Type: "audio_drama", Codec: "mp3", Description: "1992 Bible Society of Gabon", Size: "NT"},
			{ID: "PUUBSGN1DA", DataSaverID: "PUUBSGN1DA-opus16", Type: "audio", Codec: "mp3", Description: "1992 Bible Society of Gabon", Size: "NT"},
			{ID: "PUUDPIO1DA", DataSaverID: "PUUDPIO1DA-opus16", Type: "audio", Codec: "mp3", Description: "1992, Bible Society of Cameroon, Bible Society of Gabon", Size: "OT"},
		},
	},
	{
		Code: "NZB", Name: "Nzebi", NativeName: "Nzebi",
		Filesets: []Fileset{
			{ID: "NZBDPIO1DA", DataSaverID: "NZBDPIO1DA-opus16", Type: "audio", Codec: "mp3", Description: "2022, Bible Society of Gabon", Size: "OT"},
			{ID: "NZBUBSN2DA", DataSaverID: "NZBUBSN2DA-opus16", Type: "audio_drama", Codec: "mp3", Description: "1978 Bible Society of Gabon", Size: "NT"},
		},
	},
	{
		Code: "FRA", Name: "Français", NativeName: "Français",
		Filesets: []Fileset{
			{ID: "FRADPIN1DA", DataSaverID: "FRADPIN1DA-opus16", Type: "audio", Codec: "mp3", Description: "Parole de Vie - African (Davar Audio)", Size: "NT"},
			{ID: "FRADPIO1DA", DataSaverID: "FRADPIO1DA-opus16", Type: "audio", Codec: "mp3", Description: "Parole de Vie - African (Davar Audio)", Size: "OT"},
			{ID: "FRNLSNO1DA", DataSaverID: "FRNLSNO1DA-opus16", Type: "audio", Codec: "mp3", Description: "French - La Nouvelle Bible Louis Segond 1978", Size: "OT"},
			{ID: "FRNLSNN1DA16", DataSaverID: "", Type: "audio", Codec: "mp3", Description: "French - La Nouvelle Bible Louis Segond 1978", Size: "NT"},
			{ID: "FRNLSNN1DA", DataSaverID: "", Type: "audio", Codec: "mp3", Description: "French - La Nouvelle Bible Louis Segond 1978", Size: "NT"},
			{ID: "FRNPDVN2DA", DataSaverID: "", Type: "audio_drama", Codec: "mp3", Description: "Parole de Vie - African", Size: "NT"},
			{ID: "FRNPDCN2DA", DataSaverID: "", Type: "audio_drama", Codec: "mp3", Description: "Parole de Vie - Canadian", Size: "NT"},
			{ID: "FRNPDCN2DA16", DataSaverID: "", Type: "audio_drama", Codec: "mp3", Description: "Parole de Vie - Canadian", Size: "NT"},
			{ID: "FRNTLSN2DA", DataSaverID: "FRNTLSN2DA-opus16", Type: "audio_drama", Codec: "mp3", Description: "French - 1910 Louis Segond (Tresorsonore recording)", Size: "NT"},
			{ID: "FRNTLSO2DA", DataSaverID: "FRNTLSO2DA-opus16", Type: "audio_drama", Codec: "mp3", Description: "French - 1910 Louis Segond (Tresorsonore recording)", Size: "OT"},
		},
	},
	{
		Code: "BKW", Name: "Bekwel", NativeName: "Bekwel",
		Filesets: []Fileset{
			{ID: "BKWWBTP1DA", DataSaverID: "BKWWBTP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Bekwel - LORA WBT", Size: "NTP"},
		},
	},
	{
		Code: "BBG", Name: "Barama", NativeName: "Barama",
		Filesets: []Fileset{
			{ID: "BBGCIEP1DA", DataSaverID: "BBGCIEP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Varama - CIEATLG", Size: "NTP"},
		},
	},
	{
		Code: "BUW", Name: "Bubi", NativeName: "Bubi",
		Filesets: []Fileset{
			{ID: "BUWCIEP1DA", DataSaverID: "BUWCIEP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Bubi - CIEATLG", Size: "NTP"},
		},
	},
	{
		Code: "DMA", Name: "Duma", NativeName: "Duma",
		Filesets: []Fileset{
			{ID: "DMACIEP1DA", DataSaverID: "DMACIEP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Duma - CIEATLG", Size: "NTP"},
		},
	},
	{
		Code: "KEB", Name: "Kélé", NativeName: "Kélé",
		Filesets: []Fileset{
			{ID: "KEBCIEP1DA", DataSaverID: "KEBCIEP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Kele - CIEATLG", Size: "NTP"},
		},
	},
	{
		Code: "LUP", Name: "Lumbu", NativeName: "Lumbu",
		Filesets: []Fileset{
			{ID: "LUPCIEP1DA", DataSaverID: "LUPCIEP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Lumbu - CIEATLG", Size: "NTP"},
		},
	},
	{
		Code: "ZMN", Name: "Mbangwe", NativeName: "Mbangwe",
		Filesets: []Fileset{
			{ID: "ZMNCIEP1DA", DataSaverID: "ZMNCIEP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Mbangwe - CIEATLG", Size: "NTP"},
		},
	},
	{
		Code: "NMD", Name: "Ndumu", NativeName: "Ndumu",
		Filesets: []Fileset{
			{ID: "NMDCIEP1DA", DataSaverID: "NMDCIEP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Ndumu - CIEATLG", Size: "NTP"},
		},
	},
	{
		Code: "PIC", Name: "Pinji", NativeName: "Pinji",
		Filesets: []Fileset{
			{ID: "PICCIEP1DA", DataSaverID: "PICCIEP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Pinji - 2023 CIEATLG", Size: "NTP"},
		},
	},
	{
		Code: "TSV", Name: "Tsogo", NativeName: "Tsogo",
		Filesets: []Fileset{
			{ID: "TSVBSGN2DA", DataSaverID: "TSVBSGN2DA-opus16", Type: "audio_drama", Codec: "mp3", Description: "1983 Bible Society of Gabon", Size: "NT"},
			{ID: "TSVBSGN1DA", DataSaverID: "TSVBSGN1DA-opus16", Type: "audio", Codec: "mp3", Description: "1983 Bible Society of Gabon", Size: "NT"},
			{ID: "TSVDPIO1DA", DataSaverID: "TSVDPIO1DA-opus16", Type: "audio", Codec: "mp3", Description: "2023, Bible Society of Gabon", Size: "OT"},
		},
	},
	{
		Code: "VIF", Name: "Vili", NativeName: "Vili",
		Filesets: []Fileset{
			{ID: "VIFBSCN2DA", DataSaverID: "VIFBSCN2DA-opus16", Type: "audio_drama", Codec: "mp3", Description: "Bible Society of Congo and WBT", Size: "NT"},
			{ID: "VIFBSCN1DA", DataSaverID: "VIFBSCN1DA-opus16", Type: "audio", Codec: "mp3", Description: "Bible Society of Congo and WBT", Size: "NT"},
			{ID: "VIFCIEP1DA", DataSaverID: "VIFCIEP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Vili - CIEATLG", Size: "NTP"},
		},
	},
	{
		Code: "SYX", Name: "Samay", NativeName: "Samay",
		Filesets: []Fileset{
			{ID: "SYXCIEP1DA", DataSaverID: "SYXCIEP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Samay - CIEATLG", Size: "NTP"},
		},
	},
	{
		Code: "SYI", Name: "Seki", NativeName: "Seki",
		Filesets: []Fileset{
			{ID: "SYICIEP1DA", DataSaverID: "SYICIEP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Seki - CIEATLG", Size: "NTP"},
		},
	},
	{
		Code: "BNG", Name: "Benga", NativeName: "Benga",
		Filesets: []Fileset{
			{ID: "BNGCIEP1DA", DataSaverID: "BNGCIEP1DA-opus16", Type: "audio", Codec: "mp3", Description: "Benga - 2023 CIEATLG", Size: "NTP"},
		},
	},
}

// --- Allowlist lookup maps (built once at init time) ---

var (
	allowedLanguageCodes map[string]bool
	allowedFilesetIDs    map[string]bool
	allowedBibleIDs      map[string]bool
	languageByCode       map[string]*Language
)

func init() {
	allowedLanguageCodes = make(map[string]bool, len(GabonLanguages))
	allowedFilesetIDs = make(map[string]bool, len(GabonLanguages)*4)
	allowedBibleIDs = make(map[string]bool, len(GabonLanguages)*4)
	languageByCode = make(map[string]*Language, len(GabonLanguages))

	for i := range GabonLanguages {
		lang := &GabonLanguages[i]
		code := strings.ToUpper(lang.Code)
		allowedLanguageCodes[code] = true
		languageByCode[code] = lang

		for _, fs := range lang.Filesets {
			allowedFilesetIDs[fs.ID] = true
			if fs.DataSaverID != "" {
				allowedFilesetIDs[fs.DataSaverID] = true
			}
			if len(fs.ID) >= 6 {
				allowedBibleIDs[fs.ID[:6]] = true
			}
		}
	}
}

// IsAllowedLanguage returns true if the language code is in the Gabon allowlist.
// The check is case-insensitive.
func IsAllowedLanguage(code string) bool {
	return allowedLanguageCodes[strings.ToUpper(code)]
}

// IsAllowedFileset returns true if the fileset ID (including -opus16 variants)
// is in the Gabon allowlist.
func IsAllowedFileset(filesetID string) bool {
	return allowedFilesetIDs[filesetID]
}

// IsAllowedBibleID returns true if the Bible abbreviation (e.g. "FANBSG") is
// derived from an allowed fileset.
func IsAllowedBibleID(bibleID string) bool {
	return allowedBibleIDs[bibleID]
}

// FindLanguage returns the Language for an ISO code, or nil if not found.
func FindLanguage(code string) (*Language, bool) {
	lang, ok := languageByCode[strings.ToUpper(code)]
	return lang, ok
}

// AllFilesetIDs returns every allowed fileset ID (not including -opus16 variants).
func AllFilesetIDs() []string {
	ids := make([]string, 0, len(allowedFilesetIDs))
	for _, lang := range GabonLanguages {
		for _, fs := range lang.Filesets {
			ids = append(ids, fs.ID)
		}
	}
	return ids
}
