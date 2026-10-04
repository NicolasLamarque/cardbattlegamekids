package main

import (
	"context"
	"log"

	"CardsBattle/internal/lanserver"
	"CardsBattle/internal/store"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	db  *store.Store
	lan *lanserver.Server
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	db, err := store.Open()
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	if err := db.SyncSeedCharacters(); err != nil {
		log.Fatalf("sync characters: %v", err)
	}
	if err := db.SyncSeedModifiers(); err != nil {
		log.Fatalf("sync modifiers: %v", err)
	}
	a.db = db
}

// shutdown ferme proprement la base et le serveur local quand l'app se termine.
func (a *App) shutdown(ctx context.Context) {
	if a.lan != nil {
		a.lan.Stop()
	}
	if a.db != nil {
		a.db.Close()
	}
}

// GetCharacters retourne tous les personnages de la collection, stats incluses.
func (a *App) GetCharacters() ([]store.Character, error) {
	return a.db.ListCharacters()
}

// UpdateCharacterStats fixe Force/PV à la main pour un personnage (mode custom).
func (a *App) UpdateCharacterStats(id int, force int, pv int) error {
	return a.db.UpdateCharacterStats(id, force, pv)
}

// ResetCharacterStats recalcule Force/PV depuis les favoris AniList.
func (a *App) ResetCharacterStats(id int) error {
	return a.db.ResetCharacterStats(id)
}

// GetDraftModifiers retourne le catalogue des malus/bonus du Bocal.
func (a *App) GetDraftModifiers() ([]store.DraftModifier, error) {
	return a.db.ListDraftModifiers()
}

// SaveDraftModifier crée ou met à jour un malus/bonus du Bocal.
func (a *App) SaveDraftModifier(m store.DraftModifier) error {
	return a.db.SaveDraftModifier(m)
}

// DeleteDraftModifier retire un malus/bonus du catalogue.
func (a *App) DeleteDraftModifier(id string) error {
	return a.db.DeleteDraftModifier(id)
}

// StartLanServer démarre le petit serveur local et retourne son URL et un
// QR code prêt à afficher pour qu'un autre appareil du même WiFi se connecte.
func (a *App) StartLanServer() (map[string]string, error) {
	if a.lan != nil {
		qr, err := lanserver.QRCodeDataURI(a.lan.URL())
		if err != nil {
			return nil, err
		}
		return map[string]string{"url": a.lan.URL(), "qrCode": qr}, nil
	}

	s, err := lanserver.Start(assets, "frontend/dist",
		func(msg []byte) {
			wailsRuntime.EventsEmit(a.ctx, "guest:action", string(msg))
		},
		func() {
			wailsRuntime.EventsEmit(a.ctx, "guest:connected")
		},
		func() {
			wailsRuntime.EventsEmit(a.ctx, "guest:disconnected")
		},
	)
	if err != nil {
		return nil, err
	}
	a.lan = s

	qr, err := lanserver.QRCodeDataURI(s.URL())
	if err != nil {
		return nil, err
	}
	return map[string]string{"url": s.URL(), "qrCode": qr}, nil
}

// BroadcastGameState envoie l'état (déjà filtré côté hôte) à l'invité connecté.
func (a *App) BroadcastGameState(payload string) {
	if a.lan != nil {
		a.lan.Broadcast([]byte(payload))
	}
}

// StopLanServer arrête le serveur local.
func (a *App) StopLanServer() error {
	if a.lan == nil {
		return nil
	}
	err := a.lan.Stop()
	a.lan = nil
	return err
}

// ListPlayers retourne les joueurs enregistrés (nom, portefeuille, stats).
func (a *App) ListPlayers() ([]store.Player, error) {
	return a.db.ListPlayers()
}

// CreatePlayer enregistre un nouveau joueur.
func (a *App) CreatePlayer(name string) (store.Player, error) {
	return a.db.CreatePlayer(name)
}
