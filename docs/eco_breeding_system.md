# 生态养殖管理系统设计文档

---

## 一、系统命名

| 项目 | 内容 |
|------|------|
| **系统中文名** | 牧智云 |
| **系统英文名** | EcoBreed |
| **login_slogan** | 智慧养殖，生态未来 |
| **login_desc** | 牧智云是一套专业的生态养殖管理平台，为养殖户提供牲畜管理、疫苗接种、疾病防控、饲料追踪等全流程数字化解决方案 |
| **login_features** | 牲畜追溯、智能喂养、疫苗管理、疾病预警 |

---

## 二、项目概述

### 2.1 项目背景
随着现代农业的发展，传统养殖模式已难以满足规模化、标准化的生产需求。牧智云生态养殖管理系统旨在通过信息化手段，帮助养殖企业和个体养殖户实现养殖全过程的数字化管理，提升养殖效率，降低养殖风险。

### 2.2 项目目标
- 实现牲畜信息的全生命周期管理
- 规范饲料使用和喂养记录
- 建立完善的疫苗接种和疾病防控体系
- 提供养殖知识学习平台
- 实现养殖户认证和权限管理

### 2.3 技术架构
- 后端：Go + Gin + GORM
- 前端：Vue3 + TypeScript + Ant Design Vue
- 数据库：MySQL 8.0+
- 权限：RBAC + Casbin

---

## 三、用户角色说明

### 3.1 管理员（admin）
系统管理员角色，拥有系统所有功能的管理权限。

**权限范围：**
- 管理所有基础数据（牲畜分类、养殖区域、饲料信息）
- 管理所有业务数据（牲畜批次、喂养记录、疫苗接种、疾病记录）
- 发布和管理养殖知识文章
- 审核养殖户认证信息
- 发布系统公告

### 3.2 养殖户（farmer）
通过认证的养殖户角色，可以管理自己的养殖数据。

**权限范围：**
- 浏览养殖知识文章
- 管理自己的养殖区域
- 管理自己的牲畜批次信息
- 管理自己的饲料和喂养记录
- 管理自己的疫苗接种记录
- 上报和管理自己的疾病记录

**使用限制：**
- 需要通过实名认证才能使用养殖管理功能
- 只能查看和管理自己创建的数据

---

## 四、功能模块详细说明

### 4.1 牲畜分类管理（livestock_category）

**功能描述：** 管理牲畜的分类信息，如猪、牛、羊、鸡、鸭、鹅等。

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 分类名称 |
| code | string | 是 | 分类编码 |
| icon | string | 否 | 分类图标 |
| description | string | 否 | 分类描述 |
| sort | int | 否 | 排序号 |
| status | int | 否 | 状态(0-禁用 1-启用) |

**操作权限：** 仅管理员

---

### 4.2 养殖区域管理（breeding_area）

**功能描述：** 管理养殖区域/圈舍信息。

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 区域名称 |
| code | string | 是 | 区域编码 |
| area_type | string | 否 | 区域类型(圈舍/牧场/鱼塘等) |
| capacity | int | 否 | 容纳数量 |
| location | string | 否 | 位置描述 |
| description | string | 否 | 区域描述 |
| status | int | 否 | 状态(0-禁用 1-启用) |
| created_by | uint | 是 | 创建人ID |

**操作权限：** 
- 管理员：管理所有区域
- 养殖户：管理自己创建的区域（数据隔离）

---

### 4.3 牲畜批次管理（livestock_batch）

**功能描述：** 管理牲畜批次信息，同一类型的牲畜在同一区域的一批记录。

**字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| batch_no | string | 是 | 批次编号 |
| category_id | uint | 是 | 牲畜分类ID |
| area_id | uint | 是 | 养殖区域ID |
| quantity | int | 是 | 数量 |
| species | string | 否 | 牲畜品种 |
| intro | string | 否 | 简介 |
| purchase_price | decimal | 否 | 购入单价 |
| birth_date | date | 否 | 出生日期 |
| status | int | 否 | 状态(0-在养 1-出栏 2-死亡) |
| remark | string | 否 | 备注 |
| created_by | uint | 是 | 创建人ID |

```bash
livestock
{
                "id": 9,
                "batchNo": "Y10000102",
                "name": "麻鸭",
                "img": "http://111.229.89.126:9093/files/download/1758524549601-麻鸭01.jpg",
                "farmerId": 2,
                "typeId": 6,
                "areaId": 9,
                "sex": "公",
                "species": "绍兴麻鸭",
                "intro": "原产于浙江绍兴，是我国著名的蛋用型地方品种。\n特点：体型较小，羽毛以麻褐色为主，公鸭头颈墨绿色，母鸭全身麻色。\n产蛋性能优异，年产蛋量可达280-320枚，蛋壳多为青色（绿壳蛋），适合放养和稻田养鸭。",
                "breedNumber": 100,
                "birthDate": "2025-08-20",
                "status": "在养",
                "farmerName": "自然源生态农业",
                "typeName": "鸭",
                "areaName": "二号育鸭舍"
            },
```



**操作权限：** 

- 管理员：管理所有批次
- 养殖户：管理自己创建的批次（数据隔离）

---

### 4.4 饲料信息管理（feed_info）

**功能描述：** 管理饲料的基本信息。

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 饲料名称 |
| code | string | 否 | 饲料编码 |
| feed_type | string | 是 | 饲料类型(精料/粗料/添加剂等) |
| brand | string | 否 | 品牌 |
| specification | string | 否 | 规格 |
| unit | string | 否 | 单位 |
| price | decimal | 否 | 单价 |
| stock | decimal | 否 | 库存数量 |
| description | string | 否 | 饲料描述 |
| status | int | 否 | 状态(0-禁用 1-启用) |
| created_by | uint | 是 | 创建人ID |

**操作权限：** 
- 管理员：管理所有饲料
- 养殖户：管理自己创建的饲料（数据隔离）

```bash
广州动物园 https://mmapgwh.map.qq.com/shortlink/short?l=_cc66b0967a1fc6027bfb0efc4369ff8d&tempSource=pcMap


荔枝湾景区 https://mmapgwh.map.qq.com/shortlink/short?l=_509d38d6d841a6729500221242cfcd81&tempSource=pcMap
```



