# Test play file


1. Build and start players on different machines
   * `make build-player tag=player`
   * `docker run --rm -p 8080:8080 -p 8081:8081 --device /dev/snd --name player player`
   * or run player `go run cmd/player/main.go`
2. Put on dir `audio/` audio files
3. Configurate `example.go`:
   * input all players on in map `player`   
4. Start example server
   * `go run example/multi-player/server.go`