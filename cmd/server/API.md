# Server API

Воспроизведение файла
--
URI: `/player/file`

Метод: `POST`

Тело запроса:
```json
{
	"file": "string",
	"playerIP": "string",
	"playerPort": "string",
	"playerDeviceName": "string"
}
```
> file - полный путь до файла на сервере
> 
> playerIP - ip плеера, на котором будет воспроизводиться файл
> 
> playerPort - порт плеера, на который сервер будет отсылать аудио сигнал
> 
> playerDeviceName - устройство на котором будет идти воспроизведение


Тело ответа:
```json
{
	"uuid": "string",
	"channels": "string",
	"rate": "string"
}
```

> uuid - uuid хранилища в которое будет сохраняться аудио до воспроизведения
> 
> channels - количество аудиоканалов (из аудио файла)
> 
> rate - частота дискретизации (из аудио файла)