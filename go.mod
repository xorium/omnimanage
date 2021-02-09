module omnimanage

go 1.15

require (
	github.com/google/jsonapi v0.0.0-20201022225600-f822737867f6
	github.com/labstack/echo/v4 v4.1.17
	github.com/mfcochauxlaberge/jsonapi v0.19.0 // indirect
	github.com/pkg/errors v0.8.1
	gitlab.omnicube.ru/libs/omnilib v0.0.0-00010101000000-000000000000
	golang.org/x/tools v0.0.0-20191029190741-b9c20aec41a5
	gorm.io/datatypes v1.0.0
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)

replace gitlab.omnicube.ru/libs/omnilib => ../omnilib
