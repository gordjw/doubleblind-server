module github.com/gordjw/doubleblind

go 1.21.6

replace doubleblind/server => ./server

require doubleblind/server v0.0.0-00010101000000-000000000000

require (
	cloud.google.com/go/compute v1.20.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	doubleblind/models v0.0.0-00010101000000-000000000000 // indirect
	github.com/go-chi/chi v1.5.5 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gorilla/schema v1.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/ncruces/go-sqlite3 v0.12.0 // indirect
	github.com/ncruces/julianday v1.0.0 // indirect
	github.com/tetratelabs/wazero v1.6.0 // indirect
	golang.org/x/oauth2 v0.17.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace doubleblind/models => ./models
