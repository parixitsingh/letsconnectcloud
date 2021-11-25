module letsconnectcloud

replace route => ./route

require (
	github.com/gorilla/mux v1.8.0 // indirect
	route v0.0.0-00010101000000-000000000000
)

go 1.15
