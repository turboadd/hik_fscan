@echo off
setlocal enabledelayedexpansion

E:
cd \Projects\hik_fcan

echo Cleaning previous build...
if exist build (
    rmdir /s /q build
)

echo Building...
mkdir build
cd build

echo Configuring with CMake (MinGW Makefiles)
cmake -G "MinGW Makefiles" ..

echo Building with mingw2-make...
mingw32-make

echo Build complete.
pause

