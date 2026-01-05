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
	CCLogEndOutcomeSuccess CCLogEndOutcome = "success"
	CCLogEndOutcomeFailure CCLogEndOutcome = "failure"
)

func ParseCCLogEndOutcome(s string) (res CCLogEndOutcome, err error) {
	switch s {
	case "success":
		return CCLogEndOutcomeSuccess, nil
	case "failure":
		return CCLogEndOutcomeFailure, nil
	default:
		return CCLogEndOutcomeUnknown, fmt.Errorf("unknown outcome: %s", s)
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
