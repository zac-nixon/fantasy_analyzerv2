package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

type Parser struct{}

var BlackList []string = []string{"Phi", "Car", "NYG", "Den", "Ind", "Ten", "Buf", "Cin", "Dal", "Sea"}

var TeamMap map[string]string = map[string]string{
	"Jacksonville":  "Jax",
	"Denver":        "Den",
	"Minnesota":     "Min",
	"Kansas City":   "KC",
	"Arizona":       "Ari",
	"Philadelphia":  "Phi",
	"Buffalo":       "Buf",
	"Cleveland":     "Cle",
	"Seattle":       "Sea",
	"Miami":         "Mia",
	"Indianapolis":  "Ind",
	"Cincinnati":    "Cin",
	"Houston":       "Hou",
	"Tennessee":     "Ten",
	"Carolina":      "Car",
	"Washington":    "Was",
	"Chicago":       "Chi",
	"Tampa Bay":     "TB",
	"Oakland":       "Oak",
	"Green Bay":     "GB",
	"Detroit":       "Det",
	"Baltimore":     "Bal",
	"Giants":        "NYG",
	"Atlanta":       "Atl",
	"Dallas":        "Dal",
	"Rams":          "LA",
	"Chargers":      "LAC",
	"Jets":          "NYJ",
	"New Orleans":   "NO",
	"San Francisco": "SF",
	"Pittsburgh":    "Pit",
	"New England":   "NE",
}

func (p *Parser) Parse() (QBs []*Player, RBs []*Player, WRs []*Player, TEs []*Player, DSTs []*Player) {
	DSTs, err := p.parseDST()
	if err != nil {
		panic(err)
	}

	QBs, err = p.parseQB(DSTs)
	if err != nil {
		panic(err)
	}

	RBs, err = p.parseRB(DSTs)
	if err != nil {
		panic(err)
	}

	WRs, err = p.parseWR(DSTs)
	if err != nil {
		panic(err)
	}

	TEs, err = p.parseTE(DSTs)
	if err != nil {
		panic(err)
	}

	return
}

func (p *Parser) parseQB(DSTs []*Player) (QBs []*Player, err error) {
	f, err := os.Open("./csv/qb.csv")
	if err != nil {
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return
	}

	for _, record := range records {
		player := p.parseBasics(record, 14, DSTs)
		if player == nil {
			continue
		}
		if record[4] != "--" {
			var attempts, yards, interceptions, tds, rushYards, rushTds, fumbles float64
			attempts, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[4], ",", "", -1)), 64)
			if err != nil {
				return QBs, err
			}

			yards, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[6], ",", "", -1)), 64)
			if err != nil {
				return QBs, err
			}

			interceptions, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[7], ",", "", -1)), 64)
			if err != nil {
				return QBs, err
			}

			tds, err := strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[8], ",", "", -1)), 64)
			if err != nil {
				return QBs, err
			}

			rushYards, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[9], ",", "", -1)), 64)
			if err != nil {
				return QBs, err
			}

			rushTds, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[10], ",", "", -1)), 64)
			if err != nil {
				return QBs, err
			}

			fumbles, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[11], ",", "", -1)), 64)
			if err != nil {
				return QBs, err
			}
			player.passingStats = &PassingStats{
				passingAttempts:      attempts / player.games,
				passingYards:         yards / player.games,
				passingInterceptions: interceptions / player.games,
				passingTD:            tds / player.games,
			}
			player.rushingStats = &RushingStats{
				rushingYards: rushYards / player.games,
				rushingTD:    rushTds / player.games,
			}
			player.fumbles = fumbles / player.games
			player.Position = QB
			QBs = append(QBs, player)
		}
	}
	return
}

func (p *Parser) parseRB(DSTs []*Player) (RBs []*Player, err error) {
	f, err := os.Open("./csv/rb.csv")
	if err != nil {
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return
	}

	for _, record := range records {
		player := p.parseBasics(record, 13, DSTs)
		if player == nil {
			continue
		}
		if record[4] != "--" {
			var rushes, rushyd, rushtd, receptions, receivingYards, receivingtd, fumbles float64
			rushes, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[4], ",", "", -1)), 64)
			if err != nil {
				return RBs, err
			}

			rushyd, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[5], ",", "", -1)), 64)
			if err != nil {
				return RBs, err
			}

			rushtd, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[6], ",", "", -1)), 64)
			if err != nil {
				return RBs, err
			}

			receptions, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[7], ",", "", -1)), 64)
			if err != nil {
				return RBs, err
			}

			receivingYards, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[8], ",", "", -1)), 64)
			if err != nil {
				return RBs, err
			}

			receivingtd, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[9], ",", "", -1)), 64)
			if err != nil {
				return RBs, err
			}

			fumbles, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[10], ",", "", -1)), 64)
			if err != nil {
				return RBs, err
			}
			player.rushingStats = &RushingStats{
				rushingTD:       rushtd / player.games,
				rushingYards:    rushyd / player.games,
				rushingAttempts: rushes / player.games,
			}

			player.receivingStats = &ReceivingStats{
				Receptions:     receptions / player.games,
				ReceivingYards: receivingYards / player.games,
				ReceivingTD:    receivingtd / player.games,
			}

			player.fumbles = fumbles / player.games
			player.Position = RB
			RBs = append(RBs, player)
		}
	}
	return
}

