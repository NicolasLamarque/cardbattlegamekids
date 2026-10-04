package store

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"math"
)

//go:embed seed/attacks.jjk.json
var seedAttacksJSON []byte

type seedAttackEntry struct {
	CharacterID int `json:"characterId"`
	Attacks     []struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"attacks"`
}

// Attack est une technique liable à une carte au combat : elle ajoute
// ForceBonus à la Force de la carte jouée, en échange de ManaCost de mana.
type Attack struct {
	ID          int    `json:"id"`
	CharacterID int    `json:"characterId"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	ForceBonus  int    `json:"forceBonus"`
	ManaCost    int    `json:"manaCost"`
	SortOrder   int    `json:"sortOrder"`
}

// genericAttackPalette tourne pour les personnages sans technique connue —
// chaque attaque générée reste visuellement distincte.
var genericAttackPalette = []string{
	"#f87171", "#fb923c", "#facc15", "#4ade80", "#38bdf8", "#818cf8", "#e879f9", "#f472b6",
}

var genericAttackNames = []string{
	"Frappe maudite", "Onde de choc", "Lame d'énergie", "Déchaînement", "Technique secrète",
}

// SyncSeedAttacks remplit la table des attaques pour tout personnage qui
// n'en a encore aucune — jamais pour un personnage déjà pourvu (importer le
// fixture plus tard, ou l'étendre, ne touche jamais aux lignes déjà là,
// qu'elles viennent du seed ou d'une édition dans l'admin).
func (s *Store) SyncSeedAttacks() error {
	var curated []seedAttackEntry
	if err := json.Unmarshal(seedAttacksJSON, &curated); err != nil {
		return err
	}
	curatedByID := make(map[int]seedAttackEntry, len(curated))
	for _, c := range curated {
		curatedByID[c.CharacterID] = c
	}

	existing := map[int]bool{}
	rows, err := s.db.Query(`SELECT DISTINCT character_id FROM character_attacks`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		existing[id] = true
	}
	rows.Close()

	chars, err := s.db.Query(`SELECT id, force, max_mana FROM characters`)
	if err != nil {
		return err
	}
	type charInfo struct {
		id      int
		force   int
		maxMana int
	}
	var toSeed []charInfo
	for chars.Next() {
		var c charInfo
		if err := chars.Scan(&c.id, &c.force, &c.maxMana); err != nil {
			chars.Close()
			return err
		}
		if !existing[c.id] {
			toSeed = append(toSeed, c)
		}
	}
	chars.Close()

	for _, c := range toSeed {
		maxAttacks := maxAttacksFromBudget(c.maxMana)

		var names, colors []string
		if entry, ok := curatedByID[c.id]; ok {
			for _, a := range entry.Attacks {
				names = append(names, a.Name)
				colors = append(colors, a.Color)
			}
		} else {
			for i := 0; i < maxAttacks; i++ {
				names = append(names, genericAttackNames[i%len(genericAttackNames)])
				colors = append(colors, genericAttackPalette[i%len(genericAttackPalette)])
			}
		}

		if len(names) > maxAttacks {
			names = names[:maxAttacks]
			colors = colors[:maxAttacks]
		}

		for i, name := range names {
			forceBonus := int(math.Round(float64(c.force) * (0.2 + 0.15*float64(i))))
			manaCost := int(math.Round(float64(c.maxMana) * (0.15 + 0.15*float64(i))))
			if manaCost < 1 {
				manaCost = 1
			}
			_, err := s.db.Exec(`
				INSERT INTO character_attacks (character_id, name, color, force_bonus, mana_cost, sort_order)
				VALUES (:character_id, :name, :color, :force_bonus, :mana_cost, :sort_order)`,
				sql.Named("character_id", c.id),
				sql.Named("name", name),
				sql.Named("color", colors[i]),
				sql.Named("force_bonus", forceBonus),
				sql.Named("mana_cost", manaCost),
				sql.Named("sort_order", i),
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Store) ListCharacterAttacks() ([]Attack, error) {
	rows, err := s.db.Query(`
		SELECT id, character_id, name, color, force_bonus, mana_cost, sort_order
		FROM character_attacks
		ORDER BY character_id, sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Attack
	for rows.Next() {
		var a Attack
		if err := rows.Scan(&a.ID, &a.CharacterID, &a.Name, &a.Color, &a.ForceBonus, &a.ManaCost, &a.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// SaveCharacterAttack crée une attaque (ID 0) ou met à jour celle désignée
// par ID — l'admin peut donc tout changer à la main après le remplissage
// automatique.
func (s *Store) SaveCharacterAttack(a Attack) (Attack, error) {
	if a.ID == 0 {
		res, err := s.db.Exec(`
			INSERT INTO character_attacks (character_id, name, color, force_bonus, mana_cost, sort_order)
			VALUES (:character_id, :name, :color, :force_bonus, :mana_cost, :sort_order)`,
			sql.Named("character_id", a.CharacterID),
			sql.Named("name", a.Name),
			sql.Named("color", a.Color),
			sql.Named("force_bonus", a.ForceBonus),
			sql.Named("mana_cost", a.ManaCost),
			sql.Named("sort_order", a.SortOrder),
		)
		if err != nil {
			return Attack{}, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return Attack{}, err
		}
		a.ID = int(id)
		return a, nil
	}

	_, err := s.db.Exec(`
		UPDATE character_attacks
		SET name = :name, color = :color, force_bonus = :force_bonus, mana_cost = :mana_cost, sort_order = :sort_order
		WHERE id = :id`,
		sql.Named("name", a.Name),
		sql.Named("color", a.Color),
		sql.Named("force_bonus", a.ForceBonus),
		sql.Named("mana_cost", a.ManaCost),
		sql.Named("sort_order", a.SortOrder),
		sql.Named("id", a.ID),
	)
	return a, err
}

func (s *Store) DeleteCharacterAttack(id int) error {
	_, err := s.db.Exec(`DELETE FROM character_attacks WHERE id = :id`, sql.Named("id", id))
	return err
}

// ResetCharacterAttacks efface les attaques d'un personnage et les
// régénère avec les valeurs par défaut (seed curé ou générique).
func (s *Store) ResetCharacterAttacks(characterID int) error {
	if _, err := s.db.Exec(`DELETE FROM character_attacks WHERE character_id = :id`, sql.Named("id", characterID)); err != nil {
		return err
	}
	return s.SyncSeedAttacks()
}
