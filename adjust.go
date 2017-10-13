package main

import (
	"sort"
)

func AdjustPlayers(QBs []*Player, RBs []*Player, WRs []*Player, TEs []*Player, DSTs []*Player) {
	medianRushAllowed := getMedianRushAllowed(DSTs)
	medianPassAllowed := getMedianPassAllowed(DSTs)
	medianPassAttempts := getMedianPassAttempts(QBs)
	medianRushAttempts := getMedianRushAttempts(RBs)
	medianReceptionsWR := getMedianReceptions(WRs)
	medianReceptionsTE := getMedianReceptions(TEs)

	medianPointsAllowed := getMedianPointsAllowedQB(DSTs)
	adjustQBs(QBs, medianPassAttempts, medianPointsAllowed, medianPassAllowed)

	medianPointsAllowed = getMedianPointsAllowedRB(DSTs)
	adjustRBs(RBs, medianRushAllowed, medianPointsAllowed, medianRushAttempts)

	medianPointsAllowed = getMedianPointsAllowedWR(DSTs)
	adjustWRs(WRs, medianPassAllowed, medianPointsAllowed, medianReceptionsWR)

	medianPointsAllowed = getMedianPointsAllowedTE(DSTs)
	adjustTEs(TEs, medianPassAllowed, medianPointsAllowed, medianReceptionsTE)

	medianPointsAllowed = getMedianPointsAllowedOffense(DSTs)
	adjustDSTs(DSTs, medianPointsAllowed)
}

func getMedianPointsAllowedQB(DSTs []*Player) float64 {
	fs := make([]float64, 0)
	for _, d := range DSTs {
		fs = append(fs, d.defenseStats.pointsToQB)
	}
	sort.Float64s(fs)
	return fs[len(fs)/2]
}

func getMedianPointsAllowedRB(DSTs []*Player) float64 {
	fs := make([]float64, 0)
	for _, d := range DSTs {
		fs = append(fs, d.defenseStats.pointsToRB)
	}
	sort.Float64s(fs)
	return fs[len(fs)/2]
}

func getMedianPointsAllowedWR(DSTs []*Player) float64 {
	fs := make([]float64, 0)
	for _, d := range DSTs {
		fs = append(fs, d.defenseStats.pointsToWR)
	}
	sort.Float64s(fs)
	return fs[len(fs)/2]
}

func getMedianPointsAllowedTE(DSTs []*Player) float64 {
	fs := make([]float64, 0)
	for _, d := range DSTs {
		fs = append(fs, d.defenseStats.pointsToTE)
	}
	sort.Float64s(fs)
	return fs[len(fs)/2]
}

func getMedianPointsAllowedOffense(DSTs []*Player) float64 {
	fs := make([]float64, 0)
	for _, d := range DSTs {
		fs = append(fs, d.defenseStats.opposingOffensePoints)
	}
	sort.Float64s(fs)
	return fs[len(fs)/2]
}

func getMedianRushAllowed(DSTs []*Player) float64 {
	fs := make([]float64, 0)
	for _, d := range DSTs {
		fs = append(fs, d.defenseStats.rushYardsAllowed)
	}
	sort.Float64s(fs)
	return fs[len(fs)/2]
}

func getMedianPassAllowed(DSTs []*Player) float64 {
	fs := make([]float64, 0)
	for _, d := range DSTs {
		fs = append(fs, d.defenseStats.passYardsAllowed)
	}
	sort.Float64s(fs)
	return fs[len(fs)/2]
}

func getMedianPassAttempts(QBs []*Player) float64 {
	fs := make([]float64, 0)
	for _, d := range QBs {
		fs = append(fs, d.passingStats.passingAttempts)
	}
	sort.Float64s(fs)
	return fs[len(fs)/2]
}

func getMedianRushAttempts(RBs []*Player) float64 {
	fs := make([]float64, 0)
	for _, d := range RBs {
		fs = append(fs, d.rushingStats.rushingAttempts)
	}
	sort.Float64s(fs)
	return fs[len(fs)/2]
}

func getMedianReceptions(WRs []*Player) float64 {
	fs := make([]float64, 0)
	for _, d := range WRs {
		fs = append(fs, d.receivingStats.Receptions+d.receivingStats.Targets)
	}
	sort.Float64s(fs)
	return fs[len(fs)/2]
}

func adjustQBs(QBs []*Player, medianPassAttempts, medianPointsAllowed, medianPassAllowed float64) {
	for _, q := range QBs {
		opponentPassAllowed := q.oppositeObject.defenseStats.passYardsAllowed
		opponentPointsAllowed := q.oppositeObject.defenseStats.pointsToQB
		defMultipler := opponentPassAllowed / medianPassAllowed
		attemptMultipler := q.passingStats.passingAttempts / medianPassAttempts
		finalMultipler := ((defMultipler + attemptMultipler) / 2) + (opponentPointsAllowed / medianPointsAllowed)
		q.ProjectedPoints *= finalMultipler
	}
}

func adjustRBs(RBs []*Player, medianRushAllowed, medianPointsAllowed, medianRushAttempts float64) {
	for _, r := range RBs {
		opponentRushAllowed := r.oppositeObject.defenseStats.rushYardsAllowed
		defMultipler := opponentRushAllowed / medianRushAllowed
		attemptMultipler := r.rushingStats.rushingAttempts / medianRushAttempts
		opponentPointsAllowed := r.oppositeObject.defenseStats.pointsToRB
		finalMultipler := ((defMultipler + attemptMultipler) / 2) + (opponentPointsAllowed / medianPointsAllowed)
		r.ProjectedPoints *= finalMultipler
	}
}

func adjustWRs(WRs []*Player, medianPassAllowed, medianPointsAllowed, medianReceptions float64) {
	for _, p := range WRs {
		opponentPassAllowed := p.oppositeObject.defenseStats.passYardsAllowed
		opponentPointsAllowed := p.oppositeObject.defenseStats.pointsToWR
		defMultipler := opponentPassAllowed / medianPassAllowed
		attemptMultipler := p.receivingStats.Receptions / medianReceptions
		finalMultipler := ((defMultipler + attemptMultipler) / 2) + (opponentPointsAllowed / medianPointsAllowed)
		p.ProjectedPoints *= finalMultipler
	}
}

func adjustTEs(TEs []*Player, medianPassAllowed, medianPointsAllowed, medianReceptions float64) {
	for _, p := range TEs {
		opponentPassAllowed := p.oppositeObject.defenseStats.passYardsAllowed
		defMultipler := opponentPassAllowed / medianPassAllowed
		opponentPointsAllowed := p.oppositeObject.defenseStats.pointsToTE
		attemptMultipler := p.receivingStats.Receptions / medianReceptions
		finalMultipler := ((defMultipler + attemptMultipler) / 2) + (opponentPointsAllowed / medianPointsAllowed)
		p.ProjectedPoints *= finalMultipler
	}
}

func adjustDSTs(DSTs []*Player, medianPointsAllowed float64) {
	for _, d := range DSTs {
		d.ProjectedPoints *= (d.defenseStats.opposingOffensePoints / medianPointsAllowed)
	}
}
