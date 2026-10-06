# Car Dealership System

Система управления автосалоном с функциями каталога автомобилей, корзины покупок, заказов и административной панели.

## Технологии

- **Backend**: Go (Gin framework)
- **Database**: PostgreSQL
- **Templates**: Go HTML templates
- **Authentication**: 2FA через email
- **Migrations**: golang-migrate

## Требования

- Go 1.21+
- Docker & Docker Compose
- golang-migrate CLI
- Make (опционально, для удобства)

## Быстрый старт

### 1. Клонировать репозиторий

```bash
git clone <repository-url>
cd car-dealership-system
```

### 2. Установить golang-migrate

**Windows (Scoop):**
```powershell
scoop install migrate
```

**Windows (Chocolatey):**
```powershell
choco install golang-migrate
```

**Linux/Mac:**
```bash
# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Mac
brew install golang-migrate
```

### 3. Запустить базу данных

```bash
docker-compose up -d
```

Это запустит:
- PostgreSQL на порту 5432
- pgAdmin на порту 5050 (http://localhost:5050)

### 4. Применить миграции

```bash
make migrate-up
```

Или напрямую:
```bash
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/car_dealership?sslmode=disable" up
```

### 5. Настроить конфигурацию

Отредактируйте `config.yaml` и укажите настройки email для 2FA:

```yaml
email:
  smtp_host: smtp.gmail.com
  smtp_port: 587
  username: your-email@gmail.com
  password: your-app-password  # Google App Password
  from: your-email@gmail.com
```

### 6. Запустить приложение

```bash
go run cmd/server/main.go
```

Приложение будет доступно по адресу: http://localhost:8080

## Тестовые данные

После применения миграций будут созданы тестовые пользователи:

- **Админ**: `formed86@gmail.com` / `123456`
- **Клиент**: `baloban.oleg2016@yandex.ru` / `123456`

И 10 автомобилей с изображениями.

## Управление миграциями

Подробная документация по миграциям: [MIGRATIONS.md](MIGRATIONS.md)

**Основные команды:**

```bash
make migrate-up      # Применить все миграции
make migrate-down    # Откатить последнюю миграцию
make migrate-reset   # Полный сброс базы данных
make migrate-version # Показать текущую версию
make db-status       # Показать статус базы данных
```

## Структура проекта

```
.
├── cmd/
│   └── server/          # Точка входа приложения
├── internal/
│   ├── config/          # Конфигурация
│   ├── handlers/        # HTTP обработчики
│   ├── middleware/      # Middleware
│   ├── models/          # Модели данных
│   ├── repositories/    # Слой доступа к данным
│   └── services/        # Бизнес-логика
├── migrations/          # Миграции базы данных
├── pkg/
│   ├── database/        # Утилиты для работы с БД
│   ├── renderer/        # Кастомный рендерер шаблонов
│   └── session/         # Управление сессиями
├── static/              # Статические файлы (CSS, JS, изображения)
├── templates/           # HTML шаблоны
├── config.yaml          # Конфигурация приложения
├── docker-compose.yml   # Docker конфигурация
└── Makefile            # Команды для управления проектом
```

## Функционал

### Для клиентов:
- ✅ Просмотр каталога автомобилей
- ✅ Поиск и фильтрация автомобилей
- ✅ Добавление в корзину
- ✅ Оформление заказов
- ✅ История заказов
- ✅ 2FA аутентификация через email

### Для администраторов:
- ✅ Управление автомобилями (CRUD)
- ✅ Управление заказами
- ✅ Аналитика и отчеты
- ✅ Экспорт/импорт данных (JSON, XML, CSV)
- ✅ Добавление изображений к автомобилям
- ✅ Dashboard с статистикой

## API Endpoints

### Публичные
- `GET /` - Главная страница
- `GET /cars` - Каталог автомобилей
- `GET /cars/:id` - Детальная информация об автомобиле
- `GET /login` - Страница входа
- `POST /login` - Вход в систему
- `GET /register` - Страница регистрации
- `POST /register` - Регистрация
- `GET /verify` - Страница 2FA
- `POST /verify` - Подтверждение 2FA кода

### Защищенные (требуется аутентификация)
- `GET /profile` - Профиль пользователя
- `GET /cart` - Корзина
- `POST /cart/add` - Добавить в корзину
- `POST /cart/remove` - Удалить из корзины
- `GET /orders/history` - История заказов
- `POST /orders` - Создать заказ

### Админ панель (требуется роль admin)
- `GET /admin/dashboard` - Dashboard
- `GET /admin/cars` - Управление автомобилями
- `GET /admin/orders` - Управление заказами
- `GET /admin/analytics` - Аналитика
- `POST /admin/cars/update-image` - Обновить изображение автомобиля

## Разработка

### Запуск в режиме разработки

```bash
# Запустить базу данных
docker-compose up -d

# Применить миграции
make migrate-up

# Запустить приложение с hot reload (если установлен air)
air

# Или обычный запуск
go run cmd/server/main.go
```

### Создание новой миграции

```bash
make migrate-create NAME=add_new_feature
```

Это создаст два файла:
- `migrations/000007_add_new_feature.up.sql`
- `migrations/000007_add_new_feature.down.sql`

### Сброс базы данных

```bash
make migrate-reset
```

## Конфигурация

Основные настройки в `config.yaml`:

```yaml
server:
  host: 0.0.0.0
  port: 8080
  templates_path: templates
  static_path: static

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: car_dealership
  sslmode: disable

session:
  secret: your-32-character-secret-key-here
  max_age: 86400  # 24 hours

email:
  smtp_host: smtp.gmail.com
  smtp_port: 587
  username: your-email@gmail.com
  password: your-app-password
  from: your-email@gmail.com

two_factor:
  code_length: 6
  code_expiration: 5m

logging:
  level: info
  format: text
```

## Устранение проблем

### База данных не подключается

```bash
# Проверить статус контейнера
docker ps

# Проверить логи
docker logs car_dealership_db

# Перезапустить контейнер
docker-compose restart
```

### Ошибки миграций

```bash
# Проверить текущую версию
make migrate-version

# Проверить статус базы
make db-status

# Сбросить базу
make migrate-reset
```

### Email не отправляется

1. Убедитесь, что используете Google App Password (не обычный пароль)
2. Включите "Less secure app access" в настройках Google аккаунта
3. Проверьте настройки SMTP в `config.yaml`

## Лицензия

MIT

## Авторы

- Ваше имя
