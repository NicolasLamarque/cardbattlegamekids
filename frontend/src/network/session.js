// L'app native Wails (l'hôte) a window.go injecté par le runtime. Un
// navigateur qui charge la page via le serveur LAN ne l'a pas — c'est
// comme ça qu'on distingue l'hôte de l'invité, sans rien à configurer.
export const isGuest = typeof window.go === 'undefined'
