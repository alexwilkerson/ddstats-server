package ddapi

import (
	"encoding/binary"
	"errors"
)

// DeathTypes as defined by the DD API
var DeathTypes = []string{
	"FALLEN", "SWARMED", "IMPALED", "GORED", "INFESTED", "OPENED", "PURGED",
	"DESECRATED", "SACRIFICED", "EVISCERATED", "ANNIHILATED", "INTOXICATED",
	"ENVENMONATED", "INCARNATED", "DISCARNATED", "BARBED",
}

// ErrPlayerNotFound returned when player not found from the DD API
var ErrPlayerNotFound = errors.New("player not found")

// Player is the struct returned after parsing the binary data
// blob returned from the DD API.
type Player struct {
	PlayerName          string  `json:"player_name"`
	PlayerID            uint64  `json:"player_id"`
	Rank                int32   `json:"rank"`
	Time                float64 `json:"time"`
	Kills               int32   `json:"kills"`
	Gems                int32   `json:"gems"`
	DaggersHit          int32   `json:"daggers_hit"`
	DaggersFired        int32   `json:"daggers_fired"`
	Accuracy            float64 `json:"accuracy"`
	DeathType           string  `json:"death_type"`
	OverallTime         float64 `json:"overall_time"`
	OverallKills        uint64  `json:"overall_kills"`
	OverallGems         uint64  `json:"overall_gems"`
	OverallDeaths       uint64  `json:"overall_deaths"`
	OverallDaggersHit   uint64  `json:"overall_daggers_hit"`
	OverallDaggersFired uint64  `json:"overall_daggers_fired"`
	OverallAccuracy     float64 `json:"overall_accuracy"`
}

// BytesToPlayer takes a byte array and an initial offset
// and returns a Player object. Will return an error if the
// Player is not found
func BytesToPlayer(b []byte, bytePosition int) (*Player, error) {
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
	player.Time = float64(toInt32(b, bytePosition+12)) / 10000
	player.Kills = toInt32(b, bytePosition+16)
	player.Gems = toInt32(b, bytePosition+28)
	player.DaggersHit = toInt32(b, bytePosition+24)
	player.DaggersFired = toInt32(b, bytePosition+20)
	if player.DaggersFired > 0 {
		player.Accuracy = float64(player.DaggersHit) / float64(player.DaggersFired) * 100
	}
	player.DeathType = DeathTypes[toInt16(b, bytePosition+32)]
	player.OverallTime = float64(toUint64(b, bytePosition+60)) / 10000
	player.OverallKills = toUint64(b, bytePosition+44)
	player.OverallGems = toUint64(b, bytePosition+68)
	player.OverallDeaths = toUint64(b, bytePosition+36)
	player.OverallDaggersHit = toUint64(b, bytePosition+76)
	player.OverallDaggersFired = toUint64(b, bytePosition+52)
	if player.OverallDaggersFired > 0 {
		player.OverallAccuracy = float64(player.OverallDaggersHit) / float64(player.OverallDaggersFired) * 100
	}

	return &player, nil
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
