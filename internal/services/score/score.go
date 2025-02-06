package score

type Score struct {
	ScoreID      string  `json:"scoreID" bson:"scoreID"`
	Season       string  `json:"season" bson:"season"`
	ScoredAt     string  `json:"scoredAt" bson:"scoredAt"`
	ScoreDetails Details `json:"scoreDetails" bson:"scoreDetails"`
}

type Details struct {
	Nickname     string `json:"playerNickname" bson:"playerNickname"`
	Rating       int    `json:"rating" bson:"rating"`
	Wins         int    `json:"wins" bson:"wins"`
	Losses       int    `json:"losses" bson:"losses"`
	WinLoseRatio string `json:"winLoseRatio" bson:"winLoseRatio"`
	Region       string `json:"region" bson:"region"`
	WorldRank    int    `json:"rank" bson:"rank"`
}
