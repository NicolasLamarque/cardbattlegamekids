package store

import "database/sql"

type DraftModifier struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Type   string `json:"type"` // "bonus" | "malus"
	Kind   string `json:"kind"` // "money" | "slot"
	Amount int    `json:"amount"`
	Rare   bool   `json:"rare"`
}

type seedModifier struct {
	ID     string
	Label  string
	Type   string
	Kind   string
	Amount int
	Rare   bool
}

var defaultModifiers = []seedModifier{
	{"malus-money-small", "Frais de douane", "malus", "money", 25000, false},
	{"malus-money-big", "Taxe surprise", "malus", "money", 50000, false},
	{"bonus-money-small", "Pourboire", "bonus", "money", 25000, false},
	{"bonus-money-big", "Jackpot", "bonus", "money", 50000, false},
	{"bonus-slot", "Slot bonus", "bonus", "slot", 1, true},
	{"malus-slot", "Slot perdu", "malus", "slot", 1, true},
}

// SyncSeedModifiers importe le catalogue de départ des malus/bonus du Bocal
// — n'écrase jamais une ligne déjà modifiée à la main.
func (s *Store) SyncSeedModifiers() error {
	for _, m := range defaultModifiers {
		rare := 0
		if m.Rare {
			rare = 1
		}
		_, err := s.db.Exec(`
			INSERT OR IGNORE INTO draft_modifiers (id, label, type, kind, amount, rare)
			VALUES (:id, :label, :type, :kind, :amount, :rare)`,
			sql.Named("id", m.ID),
			sql.Named("label", m.Label),
			sql.Named("type", m.Type),
			sql.Named("kind", m.Kind),
			sql.Named("amount", m.Amount),
			sql.Named("rare", rare),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListDraftModifiers() ([]DraftModifier, error) {
	rows, err := s.db.Query(`SELECT id, label, type, kind, amount, rare FROM draft_modifiers ORDER BY rowid`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []DraftModifier
	for rows.Next() {
		var m DraftModifier
		var rare int
		if err := rows.Scan(&m.ID, &m.Label, &m.Type, &m.Kind, &m.Amount, &rare); err != nil {
			return nil, err
		}
		m.Rare = rare != 0
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *Store) SaveDraftModifier(m DraftModifier) error {
	rare := 0
	if m.Rare {
		rare = 1
	}
	_, err := s.db.Exec(`
		INSERT INTO draft_modifiers (id, label, type, kind, amount, rare)
		VALUES (:id, :label, :type, :kind, :amount, :rare)
		ON CONFLICT(id) DO UPDATE SET
			label = excluded.label,
			type = excluded.type,
			kind = excluded.kind,
			amount = excluded.amount,
			rare = excluded.rare`,
		sql.Named("id", m.ID),
		sql.Named("label", m.Label),
		sql.Named("type", m.Type),
		sql.Named("kind", m.Kind),
		sql.Named("amount", m.Amount),
		sql.Named("rare", rare),
	)
	return err
}

func (s *Store) DeleteDraftModifier(id string) error {
	_, err := s.db.Exec(`DELETE FROM draft_modifiers WHERE id = :id`, sql.Named("id", id))
	return err
}
