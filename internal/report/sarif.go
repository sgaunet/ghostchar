package report

import (
	"encoding/json"
	"io"

	"github.com/sgaunet/ghostchar/internal/charset"
	"github.com/sgaunet/ghostchar/internal/scanner"
)

const sarifSchema = "https://docs.oasis-open.org/sarif/sarif/v2.1.0/schemas/sarif-schema-2.1.0.json"
const sarifVersion = "2.1.0"

type sarifLog struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name  string      `json:"name"`
	Rules []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string           `json:"id"`
	ShortDescription sarifMessage     `json:"shortDescription"`
	DefaultConfig    sarifRuleConfig  `json:"defaultConfiguration"`
}

type sarifRuleConfig struct {
	Level string `json:"level"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifResult struct {
	RuleID    string           `json:"ruleId"`
	Message   sarifMessage     `json:"message"`
	Level     string           `json:"level"`
	Locations []sarifLocation  `json:"locations"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           sarifRegion           `json:"region"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
}

var sarifRules = []sarifRule{
	{
		ID:               "ghostchar/invisible",
		ShortDescription: sarifMessage{Text: "Invisible Unicode character detected"},
		DefaultConfig:    sarifRuleConfig{Level: "warning"},
	},
	{
		ID:               "ghostchar/pua",
		ShortDescription: sarifMessage{Text: "Private Use Area character detected"},
		DefaultConfig:    sarifRuleConfig{Level: "warning"},
	},
	{
		ID:               "ghostchar/bidi",
		ShortDescription: sarifMessage{Text: "Bidirectional control character detected"},
		DefaultConfig:    sarifRuleConfig{Level: "error"},
	},
}

func categoryToRuleID(cat charset.Category) string {
	return "ghostchar/" + string(cat)
}

func categoryToLevel(cat charset.Category) string {
	if cat == charset.Bidi {
		return "error"
	}
	return "warning"
}

// FormatSARIF writes the scan result in SARIF 2.1.0 format.
func FormatSARIF(w io.Writer, result *scanner.ScanResult) error {
	results := make([]sarifResult, 0, len(result.Findings))
	for _, f := range result.Findings {
		results = append(results, sarifResult{
			RuleID:  categoryToRuleID(f.Category),
			Message: sarifMessage{Text: f.CodepointHex() + " " + f.Name},
			Level:   categoryToLevel(f.Category),
			Locations: []sarifLocation{
				{
					PhysicalLocation: sarifPhysicalLocation{
						ArtifactLocation: sarifArtifactLocation{URI: f.File},
						Region: sarifRegion{
							StartLine:   f.Line,
							StartColumn: f.Column,
						},
					},
				},
			},
		})
	}

	log := sarifLog{
		Schema:  sarifSchema,
		Version: sarifVersion,
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:  "ghostchar",
						Rules: sarifRules,
					},
				},
				Results: results,
			},
		},
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(log)
}
