## Финансовый сервис

### Описание

Финансовый сервис позволяет:

- Пополнять баланс пользователя
- Переводить деньги между пользователями
- Просматривать 10 последних операций пользователя

### Запуск проекта

```bash
git clone https://github.com/EugeneKrivoshein/fin_service.git
cd fin_service
make run
```

### API Документация (Swagger)

Для просмотра документации API можно открыть Swagger UI по адресу:

```
http://localhost:8080/swagger/index.html
```

### API Эндпоинты

- **POST /deposit** — пополнение баланса пользователя
- **POST /transfer** — перевод денег между пользователями
- **GET /transactions?user\_id=1** — просмотр 10 последних операций пользователя
