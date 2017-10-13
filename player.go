package main

const (
	QB   = "QB"
	RB   = "RB"
	WR   = "WR"
	TE   = "TE"
	DST  = "DST"
	FLEX = "FLEX"
)

type PassingStats struct {
	passingTD            float64
	passingYards         float64
	passingInterceptions float64
	passingAttempts      float64
}

type RushingStats struct {
	rushingTD       float64
	rushingYards    float64
	rushingAttempts float64
}

type ReceivingStats struct {
	ReceivingTD    float64
	ReceivingYards float64
	Receptions     float64
	Targets        float64
}

type DefenseStats struct {
	sacks            float64
	interceptions    float64
	fumbleRecovery   float64
	touchdowns       float64
	safeties         float64
	pointsAllowed    float64
	passYardsAllowed float64
	rushYardsAllowed float64

	pointsToQB float64
	pointsToRB float64
	pointsToWR float64
	pointsToTE float64

	opposingOffensePoints float64
}

type Player struct {
	Position       string  `json:"position"`
	Name           string  `json:"name"`
	Opposition     string  `json:"opposition"`
	Salary         int     `json:"salary"`
	team           string  `json:"-"`
	oppositeObject *Player `json:"-"`
	games          float64 `json:"-"`
	fumbles        float64 `json:"-"`

	passingStats   *PassingStats   `json:"-"`
	rushingStats   *RushingStats   `json:"-"`
	receivingStats *ReceivingStats `json:"-"`
	defenseStats   *DefenseStats   `json:"-"`

	ProjectedPoints float64 `json:"projectedPoints"`
}

type PlayersPoints []*Player

func (ps PlayersPoints) Len() int {
	return len(ps)
}

func (ps PlayersPoints) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps PlayersPoints) Less(i, j int) bool {
	return ps[i].ProjectedPoints > ps[j].ProjectedPoints
}

type DSTRushAllowed []*Player

func (ps DSTRushAllowed) Len() int {
	return len(ps)
}

func (ps DSTRushAllowed) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps DSTRushAllowed) Less(i, j int) bool {
	return ps[i].defenseStats.rushYardsAllowed < ps[j].defenseStats.rushYardsAllowed
}

type DSTPassAllowed []*Player

func (ps DSTPassAllowed) Len() int {
	return len(ps)
}

func (ps DSTPassAllowed) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps DSTPassAllowed) Less(i, j int) bool {
	return ps[i].defenseStats.passYardsAllowed < ps[j].defenseStats.passYardsAllowed
}

type DSTPointsAllowed []*Player

func (ps DSTPointsAllowed) Len() int {
	return len(ps)
}

func (ps DSTPointsAllowed) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps DSTPointsAllowed) Less(i, j int) bool {
	return ps[i].defenseStats.pointsAllowed < ps[j].defenseStats.pointsAllowed
}

type QBPassAttempts []*Player

func (ps QBPassAttempts) Len() int {
	return len(ps)
}

func (ps QBPassAttempts) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps QBPassAttempts) Less(i, j int) bool {
	return ps[i].passingStats.passingAttempts < ps[j].passingStats.passingAttempts
}

type RBRushAttempts []*Player

func (ps RBRushAttempts) Len() int {
	return len(ps)
}

func (ps RBRushAttempts) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps RBRushAttempts) Less(i, j int) bool {
	return ps[i].rushingStats.rushingAttempts < ps[j].rushingStats.rushingAttempts
}

type PlayerReceptions []*Player

func (ps PlayerReceptions) Len() int {
	return len(ps)
}

func (ps PlayerReceptions) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps PlayerReceptions) Less(i, j int) bool {
	return ps[i].receivingStats.Receptions < ps[j].receivingStats.Receptions
}

func (p *Player) ScoreQB() {
	p.ProjectedPoints = p.scorePassing() + p.scoreRushing() - float64(p.fumbles)
}

func (p *Player) ScoreRB() {
	p.ProjectedPoints = p.scoreRushing() + p.scoreReceiving() - float64(p.fumbles)
}

func (p *Player) ScoreWR() {
	p.ProjectedPoints = p.scoreReceiving() - float64(p.fumbles)
}

func (p *Player) ScoreTE() {
	p.ProjectedPoints = p.scoreReceiving() - float64(p.fumbles)
}

func (p *Player) ScoreDST() {
	p.ProjectedPoints += float64(p.defenseStats.sacks)
	p.ProjectedPoints += float64(p.defenseStats.interceptions) * float64(2)
	p.ProjectedPoints += float64(p.defenseStats.fumbleRecovery) * float64(2)
	p.ProjectedPoints += float64(p.defenseStats.touchdowns) * float64(6)
	p.ProjectedPoints += float64(p.defenseStats.safeties) * float64(2)

	pointsAllowedScore := 0

	if p.defenseStats.pointsAllowed == 0 {
		pointsAllowedScore = 10
	} else if p.defenseStats.pointsAllowed <= 6 {
		pointsAllowedScore = 7
	} else if p.defenseStats.pointsAllowed <= 13 {
		pointsAllowedScore = 4
	} else if p.defenseStats.pointsAllowed <= 20 {
		pointsAllowedScore = 0
	} else if p.defenseStats.pointsAllowed <= 27 {
		pointsAllowedScore = -1
	} else {
		pointsAllowedScore = -4
	}

	p.ProjectedPoints += float64(pointsAllowedScore)
}

func (p *Player) scorePassing() float64 {
	var score float64
	score += float64(p.passingStats.passingTD) * float64(4)
	score += float64(p.passingStats.passingYards) * float64(.04)
	score += float64(p.passingStats.passingInterceptions) * float64(-1)
	if p.passingStats.passingYards > 300 {
		score += float64(3)
	}
	return score
}

func (p *Player) scoreRushing() float64 {
	var score float64
	score += float64(p.rushingStats.rushingTD) * float64(6)
	score += float64(p.rushingStats.rushingYards) * float64(.1)
	if p.rushingStats.rushingYards > 100 {
		score += float64(3)
	}
	return score
}

func (p *Player) scoreReceiving() float64 {
	var score float64
	score += float64(p.receivingStats.ReceivingTD) * float64(6)
	score += float64(p.receivingStats.ReceivingYards) * float64(.1)
	score += float64(p.receivingStats.Receptions)
	if p.receivingStats.ReceivingYards > 100 {
		score += float64(3)
	}
	return score
}
