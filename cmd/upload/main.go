package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	pb "github.com/beriloqueiroz/music-stream/api/proto"
	"github.com/dhowden/tag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	uploadPtr := flag.Bool("upload", false, "Fazer upload de música")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Uso: go run main.go [-upload] <music_id|filepath>")
	}

	if *uploadPtr {
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Falha ao conectar: %v", err)
		}
		defer conn.Close()

		client := pb.NewMusicServiceClient(conn)
		// upload all musics inside a folder
		err = uploadMusics(client, flag.Args()[0])
		if err != nil {
			log.Fatalf("Erro ao fazer upload das músicas: %v", err)
		}
	}
}
func uploadMusics(client pb.MusicServiceClient, folderPath string) error {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("erro ao ler o diretório: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			// Recursively process subfolders
			err = uploadMusics(client, filepath.Join(folderPath, file.Name()))
			if err != nil {
				return fmt.Errorf("erro ao processar subdiretório: %v", err)
			}
			continue
		}

		if (filepath.Ext(file.Name()) != ".mp3") && (filepath.Ext(file.Name()) != ".flac") {
			continue
		}

		err = uploadMusic(client, filepath.Join(folderPath, file.Name()))
		if err != nil {
			return fmt.Errorf("erro ao fazer upload da música: %v", err)
		}
	}
	return nil
}

type UploadOptions struct {
	Title  string
	Artist string
	Album  string
}

func uploadMusic(client pb.MusicServiceClient, filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo: %v", err)
	}
	defer file.Close()

	// Extrair metadados do MP3
	metadata, err := tag.ReadFrom(file)
	if err != nil {
		log.Printf("Aviso: não foi possível ler metadados: %v", err)
	}

	// Voltar ao início do arquivo após ler os metadados
	file.Seek(0, 0)

	stream, err := client.UploadMusic(context.Background())
	if err != nil {
		return fmt.Errorf("erro ao iniciar upload: %v", err)
	}

	// Usar valores fornecidos ou extrair dos metadados
	fileName := filepath.Base(filePath)
	title := ""
	artist := ""
	album := ""
	typeFile := metadata.FileType()

	// Se não foram fornecidos, tentar extrair do arquivo
	if title == "" || artist == "" || album == "" {
		if metadata != nil {
			if title == "" && metadata.Title() != "" {
				title = metadata.Title()
			}
			if artist == "" && metadata.Artist() != "" {
				artist = metadata.Artist()
			}
			if album == "" && metadata.Album() != "" {
				album = metadata.Album()
			}
		}
	}

	// Usar valores padrão se ainda estiverem vazios
	if title == "" {
		title = fileName
	}
	if artist == "" {
		artist = "Desconhecido"
	}
	if album == "" {
		album = "Desconhecido"
	}

	err = stream.Send(&pb.UploadRequest{
		Data: &pb.UploadRequest_Metadata{
			Metadata: &pb.MusicMetadata{
				Title:    title,
				Artist:   artist,
				Album:    album,
				Type:     string(typeFile),
				Comments: metadata.Comment(),
				AlbumArt: func() []byte {
					if metadata.Picture() != nil {
						return metadata.Picture().Data
					}
					return nil
				}(),
				AlbumArtType: metadata.Picture().Ext,
				Genre:        metadata.Genre(),
				Composer:     metadata.Composer(),
				Year:         int32(metadata.Year()),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("erro ao enviar metadata: %v", err)
	}

	// Enviar chunks do arquivo
	buffer := make([]byte, 1024*32)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erro ao ler arquivo: %v", err)
		}

		err = stream.Send(&pb.UploadRequest{
			Data: &pb.UploadRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		})
		if err != nil {
			return fmt.Errorf("erro ao enviar chunk: %v", err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("erro ao finalizar upload: %v", err)
	}

	fmt.Printf("Upload concluído! ID: %s\n", response.MusicId)
	return nil
}

func searchMusic(client pb.MusicServiceClient, query string) error {
	response, err := client.SearchMusic(context.Background(), &pb.SearchRequest{
		Query:    query,
		Page:     0,
		PageSize: 10,
	})
	if err != nil {
		return fmt.Errorf("erro ao buscar música: %v", err)
	}

	fmt.Printf("Resultados encontrados: %d\n", response.Total)
	for _, music := range response.MusicList {
		fmt.Printf("ID: %s, Título: %s, Artista: %s, Álbum: %s\n", music.Id, music.Title, music.Artist, music.Album)
	}
	return nil
}
