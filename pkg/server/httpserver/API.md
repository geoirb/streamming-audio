# Server API

Запустить воспроизведение файла
--
* URI: 
```  
/player/file
```
* Метод: 
```  
POST
```
* Тело запроса:
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

* Тело ответа:
```json
{
	"uuid": "string",
	"channels": uint32,
	"rate": uint32,
	"bitsPerSample": uint32
}
```
> uuid - uuid хранилища в которое будет сохраняться аудио до воспроизведения
> 
> channels - количество аудиоканалов (из аудио файла)
> 
> rate - частота дискретизации (из аудио файла)
>
> bitsPerSample - количество бит на семпл

* Описание:

Сервер на порт `playerPort` плеера `playerIP` начинает передавать аудио данные из файла `file`. Плеер сохранет аудио данные в хранилище `uuid` и, постепенно вычитывая из хранилища, воспроизводит на аудиоустройстве `playerDeviceName`

Приостановить воспроизведение файла
---
* URI:
```
/player/file/stop
```
* Метод:
```
POST
```
* Тело запроса:
```json
{
	"playerIP": "string",
	"playerPort": "string",
	"playerDeviceName": "string",
	"uuid": "string"
}
```
> playerIP - ip плеера, на котором будет воспроизводиться файл
> 
> playerPort - порт плеера, на который сервер будет отсылать аудио сигнал
> 
> playerDeviceName - устройство, на котором будет идти воспроизведение
> 
> uuid - хранилище, из которого будет идти воспроизведение

* Описание:

Сервер перестает передавать аудио данные на порт `playerPort` плеера `playerIP`. Плеер останавливает воспроизведение на аудиоустройстве `playerDeviceName` и очищает хранилище `uuid`

Получить состояние плеера
---
* URI:
```
/player/state
```
* Метод:
```
GET
```
* Тело запроса:
```json
{
	"playerIP": "string",
}
```
> playerIP - ip плеера

* Тело ответа:
```json
{
	"ports": "string", 
	"storages": "string", 
	"devices": "string"
}
```
> ports - порты на которых идет прием данных
> 
> storage - все uuid хранилищ имеющихся на плеере
> 
> devices - все занятые звуковые устройства 

* Описание:

Плеер `playerIP` возращает `ports` на которые он принимает данные, `storage` - uuid хранилищ на плеере и `devices` - все занятые звуковые устройства 

Начать прием данных на плеере
---
* URI:
```
/player/receive/start
```
* Метод:
```
POST
```
* Тело запроса:
```json
{
	"playerIP": "string",
	"playerPort": "string",
	"uuid": "string"
}
```
> playerIP - ip плеера
> 
> playerPort - порт плеера, на который сервер будет отсылать аудио сигнал
> 
> uuid - хранилище, куда будет сохраняться данные, необязательное поле

* Тело ответа:
```json
{
	"uuid": "string"
}
```
> uuid - хранилище, куда будет сохраняться данные

* Описание:

Плеер `playerIP` начинает прием данных на порте `playerPort` и сохраняет их в хранилище `uuid`

Завершить прием данных на плеере
---
* URI:
```
/player/receive/stop
```
* Метод:
```
POST
```
* Тело запроса:
```json
{
	"playerIP": "string",
	"playerPort": "string"
}
```
> playerIP - ip плеера
> 
> playerPort - порт плеера, на который сервер передает аудио сигнал

* Описание:
  
Плеер `playerIP` прекращает прием данных на порте `playerPort`

Запуск воспроизведения на плеере
---
* URI:
```
/player/play
```
* Метод:
```
POST
```
* Тело запроса:
```json
{
	"playerIP":"string",
	"uuid":"string",
	"playerDeviceName":"string",
	"channels":uint32,
	"rate":uint32,
	"bitsPerSample": uint32

}
```

> playerIP - ip плеера, на котором будет воспроизводиться файл
> 
> uuid - uuid хранилища, из которого будет воспроизведение
> 
> playerDeviceName - устройство на котором будет идти воспроизведение
> 
> channels - количество аудиоканалов 
> 
> rate - частота дискретизации 
> 
> bitsPerSample - количество бит на семпл

* Описание:
  
Плеер `playerIP` начинает воспроизводить аудио из хранилища `uuid` на `playerDeviceName` с количеством каналов `channels` и частотой дискретизации `rate`

Остановка воспроизведения на плеере
---
* URI:
```
/player/stop
```
* Метод:
```
POST
```
* Тело запроса:
```json
{
	"playerIP": "string",
	"playerDeviceName": "string"
}
```
> playerIP - ip плеера
> 
> playerDeviceName - устройство, на котором идет воспроизведение

* Описание:

Останавливает воспроизведение на устройстве `playerDeviceName` на плеере `playerIP`

Очистить хранилище на плеере
---
* URI:
```
/player/clearstorage
```
* Метод:
```
POST
```
* Тело запроса:
```json
{
	"playerIP": "string",
	"uuid": "string"
}
```
> playerIP - ip плеера
>
> uuid - хранилище

