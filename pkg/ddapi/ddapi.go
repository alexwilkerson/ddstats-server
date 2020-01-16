package ddapi

import (
	"encoding/binary"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	// EndpointGetUserByID is the endpoint to get a user by ID
	EndpointGetUserByID = "http://dd.hasmodai.com/backend16/get_user_by_id_public.php"
	// EndpointGetUserByRank is the endpoint to get a user by rank
	EndpointGetUserByRank = "http://dd.hasmodai.com/backend16/get_user_by_rank_public.php"
	// EndpointGetScores is the endpoint to get the leaderboard
	EndpointGetScores = "http://dd.hasmodai.com/backend16/get_scores.php"
	// EndpointGetUserSearch is the endpoint to get search for users
	EndpointGetUserSearch = "http://dd.hasmodai.com/backend16/get_user_search_public.php"
)

var (
	// ErrPlayerNotFound returned when player not found from the DD API
	ErrPlayerNotFound = errors.New("player not found")
	// ErrNoPlayersFound is returned when user search produces no users
	ErrNoPlayersFound = errors.New("no players found")
	// ErrStatusCode is returned when the Devil Daggers server responds with an error
	ErrStatusCode = errors.New("error getting a response from the Devil Daggers API")
)

// API is used as an abstraction and to inject the client into the ddapi package
type API struct {
	Client  *http.Client
	Watcher *Watcher
}

// NewAPI returns an API struct
func NewAPI(client *http.Client) *API {
	return &API{
		Client:  client,
		Watcher: NewWatcher(),
	}
}

// DeathTypes as defined by the DD API
var DeathTypes = []string{
	"FALLEN",
	"SWARMED",
	"IMPALED",
	"GORED",
	"INFESTED",
	"OPENED",
	"PURGED",
	"DESECRATED",
	"SACRIFICED",
	"EVISCERATED",
	"ANNIHILATED",
	"INTOXICATED",
	"ENVENMONATED",
	"INCARNATED",
	"DISCARNATED",
	"BARBED",
}

// Player is the struct returned after parsing the binary data
// blob returned from the DD API.
type Player struct {
	PlayerID               uint64  `json:"player_id"`
	PlayerName             string  `json:"player_name"`
	Rank                   uint32  `json:"rank"`
	GameTime               float64 `json:"game_time"`
	EnemiesKilled          uint32  `json:"enemies_killed"`
	Gems                   uint32  `json:"gems"`
	DaggersHit             uint32  `json:"daggers_hit"`
	DaggersFired           uint32  `json:"daggers_fired"`
	Accuracy               float64 `json:"accuracy"`
	DeathType              string  `json:"death_type"`
	OverallGameTime        float64 `json:"overall_game_time"`
	OverallAverageGameTime float64 `json:"overall_average_game_time"`
	OverallEnemiesKilled   uint64  `json:"overall_enemies_killed"`
	OverallGems            uint64  `json:"overall_gems"`
	OverallDeaths          uint64  `json:"overall_deaths"`
	OverallDaggersHit      uint64  `json:"overall_daggers_hit"`
	OverallDaggersFired    uint64  `json:"overall_daggers_fired"`
	OverallAccuracy        float64 `json:"overall_accuracy"`
}

// Leaderboard is a struct returned after being converted from bytes
type Leaderboard struct {
	GlobalDeaths          uint64    `json:"global_deaths"`
	GlobalEnemiesKilled   uint64    `json:"global_enemies_killed"`
	GlobalGameTime        float64   `json:"global_game_time"`
	GlobalAverageGameTime float64   `json:"global_average_game_time"`
	GlobalGems            uint64    `json:"global_gems"`
	GlobalDaggersFired    uint64    `json:"global_daggers_fired"`
	GlobalDaggersHit      uint64    `json:"global_daggers_hit"`
	GlobalAccuracy        float64   `json:"global_accuracy"`
	GlobalPlayerCount     uint32    `json:"global_player_count"`
	PlayerCount           int       `json:"player_count"`
	Players               []*Player `json:"players"`
}

