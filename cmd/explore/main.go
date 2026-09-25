package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// --- Bible Brain response shapes ---

type BBResponse struct {
	Data []BBBible `json:"data"`
}

type BBBible struct {
	Abbr     string                 `json:"abbr"`
	Name     string                 `json:"name"`
	Language string                 `json:"language"`
	Autonym  string                 `json:"autonym"`
	ISO      string                 `json:"iso"`
	Date     string                 `json:"date"`
	Filesets map[string][]BBFileset `json:"filesets"`
}

type BBFileset struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Size      string `json:"size"`
	Bitrate   string `json:"bitrate"`
	Codec     string `json:"codec"`
	Container string `json:"container"`
}

// --- Languages to explore ---

var gabonLanguages = []struct {
	Code string
	Name string
}{
	{Code: "MYE", Name: "Myène"},
	{Code: "FAN", Name: "Fang"},
	{Code: "PUU", Name: "Punu"},
	{Code: "NZB", Name: "Nzebi"},
	{Code: "FRA", Name: "Français"},
	{Code: "TEKE", Name: "Teke"},
	{Code: "BKW", Name: "Bekwel"},
	{Code: "BBG", Name: "Barama"},
	{Code: "BUW", Name: "Bubi"},
	{Code: "DMA", Name: "Duma"},
	{Code: "KEB", Name: "Kélé"},
	{Code: "LUP", Name: "Lumbu"},
	{Code: "ZMN", Name: "Mbangwe"},
	{Code: "NMD", Name: "Ndumu"},
	{Code: "NZB", Name: "Njebi"},
	{Code: "PIC", Name: "Pinji"},
	{Code: "TSV", Name: "Tsogo"},
	{Code: "VIF", Name: "Vili"},
	{Code: "SYX", Name: "Samay"},
	{Code: "BKW", Name: "Bekwel"},
	{Code: "SYI", Name: "Seki"},
	{Code: "BNG", Name: "Benga"},
}

func main() {
	godotenv.Load()

	apiKey := os.Getenv("BIBLE_BRAIN_API_KEY")
	if apiKey == "" {
		fmt.Println("ERROR: BIBLE_BRAIN_API_KEY not set in .env")
		os.Exit(1)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	fmt.Println("=== SCRIPTIA — GABON FILESET EXPLORER ===")
	fmt.Println()

	for _, lang := range gabonLanguages {
		fmt.Printf("────────────────────────────────────────\n")
		fmt.Printf("LANGUAGE: %s (%s)\n", lang.Name, lang.Code)
		fmt.Printf("────────────────────────────────────────\n")

		bibles, err := fetchBibles(client, apiKey, lang.Code)
		if err != nil {
			fmt.Printf("  ERROR: %v\n\n", err)
			continue
		}

		if len(bibles) == 0 {
			fmt.Printf("  No Bibles found — skip this language\n\n")
			continue
		}

		for _, bible := range bibles {
			fmt.Printf("\n  Bible: %s (%s)\n", bible.Name, bible.Abbr)

			prodFilesets, ok := bible.Filesets["dbp-prod"]
			if !ok || len(prodFilesets) == 0 {
				fmt.Printf("  No production filesets\n")
				continue
			}

			// Separate normal and data-saver filesets
			var normal []BBFileset
			var dataSaver []BBFileset

			for _, fs := range prodFilesets {
				if fs.Type == "video_stream" {
					continue // skip video
				}
				if strings.HasSuffix(fs.ID, "-opus16") {
					dataSaver = append(dataSaver, fs)
				} else {
					normal = append(normal, fs)
				}
			}

			fmt.Printf("  Normal quality (mp3 64kbps):\n")
			for _, fs := range normal {
				fmt.Printf("    ID: %-25s  type: %-15s  size: %s\n",
					fs.ID, fs.Type, fs.Size)
			}

			fmt.Printf("  Data Saver (opus 16kbps):\n")
			for _, fs := range dataSaver {
				fmt.Printf("    ID: %-25s  type: %-15s  size: %s\n",
					fs.ID, fs.Type, fs.Size)
			}
		}

		// Print the suggested gabon_config.go entry
		fmt.Printf("\n  >>> Suggested gabon_config.go entry:\n\n")
		printConfigEntry(lang.Code, lang.Name, bibles)
		fmt.Println()
	}

	fmt.Println("=== DONE ===")
}

func fetchBibles(client *http.Client, apiKey, langCode string) ([]BBBible, error) {
	url := fmt.Sprintf(
		"https://4.dbt.io/api/bibles?language_code=%s&page=1&limit=25&v=4&key=%s",
		langCode, apiKey,
	)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	var result BBResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing json: %w", err)
	}

	return result.Data, nil
}

func printConfigEntry(isoCode, displayName string, bibles []BBBible) {
	fmt.Printf("    {\n")
	fmt.Printf("        Code:       %q,\n", isoCode)
	fmt.Printf("        Name:       %q,\n", displayName)
	fmt.Printf("        NativeName: %q,\n", displayName)
	fmt.Printf("        Filesets: []Fileset{\n")

	for _, bible := range bibles {
		prodFilesets, ok := bible.Filesets["dbp-prod"]
		if !ok {
			continue
		}

		// Build a map: base ID → {normal, dataSaver}
		type pair struct {
			Normal    *BBFileset
			DataSaver *BBFileset
		}
		pairs := map[string]*pair{}

		for i, fs := range prodFilesets {
			if fs.Type == "video_stream" {
				continue
			}
			if strings.HasSuffix(fs.ID, "-opus16") {
				baseID := strings.TrimSuffix(fs.ID, "-opus16")
				if pairs[baseID] == nil {
					pairs[baseID] = &pair{}
				}
				pairs[baseID].DataSaver = &prodFilesets[i]
			} else {
				baseID := fs.ID
				if pairs[baseID] == nil {
					pairs[baseID] = &pair{}
				}
				pairs[baseID].Normal = &prodFilesets[i]
			}
		}

		// Prefer audio_drama over audio — print audio_drama first
		printed := map[string]bool{}
		for _, preferType := range []string{"audio_drama", "audio"} {
			for baseID, p := range pairs {
				if printed[baseID] {
					continue
				}
				fs := p.Normal
				if fs == nil {
					continue
				}
				if fs.Type != preferType {
					continue
				}
				printed[baseID] = true

				dataSaverID := ""
				if p.DataSaver != nil {
					dataSaverID = p.DataSaver.ID
				}

				fmt.Printf("            {\n")
				fmt.Printf("                // %s — %s\n", bible.Name, fs.Size)
				fmt.Printf("                ID:          %q,\n", fs.ID)
				fmt.Printf("                DataSaverID: %q,\n", dataSaverID)
				fmt.Printf("                Type:        %q,\n", fs.Type)
				fmt.Printf("                Codec:       %q,\n", fs.Codec)
				fmt.Printf("                Description: %q,\n", bible.Name)
				fmt.Printf("                Size:        %q,\n", fs.Size)
				fmt.Printf("            },\n")
			}
		}
	}

	fmt.Printf("        },\n")
	fmt.Printf("    },\n")
}
