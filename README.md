# learning-spotify-tui

# About

Create a Basic Spotify Controller within a Terminal User Interface (TUIs).
The primary focus was on using [Bubble Tea](https://github.com/charmbracelet/bubbletea),
and [Bubbles](https://github.com/charmbracelet/bubbles) to learn about the components available
and how they interact with each other


# References

- https://github.com/charmbracelet/bubbletea
- https://github.com/charmbracelet/bubbles
- https://github.com/charmbracelet/lipgloss
- https://github.com/zmb3/spotify
- ChatGPT4
 
# Technologies

- Go
- Bubble Tea
- Bubbles
- Lip Gloss
- VSCode + Vim Extension 

# How to get the Spotify Client ID and Spotify Client Secret

- Log into https://developer.spotify.com/
- After Log in, and then click on `Dashboard` 
- Click on `Create app`
- Specify the `App name`, `App description`, `Website`
- For the `Redirect URIs`, specified `http://localhost:8080/callback`, which is same as seen in `main.go`
- Select `Web API` and `Web Playback SDK`
- After creating the app, Select the given app
- Click on Basic information
- Click on `View client secret`
- Copy the `Client ID` and `Client Secret` value
- Paste them in some file `.env.spotify`
- Ensure those are exported as `SPOTIFY_CLIENT_ID` and `SPOTIFY_CLIENT_SECRET` so that it can be 
  used within the program

```
# .env.spotify
export SPOTIFY_CLIENT_ID=""
export SPOTIFY_CLIENT_SECRET=""
```

# How to use

```
# To install the go libraries
go mod tidy

# To load the spotify client id, and spotify client secret
source .env.spotify

# To run the program
go run *.go

# Then Click on the URL, to login and provide the OAuth Token
```

# Things to try out

- [ ] Searching through playlists
- [ ] Searching through tracks
- [ ] Using Refresh Token, instead of logging into the URL constantly