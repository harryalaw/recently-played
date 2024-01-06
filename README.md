### Overview

I listen to a lot of music and have been wanting to track it a bit better, mostly to understand what new albums I've enjoyed over the course of a year!
This is a project for me to be able to track the songs that I am listening to on Spotify using their API.

I've also expanded it in places to track other information - I have a Song of the Day playlist which I've recently expanded this to cover.

### Recently Played

I have deployed the code in `cmd/functions/recently-played.go` to a Lambda function which is triggered by a cron job every half hour.
It refreshes an access token, then calls Spotify's `recently-played` endpoint to get the 50 most recently played songs. It filters this list for the values that I have not previously recorded and then inserts this into a database.

### Song of the Day

In order to backup my song of the day playlist (going back to 2019) in case I ever lose access to it I have adjusted the `cmd/server/main.go` to fetch the contents of the playlist. At the moment I am manually updating the offset to only fetch the latest values but could automate this and similarly deploy it to Lambda.

### Visualisations
Currently working on a UI to visualise this information.
