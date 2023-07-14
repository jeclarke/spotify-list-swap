# Spotify List Swap

## Why
This tool makes a new playlist based on an existing one, with each track swapped for another from the same album.

## Guide

### Get the repo
````
git clone https://github.com/jeclarke/spotify-list-swap.git
````

### Create a spotify app

Create an app at [https://developer.spotify.com/my-applications/]()

Give it any name you like and set the redirect URL to:

````
http://localhost:8080/callback
````

Get the client ID and client secret and put them in the file:
````
cmd/listswap.env
````

### Build the local app
First, install the go tools used to build the app here [https://go.dev/doc/install]()

Then build the app:
````
cd spotify-list-swap/cmd
go build
````

### Find a playlist you want to duplicate
In spotify right click a playlist, share, and then "copy playlist link".

The playlist ID is the part right after playlist. E.g. in `https://open.spotify.com/playlist/76TxQb4OujgeNKZqzUOOEM?si=cd0070af88714e7c` the ID is `76TxQb4OujgeNKZqzUOOEM`.

### Run the local app
````
./listswapcmd 76TxQb4OujgeNKZqzUOOEM
````
Open the auth link as directed in the console output. This tells spotify you want to let this app create a playlist for you in your account.

Example output:

````
➜  cmd ./listswapcmd 76TxQb4OujgeNKZqzUOOEM
Please log in to Spotify by visiting the following page in your browser: https://accounts.spotify.com/authorize?client_id=69fe015becc24eeea3e8425cc12e3457&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&response_type=code&scope=user-read-private+playlist-modify-private&state=abc123
You are logged in as: a user name
2023/07/14 15:33:00 Making a new list based on "Emo" with 32 tracks
2023/07/14 15:33:01   Swapping Celebration Song from The Greatest Mistake of My Life for In Circles
2023/07/14 15:33:01   Swapping White Lies - Original from White Lies for White Lies - Acoustic
2023/07/14 15:33:01   Swapping Roses for the Dead from Hours for History
...
````

