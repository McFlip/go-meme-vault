module github.com/McFlip/go-meme-vault

go 1.21.2

require (
	github.com/a-h/templ v0.2.513
	github.com/glebarez/sqlite v1.9.0
	github.com/go-chi/chi/v5 v5.0.10
)

require (
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/image v0.13.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect
)

require (
	github.com/McFlip/go-meme-vault/internal/models v0.0.0-00010101000000-000000000000
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	gorm.io/gorm v1.25.5
)

replace github.com/McFlip/go-meme-vault/internal/models => ./internal/models
