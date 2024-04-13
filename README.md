# banner-service

# Старт сервера

`make start` или `docker-compose up -d`

# Пример запросов и ответов

## Регистрация

Пользователь

```
curl --location 'http://127.0.0.1:6050/register/user' \
--header 'Content-Type: application/json' \
--data '{
    "username": "user",
    "password": "password"
}'
```

Админ

```
curl --location 'http://127.0.0.1:6050/register/admin' \
--header 'Content-Type: application/json' \
--data '{
    "username": "user",
    "password": "password"
}'
```

## Получение токена

```
curl --location 'http://127.0.0.1:6050/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "user",
    "password": "password"
}'
```

Все запросы для работы с баннером соответвуют приглаемому [API](https://github.com/randnull/banner-service/blob/main/docs/api/api.yaml) (возможна отправка через Postman), например:

Создание баннера (в качестве токена указать тот, что получен при /login)

```
curl --location 'http://127.0.0.1:6050/banner' \
--header 'token: eyJhbGciOi5cCI6IkpXVCJ9.eyJ19I6dHJ1ZX0.xDifpQs4afRltnMJW9' \
--header 'Content-Type: application/json' \
--data '{
  "tag_ids": [1, 2, 3],
  "feature_id": 1,
  "content": {
    "title": "title",
    "text": "text",
    "url": "example.com"
  },
  "is_actitve": true
}'
```

Получение баннера (в качестве токена указать тот, что получен при /login)

```
curl --location 'http://127.0.0.1:6050/user_banner?tag_id=2&feature_id=1' \
--header 'token: eyJhbGciOi5cCI6IkpXVCJ9.eyJ19I6dHJ1ZX0.xDifpQs4afRltnMJW9'
```


Получение всех баннеров (в качестве токена указать тот, что получен при /login)

```
curl --location 'http://127.0.0.1:6050/banner?limit=1&feature_id=1' \
--header 'token: eyJhbGciOi5cCI6IkpXVCJ9.eyJ19I6dHJ1ZX0.xDifpQs4afRltnMJW9'
```

