package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	pb "github.com/beriloqueiroz/music-stream/api/proto"
	"github.com/dhowden/tag"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PlayerControls struct {
	done     chan bool
	pause    chan bool
	isPaused bool
}

func main() {

	// Flags para controlar o modo de operação
	downloadPtr := flag.Bool("download", false, "Baixar a música em vez de reproduzir")
	searchPtr := flag.Bool("search", false, "Buscar música")
	uploadPtr := flag.Bool("upload", false, "Fazer upload de música")
	outputPtr := flag.String("output", "musica_baixada.mp3", "Nome do arquivo de saída para download")

	// Novas flags para metadados
	titlePtr := flag.String("title", "", "Título da música")
	artistPtr := flag.String("artist", "", "Nome do artista")
	albumPtr := flag.String("album", "", "Nome do álbum")

	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Uso: go run main.go [-download|-upload] [-output filename] <music_id|filepath>")
	}

	// Conectar ao servidor gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Falha ao conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewMusicServiceClient(conn)

	if *searchPtr {
		err := searchMusic(client, flag.Args()[0])
		if err != nil {
			log.Fatal(err)
		}
	} else if *uploadPtr {
		err := uploadMusic(client, flag.Args()[0], &UploadOptions{
			Title:  *titlePtr,
			Artist: *artistPtr,
			Album:  *albumPtr,
		})
		if err != nil {
			log.Fatal(err)
		}
	} else if *downloadPtr {
		downloadMusic(client, flag.Args()[0], *outputPtr)
	} else {
		playMusic(client, flag.Args()[0])
	}
}

func downloadMusic(client pb.MusicServiceClient, musicID, outputFile string) {
	stream, err := client.StreamMusic(context.Background(), &pb.StreamRequest{
		MusicId: musicID,
	})
	if err != nil {
		log.Fatalf("Erro ao iniciar streaming: %v", err)
	}

	out, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo: %v", err)
	}
	defer out.Close()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Erro ao receber chunk: %v", err)
		}

		if _, err := out.Write(chunk.Data); err != nil {
			log.Fatalf("Erro ao escrever chunk: %v", err)
		}
	}

	log.Printf("Música baixada com sucesso em: %s", outputFile)
}

func playMusic(client pb.MusicServiceClient, musicID string) {
	// Configurar controles do player
	controls := &PlayerControls{
		done:     make(chan bool),
		pause:    make(chan bool),
		isPaused: false,
	}

	// Capturar sinais do terminal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar reprodução em uma goroutine
	go func() {
		tempFile, err := os.CreateTemp("", "stream*.mp3")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tempFile.Name())

		// Receber e salvar chunks
		stream, err := client.StreamMusic(context.Background(), &pb.StreamRequest{
			MusicId: musicID,
		})
		if err != nil {
			log.Fatal(err)
		}

		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if _, err := tempFile.Write(chunk.Data); err != nil {
				log.Fatal(err)
			}
		}

		// Preparar reprodução
		tempFile.Seek(0, 0)
		streamer, format, err := mp3.Decode(tempFile)
		if err != nil {
			log.Fatal(err)
		}
		defer streamer.Close()

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		done := make(chan bool)
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			done <- true
		})))

		// Loop de controle
		for {
			select {
			case <-controls.pause:
				controls.isPaused = !controls.isPaused
				if controls.isPaused {
					speaker.Lock()
				} else {
					speaker.Unlock()
				}
				fmt.Printf("\rStatus: %s", getStatus(controls.isPaused))
			case <-done:
				controls.done <- true
				return
			}
		}
	}()

	// Interface de controle
	fmt.Println("\nControles:")
	fmt.Println("p - Play/Pause")
	fmt.Println("q - Sair")
	fmt.Printf("\rStatus: Reproduzindo")

	for {
		select {
		case <-controls.done:
			return
		case sig := <-sigChan:
			fmt.Printf("\nRecebido sinal: %v\n", sig)
			return
		default:
			var input string
			fmt.Scanln(&input)
			switch input {
			case "p":
				controls.pause <- true
			case "q":
				return
			}
		}
	}
}

func getStatus(isPaused bool) string {
	if isPaused {
		return "Pausado"
	}
	return "Reproduzindo"
}

type UploadOptions struct {
	Title  string
	Artist string
	Album  string
}

func uploadMusic(client pb.MusicServiceClient, filePath string, opts *UploadOptions) error {
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
	title := opts.Title
	artist := opts.Artist
	album := opts.Album

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
				Title:  title,
				Artist: artist,
				Album:  album,
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
