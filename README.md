# TgBotUltimate

Короткий обзор структуры проекта. Это Go-приложение с несколькими режимами запуска:

- Telegram-бот для подбора квартир.
- HTTP API для health-check и ручного запуска синхронизации.
- PostgreSQL-слой с автосозданием таблиц.
- Python NER-сервис для извлечения параметров из текстовых запросов.
- Cron-задачи для синхронизации данных из внешних систем.

Основной набор процессов выбирается через переменную `STATEMENTS` в `main.go`.

## Корень проекта

- `main.go` — главный entrypoint; читает `.env`, поднимает выбранные процессы (`neuro`, `platform`, `server`, `database`, `cron`) и останавливает их по сигналам.
- `.env` — локальная конфигурация окружения.

## `database/`

Пакет работы с PostgreSQL: подключение, миграции таблиц и CRUD по основным сущностям.

- `database/main.go` — открывает пул `pgx`, создаёт таблицы и индексы при старте.

### `database/data/`

Операции с бизнес-данными каталога недвижимости.

- `projects.go` — создание, обновление и чтение проектов.
- `buildings.go` — работа со зданиями.
- `sections.go` — работа с секциями/подъездами.
- `flats.go` — работа с квартирами.
- `tags.go` — работа с тегами квартир.

### `database/users/`

Данные пользователя Telegram и его состояния.

- `users.go` — создание и обновление пользователей, хранение ручных параметров и служебных полей.
- `expert_answers.go` — ответы пользователя в экспертной системе.

### `database/messages/`

- `messages.go` — история сообщений и извлечённых параметров пользователя.

### `database/favorites/`

- `favorites.go` — избранные квартиры пользователя.

### `database/expert/`

- `questions.go` — загрузка и работа с вопросами экспертной системы.

### `database/queries/`

SQL-слой проекта.

- `queries.go` — DDL таблиц, индексы, базовые SQL-шаблоны и общий запрос на выборку квартир.
- `consts.go` — SQL-константы и связанные вспомогательные значения.

### `database/queries/helper/`

Небольшие утилиты для сборки SQL и нормализации значений.

- `createQueryForSearchFlats.go` — сборка SQL для поиска квартир по параметрам.
- `convertValuesToSQLCreate.go` — подготовка данных для `INSERT`.
- `convertValuesToSQLUpdate.go` — подготовка данных для `UPDATE`.
- `safeNil.go` — безопасная работа с nullable-значениями.

## `platform/`

Telegram-платформа: подключение к боту, команды, callback-и и сценарии общения.

- `platform/main.go` — поднимает Telegram-платформу и подключает БД.

### `platform/telegram/`

- `main.go` — long polling Telegram, регистрация пользователя, роутинг `/start`, `/help`, `/questions`, `/reload`, `/flats`, `/favorites` и обычных текстовых сообщений.

### `platform/actions/`

Основные пользовательские сценарии.

- `commands.go` — регистрирует список команд бота.
- `callback.go` — роутер callback-запросов от inline-кнопок.
- `start.go` — приветственный сценарий.
- `help.go` — справка по работе бота.
- `selection.go` — подбор квартир по тексту пользователя и показ результатов.
- `expert_system.go` — сценарий экспертной системы с вопросами и вариантами ответов.
- `manual_parameters.go` — ручная настройка параметров поиска через сообщения и callback-и.

### `platform/helper/`

- `base64.go` — вспомогательная работа с base64-данными.

## `processing/`

Промежуточная обработка данных между ботом, БД и Python NER.

- `summarize.go` — собирает итоговые параметры пользователя из истории сообщений и полей профиля, строит SQL-фильтр и форматирует карточку квартиры.

### `processing/neuro/`

Интеграция с Python NER API.

- `init.go` — запускает `training/api.py` как отдельный процесс и завершает его вместе с приложением.
- `exec.go` — отправляет текст запроса в Python API `/parse` и получает извлечённые параметры.

## `server/`

HTTP-слой поверх приложения.

- `server/main.go` — запускает HTTP-сервер на порту из `PORT`.

### `server/routes/`

- `router.go` — корневой роутер `chi`, middleware и mount `/api`.

### `server/routes/handler/`

- `router.go` — HTTP-хендлеры: `GET /api/health` и `POST /api/sync`.

### `server/routes/helper/`

- `fetch.go` — общие HTTP-клиенты/обёртки для внешних запросов.

### `server/routes/external/core/`

Синхронизация данных из внешних источников.

- `feed.go` — импорт данных недвижимости из 1C-подобного feed API.
- `strapi.go` — импорт данных из Strapi.

### `server/routes/external/helper/`

- `converter.go` — преобразование внешних DTO в внутренние типы проекта.

### `server/cron-tasks/`

- `main.go` — планировщик `gocron`; по расписанию запускает синхронизацию `Strapi`.

## `training/`

Python-часть проекта: обучение, экспорт и запуск NER-модели для разбора запросов о недвижимости.

- `realty_query_ner.py` — основная логика NER: labels, подготовка датасета, обучение `transformers`, инференс, экспорт в ONNX и CLI.
- `generate_samples.py` — генератор синтетических обучающих примеров.
- `train_data.json` — большой размеченный датасет для обучения.
- `api.py` — FastAPI-сервис, который загружает модель и отдаёт `/parse` и `/health`.

### `training/ckpt/`

Сохранённая PyTorch-версия обученной модели.

- `config.json` — конфиг архитектуры модели.
- `model.safetensors` — веса модели.
- `tokenizer.json`, `tokenizer_config.json`, `vocab.txt`, `special_tokens_map.json` — артефакты токенизатора.
- `labels.json` — список меток NER.
- `training_args.bin` — параметры обучения.
- `checkpoint-13500/`, `checkpoint-27000/` — промежуточные чекпоинты тренировки.

### `training/onnx/`

Экспорт модели для быстрого CPU-инференса.

- `model.onnx` — ONNX-модель.
- `config.json`, `tokenizer.json`, `tokenizer_config.json`, `vocab.txt`, `special_tokens_map.json` — конфиг и токенизатор для ONNX-режима.

### `training/.venv/`

Локальное Python-окружение для NER-сервиса и обучения.

## `types/`

Общие типы данных, которыми обмениваются пакеты проекта.

### `types/Action/`

- `Action.go` — единый объект контекста действия Telegram-бота: `context`, update, БД и bot instance.

### `types/Database/`

Типы для БД и доменной модели.

- `DB.go` — тип подключения к базе и связанные структуры.
- `User.go` — пользователь Telegram и его параметры.
- `Message.go` — сохранённое сообщение пользователя.
- `Query.go` — структура результата SQL-запроса по квартирам.
- `Favorite.go` — избранная квартира.
- `ExpertSystemAnswer.go` — ответ пользователя в экспертной системе.
- `IDatabase.go` — интерфейсы/контракты слоя базы данных.

### `types/Expert/`

- `Question.go` — структура вопроса экспертной системы.

### `types/Neuro/`

- `Request.go` — запрос к Python NER API.
- `Response.go` — ответ NER API; умеет принимать строки и числа и приводить `<UNK>` к пустому значению.

### `types/Sync/`

Типы для внешних синхронизаций.

- `Data.go` — общие структуры синхронизации.

### `types/Sync/Sync1C/`

- `Data.go` — полная структура ответа 1C/feed API.
- `Project.go` — структуры проекта и вложенных сущностей из 1C.

### `types/Sync/SyncStrapi/`

- `Strapi.go` — структуры ответа Strapi API.

## `errors/`

- `fetch.go` — текстовые константы ошибок для HTTP/fetch-логики.

