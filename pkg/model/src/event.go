package src

type Event struct {
	ID           int                    `jsonapi:"primary,events"`
	Type         string                 `jsonapi:"attr,type"`
	Title        string                 `jsonapi:"attr,title"`
	Time         int                    `jsonapi:"attr,time"`
	SessionId    string                 `jsonapi:"attr,sessionId"`
	SessionSlug  string                 `jsonapi:"attr,sessionSlug"`
	SessionState string                 `jsonapi:"attr,sessionState"`
	Level        string                 `jsonapi:"attr,level"`
	Ttl          int                    `jsonapi:"attr,ttl"`
	Info         map[string]interface{} `jsonapi:"attr,info"`
	Company      *Company               `jsonapi:"relation,company"`
	Location     *Location              `jsonapi:"relation,location"`
	Device       *Location              `jsonapi:"relation,device"`
	User         *User                  `jsonapi:"relation,user"`
	Session      *EventsSession         `jsonapi:"relation,session"`
}

type EventsSession struct {
	ID            int       `jsonapi:"primary,eventsSessions"`
	Title         string    `jsonapi:"attr,title"`
	State         string    `jsonapi:"attr,state"`
	Level         string    `jsonapi:"attr,level"`
	LastEventTime int       `jsonapi:"attr,lastEventTime"`
	Slug          string    `jsonapi:"attr,slug"`
	Company       *Company  `jsonapi:"relation,company"`
	Device        *Device   `jsonapi:"relation,device"`
	Location      *Location `jsonapi:"relation,location"`
	LastUser      *User     `jsonapi:"relation,lastUser"`
	Events        []*Event  `jsonapi:"relation,events"`
}
