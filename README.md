# Структура

Общая структура реализации может быть описана следующей схемой:

![readme drawio](https://github.com/Ivan010403/gRPC-server/assets/125370827/09c722d6-5c48-465e-9725-ca7d010581c0)

1. Transport layer: содержит в себе реализацию http сервера для веб интерфейса, включая все handler'ы. 
2. Service layer: содержит в себе реализацию gRPC клиента.

Все слои собираются в общий тип App, который содержит в себе конструктор для инициализации всего приложения, и именно с ним мы взаимодействуем из нашего main. 
   
# HTTP сервер

Создание и инициализация сервера находится в ```internal/app/http-server```, реализация обработчиков в ```internal/http-server/handlers```. Список эндпоинтов:

| Энжпоинт | Метод | Обработчик и описание                                      |
|----------|-------|------------------------------------------------------------|
| /	       | GET   | Отображение index.html, запуск GetFullData                 |
| /upload  | POST  | Запуск client.UploadFile() c таймаутом 10 секунд (контекст)|
| /get     | POST  | Запуск client.GetFile() c таймаутом 10 секунд (контекст)   |
| /delete  | POST  | Запуск client.DeleteFile() c таймаутом 3 секунды (контекст)|


# Service layer

1. Вся реализация бизнес-логики содержится в директории ```internal/services```. Также создаётся дополнительный интерфейс, который мы пишем в месте использования, а именно в сервисном слое:

```
type FileWork interface {
	Write([]byte, string, string) error
	Update([]byte, string, string) error
	Delete(string, string) error
	Get(string, string) ([]byte, error)
	GetFullData() ([]postgres.File, error)
}
```
2. Тип ```Cloud``` из сервисного слоя реализует этот интерфейс, именно поэтому мы передаём его в наш транспортный слой.

