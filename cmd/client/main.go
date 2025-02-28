package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/beriloqueiroz/music-stream/api/proto"
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
	outputPtr := flag.String("output", "musica_baixada.mp3", "Nome do arquivo de saída para download")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Uso: go run main.go [-download] [-output filename] <music_id>")
	}
	musicID := flag.Args()[0]

	// Conectar ao servidor gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Falha ao conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewMusicServiceClient(conn)

	if *downloadPtr {
		downloadMusic(client, musicID, *outputPtr)
	} else {
		playMusic(client, musicID)
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
