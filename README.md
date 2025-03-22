# TODO list
## Api
 - [ ] Разное
   - [ ] Добавить createdt updatedt к мемам

 - [x] Сделать Repo под каждый тип объектов
   - [x] Create meme
   - [x] Get meme
   - [x] Update meme
   - [x] Delete meme
   - [ ] List memes
   - [x] Set media
   - [x] Get media
   - [x] Create board
   - [x] Get board
   - [x] Update board
   - [x] Delete board
   - [ ] Create user
   - [ ] Get user
   - [ ] Update user
   - [ ] Delete user

 - [ ] Сделать хендлеры
   - [x] POST "/memes" - создать мем (без медиа)
   - [x] PUT "/memes/{id}" - перезаписать мем
   - [x] GET "/memes/{id}" - получить мем
   - [x] PUT "/media/{id}" - записать медиа
   - [x] GET "/media/{id}" - получить медиа
   - [ ] POST "/board" - создать доску
   - [ ] GET "/search/{id}?subject=кот&text=якот" - поиск по доске ID, с запросом {"subject": "кот", "text": якот}

 - [ ] Тестирование
   - [x] POST "/memes"
   - [x] PUT "/memes/{id}"
   - [x] GET "/memes/{id}"
   - [x] PUT "/media/{id}"
   - [x] GET "/media/{id}"
   - [ ] POST "/board"
   - [ ] GET "/search/{id}?..."


## Bot
 - [ ] Сделать клиента который будет отсылать запросы на сервер
   - [ ] POST "/memes"
   - [ ] PUT "/memes/{id}"
   - [ ] GET "/memes/{id}"
   - [ ] PUT "/media/{id}"
   - [ ] GET "/media/{id}"
   - [ ] POST "/board"
   - [ ] GET "/search/{id}?..."

 - [ ] Сделать клиента телеграмм
   - [ ] Отправка фото
   - [ ] Отправка видео
   - [ ] Отправка GIF
   - [ ] Парсинг фото
   - [ ] Парсинг видео
   - [ ] Парсинг GIF
   - [ ] Парсинг inline запросов

 - [ ] Сделать парсинг команд
   - [ ] Создание мема
   - [ ] Редактирование мема
   - [ ] Просмотр мема
   - [ ] Поиск