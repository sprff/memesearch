module tg-client

go 1.24.0

require github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1

require api-client v0.0.0

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/go-chi/chi/v5 v5.2.1 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/oapi-codegen/runtime v1.1.1 // indirect
)

replace api-client => ../api-client
