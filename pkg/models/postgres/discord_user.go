package postgres

import (
	"database/sql"
	"errors"

	"github.com/alexwilkerson/ddstats-server/pkg/models"
	"github.com/jmoiron/sqlx"
)

type DiscordUserModel struct {
	DB *sqlx.DB
}

func (du *DiscordUserModel) Upsert(discordID string, ddID int) error {
	// discordUser, err := du.Select(discordID)
	// if err != nil && !errors.Is(err, models.ErrNoDiscordUserFound) {
	// 	if errors.Is(err, models.ErrNoDiscordUserFound) {
	// 		stmt := `
	// 		INSERT INTO discord_user(discord_id, dd_id)
	// 		VALUES ($1, $2)`
	// 		_, err = du.DB.Exec(stmt)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	}
	// 	return err
	// }
	// if discordUser.Verified {
	// 	return models.ErrDiscordUserVerified
	// }
	stmt := `
		INSERT INTO discord_user(discord_id, dd_id)
		VALUES ($1, $2)
		ON CONFLICT (discord_id) DO
		UPDATE SET dd_id=$2
		WHERE discord_user.discord_id=$1`
	_, err := du.DB.Exec(stmt, discordID, ddID)
	if err != nil {
		return err
	}
	return nil
}

func (du *DiscordUserModel) Select(discordID string) (*models.DiscordUser, error) {
	var discordUser models.DiscordUser
	stmt := `
		SELECT *
		FROM discord_user
		WHERE discord_id=$1`
	err := du.DB.Get(&discordUser, stmt, discordID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoDiscordUserFound
		}
		return nil, err
	}
	return &discordUser, nil
}

func (du *DiscordUserModel) Verified(discordID string) (bool, error) {
	var verified bool
	stmt := `
		SELECT verified
		FROM discord_user
		WHERE discord_id=$1`
	err := du.DB.QueryRow(stmt, discordID).Scan(&verified)
	if err != nil {
		return false, err
	}
	return verified, nil
}
