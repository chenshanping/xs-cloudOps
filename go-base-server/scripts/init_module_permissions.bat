@echo off
REM 初始化模块权限脚本 (Windows版本)
REM 用法: init_module_permissions.bat <模块名> [数据库名] [用户名] [密码]

setlocal enabledelayedexpansion

set MODULE=%1
set DB_NAME=%2
set DB_USER=%3
set DB_PASS=%4

if "%MODULE%"=="" (
    echo 错误: 请提供模块名
    echo 用法: %0 ^<模块名^> [数据库名] [用户名] [密码]
    echo 示例: %0 productType ecobreed root yourpassword
    exit /b 1
)

if "%DB_NAME%"=="" set DB_NAME=ecobreed
if "%DB_USER%"=="" set DB_USER=root
if "%DB_PASS%"=="" (
    echo 请输入数据库密码:
    set /p DB_PASS=
)

set SQL_FILE=..\sql\%MODULE%_menu.sql

if not exist "%SQL_FILE%" (
    echo 错误: 找不到菜单SQL文件: %SQL_FILE%
    exit /b 1
)

echo =========================================
echo 初始化模块权限
echo =========================================
echo 模块名: %MODULE%
echo 数据库: %DB_NAME%
echo 用户: %DB_USER%
echo =========================================
echo.

REM 执行菜单SQL
echo 步骤1: 执行菜单SQL...
mysql -u %DB_USER% -p%DB_PASS% %DB_NAME% < "%SQL_FILE%"

if %errorlevel% neq 0 (
    echo X 菜单SQL执行失败
    exit /b 1
)
echo √ 菜单SQL执行成功
echo.

REM 分配权限给管理员角色（角色ID=1）
echo 步骤2: 分配权限给管理员角色...
echo INSERT INTO sys_role_menu (role_id, menu_id) SELECT 1, id FROM sys_menu WHERE permission LIKE '%MODULE%:%%' AND id NOT IN (SELECT menu_id FROM sys_role_menu WHERE role_id = 1) ON DUPLICATE KEY UPDATE role_id=role_id; | mysql -u %DB_USER% -p%DB_PASS% %DB_NAME%

if %errorlevel% neq 0 (
    echo X 权限分配失败
    exit /b 1
)
echo √ 权限分配成功
echo.

REM 验证权限
echo 步骤3: 验证权限...
for /f %%i in ('mysql -u %DB_USER% -p%DB_PASS% %DB_NAME% -se "SELECT COUNT(*) FROM sys_menu WHERE permission LIKE '%MODULE%:%%'"') do set PERM_COUNT=%%i
echo √ 找到 %PERM_COUNT% 个权限记录
echo.

echo =========================================
echo 初始化完成！
echo =========================================
echo.
echo 下一步操作：
echo 1. 如果需要给其他角色分配权限，请在系统中进入'角色管理'进行配置
echo 2. 用户需要重新登录才能看到新的权限按钮
echo.
echo 验证权限SQL:
echo SELECT id, name, permission FROM sys_menu WHERE permission LIKE '%MODULE%:%%';
echo.

endlocal
