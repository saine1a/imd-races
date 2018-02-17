package racelisting

type RaceSiteType int32

const (
	USSA       RaceSiteType = 0
	LIVETIMING RaceSiteType = 1
)

type RaceDefinition struct {
	RaceId     string
	Qualifier  bool
	TimingSite RaceSiteType
}

/*
//var races = []RaceDefinition{{RaceId: "178531", Qualifier: true, TimingSite: LIVETIMING}, {RaceId: "178336", Qualifier: true, TimingSite: LIVETIMING}, {RaceId: "178251", Qualifier: true, TimingSite: LIVETIMING}}
var Races = []RaceDefinition{
	{RaceId: "U0637", Qualifier: true, TimingSite: USSA}, 
	{RaceId: "U0638", Qualifier: true, TimingSite: USSA}, 
	{RaceId: "U0634", Qualifier: true, TimingSite: USSA},
	{RaceId: "U0051", Qualifier: true, TimingSite: USSA},
	{RaceId: "U0053", Qualifier: true, TimingSite: USSA}}
*/

var Races = []RaceDefinition{
	{RaceId: "U0173", Qualifier: true, TimingSite: USSA},
	{RaceId: "U0175", Qualifier: true, TimingSite: USSA},
	{RaceId: "U0177", Qualifier: true, TimingSite: USSA},
	{RaceId: "U0206", Qualifier: true, TimingSite: USSA},
	{RaceId: "U0208", Qualifier: true, TimingSite: USSA},
	{RaceId: "U0210", Qualifier: true, TimingSite: USSA},
	{RaceId: "U0004", Qualifier: true, TimingSite: USSA},
	{RaceId: "U0006", Qualifier: true, TimingSite: USSA},
	{RaceId: "U0008", Qualifier: true, TimingSite: USSA}}
