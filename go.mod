module doubleblind/main

go 1.21.6

replace doubleblind/server => ./server

require doubleblind/server v0.0.0-00010101000000-000000000000

require (
	doubleblind/models v0.0.0-00010101000000-000000000000 // indirect
	github.com/ncruces/go-sqlite3 v0.12.0 // indirect
	github.com/ncruces/julianday v1.0.0 // indirect
	github.com/tetratelabs/wazero v1.6.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)

replace doubleblind/models => ./models
