# sound-server

Streaming audio

* Server - control server [API](pkg/server/httpserver/API.md)
* Player - client for playing audio signal on audio device 
* Recorder - client for receive audio signal from audio device 

## ToDoc        

### Server

- [X] Streaming audio signal on Player from:
  - [X] .wav file
  - [X] Recorder
- [X] RPC system control
  - [X] Player
  - [X] Recorder
- [X] Record .wav file
- [X] HTTP server 
- [ ] HTTP client
- [ ] Overlay 2 tracks
  
### Player
- [X] Receive audio signal
- [X] Playing audio signal
- [X] Selecting an audio card
- [X] Storage
- [X] RPC system control
- [ ] Volume control

### Recorder

- [X] Recording audio from microphone
- [X] Streaming audio signal
- [X] RPC system control

## Запуск server

1. Скачать проект на машину, на которой будет развернут server

        git clone git@github.com:GeoIrb/sound-server.git
2. Поместите аудиофайл, который необходимо будет стриммить в папку `audio/`

3. Собрать образ сервера

        make build-server tag=IMAGE-NAME
4. Запуск сервера

        docker run -d --rm \
        -p PORT:PORT \ 
        -e ENVIRONMENTS \ 
        IMAGE-NAME

**PORT** - порт, на который будет раздача (возможно это лишнее)

**ENVIRONMENTS** - переменные окружения

- FILE=/audio/`FILE`.wav - файл для стримминга
- DST_ADDRESS="IP:PORT" - на какой IP и на какой PORT будет рассылка, по умолчанию 255.255.255.255:8080 - рассылка по всей сети на порт 8080

        make build-server server
        docker run -d --rm -p 8081:8081 -p 8082:8082 -e FILE=/audio/test.wav server

## Запуск player

1. Скачать проект на машину, на которой будет развернут player

        git clone git@github.com:GeoIrb/sound-server.git
2. Собрать образ клиент

        make build-player tag=IMAGE-NAME
3. Запуск клиента

        docker run -d --rm \
        -p 0.0.0.0:8081:8081/tcp \ 
        -p 0.0.0.0:PORT:PORT -p 0.0.0.0:PORT:PORT/udp \
        --device /dev/snd \
        -e ENVIRONMENTS \
        IMAGE-NAME

**PORT** - порт, на котором будет работать клиент

**ENVIRONMENTS** - переменные окружения

- PORT - порт, на котором будет работать клиент
- PLAYBACK_DEVICE_NAME - устройство, на котором будет воспроизводиться принятый аудио сигнал

        make build-player tag player
        docker run -d --rm -p 0.0.0.0:8081:8081/tcp -p 0.0.0.0:8082:8082 -p 0.0.0.0:8082:8082/udp --device /dev/snd player
