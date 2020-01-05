package ddapi

import (
	"encoding/binary"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

// API is used as an abstraction and to inject the client into the ddapi package
type API struct {
	Client *http.Client
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

// ErrPlayerNotFound returned when player not found from the DD API
var ErrPlayerNotFound = errors.New("player not found")

// ErrNoPlayersFound is returned when user search produces no users
var ErrNoPlayersFound = errors.New("no players found")

// Player is the struct returned after parsing the binary data
// blob returned from the DD API.
type Player struct {
	PlayerName           string  `json:"player_name"`
	PlayerID             uint64  `json:"player_id"`
	Rank                 int32   `json:"rank"`
	GameTime             float64 `json:"game_time"`
	EnemiesKilled        int32   `json:"enemies_killed"`
	Gems                 int32   `json:"gems"`
	DaggersHit           int32   `json:"daggers_hit"`
	DaggersFired         int32   `json:"daggers_fired"`
	Accuracy             float64 `json:"accuracy"`
	DeathType            string  `json:"death_type"`
	OverallTime          float64 `json:"overall_time"`
	OverallEnemiesKilled uint64  `json:"overall_enemies_killed"`
	OverallGems          uint64  `json:"overall_gems"`
	OverallDeaths        uint64  `json:"overall_deaths"`
	OverallDaggersHit    uint64  `json:"overall_daggers_hit"`
	OverallDaggersFired  uint64  `json:"overall_daggers_fired"`
	OverallAccuracy      float64 `json:"overall_accuracy"`
}

// Leaderboard is a struct returned after being converted from bytes
type Leaderboard struct {
	GlobalDeaths        uint64    `json:"global_deaths"`
	GlobalEnemiesKilled uint64    `json:"global_enemies_killed"`
	GlobalTime          float64   `json:"global_time"`
	GlobalGems          uint64    `json:"global_gems"`
	GlobalDaggersFired  uint64    `json:"global_daggers_fired"`
	GlobalDaggersHit    uint64    `json:"global_daggers_hit"`
	GlobalAccuracy      float64   `json:"global_accuracy"`
	GlobalPlayerCount   int32     `json:"global_player_count"`
	PlayerCount         int       `json:"player_count"`
	Players             []*Player `json:"players"`
}

// UserByID hits the backend DD API and returns a Player
func (api *API) UserByID(id int) (*Player, error) {
	u := "http://dd.hasmodai.com/backend16/get_user_by_id_public.php"
	form := url.Values{"uid": {strconv.Itoa(id)}}
	resp, err := api.Client.PostForm(u, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	player, err := bytesToPlayer(bodyBytes, 19)
	if err != nil {
		return nil, err
	}

	return player, nil
}

// UserByRank hits the backend DD API and returns a Player
func (api *API) UserByRank(rank int) (*Player, error) {
	u := "http://dd.hasmodai.com/backend16/get_user_by_rank_public.php"
	form := url.Values{"rank": {strconv.Itoa(rank)}}
	resp, err := api.Client.PostForm(u, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	player, err := bytesToPlayer(bodyBytes, 19)
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

	u := "http://dd.hasmodai.com/backend16/get_scores.php"
	form := url.Values{"user": {"0"}, "level": {"survival"}, "offset": {strconv.Itoa(offset)}}
	resp, err := api.Client.PostForm(u, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
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
	u := "http://dd.hasmodai.com/backend16/get_user_search_public.php"
	form := url.Values{"search": {name}}
	resp, err := api.Client.PostForm(u, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	players, err := userSearchBytesToPlayers(bodyBytes)
	if err != nil {
		return nil, err
	}

	return players, nil
}

// BytesToPlayer takes a byte array and an initial offset
// and returns a Player object. Will return an error if the
// Player is not found
func bytesToPlayer(b []byte, bytePosition int) (*Player, error) {
	var player Player

	playerNameLength := int(toInt16(b, bytePosition))
	bytePosition += 2
	player.PlayerName = string(b[bytePosition : bytePosition+playerNameLength])
	bytePosition += playerNameLength
	// just figured out this information manually...
	player.PlayerID = toUint64(b, bytePosition+4)
	if player.PlayerID == 0 {
		return nil, ErrPlayerNotFound
	}
	player.Rank = toInt32(b, bytePosition)
	player.GameTime = roundToNearest(float64(toInt32(b, bytePosition+12))/10000, 4)
	player.EnemiesKilled = toInt32(b, bytePosition+16)
	player.Gems = toInt32(b, bytePosition+28)
	player.DaggersHit = toInt32(b, bytePosition+24)
	player.DaggersFired = toInt32(b, bytePosition+20)
	if player.DaggersFired > 0 {
		player.Accuracy = roundToNearest(float64(player.DaggersHit)/float64(player.DaggersFired)*100, 2)
	}
	player.DeathType = DeathTypes[toInt16(b, bytePosition+32)]
	player.OverallTime = roundToNearest(float64(toUint64(b, bytePosition+60))/10000, 4)
	player.OverallEnemiesKilled = toUint64(b, bytePosition+44)
	player.OverallGems = toUint64(b, bytePosition+68)
	player.OverallDeaths = toUint64(b, bytePosition+36)
	player.OverallDaggersHit = toUint64(b, bytePosition+76)
	player.OverallDaggersFired = toUint64(b, bytePosition+52)
	if player.OverallDaggersFired > 0 {
		player.OverallAccuracy = roundToNearest(float64(player.OverallDaggersHit)/float64(player.OverallDaggersFired)*100, 2)
	}

	return &player, nil
}

// GetScoresBytesToLeaderboard converts the byte array from the DD API
// to a Leaderboard struct
func bytesToLeaderboard(b []byte, limit int) (*Leaderboard, error) {
	var leaderboard Leaderboard

	leaderboard.GlobalDeaths = toUint64(b, 11)
	leaderboard.GlobalEnemiesKilled = toUint64(b, 19)
	leaderboard.GlobalTime = roundToNearest(float64(toUint64(b, 35))/1000, 4)
	leaderboard.GlobalGems = toUint64(b, 43)
	leaderboard.GlobalDaggersHit = toUint64(b, 51)
	leaderboard.GlobalDaggersFired = toUint64(b, 27)
	if leaderboard.GlobalDaggersFired > 0 {
		leaderboard.GlobalAccuracy = float64(leaderboard.GlobalDaggersHit) / float64(leaderboard.GlobalDaggersFired)
	}
	leaderboard.GlobalPlayerCount = toInt32(b, 75)

	leaderboard.PlayerCount = int(toInt16(b, 59))
	if limit < leaderboard.PlayerCount {
		leaderboard.PlayerCount = limit
	}

	offset := 83
	for i := 0; i < leaderboard.PlayerCount; i++ {
		p, err := bytesToPlayer(b, offset)
		if err != nil {
			return nil, ErrPlayerNotFound
		}
		offset += len(p.PlayerName) + 90
		leaderboard.Players = append(leaderboard.Players, p)
	}

	return &leaderboard, nil
}

// UserSearchBytesToPlayers converts a byte array to a player slice
func userSearchBytesToPlayers(b []byte) ([]*Player, error) {
	playerCount := int(toInt16(b, 11))
	if playerCount < 1 {
		return nil, ErrNoPlayersFound
	}
	var players []*Player
	offset := 19
	for i := 0; i < playerCount; i++ {
		p, err := bytesToPlayer(b, offset)
		if err != nil {
			return nil, ErrPlayerNotFound
		}
		offset += len(p.PlayerName) + 90
		players = append(players, p)
	}
	return players, nil
}

func toUint64(b []byte, offset int) uint64 {
	return binary.LittleEndian.Uint64(b[offset : offset+8])
}

func toInt64(b []byte, offset int) int64 {
	return int64(binary.LittleEndian.Uint64(b[offset : offset+4]))
}

func toUint32(b []byte, offset int) uint32 {
	return binary.LittleEndian.Uint32(b[offset : offset+4])
}

func toInt32(b []byte, offset int) int32 {
	return int32(binary.LittleEndian.Uint32(b[offset : offset+4]))
}

func toInt16(b []byte, offset int) int16 {
	return int16(binary.LittleEndian.Uint16(b[offset : offset+2]))
}

func roundToNearest(f float64, numberOfDecimalPlaces int) float64 {
	multiplier := math.Pow10(numberOfDecimalPlaces)
	return math.Round(f*multiplier) / multiplier
}
