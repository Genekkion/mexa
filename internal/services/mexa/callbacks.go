package mexaservice

import (
	"context"
	chatdomain "mexa/internal/domains/chat"
)

const (
	listCasesPrefix         = "cases_list"
	casualtyCheckPrefix     = "check_casualty"
	casualtyCheckListPrefix = "list_check_casualty"
	treatStartPrefix        = "treat_start"
	treatEndPrefix          = "treat_end"
	attachCasePrefix        = "attach"
	detachCasePrefix        = "detach"
	deteriorationPrefix     = "deteriorate"
	ignorePrefix            = "ignore"
)

type Callback struct {
	Prefix  string
	Handler chatdomain.Handler
}

func (s *Service) initCallbacks() {
	s.callbacks = map[string]chatdomain.Handler{
		ignorePrefix: func(_ context.Context, _ chatdomain.Update) (err error) {
			return nil
		},
		listCasesPrefix:     s.callbackListCases,
		attachCasePrefix:    s.callbackAttachCase,
		casualtyCheckPrefix: s.callbackCasualtyCheck,
		treatStartPrefix:    s.callbackTreatStart,
		treatEndPrefix:      s.callbackTreatEnd,
		deteriorationPrefix: s.callbackDeterioration,
	}
}
