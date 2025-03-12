package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/beriloqueiroz/music-stream/internal/helper"
	"github.com/beriloqueiroz/music-stream/internal/infra/mongodb"
	rest_server "github.com/beriloqueiroz/music-stream/internal/infra/rest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// auth rest api integration test
// playlist rest api integration test

// to run : go test  ./test/integration/rest_test.go

func TestRestIntegration(t *testing.T) {
	port := "8099"
	if testing.Short() {
		t.Skip("Pulando testes de integração")
	}

	ctx := context.Background()

	testHelper := NewTestHelper()

	container, err := testHelper.StartMongoDBContainer(ctx)
	if err != nil {
		log.Fatalf("Erro ao iniciar container: %v", err)
		os.Exit(1)
	}
	defer container.Terminate(ctx)

	database, err := testHelper.ConnectToMongoDB(ctx, container)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
		os.Exit(1)
	}

	defer database.Disconnect(ctx)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "babau123"
	}

	restServer := rest_server.NewRestServer(database.Database("music-stream"), port)
	userRepo := mongodb.NewMongoUserRepository(database.Database("music-stream"))
	playlistRepo := mongodb.NewMongoPlaylistRepository(database.Database("music-stream"))
	go restServer.Start(jwtSecret, userRepo, playlistRepo)

	time.Sleep(1 * time.Second)

	//create admin user in database
	createAdminUser(database.Database("music-stream"))

	// integration tests
	token := ""
	userID := ""
	t.Run("TestLogin", func(t *testing.T) {
		// test login
		url := "http://localhost:" + port + "/api/auth/login"
		payload := map[string]interface{}{
			"email":    "admin@teste.com",
			"password": "12365478",
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, response["token"])
		token = response["token"].(string)
		userID = response["id"].(string)
	})

	inviteCode := ""

	t.Run("TestCreateInvite", func(t *testing.T) {
		// test create invite
		url := "http://localhost:" + port + "/api/invites"
		payload := map[string]interface{}{
			"email": "testuser@teste.com",
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, response["code"])
		inviteCode = response["code"].(string)
	})

	t.Run("TestCreateUserWithInviteCode", func(t *testing.T) {
		// test create user
		url := "http://localhost:" + port + "/api/auth/register"
		payload := map[string]interface{}{
			"email":       "testuser@teste.com",
			"password":    "testpassword",
			"invite_code": inviteCode,
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("TestCreateUserWithoutInviteCode", func(t *testing.T) {
		// test create user
		url := "http://localhost:" + port + "/api/auth/register"
		payload := map[string]interface{}{
			"email":       "testuser2@teste.com",
			"password":    "testpassword2",
			"invite_code": "invalidcode",
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("LoginWithNewUserInvited", func(t *testing.T) {
		// test login with new user invited
		url := "http://localhost:" + port + "/api/auth/login"
		payload := map[string]interface{}{
			"email":    "testuser@teste.com",
			"password": "testpassword",
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, response["token"])
		assert.NotNil(t, response["id"])
		userID = response["id"].(string)
		token = response["token"].(string)
	})

	playlistID := ""
	ownerID := userID
	// playlist rest api integration test
	t.Run("TestCreatePlaylist", func(t *testing.T) {
		// test create playlist
		url := "http://localhost:" + port + "/api/playlists"
		payload := map[string]interface{}{
			"name": "testplaylist",
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, response["id"])
		assert.Equal(t, "testplaylist", response["name"])
		assert.NotNil(t, response["created_at"])
		assert.NotNil(t, response["updated_at"])
		assert.Empty(t, response["musics"])
		assert.Equal(t, ownerID, response["owner_id"])
		playlistID = response["id"].(string)
	})

	t.Run("TestAddMusicInPlaylist", func(t *testing.T) {
		// test add music in playlist
		url := "http://localhost:" + port + "/api/playlists/" + playlistID + "/musics"
		payload := map[string]interface{}{
			"music_id": "66d6d6d6d6d6d6d6d6d6d6d6",
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("TestGetPlaylistAfterAddMusic", func(t *testing.T) {
		// test get playlist after add music
		url := "http://localhost:" + port + "/api/playlists/" + playlistID + "/musics"
		payload := map[string]interface{}{}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "testplaylist", response["name"])
		assert.NotNil(t, response["created_at"])
		assert.NotNil(t, response["updated_at"])
		assert.Equal(t, ownerID, response["owner_id"])
		assert.Equal(t, 1, len(response["musics"].([]interface{})))
		assert.Equal(t, "66d6d6d6d6d6d6d6d6d6d6d6", response["musics"].([]interface{})[0].(map[string]interface{})["music_id"])
	})

	t.Run("TestRemoveMusicInPlaylist", func(t *testing.T) {
		// test remove music in playlist
		url := "http://localhost:" + port + "/api/playlists/" + playlistID + "/musics/66d6d6d6d6d6d6d6d6d6d6d6"
		payload := map[string]interface{}{
			"owner_id": ownerID,
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("TestGetPlaylistAfterRemoveMusic", func(t *testing.T) {
		// test get playlist after remove music
		url := "http://localhost:" + port + "/api/playlists/" + playlistID + "/musics"
		payload := map[string]interface{}{
			"owner_id": ownerID,
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "testplaylist", response["name"])
		assert.NotNil(t, response["created_at"])
		assert.NotNil(t, response["updated_at"])
		assert.Equal(t, ownerID, response["owner_id"])
		assert.Equal(t, 0, len(response["musics"].([]interface{})))
	})

	t.Run("TestUpdatePlaylist", func(t *testing.T) {
		// test update playlist
		url := "http://localhost:" + port + "/api/playlists/" + playlistID
		payload := map[string]interface{}{
			"name": "testplaylist2",
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "testplaylist2", response["name"])
	})

	t.Run("TestGetPlaylists", func(t *testing.T) {
		// test get playlists
		url := "http://localhost:" + port + "/api/playlists"
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		var responseArray []map[string]interface{}
		err = json.Unmarshal(body, &responseArray)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, len(responseArray))
		assert.Equal(t, playlistID, responseArray[0]["id"])
	})

	t.Run("TestDeletePlaylist", func(t *testing.T) {
		// test delete playlist
		url := "http://localhost:" + port + "/api/playlists/" + playlistID
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

}

func createAdminUser(db *mongo.Database) {
	hash, err := helper.GenerateHash("12365478")
	if err != nil {
		log.Fatal(err)
	}
	db.Collection("users").InsertOne(context.Background(), bson.M{
		"id":         uuid.New(),
		"email":      "admin@teste.com",
		"password":   hash,
		"is_admin":   true,
		"created_at": time.Now(),
	})
}