// UserByID hits the backend DD API and returns a Player
func (api *API) UserByID(id int) (*Player, error) {
	form := url.Values{"uid": {strconv.Itoa(id)}}
	resp, err := api.Client.PostForm(EndpointGetUserByID, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrStatusCode
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	player, _, err := bytesToPlayer(bodyBytes, 19)
	if err != nil {
		return nil, err
	}

	if player.PlayerName == "" {
		return nil, ErrPlayerNotFound
	}

	return player, nil
}

// UserByRank hits the backend DD API and returns a Player
func (api *API) UserByRank(rank int) (*Player, error) {
	form := url.Values{"rank": {strconv.Itoa(rank)}}
	resp, err := api.Client.PostForm(EndpointGetUserByRank, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrStatusCode
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	player, _, err := bytesToPlayer(bodyBytes, 19)
	if err != nil {
		return nil, err
	}

	return player, nil
}

// GetLeaderboard takes a limit and an offset, hits the backend DD API and returns
// a Leaderboard struct
func (api *API) GetLeaderboard(limit, offset int) (*Leaderboard, error) {
	// the DD API weirdly counts users starting from 1 but internally uses a 0 index
	// this fix it to make it more readable for users.
	if offset != 0 {
		offset--
	}

	form := url.Values{"user": {"0"}, "level": {"survival"}, "offset": {strconv.Itoa(offset)}}
	resp, err := api.Client.PostForm(EndpointGetScores, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrStatusCode
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	leaderboard, err := bytesToLeaderboard(bodyBytes, limit)
	if err != nil {
		return nil, err
	}

	return leaderboard, nil
}

// UserSearch takes a user name and hits the backend DD API and returns a slice of Players
func (api *API) UserSearch(name string) ([]*Player, error) {
	// The Devil Daggers API responds with no users found if the user name is
	// longer than 16 characters. This truncates the name to 16 characters and
	// does a partial match, hopefully finding the intended user. If not, will
	// return the intended user and any other user which matches the substring.
	name = strings.TrimSpace(name)
	if len(name) > 16 {
		name = name[:16]
	}
	form := url.Values{"search": {name}}
	resp, err := api.Client.PostForm(EndpointGetUserSearch, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrStatusCode
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	players, err := userSearchBytesToPlayers(bodyBytes)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(players, func(i, j int) bool {
		return players[i].Rank < players[j].Rank
	})

	return players, nil
}

// BytesToPlayer takes a byte array and an initial offset
// and returns a Player object. Will return an error if the
// Player is not found
func bytesToPlayer(b []byte, bytePosition int) (*Player, int, error) {
	var player Player

	playerNameLength := uint8(b[bytePosition])

	bytePosition += 2
	player.PlayerName = strings.ToValidUTF8(string(b[bytePosition:bytePosition+int(playerNameLength)]), "?")
	bytePosition += int(playerNameLength)
	// just figured out this information manually...
	player.PlayerID = toUint64(b, bytePosition+4)
	if player.PlayerID == 0 {
		return nil, 0, ErrPlayerNotFound
	}
	player.Rank = toUint32(b, bytePosition)
	player.GameTime = roundToNearest(float64(toUint32(b, bytePosition+12))/10000, 4)
	player.EnemiesKilled = toUint32(b, bytePosition+16)
	player.Gems = toUint32(b, bytePosition+28)
	player.DaggersHit = toUint32(b, bytePosition+24)
	player.DaggersFired = toUint32(b, bytePosition+20)
	if player.DaggersFired > 0 {
		player.Accuracy = roundToNearest(float64(player.DaggersHit)/float64(player.DaggersFired)*100, 2)
	}
	player.DeathType = DeathTypes[toUint16(b, bytePosition+32)]
	player.OverallGameTime = roundToNearest(float64(toUint64(b, bytePosition+60))/10000, 4)
	player.OverallDeaths = toUint64(b, bytePosition+36)
	player.OverallAverageGameTime = roundToNearest(player.OverallGameTime/float64(player.OverallDeaths), 4)
	player.OverallEnemiesKilled = toUint64(b, bytePosition+44)
	player.OverallGems = toUint64(b, bytePosition+68)
	player.OverallDaggersHit = toUint64(b, bytePosition+76)
	player.OverallDaggersFired = toUint64(b, bytePosition+52)
	if player.OverallDaggersFired > 0 {
		player.OverallAccuracy = roundToNearest(float64(player.OverallDaggersHit)/float64(player.OverallDaggersFired)*100, 2)
	}

	return &player, int(playerNameLength), nil
}

// GetScoresBytesToLeaderboard converts the byte array from the DD API
// to a Leaderboard struct
func bytesToLeaderboard(b []byte, limit int) (*Leaderboard, error) {
	outfile, err := os.Create("leaderboard.bytes")
	if err != nil {
		return nil, err
	}
	outfile.Write(b)
	outfile.Close()

	var leaderboard Leaderboard
	leaderboard.Players = []*Player{} // init this so won't be nil

	leaderboard.GlobalDeaths = toUint64(b, 11)
	leaderboard.GlobalEnemiesKilled = toUint64(b, 19)
	leaderboard.GlobalGameTime = roundToNearest(float64(toUint64(b, 35))/1000, 4)
	leaderboard.GlobalAverageGameTime = roundToNearest(leaderboard.GlobalGameTime/float64(leaderboard.GlobalDeaths), 4)
	leaderboard.GlobalGems = toUint64(b, 43)
	leaderboard.GlobalDaggersHit = toUint64(b, 51)
	leaderboard.GlobalDaggersFired = toUint64(b, 27)
	if leaderboard.GlobalDaggersFired > 0 {
		leaderboard.GlobalAccuracy = float64(leaderboard.GlobalDaggersHit) / float64(leaderboard.GlobalDaggersFired) * 100
	}
	leaderboard.GlobalPlayerCount = toUint32(b, 75)

	leaderboard.PlayerCount = int(toUint16(b, 59))
	if limit < leaderboard.PlayerCount {
		leaderboard.PlayerCount = limit
	}

	offset := 83
	for i := 0; i < leaderboard.PlayerCount; i++ {
		p, playerNameLength, err := bytesToPlayer(b, offset)
		if err != nil {
			return nil, ErrPlayerNotFound
		}
		offset += playerNameLength + 90
		leaderboard.Players = append(leaderboard.Players, p)
	}

	return &leaderboard, nil
}

// UserSearchBytesToPlayers converts a byte array to a player slice
func userSearchBytesToPlayers(b []byte) ([]*Player, error) {
	playerCount := int(toUint16(b, 11))
	if playerCount < 1 {
		return nil, ErrNoPlayersFound
	}
	var players []*Player
	offset := 19
	for i := 0; i < playerCount; i++ {
		p, playerNameLength, err := bytesToPlayer(b, offset)
		if err != nil {
			return nil, ErrPlayerNotFound
		}
		offset += playerNameLength + 90
		players = append(players, p)
	}
	return players, nil
}

func toUint64(b []byte, offset int) uint64 {
	return binary.LittleEndian.Uint64(b[offset : offset+8])
}

func toUint32(b []byte, offset int) uint32 {
	return binary.LittleEndian.Uint32(b[offset : offset+4])
}

func toUint16(b []byte, offset int) uint16 {
	return binary.LittleEndian.Uint16(b[offset : offset+2])
}

func roundToNearest(f float64, numberOfDecimalPlaces int) float64 {
	multiplier := math.Pow10(numberOfDecimalPlaces)
	return math.Round(f*multiplier) / multiplier
}
