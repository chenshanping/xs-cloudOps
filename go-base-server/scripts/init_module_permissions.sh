#!/bin/bash

# 初始化模块权限脚本
# 用法: ./init_module_permissions.sh <模块名> [数据库名] [用户名]

MODULE=$1
DB_NAME=${2:-"go_base"}
DB_USER=${3:-"root"}

if [ -z "$MODULE" ]; then
    echo "错误: 请提供模块名"
    echo "用法: $0 <模块名> [数据库名] [用户名]"
    echo "示例: $0 productType go_base root"
    exit 1
fi

SQL_FILE="../sql/${MODULE}_menu.sql"

if [ ! -f "$SQL_FILE" ]; then
    echo "错误: 找不到菜单SQL文件: $SQL_FILE"
    exit 1
fi

echo "========================================="
echo "初始化模块权限"
echo "========================================="
echo "模块名: $MODULE"
echo "数据库: $DB_NAME"
echo "用户: $DB_USER"
echo "========================================="
echo ""

# 执行菜单SQL
echo "步骤1: 执行菜单SQL..."
mysql -u "$DB_USER" -p "$DB_NAME" < "$SQL_FILE"

if [ $? -eq 0 ]; then
    echo "✓ 菜单SQL执行成功"
else
    echo "✗ 菜单SQL执行失败"
    exit 1
fi

echo ""

# 分配权限给管理员角色（角色ID=1）
echo "步骤2: 分配权限给管理员角色..."
mysql -u "$DB_USER" -p "$DB_NAME" <<EOF
INSERT INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu 
WHERE permission LIKE '${MODULE}:%'
  AND id NOT IN (SELECT menu_id FROM sys_role_menu WHERE role_id = 1)
ON DUPLICATE KEY UPDATE role_id=role_id;
EOF

if [ $? -eq 0 ]; then
    echo "✓ 权限分配成功"
else
    echo "✗ 权限分配失败"
    exit 1
fi

echo ""

# 验证权限
echo "步骤3: 验证权限..."
PERM_COUNT=$(mysql -u "$DB_USER" -p "$DB_NAME" -se "SELECT COUNT(*) FROM sys_menu WHERE permission LIKE '${MODULE}:%'")
echo "✓ 找到 $PERM_COUNT 个权限记录"

echo ""
echo "========================================="
echo "初始化完成！"
echo "========================================="
echo ""
echo "下一步操作："
echo "1. 如果需要给其他角色分配权限，请在系统中进入'角色管理'进行配置"
echo "2. 用户需要重新登录才能看到新的权限按钮"
echo ""
echo "验证权限SQL:"
echo "SELECT id, name, permission FROM sys_menu WHERE permission LIKE '${MODULE}:%';"
echo ""
