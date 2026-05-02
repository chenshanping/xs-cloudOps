@echo off
setlocal
cd /d "%~dp0\.."
set GOMAXPROCS=1
go test -run TestRolePermissionSmoke -count=1 ./tests
endlocal
