package mb

type AccessTechnology int

const (
	UnknownAt    AccessTechnology = 0 << 0
	PotsAt       AccessTechnology = 1 << 0
	GSMAt        AccessTechnology = 2 << 0
	GSMCompactAt AccessTechnology = 1 << 2
	GPRSAt       AccessTechnology = 1 << 3
	EDGEAt       AccessTechnology = 1 << 4
	UMTSAt       AccessTechnology = 1 << 5
	HSDPAAt      AccessTechnology = 1 << 6
	HSUPAAt      AccessTechnology = 1 << 7
	HSPAAt       AccessTechnology = 1 << 8
	PLUSAt       AccessTechnology = 1 << 9
	OneXRTAt     AccessTechnology = 1 << 10
	EVDO0At      AccessTechnology = 1 << 11
	EVDOAAt      AccessTechnology = 1 << 12
	EVDOBAt      AccessTechnology = 1 << 13
	LTEAt        AccessTechnology = 1 << 14
	FiveGNRAt    AccessTechnology = 1 << 15
)
