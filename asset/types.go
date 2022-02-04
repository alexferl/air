package asset

import "github.com/davidbyttow/govips/v2/vips"

var StringToInterestingTypes = map[string]vips.Interesting{
	"none":      vips.InterestingNone,
	"centre":    vips.InterestingCentre,
	"entropy":   vips.InterestingEntropy,
	"attention": vips.InterestingAttention,
	"low":       vips.InterestingLow,
	"high":      vips.InterestingHigh,
	"all":       vips.InterestingAll,
}

var InterestingTypesToString = map[vips.Interesting]string{
	vips.InterestingNone:      "none",
	vips.InterestingCentre:    "centre",
	vips.InterestingEntropy:   "entropy",
	vips.InterestingAttention: "attention",
	vips.InterestingLow:       "low",
	vips.InterestingHigh:      "high",
	vips.InterestingAll:       "all",
}

var ImageTypes = invertImageTypes()

func invertImageTypes() map[string]vips.ImageType {
	m := map[string]vips.ImageType{}
	for k, v := range vips.ImageTypes {
		m[v] = k
	}
	return m
}
