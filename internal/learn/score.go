package learn

import (
	"fmt"
	"math"
	"time"

	"github.com/VieiraGabrielAlexandre/aliasgen/internal/store"
)

type Suggestion struct {
	Alias   string
	Command string
	Score   float64
	Reason  string
}

func GenerateSuggestions(stats []store.CommandStat, now time.Time) []Suggestion {
	var out []Suggestion
	for _, s := range stats {
		alias := Abbrev(s.Command)
		if alias == "" || len(alias) >= len(s.Command) {
			continue
		}
		rec := recencyWeight(now.Sub(s.LastUsed))
		saved := float64(len(s.Command) - len(alias))
		score := 0.6*float64(s.Uses) + 0.3*saved + 0.1*rec
		out = append(out, Suggestion{
			Alias:   alias,
			Command: s.Command,
			Score:   score,
			Reason:  reasonText(s.Uses, saved, rec),
		})
	}
	// dedup por alias mantendo maior score
	best := map[string]Suggestion{}
	for _, sg := range out {
		if cur, ok := best[sg.Alias]; !ok || sg.Score > cur.Score {
			best[sg.Alias] = sg
		}
	}
	out = out[:0]
	for _, v := range best {
		out = append(out, v)
	}
	return out
}

func recencyWeight(d time.Duration) float64 {
	// quanto mais recente, maior peso; curva suave
	h := d.Hours()
	return 1.0 / (1.0 + math.Log1p(h/24.0))
}

func reasonText(uses int, saved float64, rec float64) string {
	return fmt.Sprintf("freq=%d, ganho_teclas=%.0f, rec=%.2f", uses, saved, rec)
}
