This repository is a microservice.
It interacts with the main service via the G-RPC protocol. The file with the extension .proto is located in the repository https://github.com/tolseone/protos

# Теория. Ход разработки сервиса.

## Теория
G-RPC:

бинарный формат передачи данных;

Protocol Buffers;

Кодогенерация;

### Бинарный формат передачи данных
Для передачи данных между сервисами мы используем JSON. 

И в случае REST API мы передаём JSON целиком. кодируем в байты и отправляем
размер исходного сообщения == размеру отправленного сообщения! 
Целиком: значит key + value в явном виде

В случае G-RPC мы можем, пару “Email”: “some@mail.ru” представить в следующем виде:

Email: 1 - заносим это в контракт

Далее, кодируем сообщение в бинарный формат и передаём. Сервис, который получает данные, декодирует информацию, смотря в контракт!

### Protocol Buffers

Это контракт, некое описание структуры сообщения, который мы передаём между сервисами.

### Кодогенерация
Созданная людьми утилита - protoc



## НАПИСАНИЕ КОДА.
https://selectel.ru/blog/tutorials/go-grcp/?utm_source=youtube.com&utm_medium=referral&utm_campaign=help_tgbot-grcp_181123_tuzov_paid
ОПИСАНИЕ КОНТРАКТА Protocol Buffers
Создаём и наполняем файл sso.proto

По нему генерируем файлы .go
Информация по установке:
используем официальную утилиту protoc (компилятор Protocol Buffers)
https://grpc.io/docs/languages/go/quickstart/ + плагины

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

Команда для генерации:

protoc -I proto proto/sso/sso.proto –go_out=./gen/go
–go_opt=paths=source_relative –go-grpc_out=./gen/go/ –go-grpc_pot=paths=source_relative

Создание папки, куда будут генерироваться файлы:

mkdir -p gen/go

Taskfile.yaml с автоматизированной генерацией кода
task generate
lib: https://taskfile.dev/installation/
Создание папки проекта
item
├── cmd.............. Команды для запуска приложения и утилит
│	├── migrator.... Утилита для миграций базы данных
│	└── item......... Основная точка входа в сервис item
├── config........... Конфигурационные yaml-файлы
├── internal......... Внутренности проекта
│	├── app.......... Код для запуска различных компонентов приложения
│	│	└── grpc.... Запуск gRPC-сервера
│	├── config....... Загрузка конфигурации
│	├── domain
│	│	└── models.. Структуры данных и модели домена
│	├── grpc
│	│	└── auth.... gRPC-хэндлеры сервиса Auth
│	├── lib.......... Общие вспомогательные утилиты и функции
│	├── services..... Сервисный слой (бизнес-логика)
│	│	├── auth
│	│	└── permissions
│	└── storage...... Слой хранения данных
│	└── postgresql.. Реализация на PostgreSQL
├── migrations....... Миграции для базы данных
├── storage.......... Файлы хранилища PostgreSQL базы данных

ДАЛЕЕ ПО ПОРЯДКУ:

## Прописать конфиг и проинициализировать в main.go

Для запуска конфига может пригодиться:

CONFIG_PATH=./path/to/config/file.yaml myApp

myApp --config=./path/to/config/file.yaml


## Логгер

setupLogger внутри main.go

путь auth-service/internal/lib/logger/handlers/slogdiscard|slogpretty для локальной разработки и тестов.

sl.go - вспомогательная функция, чтобы добавлять какую-либо ошибку в лог

## Написание gRPC-сервера: обработка запросов
go get github.com/tolseone/protos - получение нашего же пакета
Создание файла, в котором будем описывать обработчики запросов
internal/grpc/item/server.go

Все входящие запросы будет обрабатывать serverAPI struct
поле authv1.UnimplementedAuthServer служит для того, чтобы дёргать ручки, даже если не все методы интерфейса реализованы!

создание internal/app/grpc/app.go - приложение в которое мы будем оборачивать grpc-service. Служит чтобы разгрузить файл main.go
GracefulStop() - нужен при остановке приложения. Прекращает приём новых запросов, дожидается выполнение всех текущих соединений. И уже потом происходит выход из приложения

Валидация данных в файле auth-service/internal/grpc/auth/server.go 
Сейчас напишу просто через if’ы, потом импортирую пакет validator и с ним сделаю

Здесь же объявляю интерфейс Auth и имплементирую интерфейс внутри методов

## СЕРВИСНЫЙ СЛОЙ (MODEL)
internal/services/auth/auth.go - занимается бизнес логикой и будет взаимодействовать с БД (никто кроме него не будет взаимодействовать с бд)

Создаём папку DOMAIN, где будем хранить модели данных - к этому пакету можно обратиться из любого слоя, чтобы устранить явление закольцованности и следовать принципу Dependency Injection.

## Создание миграций
Миграции - это пошаговое изменение схемы данных, которое позволяет её модифицировать или обратно откатывать.

