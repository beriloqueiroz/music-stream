package rest_server_playlist

import (
	"encoding/json"
	"net/http"

	"github.com/beriloqueiroz/music-stream/internal/application"
)

type PlaylistHandler struct {
	service *application.PlaylistService
}

func NewPlaylistHandler(service *application.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{service: service}
}

type createPlaylistRequest struct {
	Name string `json:"name"`
}

type addMusicInPlaylistRequest struct {
	MusicID string `json:"music_id"`
}

type removeMusicInPlaylistRequest struct {
	MusicID string `json:"music_id"`
}

type updatePlaylistRequest struct {
	Name string `json:"name"`
}

type deletePlaylistRequest struct {
}

type getPlaylistRequest struct {
}

func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req createPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ownerID := r.Header.Get("owner_id")
	if ownerID == "" {
		http.Error(w, "Owner ID is required", http.StatusBadRequest)
		return
	}
	playlist, err := h.service.CreatePlaylist(r.Context(), req.Name, ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(playlist)
}

// POST /api/playlists/{id}/musics
func (h *PlaylistHandler) AddMusicInPlaylist(w http.ResponseWriter, r *http.Request) {
	var req addMusicInPlaylistRequest
	playlistID := r.PathValue("id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ownerID := r.Header.Get("owner_id")
	if ownerID == "" {
		http.Error(w, "Owner ID is required", http.StatusBadRequest)
		return
	}
	err := h.service.AddMusicToPlaylist(r.Context(), playlistID, req.MusicID, ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// DELETE /api/playlists/{id}/musics/{musicId}
func (h *PlaylistHandler) RemoveMusicInPlaylist(w http.ResponseWriter, r *http.Request) {
	var req removeMusicInPlaylistRequest
	playlistID := r.PathValue("id")
	musicID := r.PathValue("musicId")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ownerID := r.Header.Get("owner_id")
	if ownerID == "" {
		http.Error(w, "Owner ID is required", http.StatusBadRequest)
		return
	}
	err := h.service.RemoveMusicFromPlaylist(r.Context(), playlistID, musicID, ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(req)
}

// PUT /api/playlists/{id}
func (h *PlaylistHandler) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req updatePlaylistRequest
	playlistID := r.PathValue("id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ownerID := r.Header.Get("owner_id")
	if ownerID == "" {
		http.Error(w, "Owner ID is required", http.StatusBadRequest)
		return
	}
	playlist, err := h.service.UpdatePlaylist(r.Context(), playlistID, req.Name, ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(playlist)
}

// DELETE /api/playlists/{id}
func (h *PlaylistHandler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	var req deletePlaylistRequest
	playlistID := r.PathValue("id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ownerID := r.Header.Get("owner_id")
	if ownerID == "" {
		http.Error(w, "Owner ID is required", http.StatusBadRequest)
		return
	}
	err := h.service.DeletePlaylist(r.Context(), playlistID, ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(req)
}

// GET /api/playlists/{id}/musics
func (h *PlaylistHandler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	var req getPlaylistRequest
	playlistID := r.PathValue("id")
	if playlistID == "" {
		http.Error(w, "Playlist ID is required", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ownerID := r.Header.Get("owner_id")
	if ownerID == "" {
		http.Error(w, "Owner ID is required", http.StatusBadRequest)
		return
	}
	playlist, err := h.service.GetPlaylist(r.Context(), playlistID, ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(playlist)
}

// GET /api/playlists
func (h *PlaylistHandler) GetPlaylists(w http.ResponseWriter, r *http.Request) {
	var req getPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ownerID := r.Header.Get("owner_id")
	if ownerID == "" {
		http.Error(w, "Owner ID is required", http.StatusBadRequest)
		return
	}
	playlists, err := h.service.GetPlaylists(r.Context(), ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(playlists)
}
