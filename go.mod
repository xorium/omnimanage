module omnimanage

go 1.15

require (
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/fatih/structs v1.1.0
	github.com/google/jsonapi v0.0.0-20201022225600-f822737867f6
	github.com/labstack/echo/v4 v4.1.17
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.6.1 // indirect
	gitlab.omnicube.ru/libs/omnilib v0.0.0-00010101000000-000000000000
	gorm.io/datatypes v1.0.0
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)

replace gitlab.omnicube.ru/libs/omnilib => ../omnilib
