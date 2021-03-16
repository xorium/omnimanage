package domain

//go:generate go run $PWD/pkg/utils/codegen/ --mode model-controller --model User --file_in $GOFILE

type User struct {
	ID            string                 `jsonapi:"primary,users" default:"1"`
	Name          string                 `jsonapi:"attr,userName" validate:"required" default:"JohnLogin"`
	Password      string                 `jsonapi:"attr,password" default:"123456"`
	FirstName     string                 `jsonapi:"attr,firstName" default:"John"`
	LastName      string                 `jsonapi:"attr,lastName" default:"Doe"`
	PhoneNumber   string                 `jsonapi:"attr,phoneNumber" default:"+79991234567"`
	Email         string                 `jsonapi:"attr,email" validate:"email" default:"krasava@omnicube.ru"`
	Image         string                 `jsonapi:"attr,image" default:"http://www.example.com/"`
	Settings      map[string]interface{} `jsonapi:"attr,settings" default:"{\"Foo\": 123}"`
	Company       *Company               `jsonapi:"relation,company" default:"{\"ID\": \"123\"}"`
	Location      *Location              `jsonapi:"relation,location" default:"{\"ID\": \"123\"}"`
	Roles         []*Role                `jsonapi:"relation,roles" default:"[{\"ID\": \"123\"}]"`
	Subscriptions []*Subscription        `jsonapi:"relation,subscriptions" default:"[{\"ID\": \"123\"}]"`
}
