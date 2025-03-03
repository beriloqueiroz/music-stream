package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	rest_server "github.com/beriloqueiroz/music-stream/cmd/server/rest"
	"github.com/beriloqueiroz/music-stream/internal/helper"
	"github.com/beriloqueiroz/music-stream/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// auth rest api integration test
// playlist rest api integration test

// to run : go test  ./test/integration/rest_test.go

func TestRestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando testes de integração")
	}

	ctx := context.Background()

	container, err := StartMongoDBContainer(ctx)
	if err != nil {
		log.Fatalf("Erro ao iniciar container: %v", err)
		os.Exit(1)
	}
	defer container.Terminate(ctx)

	database, err := ConnectToMongoDB(ctx, container)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
		os.Exit(1)
	}

	defer database.Disconnect(ctx)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "babau123"
	}

	restServer := rest_server.NewRestServer(database.Database("music-stream"))
	go restServer.Start(jwtSecret)

	time.Sleep(1 * time.Second)

	//create admin user in database
	createAdminUser(database.Database("music-stream"))

	// integration tests

	token := ""
	t.Run("TestLogin", func(t *testing.T) {
		// test login
		url := "http://localhost:8080/api/auth/login"
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

	})

	inviteCode := ""

	t.Run("TestCreateInvite", func(t *testing.T) {
		// test create invite
		url := "http://localhost:8080/api/invites"
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
		url := "http://localhost:8080/api/auth/register"
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
		url := "http://localhost:8080/api/auth/register"
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

	userID := ""

	t.Run("LoginWithNewUserInvited", func(t *testing.T) {
		// test login with new user invited
		url := "http://localhost:8080/api/auth/login"
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
	})

	playlistID := ""
	ownerID := userID
	// playlist rest api integration test
	t.Run("TestCreatePlaylist", func(t *testing.T) {
		// test create playlist
		url := "http://localhost:8080/api/playlists"
		payload := map[string]interface{}{
			"name":     "testplaylist",
			"owner_id": ownerID,
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
		url := "http://localhost:8080/api/playlists/" + playlistID + "/musics"
		payload := map[string]interface{}{
			"music_id": "66d6d6d6d6d6d6d6d6d6d6d6",
			"owner_id": ownerID,
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
		url := "http://localhost:8080/api/playlists/" + playlistID + "/musics"
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
		assert.Equal(t, 1, len(response["musics"].([]interface{})))
		assert.Equal(t, "66d6d6d6d6d6d6d6d6d6d6d6", response["musics"].([]interface{})[0].(map[string]interface{})["music_id"])
	})

	t.Run("TestRemoveMusicInPlaylist", func(t *testing.T) {
		// test remove music in playlist
		url := "http://localhost:8080/api/playlists/" + playlistID + "/musics/66d6d6d6d6d6d6d6d6d6d6d6"
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
		url := "http://localhost:8080/api/playlists/" + playlistID + "/musics"
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
		url := "http://localhost:8080/api/playlists/" + playlistID
		payload := map[string]interface{}{
			"name":     "testplaylist2",
			"owner_id": ownerID,
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
	})

}

func StartMongoDBContainer(ctx context.Context) (testcontainers.Container, error) {

	containerRequest := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27018:27018"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": "root",
			"MONGO_INITDB_DATABASE":      "music-stream",
		},
		Cmd: []string{"mongod", "--auth", "--port", "27018"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	log.Println("Container MongoDB iniciado com sucesso")

	return container, nil
}

func ConnectToMongoDB(ctx context.Context, container testcontainers.Container) (*mongo.Client, error) {
	mappedPort, err := container.MappedPort(ctx, "27018")

	if err != nil {
		log.Fatalf("Erro ao obter porta mapeada: %v", err)
		os.Exit(1)
	}

	dbURL := fmt.Sprintf("mongodb://root:root@localhost:%s/music-stream?authSource=admin", mappedPort.Port())

	db, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
		os.Exit(1)
	}

	err = db.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MongoDB conectado com sucesso")

	return db, nil
}

func createAdminUser(db *mongo.Database) {
	hash, err := helper.GenerateHash("12365478")
	if err != nil {
		log.Fatal(err)
	}
	db.Collection("users").InsertOne(context.Background(), bson.M{
		"email":      "admin@teste.com",
		"password":   hash,
		"is_admin":   true,
		"created_at": time.Now(),
	})
}

func PrintAllPlaylists(db *mongo.Database) {
	playlists, err := db.Collection("playlists").Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	for playlists.Next(context.Background()) {
		var playlist models.Playlist
		err = playlists.Decode(&playlist)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(playlist)
	}
}
