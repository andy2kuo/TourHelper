@echo off
chcp 65001 >nul
REM TourHelper Backend Admin Service startup script
REM This script starts the backend admin service with customizable service name, version, and environment

REM Check if running from service directory
if not exist "go.mod" (
    echo [ERROR] This script must be run from the 'service' directory!
    echo Current directory: %CD%
    echo.
    echo Please navigate to the service directory and run this script again.
    echo Example: cd d:\GitHub\TourHelper\service
    pause
    exit /b 1
)

REM Verify go.mod contains correct module name
findstr /C:"module github.com/andy2kuo/TourHelper" go.mod >nul
if errorlevel 1 (
    echo [ERROR] Invalid go.mod file detected!
    echo This script must be run from the TourHelper service directory.
    pause
    exit /b 1
)

echo ====================================
echo TourHelper Backend Admin Service
echo ====================================
echo.

REM Prompt user for service name (default: backend_admin_server)
set /p SERVICE_NAME="Enter Service Name [backend_admin_server]: "
if "%SERVICE_NAME%"=="" set SERVICE_NAME=backend_admin_server

REM Prompt user for service version (default: 0.0.1-dev)
set /p SERVICE_VERSION="Enter Service Version [0.0.1-dev]: "
if "%SERVICE_VERSION%"=="" set SERVICE_VERSION=0.0.1-dev

REM Prompt user for service environment (default: dev)
set /p SERVICE_ENV="Enter Service Environment [dev]: "
if "%SERVICE_ENV%"=="" set SERVICE_ENV=dev

echo.
echo ====================================
echo Configuration:
echo Service Name:    %SERVICE_NAME%
echo Service Version: %SERVICE_VERSION%
echo Service Env:     %SERVICE_ENV%
echo ====================================
echo.

REM Switch to service directory (where go.mod is located)
cd /d "%~dp0"

REM Run Go program with ldflags to set variables
set LDFLAGS=-X main.SERVICE_NAME=%SERVICE_NAME% -X main.SERVICE_VERSION=%SERVICE_VERSION% -X main.SERVICE_ENV=%SERVICE_ENV%
go run -ldflags "%LDFLAGS%" cmd\backend_admin\main.go

REM Pause if program exits with error
if errorlevel 1 (
    echo.
    echo [ERROR] Program execution failed, error code: %errorlevel%
    pause
)
