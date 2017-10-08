package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sort"
	"sync"
)

func main() {
	p := Parser{}

	QBs, RBs, WRs, TEs, DSTs := p.Parse()

	for _, q := range QBs {
		q.ScoreQB()
	}

	for _, r := range RBs {
		r.ScoreRB()
	}

	for _, w := range WRs {
		w.ScoreWR()
	}

	for _, t := range TEs {
		t.ScoreTE()
	}

	for _, d := range DSTs {
		d.ScoreDST()
	}

	QBsTemp := make([]*Player, 0)
	RBsTemp := make([]*Player, 0)
	WRsTemp := make([]*Player, 0)
	TEsTemp := make([]*Player, 0)
	DSTsTemp := make([]*Player, 0)

	for _, q := range QBs {
		if q.ProjectedPoints >= 5 {
			QBsTemp = append(QBsTemp, q)
		}
	}

	for _, r := range RBs {
		if r.ProjectedPoints >= 5 {
			RBsTemp = append(RBsTemp, r)
		}
	}

	for _, w := range WRs {
		if w.ProjectedPoints >= 5 {
			WRsTemp = append(WRsTemp, w)
		}
	}

	for _, t := range TEs {
		if t.ProjectedPoints >= 3 {
			TEsTemp = append(TEsTemp, t)
		}
	}

	for _, d := range DSTs {
		if d.Salary >= 2500 {
			DSTsTemp = append(DSTsTemp, d)
		}
	}

	QBs = QBsTemp
	RBs = RBsTemp
	WRs = WRsTemp
	TEs = TEsTemp
	DSTs = DSTsTemp

	AdjustPlayers(QBs, RBs, WRs, TEs, DSTs)

	QBsTemp = make([]*Player, 0)
	RBsTemp = make([]*Player, 0)
	WRsTemp = make([]*Player, 0)
	TEsTemp = make([]*Player, 0)
	DSTsTemp = make([]*Player, 0)

	for _, q := range QBs {
		if q.ProjectedPoints >= 15 {
			QBsTemp = append(QBsTemp, q)
			//fmt.Printf("%s - %f\n", q.Name, q.ProjectedPoints)
		}
	}

	for _, r := range RBs {
		if r.ProjectedPoints >= 5 {
			RBsTemp = append(RBsTemp, r)
			//fmt.Printf("%s - %f\n", r.Name, r.ProjectedPoints)
		}
	}

	for _, w := range WRs {
		if w.ProjectedPoints >= 8 {
			WRsTemp = append(WRsTemp, w)
			//fmt.Printf("%s - %f\n", w.Name, w.ProjectedPoints)
		}
	}

	for _, t := range TEs {
		if t.ProjectedPoints >= 6 {
			TEsTemp = append(TEsTemp, t)
			//fmt.Printf("%s - %f\n", t.Name, t.ProjectedPoints)
		}
	}

	QBs = QBsTemp
	RBs = RBsTemp
	WRs = WRsTemp
	TEs = TEsTemp

	/*
		fmt.Printf("%d\n", len(QBs))
		fmt.Printf("%d\n", len(RBs))
		fmt.Printf("%d\n", len(WRs))
		fmt.Printf("%d\n", len(TEs))
		return
	*/
	cores := runtime.NumCPU() - 1
	wgs := make([]sync.WaitGroup, cores)
	rosters := make([]Rosters, cores)

	for i := 0; i < cores; i++ {
		wgs[i].Add(1)
		go func(pid int) {
			rosters[pid] = create(pid, QBs, RBs, WRs, TEs, DSTs)
			wgs[pid].Done()
		}(i)
	}

	for i := 0; i < cores; i++ {
		wgs[i].Wait()
	}

	master := make(Rosters, 0)
	for i := 0; i < cores; i++ {
		master = append(master, rosters[i]...)
	}

	sort.Sort(master)
	if len(master) > 100 {
		master = master[:100]
	}

	/*
		for _, r := range master {
			b, err := json.Marshal(*r)
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("%s\n", string(b))
		}
	*/

	b, _ := json.Marshal(master)
	fmt.Printf("%s\n", string(b))
}

