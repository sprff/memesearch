
## Meme Search

Бот для поиска мемов по описанию. Работает через inline режим

Для запуска используется `docker compose up`
Структура необходимого .env файла:
```
MS_APISERVER_YAS3_KEY=
MS_APISERVER_YAS3_SECRET=
MS_APISERVER_DB_USER=
MS_APISERVER_DB_PASS=
MS_APISERVER_DB_NAME=
MS_TGCLIENT_BOT_TOKEN=
MS_API_CONFIG_PATH=
```

### Кодогенерация

```
cd api-client/ && go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest && go generate ./... && cd ..
cd api-server/ && go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest && go generate ./... && cd ..
```