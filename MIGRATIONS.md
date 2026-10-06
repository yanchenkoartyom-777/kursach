# Database Migrations Guide

Этот проект использует [golang-migrate](https://github.com/golang-migrate/migrate) для управления миграциями базы данных.

## Установка golang-migrate

### Windows (через Scoop)
```powershell
scoop install migrate
```

### Windows (через Chocolatey)
```powershell
choco install golang-migrate
```

### Windows (вручную)
1. Скачайте бинарник с [releases page](https://github.com/golang-migrate/migrate/releases)
2. Распакуйте и добавьте в PATH

### Linux/Mac
```bash
# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Mac (через Homebrew)
brew install golang-migrate
```

## Проверка установки

```bash
migrate -version
```

## Основные команды

### 1. Применить все миграции (Migrate Up)

```bash
make migrate-up
```

Применяет все неприменённые миграции к базе данных.

**Когда использовать:**
- При первом запуске проекта
- После добавления новых миграций
- После клонирования репозитория

### 2. Откатить последнюю миграцию (Migrate Down)

```bash
make migrate-down
```

Откатывает последнюю применённую миграцию.

**Когда использовать:**
- Когда нужно откатить последнее изменение
- При тестировании миграций

### 3. Полный сброс базы данных (Migrate Reset)

```bash
make migrate-reset
```

Откатывает все миграции и заново применяет их.

**Когда использовать:**
- Когда нужно начать с чистой базы
- При изменении структуры миграций
- Для тестирования с начальными данными

### 4. Проверить версию миграции

```bash
make migrate-version
```

Показывает текущую версию миграции в базе данных.

### 5. Принудительно установить версию

```bash
make migrate-force V=6
```

Принудительно устанавливает версию миграции (используйте осторожно!).

**Когда использовать:**
- Когда миграция застряла в "dirty" состоянии
- Для исправления ошибок версионирования

### 6. Создать новую миграцию

```bash
make migrate-create NAME=add_user_avatar
```

Создаёт пару файлов миграции (up и down) с указанным именем.

### 7. Проверить статус базы данных

```bash
make db-status
```

Показывает:
- Текущую версию миграции
- Список таблиц
- Количество записей в каждой таблице

## Структура миграций

Миграции находятся в папке `migrations/` и имеют формат:

```
000001_create_users_table.up.sql    - Применение миграции
000001_create_users_table.down.sql  - Откат миграции
```

### Текущие миграции:

1. **000001_create_users_table** - Таблица пользователей
2. **000002_create_two_factor_codes_table** - Таблица 2FA кодов
3. **000003_create_cars_table** - Таблица автомобилей
4. **000004_create_orders_tables** - Таблицы заказов и элементов заказов
5. **000005_create_cart_items_table** - Таблица корзины
6. **000006_seed_data** - Начальные данные (пользователи и автомобили)

## Примеры использования

### Первый запуск проекта

```bash
# 1. Запустить Docker контейнеры
docker-compose up -d

# 2. Применить миграции
make migrate-up

# 3. Проверить статус
make db-status
```

### Сброс базы для тестирования

```bash
# Полный сброс
make migrate-reset
```

### Откат последней миграции

```bash
# Откатить последнюю миграцию
make migrate-down

# Проверить версию
make migrate-version
```

### Создание новой миграции

```bash
# Создать новую миграцию
make migrate-create NAME=add_car_features

# Это создаст два файла:
# migrations/000007_add_car_features.up.sql
# migrations/000007_add_car_features.down.sql
```

## Тестовые данные

После применения миграций будут созданы:

**Пользователи:**
- Админ: `formed86@gmail.com` / `123456`
- Клиент: `baloban.oleg2016@yandex.ru` / `123456`

**Автомобили:**
- 10 автомобилей с изображениями
- Различные марки: Toyota, Honda, BMW, Mercedes-Benz, Audi, Tesla, Ford, Volkswagen, Mazda, Hyundai

## Устранение проблем

### Ошибка "Dirty database version"

Если миграция не удалась и база данных в "dirty" состоянии:

```bash
# Проверить текущую версию
make migrate-version

# Принудительно установить версию (замените N на номер последней успешной миграции)
make migrate-force V=N

# Попробовать снова
make migrate-up
```

### Ошибка подключения к базе данных

```bash
# Проверить, что контейнер запущен
docker ps | findstr car_dealership_db

# Если не запущен, запустить
docker-compose up -d

# Проверить логи
docker logs car_dealership_db
```

### Миграция не применяется

```bash
# Проверить текущую версию
make migrate-version

# Проверить статус базы
make db-status

# Если нужно, сбросить базу
make migrate-reset
```

## Прямое использование migrate CLI

Если нужно использовать migrate напрямую:

```bash
# Применить миграции
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/car_dealership?sslmode=disable" up

# Откатить миграции
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/car_dealership?sslmode=disable" down

# Проверить версию
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/car_dealership?sslmode=disable" version
```

## Лучшие практики

1. **Всегда создавайте down-миграции** - это позволяет откатывать изменения
2. **Тестируйте миграции** - проверяйте up и down перед коммитом
3. **Не изменяйте применённые миграции** - создавайте новые вместо этого
4. **Используйте транзакции** - добавляйте `BEGIN;` и `COMMIT;` в миграции
5. **Делайте миграции атомарными** - одна миграция = одно изменение

## Дополнительная информация

- [Документация golang-migrate](https://github.com/golang-migrate/migrate)
- [Примеры миграций](https://github.com/golang-migrate/migrate/tree/master/database/postgres)
