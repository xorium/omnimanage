package model

type Subscription struct {
	ID    int    `jsonapi:"primary,subscriptions"`
	Title string `jsonapi:"attr,title"`
	//ContactChannels map[string]interface{} `jsonapi:"attr,contactChannels"`
	//Options         map[string]interface{} `jsonapi:"attr,options"`
	//Company         *Company               `jsonapi:"relation,company"`
	UserID string
	User   *User `jsonapi:"relation,user"`
	//Rules           []*Rule                `jsonapi:"relation,rules"`
}