---

### 4.5 喂养记录管理（feeding_record）

**功能描述：** 记录每次的喂养情况。

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| batch_id | uint | 是 | 牲畜批次ID |
| feed_id | uint | 是 | 饲料ID |
| feed_amount | decimal | 是 | 喂养数量 |
| feed_time | datetime | 是 | 喂养时间 |
| operator | string | 否 | 操作人员 |
| remark | string | 否 | 备注 |
| created_by | uint | 是 | 创建人ID |

**操作权限：** 
- 管理员：管理所有记录
- 养殖户：管理自己创建的记录（数据隔离）

---

### 4.6 疫苗接种管理（vaccination_record）

**功能描述：** 管理牲畜的疫苗接种记录。

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| batch_id | uint | 是 | 牲畜批次ID |
| vaccine_name | string | 是 | 疫苗名称 |
| vaccine_type | string | 否 | 疫苗类型 |
| vaccine_batch | string | 否 | 疫苗批号 |
| dose | string | 否 | 剂量 |
| injection_method | string | 否 | 接种方式 |
| vaccination_date | date | 是 | 接种日期 |
| next_date | date | 否 | 下次接种日期 |
| veterinarian | string | 否 | 接种兽医 |
| cost | decimal | 否 | 费用 |
| remark | string | 否 | 备注 |
| created_by | uint | 是 | 创建人ID |

**操作权限：** 
- 管理员：管理所有记录
- 养殖户：管理自己创建的记录（数据隔离）

---

### 4.7 疾病记录管理（disease_record）

**功能描述：** 记录和上报牲畜疾病情况。

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| batch_id | uint | 是 | 牲畜批次ID |
| disease_name | string | 是 | 疾病名称 |
| symptoms | text | 是 | 症状描述 |
| affected_count | int | 否 | 发病数量 |
| death_count | int | 否 | 死亡数量 |
| discover_date | date | 是 | 发现日期 |
| treatment | text | 否 | 治疗方案 |
| treatment_cost | decimal | 否 | 治疗费用 |
| recover_date | date | 否 | 康复日期 |
| status | int | 否 | 状态(0-治疗中 1-已康复 2-已处理) |
| remark | string | 否 | 备注 |
| created_by | uint | 是 | 创建人ID |

**操作权限：** 
- 管理员：管理所有记录
- 养殖户：上报和管理自己的记录（数据隔离）

---

### 4.8 养殖知识管理（breeding_knowledge）

**功能描述：** 管理员发布养殖知识文章，供养殖户学习。

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 文章标题 |
| category | string | 否 | 文章分类 |
| cover | string | 否 | 封面图片 |
| summary | string | 否 | 文章摘要 |
| content | text | 是 | 文章内容 |
| author | string | 否 | 作者 |
| view_count | int | 否 | 浏览次数 |
| sort | int | 否 | 排序号 |
| status | int | 否 | 状态(0-草稿 1-发布) |
| created_by | uint | 是 | 创建人ID |

**操作权限：** 
- 管理员：发布和管理文章
- 养殖户：仅浏览（前台接口）

---

### 4.9 养殖户认证管理（farmer_certification）

**功能描述：** 养殖户提交认证申请，管理员审核。

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| real_name | string | 是 | 真实姓名 |
| id_card | string | 是 | 身份证号 |
| id_card_front | string | 是 | 身份证正面照 |
| id_card_back | string | 是 | 身份证背面照 |
| phone | string | 是 | 联系电话 |
| address | string | 是 | 养殖地址 |
| farm_name | string | 否 | 养殖场名称 |
| farm_scale | string | 否 | 养殖规模 |
| business_license | string | 否 | 营业执照 |
| created_by | uint | 是 | 申请人用户ID |

**审核字段：**
| 字段 | 类型 | 说明 |
|------|------|------|
| audit_status | int | 审批状态(0-待审核 1-通过 2-拒绝) |
| audit_remark | string | 审批备注 |
| audit_time | datetime | 审批时间 |
| audit_by | uint | 审批人ID |

**操作权限：** 
- 管理员：审核认证申请
- 养殖户：提交认证申请，查看自己的认证状态

---

### 4.10 公告信息管理（notice）

**功能描述：** 管理员发布系统公告。

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 公告标题 |
| content | text | 是 | 公告内容 |
| notice_type | string | 否 | 公告类型(通知/警告/活动等) |
| priority | int | 否 | 优先级 |
| start_time | datetime | 否 | 生效时间 |
| end_time | datetime | 否 | 失效时间 |
| status | int | 否 | 状态(0-草稿 1-发布) |
| created_by | uint | 是 | 创建人ID |

**操作权限：** 仅管理员

---

## 五、数据字典

### 5.1 状态码定义

**通用状态：**
| 值 | 说明 |
|----|------|
| 0 | 禁用/关闭 |
| 1 | 启用/开启 |

**审核状态（audit_status）：**
| 值 | 说明 |
|----|------|
| 0 | 待审核 |
| 1 | 审核通过 |
| 2 | 审核拒绝 |

**牲畜状态：**
| 值 | 说明 |
|----|------|
| 0 | 在养 |
| 1 | 出栏 |
| 2 | 死亡 |

**疾病状态：**
| 值 | 说明 |
|----|------|
| 0 | 治疗中 |
| 1 | 已康复 |
| 2 | 已处理 |

**文章状态：**
| 值 | 说明 |
|----|------|
| 0 | 草稿 |
| 1 | 已发布 |

### 5.2 饲料类型

| 编码 | 说明 |
|------|------|
| concentrate | 精料 |
| roughage | 粗料 |
| additive | 添加剂 |
| premix | 预混料 |
| complete | 全价料 |

### 5.3 养殖区域类型

| 编码 | 说明 |
|------|------|
| barn | 圈舍 |
| pasture | 牧场 |
| pond | 鱼塘 |
| coop | 禽舍 |
| other | 其他 |

### 5.4 公告类型

