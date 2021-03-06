module github.com/Teeworlds-Server-Moderation/monitor-zcatch

go 1.15

replace github.com/Teeworlds-Server-Moderation/monitor-zcatch/ => ./

// https://pkg.go.dev/github.com/Teeworlds-Server-Moderation/common@v0.7.x
require (
	github.com/Teeworlds-Server-Moderation/common v0.7.2
	github.com/jxsl13/simple-configo v1.2.6
	github.com/jxsl13/twapi v1.1.2
)