func (p *Parser) parseWR(DSTs []*Player) (WRs []*Player, err error) {
	f, err := os.Open("./csv/wr.csv")
	if err != nil {
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return
	}

	for _, record := range records {
		player := p.parseBasics(record, 13, DSTs)
		if player == nil {
			continue
		}
		if record[4] != "--" {
			var targets, rushyd, rushtd, receptions, receivingYards, receivingtd, fumbles float64
			receptions, err := strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[4], ",", "", -1)), 64)
			if err != nil {
				return WRs, err
			}

			receivingYards, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[5], ",", "", -1)), 64)
			if err != nil {
				return WRs, err
			}

			targets, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[6], ",", "", -1)), 64)
			if err != nil {
				return WRs, err
			}

			receivingtd, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[7], ",", "", -1)), 64)
			if err != nil {
				return WRs, err
			}

			rushyd, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[8], ",", "", -1)), 64)
			if err != nil {
				return WRs, err
			}

			rushtd, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[9], ",", "", -1)), 64)
			if err != nil {
				return WRs, err
			}

			fumbles, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[10], ",", "", -1)), 64)
			if err != nil {
				return WRs, err
			}

			player.rushingStats = &RushingStats{
				rushingTD:    rushtd / player.games,
				rushingYards: rushyd / player.games,
			}

			player.receivingStats = &ReceivingStats{
				Receptions:     receptions / player.games,
				ReceivingYards: receivingYards / player.games,
				ReceivingTD:    receivingtd / player.games,
				Targets:        targets / player.games,
			}

			player.fumbles = fumbles / player.games
			player.Position = WR
			WRs = append(WRs, player)
		}
	}
	return
}

func (p *Parser) parseTE(DSTs []*Player) (TEs []*Player, err error) {
	f, err := os.Open("./csv/te.csv")
	if err != nil {
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return
	}

	for _, record := range records {
		player := p.parseBasics(record, 13, DSTs)
		if player == nil {
			continue
		}
		if record[4] != "--" {
			var targets, rushyd, rushtd, receptions, receivingYards, receivingtd, fumbles float64
			receptions, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[4], ",", "", -1)), 64)
			if err != nil {
				return TEs, err
			}

			receivingYards, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[5], ",", "", -1)), 64)
			if err != nil {
				return TEs, err
			}

			targets, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[6], ",", "", -1)), 64)
			if err != nil {
				return TEs, err
			}

			receivingtd, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[7], ",", "", -1)), 64)
			if err != nil {
				return TEs, err
			}

			rushyd, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[8], ",", "", -1)), 64)
			if err != nil {
				return TEs, err
			}

			rushtd, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[9], ",", "", -1)), 64)
			if err != nil {
				return TEs, err
			}

			fumbles, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[10], ",", "", -1)), 64)
			if err != nil {
				return TEs, err
			}

			player.rushingStats = &RushingStats{
				rushingTD:    rushtd / player.games,
				rushingYards: rushyd / player.games,
			}

			player.receivingStats = &ReceivingStats{
				Receptions:     receptions / player.games,
				ReceivingYards: receivingYards / player.games,
				ReceivingTD:    receivingtd / player.games,
				Targets:        targets / player.games,
			}

			player.fumbles = fumbles / player.games
			player.Position = TE
			TEs = append(TEs, player)
		}
	}
	return
}