| 编码 | 说明 |
|------|------|
| notice | 通知 |
| warning | 警告 |
| activity | 活动 |
| policy | 政策 |

---

## 六、API接口规划

### 6.1 管理员接口

#### 6.1.1 牲畜分类
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/livestock_category | 分类列表 |
| GET | /api/v1/livestock_category/:id | 分类详情 |
| POST | /api/v1/livestock_category | 创建分类 |
| PUT | /api/v1/livestock_category/:id | 更新分类 |
| DELETE | /api/v1/livestock_category/:id | 删除分类 |
| GET | /api/v1/livestock_category/options | 分类选项 |

#### 6.1.2 养殖区域
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/breeding_area | 区域列表 |
| GET | /api/v1/breeding_area/:id | 区域详情 |
| POST | /api/v1/breeding_area | 创建区域 |
| PUT | /api/v1/breeding_area/:id | 更新区域 |
| DELETE | /api/v1/breeding_area/:id | 删除区域 |
| GET | /api/v1/breeding_area/options | 区域选项 |

#### 6.1.3 牲畜批次
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/livestock_batch | 批次列表 |
| GET | /api/v1/livestock_batch/:id | 批次详情 |
| POST | /api/v1/livestock_batch | 创建批次 |
| PUT | /api/v1/livestock_batch/:id | 更新批次 |
| DELETE | /api/v1/livestock_batch/:id | 删除批次 |
| GET | /api/v1/livestock_batch/options | 批次选项 |

#### 6.1.4 饲料信息
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/feed_info | 饲料列表 |
| GET | /api/v1/feed_info/:id | 饲料详情 |
| POST | /api/v1/feed_info | 创建饲料 |
| PUT | /api/v1/feed_info/:id | 更新饲料 |
| DELETE | /api/v1/feed_info/:id | 删除饲料 |
| GET | /api/v1/feed_info/options | 饲料选项 |

#### 6.1.5 喂养记录
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/feeding_record | 记录列表 |
| GET | /api/v1/feeding_record/:id | 记录详情 |
| POST | /api/v1/feeding_record | 创建记录 |
| PUT | /api/v1/feeding_record/:id | 更新记录 |
| DELETE | /api/v1/feeding_record/:id | 删除记录 |

#### 6.1.6 疫苗接种
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/vaccination_record | 接种列表 |
| GET | /api/v1/vaccination_record/:id | 接种详情 |
| POST | /api/v1/vaccination_record | 创建接种 |
| PUT | /api/v1/vaccination_record/:id | 更新接种 |
| DELETE | /api/v1/vaccination_record/:id | 删除接种 |

#### 6.1.7 疾病记录
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/disease_record | 记录列表 |
| GET | /api/v1/disease_record/:id | 记录详情 |
| POST | /api/v1/disease_record | 创建记录 |
| PUT | /api/v1/disease_record/:id | 更新记录 |
| DELETE | /api/v1/disease_record/:id | 删除记录 |

#### 6.1.8 养殖知识
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/breeding_knowledge | 文章列表 |
| GET | /api/v1/breeding_knowledge/:id | 文章详情 |
| POST | /api/v1/breeding_knowledge | 创建文章 |
| PUT | /api/v1/breeding_knowledge/:id | 更新文章 |
| DELETE | /api/v1/breeding_knowledge/:id | 删除文章 |

#### 6.1.9 养殖户认证
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/farmer_certification | 认证列表 |
| GET | /api/v1/farmer_certification/:id | 认证详情 |
| PUT | /api/v1/farmer_certification/:id/audit | 审核认证 |

#### 6.1.10 公告管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/notice | 公告列表 |
| GET | /api/v1/notice/:id | 公告详情 |
| POST | /api/v1/notice | 创建公告 |
| PUT | /api/v1/notice/:id | 更新公告 |
| DELETE | /api/v1/notice/:id | 删除公告 |

### 6.2 养殖户接口（前台）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/breeding_knowledge/frontend | 养殖知识列表 |
| GET | /api/v1/breeding_knowledge/frontend/:id | 养殖知识详情 |
| GET | /api/v1/notice/frontend | 公告列表 |
| GET | /api/v1/notice/frontend/:id | 公告详情 |
| GET | /api/v1/farmer_certification/my | 我的认证信息 |
| POST | /api/v1/farmer_certification | 提交认证申请 |
| PUT | /api/v1/farmer_certification/my | 更新认证信息 |

---

## 七、数据库设计（MySQL SQL脚本）

