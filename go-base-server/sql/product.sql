-- 产品信息 建表SQL
-- 生成时间: 2026-02-14 23:32:38
-- 模块: product

CREATE TABLE IF NOT EXISTS `product` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `type_id` INT UNSIGNED NOT NULL COMMENT '产品类型',
  `name` VARCHAR(255) NOT NULL COMMENT '产品名称',
  `num` INT NULL COMMENT '产品数量',
  `price` DOUBLE NULL COMMENT '产品单价',
  `status` VARCHAR(255) NOT NULL COMMENT '状态',
  `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_name` (`name`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品信息';
