package lanserver

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	httpServer *http.Server
	upgrader   websocket.Upgrader
	url        string

	onGuestMessage    func([]byte)
	onGuestConnect    func()
	onGuestDisconnect func()

	mu      sync.Mutex
	clients map[*websocket.Conn]bool
}

// localLanIP trouve l'adresse IPv4 du vrai réseau WiFi/Ethernet — celle que
// les autres appareils du même réseau peuvent joindre. La machine a souvent
// plusieurs adaptateurs virtuels (Docker, WSL, VPN, Hyper-V) avec leur propre
// plage privée ; 192.168.x.x est la plage domestique la plus fiable, donc on
// la priorise avant de retomber sur 10.x.x.x ou 172.16-31.x.x.
func localLanIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	var candidates []net.IP
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}
		ip4 := ipNet.IP.To4()
		if ip4 == nil || ip4.IsLinkLocalUnicast() || !isPrivateIPv4(ip4) {
			continue
		}
		candidates = append(candidates, ip4)
	}

	for _, ip := range candidates {
		if ip[0] == 192 && ip[1] == 168 {
			return ip.String(), nil
		}
	}
	if len(candidates) > 0 {
		return candidates[0].String(), nil
	}
	return "", fmt.Errorf("aucune adresse IPv4 de réseau local trouvée")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip[0] == 10 ||
		(ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31) ||
		(ip[0] == 192 && ip[1] == 168)
}

// Start sert le frontend déjà construit et ouvre un endpoint WebSocket, sur
// une IP et un port choisis automatiquement, accessibles depuis le réseau local.
// onGuestMessage est appelé pour chaque message reçu d'un invité,
// onGuestConnect dès qu'un invité se connecte, et onGuestDisconnect quand il
// part — c'est l'hôte (via app.go) qui décide quoi en faire (renvoyer l'état
// actuel tout de suite pour ne pas dépendre du timing, afficher un avis, etc.).
func Start(
	assets embed.FS,
	assetsDir string,
	onGuestMessage func([]byte),
	onGuestConnect func(),
	onGuestDisconnect func(),
) (*Server, error) {
	ip, err := localLanIP()
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, err
	}

	sub, err := fs.Sub(assets, assetsDir)
	if err != nil {
		listener.Close()
		return nil, err
	}

	s := &Server{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		onGuestMessage:    onGuestMessage,
		onGuestConnect:    onGuestConnect,
		onGuestDisconnect: onGuestDisconnect,
		clients:           make(map[*websocket.Conn]bool),
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(sub)))
	mux.HandleFunc("/ws", s.handleWS)

	s.httpServer = &http.Server{Handler: mux}
	s.url = fmt.Sprintf("http://%s:%d", ip, listener.Addr().(*net.TCPAddr).Port)

	go s.httpServer.Serve(listener)

	return s, nil
}

func (s *Server) URL() string {
	return s.url
}

func (s *Server) Stop() error {
	return s.httpServer.Shutdown(context.Background())
}

// Broadcast envoie l'état de la partie (déjà filtré côté hôte) à tous les
// invités connectés.
func (s *Server) Broadcast(msg []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for conn := range s.clients {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			conn.Close()
			delete(s.clients, conn)
		}
	}
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	if s.onGuestConnect != nil {
		s.onGuestConnect()
	}

	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
		if s.onGuestDisconnect != nil {
			s.onGuestDisconnect()
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if s.onGuestMessage != nil {
			s.onGuestMessage(msg)
		}
	}
}
