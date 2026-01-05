package mexadomain

import "time"

type CaseDeteriorationId = int

type CadetDeterioration struct {
	Id        CaseDeteriorationId
	CadetId   CasualtyId
	CreatedAt time.Time
	Value     string
}
