package main

import (
	"sort"
)

func AdjustPlayers(QBs []*Player, RBs []*Player, WRs []*Player, TEs []*Player, DSTs []*Player) {
	medianRushAllowed := getMedianRushAllowed(DSTs)
	medianPassAllowed := getMedianPassAllowed(DSTs)
	medianPointsAllowed := getMedianPointsAllowed(DSTs)
	medianPassAttempts := getMedianPassAttempts(QBs)
	medianRushAttempts := getMedianRushAttempts(RBs)
	medianReceptions := getMedianReceptions(WRs)
	adjustQBs(QBs, medianPassAttempts, medianPointsAllowed, medianPassAllowed)
	adjustRBs(RBs, medianRushAllowed, medianPointsAllowed, medianRushAttempts)
	adjustReceivers(WRs, TEs, medianPassAllowed, medianPointsAllowed, medianReceptions)
}

func getMedianRushAllowed(DSTs []*Player) float64 {
	sort.Sort(DSTRushAllowed(DSTs))
	return DSTs[len(DSTs)/2].defenseStats.rushYardsAllowed
}

func getMedianPassAllowed(DSTs []*Player) float64 {
	sort.Sort(DSTPassAllowed(DSTs))
	return DSTs[len(DSTs)/2].defenseStats.passYardsAllowed
}

func getMedianPointsAllowed(DSTs []*Player) float64 {
	sort.Sort(DSTPointsAllowed(DSTs))
	return DSTs[len(DSTs)/2].defenseStats.pointsAllowed
}

func getMedianPassAttempts(QBs []*Player) float64 {
	sort.Sort(QBPassAttempts(QBs))
	return QBs[len(QBs)/2].passingStats.passingAttempts
}

func getMedianRushAttempts(RBs []*Player) float64 {
	sort.Sort(RBRushAttempts(RBs))
	return RBs[len(RBs)/2].rushingStats.rushingAttempts
}

func getMedianReceptions(WRs []*Player) float64 {
	sort.Sort(RBRushAttempts(WRs))
	return WRs[len(WRs)/2].receivingStats.Receptions
}

func adjustQBs(QBs []*Player, medianPassAttempts, medianPointsAllowed, medianPassAllowed float64) {
	for _, q := range QBs {
		opponentPassAllowed := q.oppositeObject.defenseStats.passYardsAllowed
		opponentPointsAllowed := q.oppositeObject.defenseStats.pointsAllowed
		defMultipler := opponentPassAllowed/medianPassAllowed + opponentPointsAllowed/medianPointsAllowed
		attemptMultipler := q.passingStats.passingAttempts / medianPassAttempts
		finalMultipler := (defMultipler + attemptMultipler) / 2
		q.ProjectedPoints *= finalMultipler
	}
}

func adjustRBs(RBs []*Player, medianRushAllowed, medianPointsAllowed, medianRushAttempts float64) {
	for _, r := range RBs {
		opponentRushAllowed := r.oppositeObject.defenseStats.rushYardsAllowed
		defMultipler := opponentRushAllowed / medianRushAllowed
		attemptMultipler := r.rushingStats.rushingAttempts / medianRushAttempts
		finalMultipler := (defMultipler + attemptMultipler) / 2
		r.ProjectedPoints *= finalMultipler
	}
}

func adjustReceivers(WRs, TEs []*Player, medianPassAllowed, medianPointsAllowed, medianReceptions float64) {
	for _, p := range WRs {
		opponentPassAllowed := p.oppositeObject.defenseStats.passYardsAllowed
		defMultipler := opponentPassAllowed / medianPassAllowed
		attemptMultipler := p.receivingStats.Receptions / medianReceptions
		finalMultipler := (defMultipler + attemptMultipler) / 2
		p.ProjectedPoints *= finalMultipler
	}

	for _, p := range TEs {
		opponentPassAllowed := p.oppositeObject.defenseStats.passYardsAllowed
		defMultipler := opponentPassAllowed / medianPassAllowed
		attemptMultipler := p.receivingStats.Receptions / medianReceptions
		finalMultipler := (defMultipler + attemptMultipler) / 2
		p.ProjectedPoints *= finalMultipler
	}
}
