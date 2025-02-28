mkdir -p music-stream/{cmd,internal,pkg,api}
cd music-stream
go mod init github.com/beriloqueiroz/music-stream 
go mod tidy

sudo chmod -R 755 ./docker/mongo/mongodb_data/
docker compose up -d