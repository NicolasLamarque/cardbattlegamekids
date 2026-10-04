package store

import "database/sql"

type Player struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	WalletBalance int    `json:"walletBalance"`
	GamesWon      int    `json:"gamesWon"`
	GamesPlayed   int    `json:"gamesPlayed"`
}

func (s *Store) ListPlayers() ([]Player, error) {
	rows, err := s.db.Query(`SELECT id, name, wallet_balance, games_won, games_played FROM players ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Player
	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.ID, &p.Name, &p.WalletBalance, &p.GamesWon, &p.GamesPlayed); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) CreatePlayer(name string) (Player, error) {
	res, err := s.db.Exec(`INSERT INTO players (name) VALUES (:name)`, sql.Named("name", name))
	if err != nil {
		return Player{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Player{}, err
	}
	return Player{ID: int(id), Name: name}, nil
}
