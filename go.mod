module github.com/Teeworlds-Server-Moderation/monitor-zcatch

go 1.15

replace github.com/Teeworlds-Server-Moderation/monitor-zcatch/ => ./

// https://pkg.go.dev/github.com/Teeworlds-Server-Moderation/common@v0.7.x
require (
	github.com/Teeworlds-Server-Moderation/common v0.7.5
	github.com/jxsl13/simple-configo v1.18.0
	github.com/jxsl13/twapi v1.2.1
)
