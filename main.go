package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sort"
	"sync"
)

var globalIndex int
var mutex sync.Mutex

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

	for _, q := range QBs {
		if q.ProjectedPoints >= 15 {
			QBsTemp = append(QBsTemp, q)
		}
	}

	for _, r := range RBs {
		if r.ProjectedPoints >= 8 {
			RBsTemp = append(RBsTemp, r)
		}
	}

	for _, w := range WRs {
		if w.ProjectedPoints >= 10 {
			WRsTemp = append(WRsTemp, w)
		}
	}

	for _, t := range TEs {
		if t.ProjectedPoints >= 7 {
			TEsTemp = append(TEsTemp, t)
		}
	}

	QBs = QBsTemp
	RBs = RBsTemp
	WRs = WRsTemp
	TEs = TEsTemp

	AdjustPlayers(QBs, RBs, WRs, TEs, DSTs)

	QBsTemp = make([]*Player, 0)
	RBsTemp = make([]*Player, 0)
	WRsTemp = make([]*Player, 0)
	TEsTemp = make([]*Player, 0)

	for _, q := range QBs {
		if q.ProjectedPoints >= 10 {
			QBsTemp = append(QBsTemp, q)
			//fmt.Printf("%s - %f\n", q.Name, q.ProjectedPoints)
		}
	}

	//fmt.Printf("----\n")
	for _, r := range RBs {
		if r.ProjectedPoints >= 4 {
			RBsTemp = append(RBsTemp, r)
			//fmt.Printf("%s - %f\n", r.Name, r.ProjectedPoints)
		}
	}
	//fmt.Printf("----\n")
	for _, w := range WRs {
		if w.ProjectedPoints >= 4 {
			WRsTemp = append(WRsTemp, w)
			//fmt.Printf("%s - %f\n", w.Name, w.ProjectedPoints)
		}
	}
	//fmt.Printf("----\n")
	for _, t := range TEs {
		if t.ProjectedPoints >= 3 {
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

	sort.Sort(PlayersPoints(QBs))
	sort.Sort(PlayersPoints(RBs))
	sort.Sort(PlayersPoints(WRs))
	sort.Sort(PlayersPoints(TEs))
	sort.Sort(PlayersPoints(DSTs))
	DSTs = DSTs[0:8]

	cores := runtime.NumCPU()
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

	b, _ := json.Marshal(master)
	fmt.Printf("%s\n", string(b))
}

func create(pid int, QBs, RBs, WRs, TEs, DSTs []*Player) Rosters {
	candidates := make(Rosters, 0, 100000)
	var q *Player
	for true {
		mutex.Lock()
		if globalIndex < len(QBs) {
			q = QBs[globalIndex]
			globalIndex += 1
			mutex.Unlock()
		} else {
			mutex.Unlock()
			break
		}
		roster := &Roster{}
		roster.addPlayer(q, 0, false)
		roster.popPlayer(WR, 0)
		for wr1I, wr1 := range WRs { // wr1
			added := roster.addPlayer(wr1, 0, false)
			if !added {
				continue
			}
			roster.popPlayer(RB, 0)
			for rb1I, rb1 := range RBs { // rb1
				added := roster.addPlayer(rb1, 0, false)
				if !added {
					continue
				}
				roster.popPlayer(RB, 1)
				for rb2I, rb2 := range RBs[rb1I+1:] { //rb2
					if rb2.team == rb1.team {
						continue
					}
					added := roster.addPlayer(rb2, 1, false)
					if !added {
						continue
					}
					roster.popPlayer(WR, 1)
					for wr2I, wr2 := range WRs[wr1I+1:] { // wr2
						if wr2.team == wr1.team {
							continue
						}
						added := roster.addPlayer(wr2, 1, false)
						if !added {
							continue
						}
						roster.popPlayer(WR, 2)
						for wr3I, wr3 := range WRs[wr2I+1:] { // wr3
							if wr3.team == wr1.team || wr3.team == wr2.team {
								continue
							}
							added := roster.addPlayer(wr3, 2, false)
							if !added {
								continue
							}
							roster.popPlayer(TE, 0)
							for _, te := range TEs { // te
								added := roster.addPlayer(te, 0, false)
								if !added {
									continue
								}
								for _, dst := range DSTs { // dst
									roster.popPlayer(DST, 0)
									added := roster.addPlayer(dst, 0, false)
									if !added {
										continue
									}
									addFLEX(roster, RBs[rb2I+1:], WRs[wr3I+1:], &candidates)
									roster.popPlayer(DST, 0)
								} // dst loop
								roster.popPlayer(TE, 0)
							} // te loop
							roster.popPlayer(WR, 2)
						} // wr 3 loop
						roster.popPlayer(WR, 1)
					} // wr 2 loop
					roster.popPlayer(RB, 1)
				} // rb 2 loop
				roster.popPlayer(RB, 0)
			} // rb 1 loop
			roster.popPlayer(WR, 0)
		} //wr 1 loop
	}

	return candidates
}

func addFLEX(roster *Roster, RBs, WRs []*Player, candidates *Rosters) {
	var l int
	if len(RBs) > len(WRs) {
		l = len(WRs)
	} else {
		l = len(RBs)
	}

	var rbI, wrI int
	for i := 0; i < l; i++ {
		var f *Player
		var rb *Player
		for rbI < len(RBs) && (RBs[rbI].team == roster.RB1.team || RBs[rbI].team == roster.RB2.team) {
			if rbI >= len(RBs) {
				break
			}
			rbI += 1
		}

		for wrI < len(WRs) && (WRs[wrI].team == roster.WR1.team || WRs[wrI].team == roster.WR2.team || WRs[wrI].team == roster.WR3.team) {
			if wrI >= len(WRs) {
				break
			}
			wrI += 1
		}

		if rbI < len(RBs) {
			rb = RBs[rbI]
		}

		var wr *Player
		if wrI < len(WRs) {
			wr = WRs[wrI]
		}

		if rb == nil && wr == nil {
			break
		}
		if wr == nil {
			f = rb
			rbI++
		} else if rb == nil {
			f = wr
			wrI++
		} else if rb.ProjectedPoints < wr.ProjectedPoints {
			f = wr
			wrI++
		} else {
			f = rb
			rbI++
		}
		roster.popPlayer(FLEX, 0)
		if roster.addPlayer(f, 0, true) {
			if len(*candidates) < 100 || roster.Points >= (*candidates)[len(*candidates)-1].Points {
				cpy := roster.Copy()
				*candidates = append(*candidates, cpy)
				sort.Sort(candidates)
				if len(*candidates) > 100000 {
					*candidates = (*candidates)[:100]
				}
			}
			roster.popPlayer(FLEX, 0)
			break
		}
	}
}
