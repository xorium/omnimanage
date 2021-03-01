package web

type Rule struct {
	ID                string                 `jsonapi:"primary,rules"`
	Title             string                 `jsonapi:"attr,title"`
	Slug              string                 `jsonapi:"attr,slug"`
	Expression        map[string]interface{} `jsonapi:"attr,expression"`
	Duration          int                    `jsonapi:"attr,duration"`
	EventLevel        string                 `jsonapi:"attr,eventLevel"`
	EventSessionState string                 `jsonapi:"attr,eventSessionState"`
	RuleGroup         string                 `jsonapi:"attr,ruleGroup"`
	Company           *Company               `jsonapi:"relation,company"`
	Devices           []*Device              `jsonapi:"relation,devices"`
	Params            []*Parameter           `jsonapi:"relation,params"`
}
