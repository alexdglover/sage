module github.com/alexdglover/sage

go 1.22.4

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	gorm.io/driver/sqlite v1.5.6 // indirect
	gorm.io/gorm v1.25.10 // indirect
)

replace github.com/alexdglover/sage/internal => ./internal

replace github.com/alexdglover/sage/internal/models => ./internal/models
