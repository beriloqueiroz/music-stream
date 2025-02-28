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
	PlaylistID string `json:"playlist_id"`
	MusicID    string `json:"music_id"`
	OwnerID    string `json:"owner_id"`
}

type removeMusicInPlaylistRequest struct {
	PlaylistID string `json:"playlist_id"`
	MusicID    string `json:"music_id"`
	OwnerID    string `json:"owner_id"`
}

type updatePlaylistRequest struct {
	PlaylistID string `json:"playlist_id"`
	Name       string `json:"name"`
	OwnerID    string `json:"owner_id"`
}

type deletePlaylistRequest struct {
	PlaylistID string `json:"playlist_id"`
	OwnerID    string `json:"owner_id"`
}

type getPlaylistRequest struct {
	PlaylistID string `json:"playlist_id"`
	OwnerID    string `json:"owner_id"`
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
	json.NewEncoder(w).Encode(playlist)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) AddMusicInPlaylist(w http.ResponseWriter, r *http.Request) {
	var req addMusicInPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.AddMusicToPlaylist(r.Context(), req.PlaylistID, req.MusicID, req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(req)
}

func (h *Handler) RemoveMusicInPlaylist(w http.ResponseWriter, r *http.Request) {
	var req removeMusicInPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.RemoveMusicFromPlaylist(r.Context(), req.PlaylistID, req.MusicID, req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(req)
}

func (h *Handler) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req updatePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playlist, err := h.service.UpdatePlaylist(r.Context(), req.PlaylistID, req.Name, req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(playlist)
}

func (h *Handler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	var req deletePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.DeletePlaylist(r.Context(), req.PlaylistID, req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(req)
}

func (h *Handler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	var req getPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playlist, err := h.service.GetPlaylist(r.Context(), req.PlaylistID, req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(playlist)
}
