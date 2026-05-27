# ITStep Social 🚀

Добро пожаловать в репозиторий проекта **ITStep Social**! Это фуллстек-приложение (социальная платформа/сеть), состоящее из современного фронтенда на Nuxt 3 и производительного бэкенда на Go (Gin).

## 🛠 Технологический стек

### Frontend
- **Фреймворк:** [Nuxt 3](https://nuxt.com/) / Vue 3
- **Стилизация:** [Tailwind CSS](https://tailwindcss.com/)
- **Управление состоянием:** [Pinia](https://pinia.vuejs.org/)
- **Авторизация / BaaS:** Firebase

### Backend
- **Язык:** Go 1.25
- **Веб-фреймворк:** [Gin](https://gin-gonic.com/)
- **Базы данных и хранилище:** MongoDB, Google Cloud Firestore, Firebase, PostgreSQL (разворачивается через Docker)

### Инфраструктура
- Docker & Docker Compose (для локального развертывания БД)

---

## 📁 Структура проекта

- `/frontend` — клиентская часть приложения (Nuxt 3).
- `/backend` — серверная часть приложения (Go, Gin).
- `docker-compose.yml` — конфигурация для запуска инфраструктуры (PostgreSQL).
- `run.bat` — скрипт для быстрого запуска проекта.
- `Future_Roadmap.canvas` / `ITSTEP_Project_Map.canvas` — карты проекта и планы развития (Obsidian Canvas).

---

## 🚀 Как запустить проект локально

### Требования
- Установленный **Node.js** (для фронтенда)
- Установленный **Go** (для бэкенда)
- **Docker** и **Docker Compose** (для базы данных)

### Шаги для запуска

1. **Клонируйте репозиторий:**
   ```bash
   git clone <URL_вашего_репозитория>
   cd IT
   ```

2. **Запустите инфраструктуру (Базу данных):**
   ```bash
   docker-compose up -d
   ```

3. **Запуск Backend (Go):**
   Перейдите в папку бэкенда, скопируйте `.env.example` в `.env` (при необходимости укажите ключи) и запустите сервер:
   ```bash
   cd backend
   go run cmd/app/main.go
   ```

4. **Запуск Frontend (Nuxt 3):**
   В новом окне терминала перейдите в папку фронтенда, установите зависимости и запустите dev-сервер:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

5. **Альтернативный быстрый запуск:**
   Вы можете воспользоваться скриптом `run.bat` в корне проекта (для Windows):
   ```cmd
   .\run.bat
   ```

---

## 🔑 Конфигурация и переменные окружения

Для корректной работы приложения вам потребуется настроить переменные окружения.

- В папке `/backend` переименуйте `.env.example` в `.env` и добавьте ваши ключи (например, URL базы данных, ключи Firebase). Также убедитесь, что добавлен файл сервисного аккаунта Firebase (например, `serviceAccountKey.json`).
- В папке `/frontend` переименуйте `.env.example` в `.env` и укажите публичные ключи для Firebase и URL бэкенда.

---

## 📝 Планы и развитие (Roadmap)
Начальное развитие проекта, планирование архитектуры и новые фичи описаны в файлах. По этим Майндкартам пошагово реализовывали проект начиная с 1 мая 2026 года.`Future_Roadmap.canvas` и `ITSTEP_Project_Map.canvas` (открываются через [Obsidian](https://obsidian.md/)).

FIGMA:
https://www.figma.com/board/ngd32NaMB333zekpASws5y/FigJam-basics?node-id=0-1&t=awFDYOwX6yXbiAnE-0