```sql
-- =====================================================
-- 牧智云生态养殖管理系统 数据库设计
-- 生成时间: 2026-02-06
-- 数据库: MySQL 8.0+
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 1. 牲畜分类表
-- ----------------------------
DROP TABLE IF EXISTS `livestock_category`;
CREATE TABLE `livestock_category` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
  `code` VARCHAR(50) NOT NULL COMMENT '分类编码',
  `icon` VARCHAR(255) NULL DEFAULT '' COMMENT '分类图标',
  `description` VARCHAR(500) NULL DEFAULT '' COMMENT '分类描述',
  `sort` INT NULL DEFAULT 0 COMMENT '排序号',
  `status` TINYINT NULL DEFAULT 1 COMMENT '状态(0-禁用 1-启用)',
  `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_code` (`code`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='牲畜分类表';

-- ----------------------------
-- 2. 养殖区域表
-- ----------------------------
DROP TABLE IF EXISTS `breeding_area`;
CREATE TABLE `breeding_area` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` VARCHAR(100) NOT NULL COMMENT '区域名称',
  `code` VARCHAR(50) NOT NULL COMMENT '区域编码',
  `area_type` VARCHAR(50) NULL DEFAULT 'barn' COMMENT '区域类型(barn-圈舍 pasture-牧场 pond-鱼塘 coop-禽舍 other-其他)',
  `capacity` INT NULL DEFAULT 0 COMMENT '容纳数量',
  `location` VARCHAR(255) NULL DEFAULT '' COMMENT '位置描述',
  `description` VARCHAR(500) NULL DEFAULT '' COMMENT '区域描述',
  `status` TINYINT NULL DEFAULT 1 COMMENT '状态(0-禁用 1-启用)',
  `created_by` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '创建人ID',
  `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_code` (`code`),
  INDEX `idx_created_by` (`created_by`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='养殖区域表';

-- ----------------------------
-- 3. 牲畜批次表
-- ----------------------------
DROP TABLE IF EXISTS `livestock_batch`;
CREATE TABLE `livestock_batch` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `batch_no` VARCHAR(50) NOT NULL COMMENT '批次编号',
  `category_id` BIGINT UNSIGNED NOT NULL COMMENT '牲畜分类ID',
  `area_id` BIGINT UNSIGNED NOT NULL COMMENT '养殖区域ID',
  `quantity` INT NOT NULL DEFAULT 0 COMMENT '数量',
  `source` VARCHAR(50) NULL DEFAULT 'purchase' COMMENT '来源(self-自繁 purchase-购入)',
  `purchase_date` DATE NULL COMMENT '购入日期',
  `purchase_price` DECIMAL(10,2) NULL DEFAULT 0.00 COMMENT '购入单价',
  `birth_date` DATE NULL COMMENT '出生日期',
  `status` TINYINT NULL DEFAULT 0 COMMENT '状态(0-在养 1-出栏 2-死亡)',
  `remark` VARCHAR(500) NULL DEFAULT '' COMMENT '备注',
  `created_by` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '创建人ID',
  `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_batch_no` (`batch_no`),
  INDEX `idx_category_id` (`category_id`),
  INDEX `idx_area_id` (`area_id`),
  INDEX `idx_created_by` (`created_by`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='牲畜批次表';

-- ----------------------------
-- 4. 饲料信息表
-- ----------------------------
DROP TABLE IF EXISTS `feed_info`;
CREATE TABLE `feed_info` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` VARCHAR(100) NOT NULL COMMENT '饲料名称',
  `code` VARCHAR(50) NULL DEFAULT '' COMMENT '饲料编码',
  `feed_type` VARCHAR(50) NOT NULL COMMENT '饲料类型(concentrate-精料 roughage-粗料 additive-添加剂 premix-预混料 complete-全价料)',
  `brand` VARCHAR(100) NULL DEFAULT '' COMMENT '品牌',
  `specification` VARCHAR(100) NULL DEFAULT '' COMMENT '规格',
  `unit` VARCHAR(20) NULL DEFAULT 'kg' COMMENT '单位',
  `price` DECIMAL(10,2) NULL DEFAULT 0.00 COMMENT '单价',
  `stock` DECIMAL(10,2) NULL DEFAULT 0.00 COMMENT '库存数量',
  `description` VARCHAR(500) NULL DEFAULT '' COMMENT '饲料描述',
  `status` TINYINT NULL DEFAULT 1 COMMENT '状态(0-禁用 1-启用)',
  `created_by` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '创建人ID',
  `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  INDEX `idx_feed_type` (`feed_type`),
  INDEX `idx_created_by` (`created_by`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='饲料信息表';

-- ----------------------------
-- 5. 喂养记录表
-- ----------------------------
DROP TABLE IF EXISTS `feeding_record`;
CREATE TABLE `feeding_record` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `batch_id` BIGINT UNSIGNED NOT NULL COMMENT '牲畜批次ID',
  `feed_id` BIGINT UNSIGNED NOT NULL COMMENT '饲料ID',
  `feed_amount` DECIMAL(10,2) NOT NULL COMMENT '喂养数量',
  `feed_time` DATETIME NOT NULL COMMENT '喂养时间',
  `operator` VARCHAR(50) NULL DEFAULT '' COMMENT '操作人员',
  `remark` VARCHAR(500) NULL DEFAULT '' COMMENT '备注',
  `created_by` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '创建人ID',
  `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  INDEX `idx_batch_id` (`batch_id`),
  INDEX `idx_feed_id` (`feed_id`),
  INDEX `idx_feed_time` (`feed_time`),
  INDEX `idx_created_by` (`created_by`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='喂养记录表';

-- ----------------------------
-- 6. 疫苗接种记录表
-- ----------------------------
DROP TABLE IF EXISTS `vaccination_record`;
CREATE TABLE `vaccination_record` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `batch_id` BIGINT UNSIGNED NOT NULL COMMENT '牲畜批次ID',
  `vaccine_name` VARCHAR(100) NOT NULL COMMENT '疫苗名称',
  `vaccine_type` VARCHAR(50) NULL DEFAULT '' COMMENT '疫苗类型',
  `vaccine_batch` VARCHAR(50) NULL DEFAULT '' COMMENT '疫苗批号',
  `dose` VARCHAR(50) NULL DEFAULT '' COMMENT '剂量',
  `injection_method` VARCHAR(50) NULL DEFAULT '' COMMENT '接种方式',
  `vaccination_date` DATE NOT NULL COMMENT '接种日期',
  `next_date` DATE NULL COMMENT '下次接种日期',
  `veterinarian` VARCHAR(50) NULL DEFAULT '' COMMENT '接种兽医',
  `cost` DECIMAL(10,2) NULL DEFAULT 0.00 COMMENT '费用',
  `remark` VARCHAR(500) NULL DEFAULT '' COMMENT '备注',
  `created_by` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '创建人ID',
  `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  INDEX `idx_batch_id` (`batch_id`),
  INDEX `idx_vaccination_date` (`vaccination_date`),
  INDEX `idx_next_date` (`next_date`),
  INDEX `idx_created_by` (`created_by`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='疫苗接种记录表';

-- ----------------------------
-- 7. 疾病记录表
-- ----------------------------
DROP TABLE IF EXISTS `disease_record`;
CREATE TABLE `disease_record` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `batch_id` BIGINT UNSIGNED NOT NULL COMMENT '牲畜批次ID',
  `disease_name` VARCHAR(100) NOT NULL COMMENT '疾病名称',
  `symptoms` TEXT NOT NULL COMMENT '症状描述',
  `affected_count` INT NULL DEFAULT 0 COMMENT '发病数量',
  `death_count` INT NULL DEFAULT 0 COMMENT '死亡数量',
  `discover_date` DATE NOT NULL COMMENT '发现日期',
  `treatment` TEXT NULL COMMENT '治疗方案',
  `treatment_cost` DECIMAL(10,2) NULL DEFAULT 0.00 COMMENT '治疗费用',
  `recover_date` DATE NULL COMMENT '康复日期',
  `status` TINYINT NULL DEFAULT 0 COMMENT '状态(0-治疗中 1-已康复 2-已处理)',
  `remark` VARCHAR(500) NULL DEFAULT '' COMMENT '备注',
  `created_by` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '创建人ID',
  `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  INDEX `idx_batch_id` (`batch_id`),
  INDEX `idx_discover_date` (`discover_date`),
  INDEX `idx_status` (`status`),
  INDEX `idx_created_by` (`created_by`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='疾病记录表';

-- ----------------------------
-- 8. 养殖知识表
-- ----------------------------
DROP TABLE IF EXISTS `breeding_knowledge`;
CREATE TABLE `breeding_knowledge` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `title` VARCHAR(200) NOT NULL COMMENT '文章标题',
  `category` VARCHAR(50) NULL DEFAULT '' COMMENT '文章分类',
  `cover_file_id` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '封面图片文件ID',
  `summary` VARCHAR(500) NULL DEFAULT '' COMMENT '文章摘要',
  `content` LONGTEXT NOT NULL COMMENT '文章内容',
  `author` VARCHAR(50) NULL DEFAULT '' COMMENT '作者',
  `view_count` INT NULL DEFAULT 0 COMMENT '浏览次数',
  `sort` INT NULL DEFAULT 0 COMMENT '排序号',
  `status` TINYINT NULL DEFAULT 0 COMMENT '状态(0-草稿 1-发布)',
  `created_by` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '创建人ID',
  `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  INDEX `idx_category` (`category`),
  INDEX `idx_status` (`status`),
  INDEX `idx_created_by` (`created_by`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='养殖知识表';

-- ----------------------------
-- 9. 养殖户认证表
-- ----------------------------
DROP TABLE IF EXISTS `farmer_certification`;
CREATE TABLE `farmer_certification` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `real_name` VARCHAR(50) NOT NULL COMMENT '真实姓名',
  `id_card` VARCHAR(18) NOT NULL COMMENT '身份证号',
  `id_card_front_file_id` BIGINT UNSIGNED NOT NULL COMMENT '身份证正面照文件ID',
  `id_card_back_file_id` BIGINT UNSIGNED NOT NULL COMMENT '身份证背面照文件ID',
  `phone` VARCHAR(20) NOT NULL COMMENT '联系电话',
  `address` VARCHAR(255) NOT NULL COMMENT '养殖地址',
  `farm_name` VARCHAR(100) NULL DEFAULT '' COMMENT '养殖场名称',
  `farm_scale` VARCHAR(100) NULL DEFAULT '' COMMENT '养殖规模',
  `business_license_file_id` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '营业执照文件ID',
  `created_by` BIGINT UNSIGNED NOT NULL COMMENT '申请人用户ID',
  `audit_status` TINYINT NULL DEFAULT 0 COMMENT '审批状态(0-待审核 1-通过 2-拒绝)',
  `audit_remark` VARCHAR(500) NULL DEFAULT '' COMMENT '审批备注',
  `audit_time` DATETIME NULL COMMENT '审批时间',
  `audit_by` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '审批人ID',
  `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_created_by` (`created_by`),
  INDEX `idx_audit_status` (`audit_status`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='养殖户认证表';

-- ----------------------------
-- 10. 公告信息表
-- ----------------------------
DROP TABLE IF EXISTS `notice`;
CREATE TABLE `notice` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `title` VARCHAR(200) NOT NULL COMMENT '公告标题',
  `content` LONGTEXT NOT NULL COMMENT '公告内容',
  `notice_type` VARCHAR(50) NULL DEFAULT 'notice' COMMENT '公告类型(notice-通知 warning-警告 activity-活动 policy-政策)',
  `priority` INT NULL DEFAULT 0 COMMENT '优先级',
  `start_time` DATETIME NULL COMMENT '生效时间',
  `end_time` DATETIME NULL COMMENT '失效时间',
  `status` TINYINT NULL DEFAULT 0 COMMENT '状态(0-草稿 1-发布)',
  `created_by` BIGINT UNSIGNED NULL DEFAULT 0 COMMENT '创建人ID',
  `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  INDEX `idx_notice_type` (`notice_type`),
  INDEX `idx_status` (`status`),
  INDEX `idx_start_time` (`start_time`),
  INDEX `idx_end_time` (`end_time`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告信息表';

-- ----------------------------
-- 初始化数据
-- ----------------------------

-- 牲畜分类初始数据
INSERT INTO `livestock_category` (`name`, `code`, `description`, `sort`, `status`) VALUES
('猪', 'pig', '家猪养殖', 1, 1),
('牛', 'cattle', '肉牛、奶牛养殖', 2, 1),
('羊', 'sheep', '绵羊、山羊养殖', 3, 1),
('鸡', 'chicken', '肉鸡、蛋鸡养殖', 4, 1),
('鸭', 'duck', '肉鸭、蛋鸭养殖', 5, 1),
('鹅', 'goose', '家鹅养殖', 6, 1),
('兔', 'rabbit', '肉兔、毛兔养殖', 7, 1),
('鱼', 'fish', '淡水鱼、海水鱼养殖', 8, 1);

-- 养殖户角色
INSERT INTO `sys_role` (`name`, `code`, `sort`, `status`, `remark`, `created_at`, `updated_at`) VALUES
('养殖户', 'farmer', 10, 1, '通过认证的养殖户角色', NOW(), NOW());

SET FOREIGN_KEY_CHECKS = 1;
```

---

## 八、代码生成器使用说明

本项目内置代码生成器，可以根据数据库表结构自动生成后端和前端代码。

### 8.1 使用方式

1. 访问系统管理后台
2. 进入「开发工具」->「代码生成」
3. 选择要生成代码的数据库表
4. 配置字段属性和关联关系
5. 预览生成的代码
6. 确认后执行生成

### 8.2 各模块配置示例

#### 8.2.1 牲畜分类（livestock_category）

```json
{
  "table_name": "livestock_category",
  "module_name": "livestock_category",
  "description": "牲畜分类",
  "generate_backend": true,
  "generate_frontend": true,
  "generate_sql": false,
  "has_created_at": true,
  "has_updated_at": true,
  "has_deleted_at": true,
  "has_created_by": false,
  "data_isolation": false,
  "columns": [
    {"column_name": "name", "field_name": "Name", "field_type": "string", "json_name": "name", "ts_type": "string", "comment": "分类名称", "is_required": true, "is_searchable": true, "search_type": "like", "is_list_visible": true, "is_form_visible": true, "form_type": "input"},
    {"column_name": "code", "field_name": "Code", "field_type": "string", "json_name": "code", "ts_type": "string", "comment": "分类编码", "is_required": true, "is_searchable": true, "is_list_visible": true, "is_form_visible": true, "is_unique": true, "form_type": "input"},
    {"column_name": "icon", "field_name": "Icon", "field_type": "string", "json_name": "icon", "ts_type": "string", "comment": "分类图标", "is_list_visible": false, "is_form_visible": true, "form_type": "input"},
    {"column_name": "description", "field_name": "Description", "field_type": "string", "json_name": "description", "ts_type": "string", "comment": "分类描述", "is_list_visible": false, "is_form_visible": true, "form_type": "textarea"},
    {"column_name": "sort", "field_name": "Sort", "field_type": "int", "json_name": "sort", "ts_type": "number", "comment": "排序号", "is_list_visible": true, "is_form_visible": true, "is_sortable": true, "form_type": "number", "default_value": "0"},
    {"column_name": "status", "field_name": "Status", "field_type": "int", "json_name": "status", "ts_type": "number", "comment": "状态", "is_searchable": true, "search_type": "eq", "is_list_visible": true, "is_form_visible": true, "form_type": "switch", "switch_values": {"active_value": 1, "inactive_value": 0, "active_text": "启用", "inactive_text": "禁用"}, "default_value": "1"}
  ],
  "menu_config": {
    "parent_id": 0,
    "menu_name": "牲畜分类",
    "menu_icon": "AppstoreOutlined",
    "menu_sort": 1,
    "permission": "livestock_category"
  }
}
```

#### 8.2.2 养殖区域（breeding_area）- 数据隔离

```json
{
  "table_name": "breeding_area",
  "module_name": "breeding_area",
  "description": "养殖区域",
  "generate_backend": true,
  "generate_frontend": true,
  "generate_sql": false,
  "has_created_at": true,
  "has_updated_at": true,
  "has_deleted_at": true,
  "has_created_by": true,
  "data_isolation": true,
  "admin_role_ids": "1",
  "columns": [
    {"column_name": "name", "field_name": "Name", "field_type": "string", "json_name": "name", "ts_type": "string", "comment": "区域名称", "is_required": true, "is_searchable": true, "search_type": "like", "is_list_visible": true, "is_form_visible": true, "form_type": "input"},
    {"column_name": "code", "field_name": "Code", "field_type": "string", "json_name": "code", "ts_type": "string", "comment": "区域编码", "is_required": true, "is_searchable": true, "is_list_visible": true, "is_form_visible": true, "is_unique": true, "form_type": "input"},
    {"column_name": "area_type", "field_name": "AreaType", "field_type": "string", "json_name": "area_type", "ts_type": "string", "comment": "区域类型", "is_searchable": true, "search_type": "eq", "is_list_visible": true, "is_form_visible": true, "form_type": "select", "select_options": [{"label": "圈舍", "value": "barn"}, {"label": "牧场", "value": "pasture"}, {"label": "鱼塘", "value": "pond"}, {"label": "禽舍", "value": "coop"}, {"label": "其他", "value": "other"}], "default_value": "'barn'"},
    {"column_name": "capacity", "field_name": "Capacity", "field_type": "int", "json_name": "capacity", "ts_type": "number", "comment": "容纳数量", "is_list_visible": true, "is_form_visible": true, "form_type": "number", "default_value": "0"},
    {"column_name": "location", "field_name": "Location", "field_type": "string", "json_name": "location", "ts_type": "string", "comment": "位置描述", "is_list_visible": false, "is_form_visible": true, "form_type": "input"},
    {"column_name": "description", "field_name": "Description", "field_type": "string", "json_name": "description", "ts_type": "string", "comment": "区域描述", "is_list_visible": false, "is_form_visible": true, "form_type": "textarea"},
    {"column_name": "status", "field_name": "Status", "field_type": "int", "json_name": "status", "ts_type": "number", "comment": "状态", "is_searchable": true, "search_type": "eq", "is_list_visible": true, "is_form_visible": true, "form_type": "switch", "switch_values": {"active_value": 1, "inactive_value": 0, "active_text": "启用", "inactive_text": "禁用"}, "default_value": "1"}
  ],
  "menu_config": {
    "parent_id": 0,
    "menu_name": "养殖区域",
    "menu_icon": "HomeOutlined",
    "menu_sort": 2,
    "permission": "breeding_area"
  }
}
```

#### 8.2.3 牲畜批次（livestock_batch）- 带关联关系

```json
{
  "table_name": "livestock_batch",
  "module_name": "livestock_batch",
  "description": "牲畜批次",
  "generate_backend": true,
  "generate_frontend": true,
  "generate_sql": false,
  "has_created_at": true,
  "has_updated_at": true,
  "has_deleted_at": true,
  "has_created_by": true,
  "data_isolation": true,
  "admin_role_ids": "1",
  "columns": [
    {"column_name": "batch_no", "field_name": "BatchNo", "field_type": "string", "json_name": "batch_no", "ts_type": "string", "comment": "批次编号", "is_required": true, "is_searchable": true, "search_type": "like", "is_list_visible": true, "is_form_visible": true, "is_unique": true, "form_type": "input"},
    {"column_name": "category_id", "field_name": "CategoryId", "field_type": "uint", "json_name": "category_id", "ts_type": "number", "comment": "牲畜分类ID", "is_required": true, "is_form_visible": true, "form_type": "select"},
    {"column_name": "area_id", "field_name": "AreaId", "field_type": "uint", "json_name": "area_id", "ts_type": "number", "comment": "养殖区域ID", "is_required": true, "is_form_visible": true, "form_type": "select"},
    {"column_name": "quantity", "field_name": "Quantity", "field_type": "int", "json_name": "quantity", "ts_type": "number", "comment": "数量", "is_required": true, "is_list_visible": true, "is_form_visible": true, "form_type": "number"},
    {"column_name": "source", "field_name": "Source", "field_type": "string", "json_name": "source", "ts_type": "string", "comment": "来源", "is_list_visible": true, "is_form_visible": true, "form_type": "select", "select_options": [{"label": "自繁", "value": "self"}, {"label": "购入", "value": "purchase"}], "default_value": "'purchase'"},
    {"column_name": "purchase_date", "field_name": "PurchaseDate", "field_type": "time.Time", "json_name": "purchase_date", "ts_type": "string", "comment": "购入日期", "is_list_visible": true, "is_form_visible": true, "form_type": "date"},
    {"column_name": "purchase_price", "field_name": "PurchasePrice", "field_type": "float64", "json_name": "purchase_price", "ts_type": "number", "comment": "购入单价", "db_type": "DECIMAL", "is_list_visible": true, "is_form_visible": true, "form_type": "number"},
    {"column_name": "birth_date", "field_name": "BirthDate", "field_type": "time.Time", "json_name": "birth_date", "ts_type": "string", "comment": "出生日期", "is_form_visible": true, "form_type": "date"},
    {"column_name": "status", "field_name": "Status", "field_type": "int", "json_name": "status", "ts_type": "number", "comment": "状态", "is_searchable": true, "search_type": "eq", "is_list_visible": true, "is_form_visible": true, "form_type": "select", "select_options": [{"label": "在养", "value": 0}, {"label": "出栏", "value": 1}, {"label": "死亡", "value": 2}], "default_value": "0"},
    {"column_name": "remark", "field_name": "Remark", "field_type": "string", "json_name": "remark", "ts_type": "string", "comment": "备注", "is_form_visible": true, "form_type": "textarea"}
  ],
  "relations": [
    {"relation_type": "belongsTo", "related_table": "livestock_category", "related_model": "LivestockCategory", "foreign_key": "category_id", "reference_key": "id", "display_field": "name", "comment": "牲畜分类", "use_options_api": true},
    {"relation_type": "belongsTo", "related_table": "breeding_area", "related_model": "BreedingArea", "foreign_key": "area_id", "reference_key": "id", "display_field": "name", "comment": "养殖区域", "use_options_api": true}
  ],
  "menu_config": {
    "parent_id": 0,
    "menu_name": "牲畜批次",
    "menu_icon": "UnorderedListOutlined",
    "menu_sort": 3,
    "permission": "livestock_batch"
  }
}
```

#### 8.2.4 养殖户认证（farmer_certification）- 带审批

```json
{
  "table_name": "farmer_certification",
  "module_name": "farmer_certification",
  "description": "养殖户认证",
  "generate_backend": true,
  "generate_frontend": true,
  "generate_sql": false,
  "has_created_at": true,
  "has_updated_at": true,
  "has_deleted_at": true,
  "has_created_by": true,
  "has_audit": true,
  "data_isolation": true,
  "admin_role_ids": "1",
  "generate_frontend_api": true,
  "columns": [
    {"column_name": "real_name", "field_name": "RealName", "field_type": "string", "json_name": "real_name", "ts_type": "string", "comment": "真实姓名", "is_required": true, "is_searchable": true, "search_type": "like", "is_list_visible": true, "is_form_visible": true, "form_type": "input"},
    {"column_name": "id_card", "field_name": "IdCard", "field_type": "string", "json_name": "id_card", "ts_type": "string", "comment": "身份证号", "is_required": true, "is_list_visible": true, "is_form_visible": true, "form_type": "input"},
    {"column_name": "id_card_front", "field_name": "IdCardFront", "field_type": "string", "json_name": "id_card_front", "ts_type": "string", "comment": "身份证正面照", "is_required": true, "is_form_visible": true, "form_type": "image"},
    {"column_name": "id_card_back", "field_name": "IdCardBack", "field_type": "string", "json_name": "id_card_back", "ts_type": "string", "comment": "身份证背面照", "is_required": true, "is_form_visible": true, "form_type": "image"},
    {"column_name": "phone", "field_name": "Phone", "field_type": "string", "json_name": "phone", "ts_type": "string", "comment": "联系电话", "is_required": true, "is_searchable": true, "is_list_visible": true, "is_form_visible": true, "form_type": "input"},
    {"column_name": "address", "field_name": "Address", "field_type": "string", "json_name": "address", "ts_type": "string", "comment": "养殖地址", "is_required": true, "is_list_visible": true, "is_form_visible": true, "form_type": "input"},
    {"column_name": "farm_name", "field_name": "FarmName", "field_type": "string", "json_name": "farm_name", "ts_type": "string", "comment": "养殖场名称", "is_list_visible": true, "is_form_visible": true, "form_type": "input"},
    {"column_name": "farm_scale", "field_name": "FarmScale", "field_type": "string", "json_name": "farm_scale", "ts_type": "string", "comment": "养殖规模", "is_list_visible": true, "is_form_visible": true, "form_type": "input"},
    {"column_name": "business_license", "field_name": "BusinessLicense", "field_type": "string", "json_name": "business_license", "ts_type": "string", "comment": "营业执照", "is_form_visible": true, "form_type": "image"}
  ],
  "menu_config": {
    "parent_id": 0,
    "menu_name": "认证审核",
    "menu_icon": "SafetyCertificateOutlined",
    "menu_sort": 9,
    "permission": "farmer_certification"
  }
}
```

#### 8.2.5 养殖知识（breeding_knowledge）- 带前台接口

```json
{
  "table_name": "breeding_knowledge",
  "module_name": "breeding_knowledge",
  "description": "养殖知识",
  "generate_backend": true,
  "generate_frontend": true,
  "generate_sql": false,
  "has_created_at": true,
  "has_updated_at": true,
  "has_deleted_at": true,
  "has_created_by": true,
  "generate_frontend_api": true,
  "columns": [
    {"column_name": "title", "field_name": "Title", "field_type": "string", "json_name": "title", "ts_type": "string", "comment": "文章标题", "is_required": true, "is_searchable": true, "search_type": "like", "is_list_visible": true, "is_form_visible": true, "form_type": "input"},
    {"column_name": "category", "field_name": "Category", "field_type": "string", "json_name": "category", "ts_type": "string", "comment": "文章分类", "is_searchable": true, "search_type": "eq", "is_list_visible": true, "is_form_visible": true, "form_type": "select", "select_options": [{"label": "养殖技术", "value": "technology"}, {"label": "疾病防治", "value": "disease"}, {"label": "饲料营养", "value": "feed"}, {"label": "经营管理", "value": "management"}]},
    {"column_name": "cover", "field_name": "Cover", "field_type": "string", "json_name": "cover", "ts_type": "string", "comment": "封面图片", "is_list_visible": true, "is_form_visible": true, "form_type": "image"},
    {"column_name": "summary", "field_name": "Summary", "field_type": "string", "json_name": "summary", "ts_type": "string", "comment": "文章摘要", "is_list_visible": true, "is_form_visible": true, "form_type": "textarea"},
    {"column_name": "content", "field_name": "Content", "field_type": "string", "json_name": "content", "ts_type": "string", "comment": "文章内容", "db_type": "LONGTEXT", "is_required": true, "is_form_visible": true, "form_type": "editor"},
    {"column_name": "author", "field_name": "Author", "field_type": "string", "json_name": "author", "ts_type": "string", "comment": "作者", "is_list_visible": true, "is_form_visible": true, "form_type": "input"},
    {"column_name": "view_count", "field_name": "ViewCount", "field_type": "int", "json_name": "view_count", "ts_type": "number", "comment": "浏览次数", "is_list_visible": true, "is_sortable": true, "default_value": "0"},
    {"column_name": "sort", "field_name": "Sort", "field_type": "int", "json_name": "sort", "ts_type": "number", "comment": "排序号", "is_list_visible": true, "is_form_visible": true, "is_sortable": true, "form_type": "number", "default_value": "0"},
    {"column_name": "status", "field_name": "Status", "field_type": "int", "json_name": "status", "ts_type": "number", "comment": "状态", "is_searchable": true, "search_type": "eq", "is_list_visible": true, "is_form_visible": true, "form_type": "switch", "switch_values": {"active_value": 1, "inactive_value": 0, "active_text": "发布", "inactive_text": "草稿"}, "default_value": "0"}
  ],
  "menu_config": {
    "parent_id": 0,
    "menu_name": "养殖知识",
    "menu_icon": "ReadOutlined",
    "menu_sort": 8,
    "permission": "breeding_knowledge"
  }
}
```

### 8.3 生成步骤

1. **执行数据库脚本**
   - 先在MySQL中执行上述SQL脚本，创建所有数据表

2. **使用代码生成器**
   - 登录管理后台 -> 开发工具 -> 代码生成
   - 选择数据表（如 livestock_category）
   - 系统会自动解析字段信息
   - 根据需求配置：
     - 字段的搜索、列表、表单显示
     - 表单组件类型
     - 关联关系
     - 数据隔离和审批功能
   - 点击「预览」查看生成的代码
   - 确认无误后点击「生成」

3. **注册路由**
   - 生成的路由文件在 `router/modules/` 目录
   - 需要在 `router/router.go` 中导入并注册

4. **执行菜单SQL**
   - 生成的菜单SQL在 `sql/` 目录
   - 执行后即可在后台看到对应菜单

5. **分配权限**
   - 在角色管理中为相应角色分配新模块的菜单和API权限

### 8.4 推荐生成顺序

1. `livestock_category` - 牲畜分类（基础数据，无依赖）
2. `breeding_area` - 养殖区域（基础数据，无依赖）
3. `feed_info` - 饲料信息（基础数据，无依赖）
4. `livestock_batch` - 牲畜批次（依赖分类和区域）
5. `feeding_record` - 喂养记录（依赖批次和饲料）
6. `vaccination_record` - 疫苗接种（依赖批次）
7. `disease_record` - 疾病记录（依赖批次）
8. `breeding_knowledge` - 养殖知识
9. `farmer_certification` - 养殖户认证
10. `notice` - 公告信息

---

## 九、系统配置

### 9.1 前端配置更新

在 `go-base-web/.env.development` 中更新系统信息：

```env
# 系统名称
VITE_APP_TITLE=牧智云
```

### 9.2 登录页配置

在系统配置中添加以下配置项：

| 配置键 | 配置值 |
|--------|--------|
| login_slogan | 智慧养殖，生态未来 |
| login_desc | 牧智云是一套专业的生态养殖管理平台，为养殖户提供牲畜管理、疫苗接种、疾病防控、饲料追踪等全流程数字化解决方案 |
| login_features | ["牲畜追溯", "智能喂养", "疫苗管理", "疾病预警"] |

---

## 十、附录

### 10.1 图标推荐

| 模块 | 推荐图标 |
|------|----------|
| 牲畜分类 | AppstoreOutlined |
| 养殖区域 | HomeOutlined |
| 牲畜批次 | UnorderedListOutlined |
| 饲料信息 | ShoppingOutlined |
| 喂养记录 | ScheduleOutlined |
| 疫苗接种 | MedicineBoxOutlined |
| 疾病记录 | AlertOutlined |
| 养殖知识 | ReadOutlined |
| 认证审核 | SafetyCertificateOutlined |
| 公告信息 | NotificationOutlined |

### 10.2 后续扩展建议

1. **数据统计看板** - 养殖数据可视化
2. **预警提醒** - 疫苗到期、库存不足提醒
3. **养殖日志** - 每日养殖工作记录
4. **成本核算** - 养殖成本和收益分析
5. **出栏管理** - 牲畜出栏和销售记录
6. **消息通知** - 站内消息和短信通知
