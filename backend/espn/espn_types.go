package espn

/*
Ce fichier contient les types qui vont stocker les données retournées par l'API. 
*/

// Note: il ya certains champs qui ne seront pas exploiter peut être, mais y'aurait	
// moyen qu'une version future de l'application le fasse donc je laisse.

type ESPNScoreboardResponse struct {
	Leagues []ESPNLeague         `json:"leagues"`
	Season  ESPNScoreboardSeason `json:"season"`
	Events  []ESPNEvent          `json:"events"`
}

type ESPNScoreboardSeason struct {
	Year int `json:"year"`
	Type int `json:"type"`
}

type ESPNLeague struct {
	ID           string             `json:"id"`
	UID          string             `json:"uid"`
	Name         string             `json:"name"`
	DisplayName  string             `json:"displayName"`
	Abbreviation string             `json:"abbreviation"`
	ShortName    string             `json:"shortName"`
	Slug         string             `json:"slug"`
	Season       ESPNLeagueSeason   `json:"season"`
	Calendar     []ESPNCalendarItem `json:"calendar"`
}

type ESPNLeagueSeason struct {
	Year        int                  `json:"year"`
	StartDate   string               `json:"startDate"`
	EndDate     string               `json:"endDate"`
	DisplayName string               `json:"displayName"`
}

type ESPNCalendarItem struct {
	Label     string  `json:"label"`
	StartDate string  `json:"startDate"`
	EndDate   string  `json:"endDate"`
	Event     ESPNRef `json:"event"`
}

type ESPNRef struct {
	Ref string `json:"$ref"`
}

type ESPNEvent struct {
	ID           string            `json:"id"`
	UID          string            `json:"uid"`
	Date         string            `json:"date"`
	Name         string            `json:"name"`
	ShortName    string            `json:"shortName"`
	Season       ESPNEventSeason   `json:"season"`
	Competitions []ESPNCompetition `json:"competitions"`
	Venues       []ESPNVenue       `json:"venues"`
	Status       ESPNStatus        `json:"status"`
}

type ESPNEventSeason struct {
	Year int    `json:"year"`
	Type int    `json:"type"`
	Slug string `json:"slug"`
}

type ESPNCompetition struct {
	ID          string              `json:"id"`
	UID         string              `json:"uid"`
	Date        string              `json:"date"`
	EndDate     string              `json:"endDate"`
	StartDate   string              `json:"startDate"`
	Type        ESPNCompetitionType `json:"type"`
	Venue       ESPNVenue           `json:"venue"`
	Competitors []ESPNCompetitor    `json:"competitors"`
	Status      ESPNStatus          `json:"status"`
}

type ESPNCompetitionType struct {
	ID           string `json:"id"`
	Abbreviation string `json:"abbreviation"`
}

type ESPNVenue struct {
	ID       string      `json:"id"`
	FullName string      `json:"fullName"`
	Address  ESPNAddress `json:"address"`
}

type ESPNAddress struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}

type ESPNStatus struct {
	Clock        float64        `json:"clock"`
	DisplayClock string         `json:"displayClock"`
	Period       int            `json:"period"`
	Type         ESPNStatusType `json:"type"`
}

type ESPNStatusType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
	Detail      string `json:"detail"`
	ShortDetail string `json:"shortDetail"`
}

type ESPNCompetitor struct {
	ID      string       `json:"id"`
	UID     string       `json:"uid"`
	Type    string       `json:"type"`
	Order   int          `json:"order"`
	Winner  bool         `json:"winner"`
	Athlete ESPNAthlete  `json:"athlete"`
	Records []ESPNRecord `json:"records"`
}

type ESPNAthlete struct {
	FullName    string   `json:"fullName"`
	DisplayName string   `json:"displayName"`
	ShortName   string   `json:"shortName"`
	Flag        ESPNFlag `json:"flag"`
}

type ESPNFlag struct {
	Href string `json:"href"`
	Alt  string `json:"alt"`
}

type ESPNRecord struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Type         string `json:"type"`
	Summary      string `json:"summary"`
}
