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
    1. "https://neumann-payment-api.herokuapp.com/payment", Method: POST - создает транзакцию, request body params: {"user_id": type int, "amount": type decimal, "user_email": type varchar, "currency": type varchar}

    2. "https://neumann-payment-api.herokuapp.com/payments/{id}/status", Method: PUT - обновляет статус транзакции по ее id (из за абстрактно написанного задания - перенес выбор статуса с сервиса на клиент (клиент отдает статусы - "УСПЕХ", "НЕУСПЕХ"), так же не понятно с применением к запросу "авторизации" - если речь идет о защите эндпоинта по авторизованному юзеру (брать из заголовка токен и сверять его через сервис авторизации по логике), то почему только одного???), request body params: {"status": enum valid_status}

    3. "https://neumann-payment-api.herokuapp.com/payments/{id}/status", Method: GET - возвращает статус транзакции по ее id

    4. "https://neumann-payment-api.herokuapp.com/payments/user/{id}", Method: GET - возвращает транзакции пользователя по его id

    5. "https://neumann-payment-api.herokuapp.com/payments/user?email=...", Method: GET - возвращает транзакции пользователя по его email

    5. "https://neumann-payment-api.herokuapp.com/payments/{id}", Method: PUT - отменяет транзакцию транзакцию по ее id
