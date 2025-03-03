package playlist

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
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type addMusicInPlaylistRequest struct {
	MusicID string `json:"music_id"`
	OwnerID string `json:"owner_id"`
}

type removeMusicInPlaylistRequest struct {
	MusicID string `json:"music_id"`
	OwnerID string `json:"owner_id"`
}

type updatePlaylistRequest struct {
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type deletePlaylistRequest struct {
	OwnerID string `json:"owner_id"`
}

type getPlaylistRequest struct {
	OwnerID string `json:"owner_id"`
}

func (h *Handler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req createPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playlist, err := h.service.CreatePlaylist(r.Context(), req.Name, req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(playlist)
}

// POST /api/playlists/{id}/musics
func (h *Handler) AddMusicInPlaylist(w http.ResponseWriter, r *http.Request) {
	var req addMusicInPlaylistRequest
	playlistID := r.PathValue("id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.AddMusicToPlaylist(r.Context(), playlistID, req.MusicID, req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// DELETE /api/playlists/{id}/musics/{musicId}
func (h *Handler) RemoveMusicInPlaylist(w http.ResponseWriter, r *http.Request) {
	var req removeMusicInPlaylistRequest
	playlistID := r.PathValue("id")
	musicID := r.PathValue("musicId")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.RemoveMusicFromPlaylist(r.Context(), playlistID, musicID, req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(req)
}

// PUT /api/playlists/{id}
func (h *Handler) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req updatePlaylistRequest
	playlistID := r.PathValue("id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playlist, err := h.service.UpdatePlaylist(r.Context(), playlistID, req.Name, req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(playlist)
}

// DELETE /api/playlists/{id}
func (h *Handler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	var req deletePlaylistRequest
	playlistID := r.PathValue("id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.DeletePlaylist(r.Context(), playlistID, req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(req)
}

// GET /api/playlists/{id}/musics
func (h *Handler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
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
	playlist, err := h.service.GetPlaylist(r.Context(), playlistID, req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(playlist)
}

// GET /api/playlists
func (h *Handler) GetPlaylists(w http.ResponseWriter, r *http.Request) {
	var req getPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playlists, err := h.service.GetPlaylists(r.Context(), req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(playlists)
}