func create(pid int, QBs, RBs, WRs, TEs, DSTs []*Player) Rosters {
	candidates := make(Rosters, 0)

	for i, _ := range QBs {
		if i+pid >= len(QBs) {
			if len(candidates) > 100 {
				candidates = candidates[:100]
			}
			break
		}
		q := QBs[i+pid]
		roster := &Roster{}
		roster.addPlayer(q, false)
		addWR(roster, RBs, WRs, TEs, DSTs, []string{}, &candidates)
		if len(candidates) > 100 {
			sort.Sort(candidates)
			candidates = candidates[:100]
		}
	}

	return candidates
}

func addWR(roster *Roster, RBs, WRs, TEs, DSTs []*Player, teamsNotToChoose []string, candidates *Rosters) {
	for _, wr := range WRs {
		found := false
		for _, team := range teamsNotToChoose {
			if team == wr.team {
				found = true
				break
			}
		}
		if found {
			continue
		}
		if roster.addPlayer(wr, false) {
			teamsNotToChoose = append(teamsNotToChoose, wr.team)

			if len(roster.WRs) >= WRLIMIT {
				addRB(roster, RBs, WRs, TEs, DSTs, teamsNotToChoose, candidates)
			} else {
				addWR(roster, RBs, WRs, TEs, DSTs, teamsNotToChoose, candidates)
			}

			roster.popPlayer(WR)
			teamsNotToChoose = teamsNotToChoose[:len(teamsNotToChoose)-1]
		}
	}
}

func addRB(roster *Roster, RBs, WRs, TEs, DSTs []*Player, teamsNotToChoose []string, candidates *Rosters) {
	for _, rb := range RBs {
		found := false
		for _, team := range teamsNotToChoose {
			if team == rb.team {
				found = true
				break
			}
		}
		if found {
			continue
		}
		if roster.addPlayer(rb, false) {
			teamsNotToChoose = append(teamsNotToChoose, rb.team)
			if len(roster.RBs) >= RBLIMIT {
				addTE(roster, RBs, WRs, TEs, DSTs, teamsNotToChoose, candidates)
			} else {
				addRB(roster, RBs, WRs, TEs, DSTs, teamsNotToChoose, candidates)
			}
			roster.popPlayer(RB)
			teamsNotToChoose = teamsNotToChoose[:len(teamsNotToChoose)-1]
		}
	}
}

func addTE(roster *Roster, RBs, WRs, TEs, DSTs []*Player, teamsNotToChoose []string, candidates *Rosters) {
	for _, te := range TEs {
		if roster.addPlayer(te, false) {
			addFLEX(roster, RBs, WRs, TEs, DSTs, teamsNotToChoose, candidates)
			roster.popPlayer(TE)
		}
	}
}

func addFLEX(roster *Roster, RBs, WRs, TEs, DSTs []*Player, teamsNotToChoose []string, candidates *Rosters) {
	flexOptions := make([]*Player, 0)
	temp := make(PlayersPoints, 0)
	for _, p := range RBs {
		flexOptions = append(flexOptions, p)
	}
	for _, p := range WRs {
		flexOptions = append(flexOptions, p)
	}
	sort.Sort(temp)
	for _, f := range temp {
		found := false
		for _, team := range teamsNotToChoose {
			if team == f.team {
				found = true
				break
			}
		}
		if found {
			continue
		}

		if roster.canAfford(f) {
			flexOptions = append(flexOptions, f)
		}

		if len(flexOptions) > 2 {
			break
		}
	}

	for _, p := range flexOptions {
		if roster.addPlayer(p, true) {
			addDST(roster, RBs, WRs, TEs, DSTs, teamsNotToChoose, candidates)
			roster.popPlayer(FLEX)
		}
	}
}

func addDST(roster *Roster, RBs, WRs, TEs, DSTs []*Player, teamsNotToChoose []string, candidates *Rosters) {
	for _, d := range DSTs {
		if roster.addPlayer(d, false) {
			if roster.Points >= 140 {
				cpy := roster.Copy()
				*candidates = append(*candidates, cpy)
				sort.Sort(candidates)
				if len(*candidates) > 100 {
					*candidates = (*candidates)[:100]
				}
			}
			roster.popPlayer(DST)
		}
	}
}
