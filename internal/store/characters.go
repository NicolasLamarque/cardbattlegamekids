package store

import (
	"database/sql"
	_ "embed"
	"encoding/json"
)

//go:embed seed/characters.jjk.json
var seedCharactersJSON []byte

type seedCharacter struct {
	ID   int `json:"id"`
	Name struct {
		Full   string `json:"full"`
		Native string `json:"native"`
	} `json:"name"`
	Image struct {
		Large  string `json:"large"`
		Medium string `json:"medium"`
		Local  string `json:"local"`
	} `json:"image"`
	Favourites  int    `json:"favourites"`
	Role        string `json:"role"`
	SeriesID    int    `json:"seriesId"`
	SeriesTitle string `json:"seriesTitle"`
	SeriesColor string `json:"seriesColor"`
}

type Character struct {
	ID                 int    `json:"id"`
	NameFull           string `json:"nameFull"`
	NameNative         string `json:"nameNative"`
	ImageLarge         string `json:"imageLarge"`
	ImageMedium        string `json:"imageMedium"`
	ImageLocal         string `json:"imageLocal"`
	Favourites         int    `json:"favourites"`
	Role               string `json:"role"`
	SeriesID           int    `json:"seriesId"`
	SeriesTitle        string `json:"seriesTitle"`
	SeriesColor        string `json:"seriesColor"`
	Force              int    `json:"force"`
	PV                 int    `json:"pv"`
	IsCustomStats      bool   `json:"isCustomStats"`
	CustomData         string `json:"customData"`
	MaxMana            int    `json:"maxMana"`
	ManaRegen          int    `json:"manaRegen"`
	MaxAttacksOverride *int   `json:"maxAttacksOverride"`
	MaxAttacks         int    `json:"maxAttacks"`
}

// SyncSeedCharacters importe les personnages du fixture AniList embarqué qui
// ne sont pas encore dans la base — les lignes existantes ne sont jamais
// touchées, donc élargir le fixture n'écrase rien. Force et PV sont calculés
// une seule fois ici, avec la même formule que game/stats.js, puis stockés
// pour de bon.
func (s *Store) SyncSeedCharacters() error {
	var seed []seedCharacter
	if err := json.Unmarshal(seedCharactersJSON, &seed); err != nil {
		return err
	}

	for _, c := range seed {
		ratio := ratioFromID(c.ID)
		force, pv := statsFromFavourites(c.Favourites, ratio)
		budget := budgetFromFavourites(c.Favourites)
		maxMana := maxManaFromBudget(budget)
		manaRegen := defaultManaRegen(maxMana)

		_, err := s.db.Exec(`
			INSERT OR IGNORE INTO characters (
				id, name_full, name_native, image_large, image_medium, image_local,
				favourites, role, series_id, series_title, series_color, force, pv,
				is_custom_stats, custom_data, max_mana, mana_regen
			) VALUES (
				:id, :name_full, :name_native, :image_large, :image_medium, :image_local,
				:favourites, :role, :series_id, :series_title, :series_color, :force, :pv,
				0, '{}', :max_mana, :mana_regen
			)`,
			sql.Named("id", c.ID),
			sql.Named("name_full", c.Name.Full),
			sql.Named("name_native", c.Name.Native),
			sql.Named("image_large", c.Image.Large),
			sql.Named("image_medium", c.Image.Medium),
			sql.Named("image_local", c.Image.Local),
			sql.Named("favourites", c.Favourites),
			sql.Named("role", c.Role),
			sql.Named("series_id", c.SeriesID),
			sql.Named("series_title", c.SeriesTitle),
			sql.Named("series_color", c.SeriesColor),
			sql.Named("force", force),
			sql.Named("pv", pv),
			sql.Named("max_mana", maxMana),
			sql.Named("mana_regen", manaRegen),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateCharacterStats écrase Force/PV à la main et marque le personnage
// comme "custom" — il ne sera plus recalculé automatiquement.
func (s *Store) UpdateCharacterStats(id, force, pv int) error {
	_, err := s.db.Exec(`
		UPDATE characters
		SET force = :force, pv = :pv, is_custom_stats = 1
		WHERE id = :id`,
		sql.Named("force", force),
		sql.Named("pv", pv),
		sql.Named("id", id),
	)
	return err
}

// ResetCharacterStats recalcule Force/PV depuis les favoris avec la formule
// standard et retire le marqueur "custom".
func (s *Store) ResetCharacterStats(id int) error {
	var favourites int
	err := s.db.QueryRow(`SELECT favourites FROM characters WHERE id = :id`, sql.Named("id", id)).Scan(&favourites)
	if err != nil {
		return err
	}

	ratio := ratioFromID(id)
	force, pv := statsFromFavourites(favourites, ratio)

	_, err = s.db.Exec(`
		UPDATE characters
		SET force = :force, pv = :pv, is_custom_stats = 0
		WHERE id = :id`,
		sql.Named("force", force),
		sql.Named("pv", pv),
		sql.Named("id", id),
	)
	return err
}

func (s *Store) ListCharacters() ([]Character, error) {
	rows, err := s.db.Query(`
		SELECT id, name_full, name_native, image_large, image_medium, image_local,
			favourites, role, series_id, series_title, series_color, force, pv,
			is_custom_stats, custom_data, max_mana, mana_regen, max_attacks_override
		FROM characters
		ORDER BY favourites DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Character
	for rows.Next() {
		var c Character
		var isCustom int
		var maxAttacksOverride sql.NullInt64
		if err := rows.Scan(
			&c.ID, &c.NameFull, &c.NameNative, &c.ImageLarge, &c.ImageMedium, &c.ImageLocal,
			&c.Favourites, &c.Role, &c.SeriesID, &c.SeriesTitle, &c.SeriesColor, &c.Force, &c.PV,
			&isCustom, &c.CustomData, &c.MaxMana, &c.ManaRegen, &maxAttacksOverride,
		); err != nil {
			return nil, err
		}
		c.IsCustomStats = isCustom != 0
		budget := budgetFromFavourites(c.Favourites)
		c.MaxAttacks = maxAttacksFromBudget(budget)
		if maxAttacksOverride.Valid {
			v := int(maxAttacksOverride.Int64)
			c.MaxAttacksOverride = &v
			c.MaxAttacks = v
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// UpdateCharacterManaSettings règle la recharge de mana par manche et, en
// option, force le nombre max d'attaques (nil pour revenir au calcul auto).
func (s *Store) UpdateCharacterManaSettings(id, manaRegen int, maxAttacksOverride *int) error {
	_, err := s.db.Exec(`
		UPDATE characters
		SET mana_regen = :mana_regen, max_attacks_override = :max_attacks_override
		WHERE id = :id`,
		sql.Named("mana_regen", manaRegen),
		sql.Named("max_attacks_override", maxAttacksOverride),
		sql.Named("id", id),
	)
	return err
}
