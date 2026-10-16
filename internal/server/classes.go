package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"projectdnd/internal/rules"
)

// SRDAttribution is the CC BY 4.0 credit the class data needs, shown by the client (docs/features/rules-v2.md).
const SRDAttribution = `This work includes material from the System Reference Document 5.2.1 ("SRD 5.2.1") by Wizards of the Coast LLC, available at https://www.dndbeyond.com/srd. The SRD 5.2.1 is licensed under the Creative Commons Attribution 4.0 International License, available at https://creativecommons.org/licenses/by/4.0/legalcode.`

// classesDoc is the GET /api/classes body. Tables are indexed by class level - 1.
type classesDoc struct {
	Classes     []rules.Class     `json:"classes"`
	Skills      []rules.Skill     `json:"skills"`
	Proficiency []int             `json:"proficiency"`
	FullSlots   [][9]int          `json:"fullSlots"`
	HalfSlots   [][9]int          `json:"halfSlots"`
	PactSlots   []rules.PactSlots `json:"pactSlots"`
	Attribution string            `json:"attribution"`
}

func classesJSON() []byte {
	d := classesDoc{Classes: rules.Classes(), Skills: rules.Skills(), Attribution: SRDAttribution}
	full, half, pact := rules.Class{Caster: "full"}, rules.Class{Caster: "half"}, rules.Class{Caster: "pact"}
	for lv := rules.MinLevel; lv <= rules.MaxLevel; lv++ {
		d.Proficiency = append(d.Proficiency, rules.ProfBonus(lv))
		d.FullSlots = append(d.FullSlots, rules.SlotsMax(full, lv))
		d.HalfSlots = append(d.HalfSlots, rules.SlotsMax(half, lv))
		d.PactSlots = append(d.PactSlots, rules.PactMax(pact, lv))
	}
	b, err := json.Marshal(d)
	if err != nil {
		panic(err) // static data of known types: cannot fail
	}
	return b
}

// handleClasses serves the static class data (not in RoomState: ~136 KB). It is built once; the ETag lets the
// browser revalidate cheaply (no-cache), so a new binary with new data is never hidden by a stale cache.
func handleClasses() http.HandlerFunc {
	body := classesJSON()
	sum := sha256.Sum256(body)
	etag := `"` + hex.EncodeToString(sum[:8]) + `"`
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("ETag", etag)
		http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(body))
	}
}