* Описание:
 
Очищает хранилище `uuid` на плеере `playerIP` 

Начать запись аудио в файл
---
* URI:
```
/recoder/file/start
```
* Метод:
```
POST
```
* Тело запроса:
```json
{
	"recorderIP": "string",
	"recorderDeviceName": "string",
	"channels": uint32,
	"rate": uint32,
	"receivePort": "string",
	"file": "string"
}
```
>recorderIP - ip рекордера
>
>recorderDeviceName - устройство записи
>
>channels - количество аудиопотоков
>
>rate - частота дискретизации 
>
>receivePort - порт сервера на который рекордер отправляет аудиосигнал 
>
>file - имя файла для записи

* Описание:

Начинает запись аудио с рекордера `recorderIP` в wav файл `file`

Остановить запись аудио в файл
---
* URI:
```
/recoder/file/stop
```
* Метод:
```
POST
```
* Тело запроса:
```json
{
	"recorderIP": "string",
	"recorderDeviceName": "string",
	"receivePort": "string"
}
```

>recorderIP - ip рекордера
>
>recorderDeviceName - устройство записи
>
>receivePort - порт сервера на который рекордер отправляет аудиосигнал 

* Описание:

Остановка записи аудио в файл
  
Начать передачу аудио с рекордера на плеер
---
* URI:
```
/recoder/player/play
```
* Метод:
```
POST
```
* Тело запроса:
```json
{
	"playerIP": "string",
	"playerPort": "string",
	"playerDeviceName": "string",
	"channels": uint32,
	"rate": uint32,
	"recorderIP": "string",
	"recorderDeviceName": "string"
}
```

>playerIP - ip плеера
> 
>playerPort - порт плеера, на который рекордер передает аудио сигнал
>
>playerDeviceName - устройство воспроизведение
>
>channels - количество аудиопотоков
>
>rate - частота дискретизации
> 
>recorderIP - ip рекордера
>
>recorderDeviceName - устройство записи

* Тело ответа:
```json
{
	"uuid": "string"
}
```
>uuid - хранилище на плеере с данными с рекордера

* Описание:

Рекордер `recorderIP` начинает получать аудио с устройства `recorderDeviceName` и отправляет его на плеер `playerIP` на порт `playerPort`. Плеер сохранет аудио данные в хранилище `uuid` и, постепенно вычитывая из хранилища, воспроизводит на аудиоустройстве `playerDeviceName`

Завершить передачу данных с рекордера на плеер
---
* URI:
```
/recoder/player/stop
```
* Метод:
```
POST
```
* Тело запроса:
```json
{
	"playerIP": "string",
	"playerPort": "string",
	"playerDeviceName": "string",
	"uuid": "string",
	"recorderIP": "string",
	"recorderDeviceName": "string"
}
```

>playerIP - ip плеера
> 
>playerPort - порт плеера, на который рекордер передает аудио сигнал
>
>playerDeviceName - устройство воспроизведение
>
>uuid - хранилище на плеере с данными с рекордера
> 
>recorderIP - ip рекордера
>
>recorderDeviceName - устройство записи

* Описание:

Останавливает получение аудио на рекордере `recorderIP` и передачу на плеер `playerIP`. Плеер прекращает воспроизведение аудио и очищает хранилище `uuid`

Получить состояние рекордера
---
* URI:
```
/recorder/state
```
* Метод:
```
GET
```
* Тело запроса:
```json
{
	"recorderIP": "string",
}
```
> recorderIP - ip рекордера

* Тело ответа:
```json
{
	"devices": "string"
}
```
> devices - все занятые записывающие устройства 

* Описание:

Рекордер `recorderIP` возращает `devices` - все занятые записывающие устройства 

Запустить рекордер
---
* URI:
```
/recoder/start
```
* Метод:
```
POST
```
* Тело запроса:
```json
{
	"recorderIP": "string",
	"recorderDeviceName": "string",
	"channels": uint32,
	"rate": uint32,
	"dstAddr": "string"
}
```
>recorderIP - ip рекордера
>
>recorderDeviceName - устройство записи
>
>channels - количество аудиопотоков
>
>rate - частота дискретизации
>
>dstAddr - адрес, на который необходимо отправлять аудио 

* Описание:

Запускает получение аудио с устройства `recorderDeviceName` на рекордере `recorderIP` и передает на адрес `dstAddr`

Остановить рекордер 
---
* URI:
```
/recoder/stop
```
* Метод:
```
POST
```
* Тело запроса:
```json
{
	"recorderIP": "string",
	"recorderDeviceName"`: "string"
}
```

>recorderIP - ip рекордера
>
>recorderDeviceName - устройство записи

* Описание:
  
Останавливает получение аудио с устройства `recorderDeviceName` на рекордере `recorderIP` и передачу 
