package main

const (
	QBLIMIT   = 1
	RBLIMIT   = 2
	WRLIMIT   = 3
	TELIMIT   = 1
	DSTLIMIT  = 1
	FLEXLIMIT = 1
	CAP       = 50000
)

type Roster struct {
	QBs    []*Player `json:"QB"`
	RBs    []*Player `json:"RBs"`
	WRs    []*Player `json:"WRs"`
	TEs    []*Player `json:"TE"`
	FLEXs  []*Player `json:"FLEX"`
	DSTs   []*Player `json:"DSTs"`
	Spent  int       `json:"Spent"`
	Points float64   `json:"Points"`
}

type Rosters []*Roster

func (rs Rosters) Len() int {
	return len(rs)
}

func (rs Rosters) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs Rosters) Less(i, j int) bool {
	return rs[i].Points > rs[j].Points
}

func (r *Roster) canAfford(p *Player) bool {
	return p.Salary+r.Spent <= CAP
}

func (r *Roster) addPlayer(p *Player, flex bool) bool {
	if !r.canAfford(p) {
		return false
	}

	r.Spent += p.Salary
	r.Points += p.ProjectedPoints

	if flex {
		r.FLEXs = append(r.FLEXs, p)
		return true
	}

	if p.Position == QB {
		r.QBs = append(r.QBs, p)
	}

	if p.Position == RB {
		r.RBs = append(r.RBs, p)
	}

	if p.Position == WR {
		r.WRs = append(r.WRs, p)
	}

	if p.Position == TE {
		r.TEs = append(r.TEs, p)
	}

	if p.Position == DST {
		r.DSTs = append(r.DSTs, p)
	}

	return true
}

func (r *Roster) popPlayer(position string) {
	var p *Player
	if position == QB {
		p, r.QBs = r.QBs[len(r.QBs)-1], r.QBs[:len(r.QBs)-1]
	} else if position == RB {
		p, r.RBs = r.RBs[len(r.RBs)-1], r.RBs[:len(r.RBs)-1]
	} else if position == WR {
		p, r.WRs = r.WRs[len(r.WRs)-1], r.WRs[:len(r.WRs)-1]
	} else if position == TE {
		p, r.TEs = r.TEs[len(r.TEs)-1], r.TEs[:len(r.TEs)-1]
	} else if position == DST {
		p, r.DSTs = r.DSTs[len(r.DSTs)-1], r.DSTs[:len(r.DSTs)-1]
	} else {
		p, r.FLEXs = r.FLEXs[len(r.FLEXs)-1], r.FLEXs[:len(r.FLEXs)-1]
	}
	r.Points -= p.ProjectedPoints
	r.Spent -= p.Salary
}

func (r *Roster) isValid() bool {
	return len(r.QBs) == QBLIMIT && len(r.RBs) == RBLIMIT && len(r.WRs) == WRLIMIT && len(r.TEs) == TELIMIT && len(r.DSTs) == DSTLIMIT && r.Spent <= CAP
}

func (r *Roster) Copy() *Roster {
	rosterCopy := &Roster{}
	rosterCopy.Spent = r.Spent
	rosterCopy.Points = r.Points
	rosterCopy.QBs = make([]*Player, 0)
	rosterCopy.RBs = make([]*Player, 0)
	rosterCopy.WRs = make([]*Player, 0)
	rosterCopy.TEs = make([]*Player, 0)
	rosterCopy.DSTs = make([]*Player, 0)

	for _, q := range r.QBs {
		rosterCopy.QBs = append(rosterCopy.QBs, q)
	}

	for _, rb := range r.RBs {
		rosterCopy.RBs = append(rosterCopy.RBs, rb)
	}

	for _, wr := range r.WRs {
		rosterCopy.WRs = append(rosterCopy.WRs, wr)
	}

	for _, te := range r.TEs {
		rosterCopy.TEs = append(rosterCopy.TEs, te)
	}

	for _, dst := range r.DSTs {
		rosterCopy.DSTs = append(rosterCopy.DSTs, dst)
	}

	for _, f := range r.FLEXs {
		rosterCopy.FLEXs = append(rosterCopy.FLEXs, f)
	}
	return rosterCopy
}
