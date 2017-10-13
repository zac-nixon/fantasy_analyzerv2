package main

const (
	CAP = 50000
)

type Roster struct {
	QB     *Player `json:"QB"`
	RB1    *Player `json:"RB1"`
	RB2    *Player `json:"RB2"`
	WR1    *Player `json:"WR1"`
	WR2    *Player `json:"WR2"`
	WR3    *Player `json:"WR3"`
	TE     *Player `json:"TE"`
	FLEX   *Player `json:"FLEX"`
	DST    *Player `json:"DST"`
	Spent  int     `json:"Spent"`
	Points float64 `json:"Points"`
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

func (r *Roster) addPlayer(p *Player, index int, flex bool) bool {
	if !r.canAfford(p) {
		return false
	}

	r.Spent += p.Salary
	r.Points += p.ProjectedPoints

	if flex {
		r.FLEX = p
		return true
	}

	if p.Position == QB {
		r.QB = p
	}

	if p.Position == RB {
		if index == 0 {
			r.RB1 = p
		} else {
			r.RB2 = p
		}
	}

	if p.Position == WR {
		if index == 0 {
			r.WR1 = p
		} else if index == 1 {
			r.WR2 = p
		} else {
			r.WR3 = p
		}
	}

	if p.Position == TE {
		r.TE = p
	}

	if p.Position == DST {
		r.DST = p
	}

	return true
}

func (r *Roster) popPlayer(position string, index int) {
	var p *Player

	if position == QB {
		p = r.QB
		r.QB = nil
	}

	if position == RB {
		if index == 0 {
			p = r.RB1
			r.RB1 = nil
		} else {
			p = r.RB2
			r.RB2 = nil
		}
	}

	if position == WR {
		if index == 0 {
			p = r.WR1
			r.WR1 = nil
		} else if index == 1 {
			p = r.WR2
			r.WR2 = nil
		} else {
			p = r.WR3
			r.WR3 = nil
		}
	}

	if position == TE {
		p = r.TE
		r.TE = nil
	}

	if position == DST {
		p = r.DST
		r.DST = nil
	}

	if position == FLEX {
		p = r.FLEX
		r.FLEX = nil
	}
	if p == nil {
		return
	}
	r.Points -= p.ProjectedPoints
	r.Spent -= p.Salary
}

func (r *Roster) Copy() *Roster {
	rosterCopy := &Roster{}
	rosterCopy.Spent = r.Spent
	rosterCopy.Points = r.Points

	rosterCopy.QB = r.QB
	rosterCopy.RB1 = r.RB1
	rosterCopy.RB2 = r.RB2
	rosterCopy.WR1 = r.WR1
	rosterCopy.WR2 = r.WR2
	rosterCopy.WR3 = r.WR3
	rosterCopy.TE = r.TE
	rosterCopy.DST = r.DST
	rosterCopy.FLEX = r.FLEX
	return rosterCopy
}

func (r *Roster) Equal(o *Roster) bool {
	if o == nil {
		return false
	}
	return o.QB.Name == r.QB.Name && o.RB1.Name == r.RB1.Name && o.RB2.Name == r.RB2.Name && o.WR1.Name == r.WR1.Name && o.WR2.Name == r.WR2.Name && o.WR3.Name == r.WR3.Name && o.TE.Name == r.TE.Name && o.DST.Name == r.DST.Name && o.FLEX.Name == r.FLEX.Name
}
