# ⭐⭐🎄 Advent Of Code Discord Bot 🎄⭐⭐

A discord bot written in Go to provide leaderboards and nightly notifications to your server!

## Installation

### Self Hosting

1. Set `$AOC_BOT_TOKEN` to discord bot secret token

2. `go run main.go`

### Docker



## Usage 

### General
`/leaderboard` - Show the current state of the leaderboard

### Admin
`/configure <channelId> <roleId> <leaderboardId> <sessionToken>` - Configure the bot to work on your server. Provide the desired channel id, role id to ping, the id of your AoC leaderboard and a valid session token of one member that has access to the private leaderboard

`/start-notifications` - Subscribe to nightly notifications

`/stop-notifications` - Unsubscribe from nightly notifications

`/check-notifications` - Check your servers current notification status

