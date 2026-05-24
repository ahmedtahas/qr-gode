package qrgode

import "fmt"

// Practical safe ratios of logo-area / QR-area at each error-correction level.
// These are below the theoretical recovery capacities (7/15/25/30%) to leave
// margin for finder/timing/format modules and real-world scanner variance.
var safeLogoAreaFraction = map[ErrorCorrectionLevel]float64{
	LevelL: 0.05,
	LevelM: 0.10,
	LevelQ: 0.18,
	LevelH: 0.23,
}

var eclName = map[ErrorCorrectionLevel]string{
	LevelL: "L", LevelM: "M", LevelQ: "Q", LevelH: "H",
}

// ScannabilityWarnings inspects the configuration for combinations that
// commonly produce unscannable QR codes — currently, a logo that occupies
// more of the QR area than the chosen error-correction level can recover.
//
// Warnings are heuristic and free of side effects: nothing is printed. Call
// this before SVG()/PNG() if you want to surface them to your users.
//
// Returns nil if no logo is configured or the configuration is safe.
func (q *QRCode) ScannabilityWarnings() []string {
	if q.config == nil || q.config.Logo == nil {
		return nil
	}
	logo := q.config.Logo
	if logo.Image == nil && logo.Path == "" {
		return nil
	}

	logoFraction := q.estimatedLogoFraction()
	limit, ok := safeLogoAreaFraction[q.config.ErrorCorrection]
	if !ok {
		return nil
	}
	if logoFraction <= limit {
		return nil
	}

	return []string{fmt.Sprintf(
		"logo covers ~%.0f%% of the QR area (safe limit ~%.0f%% at ECL %s); "+
			"raise ErrorCorrection (e.g. LevelH) or shrink the logo to keep the code scannable",
		logoFraction*100, limit*100, eclName[q.config.ErrorCorrection],
	)}
}

// estimatedLogoFraction returns the logo's projected area as a fraction of
// the QR's pixel area. For auto-sized logos we use the configured midpoint
// (~22.5% per dimension) as the worst-case estimate.
func (q *QRCode) estimatedLogoFraction() float64 {
	logo := q.config.Logo
	qrSide := float64(q.config.Size)
	if qrSide <= 0 {
		return 0
	}

	var w, h float64
	switch {
	case logo.Width > 0 && logo.Height > 0:
		w, h = float64(logo.Width), float64(logo.Height)
	case logo.Width > 0:
		w = float64(logo.Width)
		h = w // worst case for the unknown dimension
	case logo.Height > 0:
		h = float64(logo.Height)
		w = h
	default:
		// Auto-sized: midpoint of [logoMinSize, logoMaxSize] on the longer side.
		mid := (logoMinSize + logoMaxSize) / 2
		w = qrSide * mid
		h = qrSide * mid
	}
	return (w * h) / (qrSide * qrSide)
}
