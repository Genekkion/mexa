package mexadomain

import "fmt"

type CCLogId = int

type CCLog struct {
	Id         CCLogId
	CasualtyId CasualtyId
	Type       CCLogType

	Value CCLogValue
}

type CCLogValue struct {
	Outcome *CCLogEndOutcome `json:"outcome,omitempty"`
}

type CCLogType string

const (
	CCLogTypeTreatStart CCLogType = "treat_start"
	CCLogTypeTreatEnd   CCLogType = "treat_end"
)

type CCLogEndOutcome string

const (
	CCLogEndOutcomeUnknown CCLogEndOutcome = ""
	CCLogEndOutcomeP1      CCLogEndOutcome = "p1"
	CCLogEndOutcomeP2      CCLogEndOutcome = "p2"
	CCLogEndOutcomeP3      CCLogEndOutcome = "p3"
	CCLogEndOutcomeP4      CCLogEndOutcome = "p4"
)

func (s CCLogEndOutcome) String() string {
	switch s {
	case CCLogEndOutcomeP1:
		return "P1"
	case CCLogEndOutcomeP2:
		return "P2"
	case CCLogEndOutcomeP3:
		return "P3"
	case CCLogEndOutcomeP4:
		return "P4"
	default:
		return "Unknown"
	}
}

func ParsePValue(v int) (res CCLogEndOutcome, err error) {
	switch v {
	case 1:
		return CCLogEndOutcomeP1, nil
	case 2:
		return CCLogEndOutcomeP2, nil
	case 3:
		return CCLogEndOutcomeP3, nil
	case 4:
		return CCLogEndOutcomeP4, nil
	default:
		return CCLogEndOutcomeUnknown, fmt.Errorf("unknown p value: %d", v)
	}
}

func NewCCLogTreatStart(casualtyId CasualtyId) CCLog {
	return CCLog{
		CasualtyId: casualtyId,
		Type:       CCLogTypeTreatStart,
		Value:      CCLogValue{},
	}
}

func NewCCLogTreatEnd(casualtyId CasualtyId, outcome CCLogEndOutcome) CCLog {
	return CCLog{
		CasualtyId: casualtyId,
		Type:       CCLogTypeTreatEnd,
		Value: CCLogValue{
			Outcome: &outcome,
		},
	}
}