func (p *Parser) parseDST() (DSTs []*Player, err error) {
	f, err := os.Open("./csv/dst.csv")
	if err != nil {
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return
	}

	def_vs_qb_file, err := os.Open("./csv/def_vs_qb.csv")

	if err != nil {
		return
	}
	defer def_vs_qb_file.Close()

	def_vs_rb_file, err := os.Open("./csv/def_vs_rb.csv")

	if err != nil {
		return
	}
	defer def_vs_rb_file.Close()

	def_vs_wr_file, err := os.Open("./csv/def_vs_wr.csv")

	if err != nil {
		return
	}
	defer def_vs_wr_file.Close()

	def_vs_te_file, err := os.Open("./csv/def_vs_te.csv")

	if err != nil {
		return
	}
	defer def_vs_te_file.Close()

	offense_rank_file, err := os.Open("./csv/offense_rank.csv")

	if err != nil {
		return
	}
	defer offense_rank_file.Close()

	def_vs_qb := make(map[string]float64)
	def_vs_wr := make(map[string]float64)
	def_vs_rb := make(map[string]float64)
	def_vs_te := make(map[string]float64)
	offense_rank := make(map[string]float64)

	r2 := csv.NewReader(def_vs_qb_file)

	records2, err := r2.ReadAll()
	if err != nil {
		return
	}

	for _, record := range records2 {
		for k, v := range TeamMap {
			if strings.Contains(record[1], k) {
				record[1] = v
				break
			}
		}
		f, _ := strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[4], ",", "", -1)), 64)
		def_vs_qb[record[1]] = float64(f)
	}

	r2 = csv.NewReader(def_vs_wr_file)

	records2, err = r2.ReadAll()
	if err != nil {
		return
	}

	for _, record := range records2 {
		for k, v := range TeamMap {
			if strings.Contains(record[1], k) {
				record[1] = v
				break
			}
		}
		f, _ := strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[4], ",", "", -1)), 64)
		def_vs_wr[record[1]] = float64(f)
	}

	r2 = csv.NewReader(def_vs_rb_file)

	records2, err = r2.ReadAll()
	if err != nil {
		return
	}

	for _, record := range records2 {
		for k, v := range TeamMap {
			if strings.Contains(record[1], k) {
				record[1] = v
				break
			}
		}
		f, _ := strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[4], ",", "", -1)), 64)
		def_vs_rb[record[1]] = float64(f)
	}

	r2 = csv.NewReader(def_vs_te_file)

	records2, err = r2.ReadAll()
	if err != nil {
		return
	}

	for _, record := range records2 {
		for k, v := range TeamMap {
			if strings.Contains(record[1], k) {
				record[1] = v
				break
			}
		}
		f, _ := strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[4], ",", "", -1)), 64)
		def_vs_te[record[1]] = float64(f)
	}

	r2 = csv.NewReader(offense_rank_file)
	records2, err = r2.ReadAll()
	if err != nil {
		return
	}

	for _, record := range records2 {
		for k, v := range TeamMap {
			if strings.Contains(record[1], k) {
				record[1] = v
				break
			}
		}
		f, _ := strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[4], ",", "", -1)), 64)
		offense_rank[record[1]] = float64(f)
	}

	for _, record := range records {
		player := p.parseBasics(record, 15, []*Player{})
		if player == nil {
			continue
		}

		for k, v := range TeamMap {
			if strings.Contains(player.Name, k) {
				player.Name = v

				break
			}
		}

		var sacks, interceptions, fumbles, safeties, touchdowns, pointsAllowed, passyd, rushyd float64
		player.Position = DST
		sacks, err := strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[4], ",", "", -1)), 64)
		if err != nil {
			return DSTs, err
		}

		interceptions, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[5], ",", "", -1)), 64)
		if err != nil {
			return DSTs, err
		}

		fumbles, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[6], ",", "", -1)), 64)
		if err != nil {
			return DSTs, err
		}

		safeties, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[7], ",", "", -1)), 64)
		if err != nil {
			return DSTs, err
		}

		touchdowns, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[9], ",", "", -1)), 64)
		if err != nil {
			return DSTs, err
		}

		pointsAllowed, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[10], ",", "", -1)), 64)
		if err != nil {
			return DSTs, err
		}

		passyd, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[11], ",", "", -1)), 64)
		if err != nil {
			return DSTs, err
		}

		rushyd, err = strconv.ParseFloat(strings.TrimSpace(strings.Replace(record[12], ",", "", -1)), 64)
		if err != nil {
			return DSTs, err
		}

		player.defenseStats = &DefenseStats{
			sacks:                 sacks / player.games,
			safeties:              safeties / player.games,
			touchdowns:            touchdowns / player.games,
			pointsAllowed:         pointsAllowed / player.games,
			passYardsAllowed:      passyd / player.games,
			rushYardsAllowed:      rushyd / player.games,
			interceptions:         interceptions / player.games,
			fumbleRecovery:        fumbles / player.games,
			pointsToQB:            def_vs_qb[player.Name],
			pointsToRB:            def_vs_rb[player.Name],
			pointsToWR:            def_vs_wr[player.Name],
			pointsToTE:            def_vs_te[player.Name],
			opposingOffensePoints: offense_rank[player.Name],
		}

		DSTs = append(DSTs, player)
	}
	return
}
func (p *Parser) parseBasics(row []string, salaryPos int, DSTs []*Player) *Player {
	name := strings.TrimSpace(row[0])
	var team string
	if strings.Contains(name, ",") {
		chunks := strings.Split(name, ",")
		name = strings.TrimSpace(chunks[0])
		team = strings.TrimSpace(chunks[1])
	}
	opposition := strings.TrimSpace(row[1])

	o := strings.Replace(opposition, "@", "", 1)
	for _, v := range BlackList {
		if v == o {
			return nil
		}
	}

	var DST *Player
	for _, v := range DSTs {
		if v.Name == o {
			DST = v
			break
		}
	}

	var games float64
	var err error
	if strings.TrimSpace(row[3]) != "--" {
		games, err = strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
		if err != nil {
			panic(err)
		}
	}
	salary, err := strconv.Atoi(strings.TrimSpace(row[salaryPos]))
	player := Player{}
	player.Name = name
	player.team = team
	player.Opposition = opposition
	player.Salary = salary
	player.games = games
	player.oppositeObject = DST
	return &player
}
