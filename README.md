## Эмулятор платежного сервиса

## Организация структуры кода
    1. https://github.com/golang-standards/project-layout
    2. https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
    3. https://habr.com/ru/post/181772/
    4. https://github.com/evrone/go-clean-template

## Что нужно для запуска проекта
    1. Docker
    2. Docker-compose (version by - support 3.9 verion file)

## Запуск проекта
    1. sudo make compose-up

## .env - показательный, в реальном проекте данный файл должен находится в .gitignore

## API
    1. "payments", Method: POST - создает платеж, request body params: {"user_id": type int, "amount": type decimal, "user_email": type varchar, "currency": type varchar}

    2. "payments/{id}/status", Method: PUT - обновляет статус платежа по ее id, request body params: {"status": enum valid_status}

    3. "payments/{id}/status", Method: GET - возвращает статус платежа по ее id

    4. "payments/user/{id}", Method: GET - возвращает платеж пользователя по его id

    5. "payments/user?email=...", Method: GET - возвращает платеж пользователя по его email

    5. "payments/{id}", Method: PUT - отменяет платеж по ее id
