@echo off
chcp 65001 >nul
REM TourHelper Frontend Service startup script
REM This script starts the frontend service with customizable service name, version, and environment

echo ====================================
echo TourHelper Frontend Service
echo ====================================
echo.

REM Prompt user for service name (default: tour_helper)
set /p SERVICE_NAME="Enter Service Name [tour_helper]: "
if "%SERVICE_NAME%"=="" set SERVICE_NAME=tour_helper

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

REM Switch to project root directory
cd /d "%~dp0.."

REM Run Go program with ldflags to set variables
set LDFLAGS=-X main.SERVICE_NAME=%SERVICE_NAME% -X main.SERVICE_VERSION=%SERVICE_VERSION% -X main.SERVICE_ENV=%SERVICE_ENV%
go run -ldflags "%LDFLAGS%" service\cmd\frontend\main.go

REM Pause if program exits with error
if errorlevel 1 (
    echo.
    echo [ERROR] Program execution failed, error code: %errorlevel%
    pause
)
