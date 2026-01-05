package mexadomain

import (
	"encoding/json"
	"fmt"
	"mexa/internal/utils"
	"strings"
	"time"

	cases2 "golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	caser = cases2.Title(language.English)
)

type CaseId = int

type Case struct {
	Id_        int
	Id         CaseId
	CreatedAt  time.Time
	ExerciseId ExerciseId

	CaseValue
}

func (c Case) TgMd2() string {
	tmpl := `
*Case ID*
{id}

*Case Summary*
{summary}

*Category*
{category}

*Scenario*
{scenario}

<divider>

__*Basic Management*__

*Care Under Fire*
{care_under_fire}

*30\-2\-CAN DO*
{30-2_can_do}

*Tactical Field Care*
{tactical_field_care}

<divider>

__*Advanced Management*__

*Temperature*
{temperature}

*BP*
{bp}

*HR*
{hr}

*RR*
{rr}

*SpO2*
{spo2}

*AVPU*
{avpu}

*Pain Score*
{pain_score}

<divider>

__*Primary Survey*__

*Airway*
{airway}

*Breathing*
{breathing}

*Circulation*
{circulation}

*Disability*
{disability}

*Exposure*
{exposure}

<divider>

*Assessment and Plan*
{assessment_and_plan}

<divider>

*Prolonged Field Care*
{prolonged_field_care}

`

	m := map[string]string{
		"<divider>": strings.Repeat("-", 30),

		"{id}": fmt.Sprintf("%d", c.Id),

		"{summary}":  c.Summary,
		"{category}": caser.String(string(c.Category)),
		"{scenario}": c.Scenario,

		"{care_under_fire}":     strings.Join(c.BasicManagement.CareUnderFire, "\n"),
		"{30-2_can_do}":         strings.Join(c.BasicManagement.CanDo_30_2, "\n"),
		"{tactical_field_care}": strings.Join(c.BasicManagement.TacticalFieldCare, "\n"),

		"{temperature}": fmt.Sprintf("%.1f", c.AdvancedManagement.Temperature),
		"{bp}":          c.AdvancedManagement.BP,
		"{hr}":          fmt.Sprintf("%d", c.AdvancedManagement.HR),
		"{rr}":          fmt.Sprintf("%d", c.AdvancedManagement.RR),
		"{spo2}":        fmt.Sprintf("%d", c.AdvancedManagement.SpO2),
		"{avpu}":        c.AdvancedManagement.Avpu,
		"{pain_score}":  "N.A.",

		"{airway}":      c.PrimarySurvey.Airway,
		"{breathing}":   c.PrimarySurvey.Breathing,
		"{circulation}": c.PrimarySurvey.Circulation,
		"{disability}":  c.PrimarySurvey.Disability,
		"{exposure}":    c.PrimarySurvey.Exposure,

		"{assessment_and_plan}":  strings.Join(c.AssessmentAndPlan, "\n"),
		"{prolonged_field_care}": strings.Join(c.ProlongedFieldCare, "\n"),
	}
	if c.AdvancedManagement.PainScore > 0 {
		m["{pain_score}"] = fmt.Sprintf("%d/10", c.AdvancedManagement.PainScore)
	}
	rep := map[string]string{
		`<b\>`: "*",
	}
	for k, v := range m {
		vv := utils.EscapeMd2(v)
		for x, y := range rep {
			vv = strings.ReplaceAll(vv, x, y)
		}
		tmpl = strings.ReplaceAll(tmpl, k, vv)
	}

	return tmpl
}

type CaseValue struct {
	Summary  string   `json:"summary"`
	Category Category `json:"category"`
	Scenario string   `json:"scenario"`

	BasicManagement    BasicManagement    `json:"basic_management"`
	AdvancedManagement AdvancedManagement `json:"advanced_management"`
	PrimarySurvey      PrimarySurvey      `json:"primary_survey"`

	AssessmentAndPlan  []string `json:"assessment_and_plan"`
	ProlongedFieldCare []string `json:"prolonged_field_care"`
}

func (v CaseValue) Json() (b []byte, err error) {
	b, err = json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type Category string

const (
	CategorySevere       = "severe"
	CategoryIntermediate = "intermediate"
	CategoryLight        = "light"
)

type BasicManagement struct {
	CareUnderFire     []string `json:"care_under_fire"`
	CanDo_30_2        []string `json:"30-2_can_do"`
	TacticalFieldCare []string `json:"tactical_field_care"`
}

type AdvancedManagement struct {
	Temperature float64 `json:"temperature"`
	BP          string  `json:"bp"`
	HR          int     `json:"hr"`
	RR          int     `json:"rr"`
	SpO2        int     `json:"spo2"`
	Avpu        string  `json:"avpu"`
	PainScore   int     `json:"pain_score"`
}

type PrimarySurvey struct {
	Airway      string `json:"airway"`
	Breathing   string `json:"breathing"`
	Circulation string `json:"circulation"`
	Disability  string `json:"disability"`
	Exposure    string `json:"exposure"`
}
