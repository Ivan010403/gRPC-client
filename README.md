# Структура

Общая структура реализации может быть описана следующей схемой:

![цу drawio](https://github.com/Ivan010403/gRPC-client/assets/125370827/ddc6d063-610e-4f2e-ba54-2ee6a7bc4d18)

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


# Service layer (gRPC client)

1. Вся реализация клиента содержится в директории ```internal/app/grpcClient```. Для начала мы конструктором создаём новое подключение и клиент, используя это подключение:

```
conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials())...

proto.NewCloudClient(conn)
```

Метод NewCloudClient был сгенерирован утилитой protoc и используется для создания объекта типа proto.CloudClient, через который мы можем использовать функции из созданного контракта.

2. Наш созданный тип Client имеет 4 метода для реализации бизнес-логики

```
func (c *Client) UploadFile(ctx context.Context, data []byte, name, format_file string) (string, error)

func (c *Client) DeleteFile(ctx context.Context, name, format_file string) (string, error)

func (c *Client) GetFile(ctx context.Context, name, format_file string) ([]byte, error)

func (c *Client) GetFullData(ctx context.Context) ([]struct {
	Name          string
	Creation_date string
	Update_date   string
}, error)
```


