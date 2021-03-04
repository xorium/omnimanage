package web

type User struct {
	ID            string                 `jsonapi:"primary,users"`
	Name          string                 `jsonapi:"attr,userName" validate:"required"`
	Password      string                 `jsonapi:"attr,password"`
	FirstName     string                 `jsonapi:"attr,firstName"`
	LastName      string                 `jsonapi:"attr,lastName"`
	PhoneNumber   string                 `jsonapi:"attr,phoneNumber"`
	Email         string                 `jsonapi:"attr,email" validate:"email"`
	Image         string                 `jsonapi:"attr,image"`
	Settings      map[string]interface{} `jsonapi:"attr,settings"`
	Company       *Company               `jsonapi:"relation,company"`
	Location      *Location              `jsonapi:"relation,location"`
	Roles         []*Role                `jsonapi:"relation,roles"`
	Subscriptions []*Subscription        `jsonapi:"relation,subscriptions"`
}
