package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify"
)

var (
	leftPaneStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true)
	rightPaneStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true)
	bottomPaneStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Align(lipgloss.Center)
)

type model struct {
	client            *spotify.Client
	playlists         []spotify.SimplePlaylist
	tracks            []spotify.FullTrack
	playlistsTable    table.Model
	tracksTable       table.Model
	focused           string
	selectedPlaylist  int
	selectedTrackName string
	status            string
	volume            int
	width             int
	height            int
}

func convertPlaylistsToRows(playlists []spotify.SimplePlaylist) []table.Row {
	rows := []table.Row{}
	for _, playlist := range playlists {
		rows = append(rows, table.Row{
			playlist.Name,
		})
	}

	return rows
}
func convertTracksToRows(tracks []spotify.FullTrack) []table.Row {
	rows := []table.Row{}
	for _, track := range tracks {
		rows = append(rows, table.Row{
			track.Name,
		})
	}
	return rows
}

func focusedTableStyle() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.Border(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color("#1ed760"))
	s.Selected = s.Selected.Background(lipgloss.Color("#1ed760")).Foreground(lipgloss.Color("0"))
	return s
}

func unfocusedTableStyle() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	s.Selected = s.Selected.Background(lipgloss.Color("240")).Foreground(lipgloss.Color("0"))
	return s
}

func renderControls(m *model) string {
	return fmt.Sprintf("Song: %s || Status: %s || Volume %d%% \n[p] Play/Pause [up/down] Navigate [tab] Switch Tracks/Playlists [enter] Select [+] Volume Up [-] Volume Down [q] Quit", m.selectedTrackName, m.status, m.volume)
}

func (m *model) togglePlayPause() tea.Cmd {
	if m.status == "playing" {
		err := m.client.Pause()
		if err != nil {
			log.Println("Error pausing track:", err)
		}
		m.status = "paused"
	} else {
		err := m.client.Play()
		if err != nil {
			log.Println("Error playing track:", err)
		}
		m.status = "playing"
	}

	return nil
}

func (m *model) adjustVolume(change int) tea.Cmd {
	m.volume += change
	if m.volume > 100 {
		m.volume = 100
	}
	if m.volume < 0 {
		m.volume = 0
	}

	err := m.client.Volume(m.volume)
	if err != nil {
		log.Println("Error adjusting volume:", err)
	}

	return nil
}

func (m *model) playSelectedTrack(trackURI spotify.URI) {
	err := m.client.PlayOpt(&spotify.PlayOptions{
		URIs: []spotify.URI{
			spotify.URI(trackURI),
		},
	})

	if err != nil {
		log.Println("Error playing track:", err)
	}
	m.status = "playing"

}
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		leftPaneWidth := int(0.35 * float64(m.width))
		rightPaneWidth := m.width - leftPaneWidth
		m.playlistsTable.SetWidth(leftPaneWidth)
		m.tracksTable.SetWidth(rightPaneWidth)
	case tea.KeyMsg:
		switch msg.String() {
		case "tab": // Switch focus between tables
			if m.focused == "playlists" {
				m.focused = "tracks"
				m.playlistsTable.SetStyles(unfocusedTableStyle())
				m.tracksTable.SetStyles(focusedTableStyle())
			} else {
				m.focused = "playlists"
				m.playlistsTable.SetStyles(focusedTableStyle())
				m.tracksTable.SetStyles(unfocusedTableStyle())
			}
		case "p":
			return m, m.togglePlayPause()
		case "enter":
			if m.focused == "playlists" {
				playlistID := m.playlists[m.playlistsTable.Cursor()].ID
				tracks, err := loadTracksForPlaylist(m.client, playlistID)
				if err != nil {
					return m, nil
				}
				m.tracks = tracks
				m.tracksTable.SetRows(convertTracksToRows(tracks))
				m.tracksTable.GotoTop()
			} else if m.focused == "tracks" {
				trackURI := m.tracks[m.tracksTable.Cursor()].URI
				m.selectedTrackName = m.tracks[m.tracksTable.Cursor()].Name
				m.playSelectedTrack(trackURI)

			}
		case "+":
			return m, m.adjustVolume(10)
		case "-":
			return m, m.adjustVolume(-10)
		case "q":
			return m, tea.Quit
		}

		// Update the focused table with navigation commands
		if m.focused == "playlists" {
			m.playlistsTable.Focus()
			m.tracksTable.Blur()
			m.playlistsTable, cmd = m.playlistsTable.Update(msg)
		} else if m.focused == "tracks" {
			m.playlistsTable.Blur()
			m.tracksTable.Focus()
			m.tracksTable, cmd = m.tracksTable.Update(msg)
		}
		return m, cmd
	}

	return m, nil
}

func (m *model) View() string {
	leftPaneWidth := int(0.35 * float64(m.width))
	rightPaneWidth := m.width - leftPaneWidth

	bottomPaneHeight := int(0.2 * float64(m.height))
	topPaneHeight := m.height - bottomPaneHeight

	m.playlistsTable.SetWidth(leftPaneWidth)
	m.tracksTable.SetWidth(rightPaneWidth)

	playlistsColumn := m.playlistsTable.Columns()
	for i := range playlistsColumn {
		playlistsColumn[i].Width = leftPaneWidth/2 - 4
	}
	m.playlistsTable.SetColumns(playlistsColumn)

	tracksColumn := m.tracksTable.Columns()
	for i := range tracksColumn {
		tracksColumn[i].Width = rightPaneWidth/2 - 4
	}
	m.tracksTable.SetColumns(tracksColumn)
	leftPane := leftPaneStyle.Height(topPaneHeight).Width(leftPaneWidth / 2).MaxWidth(m.width / 2).Render(m.playlistsTable.View())
	rightPane := rightPaneStyle.Height(topPaneHeight).Width(rightPaneWidth / 2).MaxWidth(m.width / 2).Render(m.tracksTable.View())

	bottomPane := bottomPaneStyle.Width(m.width / 2).Render(renderControls(m))

	topPane := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPane,
		rightPane,
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		topPane,
		bottomPane,
	)

}

func (m *model) Init() tea.Cmd {

	playerState, err := getUserPlayerState(m.client)
	if err != nil {
		log.Println("Error getting player state:", err)
	}

	if playerState.Playing {
		m.status = "playing"
	} else {
		m.status = "paused"
	}
	m.volume = playerState.Device.Volume

	playlists, err := initializePlaylists(m.client)
	if err != nil {
		log.Fatal(err)
	}
	m.playlists = playlists
	m.playlistsTable = table.New(
		table.WithColumns([]table.Column{{Title: "Playlists", Width: 40}}),
		table.WithRows(convertPlaylistsToRows(m.playlists)),
		table.WithFocused(true),
	)
	m.playlistsTable.SetStyles(focusedTableStyle())

	tracks, err := loadTracksForPlaylist(m.client, m.playlists[m.selectedPlaylist].ID)
	if err != nil {
		return nil
	}
	m.tracks = tracks
	m.tracksTable = table.New(
		table.WithColumns([]table.Column{{Title: "Tracks", Width: 60}}),
		table.WithRows(convertTracksToRows(m.tracks)),
	)
	m.tracksTable.SetStyles(unfocusedTableStyle())
	m.focused = "playlists"

	return nil
}

func main() {

	http.HandleFunc("/callback", spotifyAuthentication)
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	log.Println("Log into Spotify by visiting the URL:", url)
	spotifyClient := <-ch

	log.Println("Initialized the Spotify client")

	user, err := spotifyClient.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Logged in as:", user.ID)
	m := &model{client: spotifyClient}
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal("error occurred:", err)
	}
}
