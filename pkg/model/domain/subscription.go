package domain

type Subscription struct {
	ID              string                 `jsonapi:"primary,subscriptions"`
	Title           string                 `jsonapi:"attr,title"`
	ContactChannels map[string]interface{} `jsonapi:"attr,contactChannels"`
	Options         map[string]interface{} `jsonapi:"attr,options"`
	Company         *Company               `jsonapi:"relation,company"`
	User            *User                  `jsonapi:"relation,user"`
	Rules           []*Rule                `jsonapi:"relation,rules"`
}
