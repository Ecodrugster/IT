@echo off
title ITSTEP Social Network Manager
cls

echo [1/3] Closing existing processes...
:: Завершаем процесс бэкенда (если он скомпилирован как main.exe)
taskkill /F /IM main.exe /T 2>nul
:: Завершаем процессы Node.js (фронтенд)
taskkill /F /IM node.exe /T 2>nul

echo [2/3] Starting Backend (Go)...
start "Backend Server" cmd /c "cd backend && go run cmd/app/main.go"

echo [3/3] Starting Frontend (Nuxt)...
start "Frontend Server" cmd /c "cd frontend && npm run dev"

echo.
echo ========================================
echo   ALL SYSTEMS RESTARTED SUCCESSFULLY!
echo ========================================
echo.
pause
