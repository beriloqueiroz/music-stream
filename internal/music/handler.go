package music

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createPlaylistRequest struct {
	Name string `json:"name"`
}

type addMusicInPlaylistRequest struct {
	PlaylistID string `json:"playlist_id"`
	MusicID    string `json:"music_id"`
}

type removeMusicInPlaylistRequest struct {
	PlaylistID string `json:"playlist_id"`
	MusicID    string `json:"music_id"`
}

type updatePlaylistRequest struct {
	PlaylistID string `json:"playlist_id"`
	Name       string `json:"name"`
}

type deletePlaylistRequest struct {
	PlaylistID string `json:"playlist_id"`
}

type getPlaylistRequest struct {
	PlaylistID string `json:"playlist_id"`
}

func (h *Handler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req createPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) AddMusicInPlaylist(w http.ResponseWriter, r *http.Request) {
	var req addMusicInPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) RemoveMusicInPlaylist(w http.ResponseWriter, r *http.Request) {
	var req removeMusicInPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req updatePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	var req deletePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	var req getPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
