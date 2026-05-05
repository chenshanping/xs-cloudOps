/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80034 (8.0.34)
 Source Host           : localhost:3306
 Source Schema         : go-base

 Target Server Type    : MySQL
 Target Server Version : 80034 (8.0.34)
 File Encoding         : 65001

 Date: 05/05/2026 23:10:04
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ai_conversations
-- ----------------------------
DROP TABLE IF EXISTS `ai_conversations`;
CREATE TABLE `ai_conversations`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `user_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '用户ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '对话标题',
  `model` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '使用的模型',
  `context_cleared_at` datetime(3) NULL DEFAULT NULL COMMENT '上下文清空时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_ai_conversations_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_ai_conversations_user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of ai_conversations
-- ----------------------------

-- ----------------------------
-- Table structure for ai_messages
-- ----------------------------
DROP TABLE IF EXISTS `ai_messages`;
CREATE TABLE `ai_messages`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `conversation_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '对话ID',
  `role` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '角色(user/assistant/system)',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '消息内容',
  `reasoning_content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '思考过程',
  `file_ids` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '附件文件ID(JSON数组)',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_ai_messages_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_ai_messages_conversation_id`(`conversation_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of ai_messages
-- ----------------------------

-- ----------------------------
-- Table structure for ai_providers
-- ----------------------------
DROP TABLE IF EXISTS `ai_providers`;
CREATE TABLE `ai_providers`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '平台名称',
  `api_key` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '平台API Key',
  `base_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '平台Base URL',
  `models_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '模型配置JSON',
  `is_default` tinyint(1) NULL DEFAULT 0 COMMENT '是否默认平台',
  `sort` bigint NULL DEFAULT 0 COMMENT '排序',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_ai_providers_name`(`name` ASC) USING BTREE,
  INDEX `idx_ai_providers_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 38 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ai_providers
-- ----------------------------
INSERT INTO `ai_providers` VALUES (35, '2026-05-01 23:30:23.996', '2026-05-01 23:30:23.996', NULL, 'deepseek', 'sk-e6c0478a770a49ad8650884ce456c801', 'https://api.deepseek.com', '[{\"id\":\"deepseek-v4-flash\",\"name\":\"deepseek-v4-flash\",\"description\":\"\",\"is_thinking\":false,\"support_vision\":false,\"support_tools\":true,\"search_strategy\":\"none\",\"support_embedding\":false,\"support_rerank\":false,\"is_free\":false,\"tags\":[\"tool\"]},{\"id\":\"deepseek-v4-pro\",\"name\":\"deepseek-v4-pro\",\"description\":\"\",\"is_thinking\":false,\"support_vision\":false,\"support_tools\":true,\"search_strategy\":\"none\",\"support_embedding\":false,\"support_rerank\":false,\"is_free\":false,\"tags\":[\"tool\"]}]', 0, 0);
INSERT INTO `ai_providers` VALUES (36, '2026-05-01 23:30:23.996', '2026-05-01 23:30:23.996', NULL, '阿里云百炼', 'sk-93a09f74fb7e49dc951d8b30b06e5a04', 'https://dashscope.aliyuncs.com/compatible-mode/v1', '[{\"id\":\"deepseek-v4-flash\",\"name\":\"deepseek-v4-flash\",\"description\":\"\",\"is_thinking\":false,\"support_vision\":false,\"support_tools\":true,\"search_strategy\":\"none\",\"support_embedding\":false,\"support_rerank\":false,\"is_free\":false,\"tags\":[\"tool\"]},{\"id\":\"deepseek-v4-pro\",\"name\":\"deepseek-v4-pro\",\"description\":\"\",\"is_thinking\":false,\"support_vision\":false,\"support_tools\":true,\"search_strategy\":\"none\",\"support_embedding\":false,\"support_rerank\":false,\"is_free\":false,\"tags\":[\"tool\"]},{\"id\":\"qwen3.6-plus-2026-04-02\",\"name\":\"qwen3.6-plus-2026-04-02\",\"description\":\"\",\"is_thinking\":true,\"support_vision\":true,\"support_tools\":true,\"search_strategy\":\"none\",\"support_embedding\":false,\"support_rerank\":false,\"is_free\":false,\"tags\":[\"reasoning\",\"vision\",\"tool\"]}]', 0, 1);
INSERT INTO `ai_providers` VALUES (37, '2026-05-01 23:30:23.996', '2026-05-01 23:30:23.996', NULL, '小米', 'sk-css8r1opz28oue1zatk62rob40woewcy9q90xkpqutw8ylbd', 'https://api.xiaomimimo.com/v1', '[{\"id\":\"mimo-v2-flash\",\"name\":\"mimo-v2-flash\",\"description\":\"\",\"is_thinking\":false,\"support_vision\":false,\"support_tools\":true,\"search_strategy\":\"none\",\"support_embedding\":false,\"support_rerank\":false,\"is_free\":false,\"tags\":[\"tool\"]},{\"id\":\"mimo-v2-omni\",\"name\":\"mimo-v2-omni\",\"description\":\"\",\"is_thinking\":false,\"support_vision\":true,\"support_tools\":true,\"search_strategy\":\"none\",\"support_embedding\":false,\"support_rerank\":false,\"is_free\":false,\"tags\":[\"vision\",\"tool\"]},{\"id\":\"mimo-v2.5\",\"name\":\"mimo-v2.5\",\"description\":\"\",\"is_thinking\":false,\"support_vision\":true,\"support_tools\":true,\"search_strategy\":\"builtin\",\"support_embedding\":false,\"support_rerank\":false,\"is_free\":false,\"tags\":[\"vision\",\"search\",\"tool\"]}]', 1, 2);

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v0` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v1` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v2` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v3` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v4` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `v5` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_casbin_rule`(`ptype` ASC, `v0` ASC, `v1` ASC, `v2` ASC, `v3` ASC, `v4` ASC, `v5` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 294 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
INSERT INTO `casbin_rule` VALUES (156, 'p', 'admin', '/api/v1/ai/chat', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (157, 'p', 'admin', '/api/v1/ai/chat/stream', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (204, 'p', 'admin', '/api/v1/ai/config', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (205, 'p', 'admin', '/api/v1/ai/config', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (149, 'p', 'admin', '/api/v1/ai/conversations', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (150, 'p', 'admin', '/api/v1/ai/conversations', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (153, 'p', 'admin', '/api/v1/ai/conversations/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (151, 'p', 'admin', '/api/v1/ai/conversations/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (152, 'p', 'admin', '/api/v1/ai/conversations/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (159, 'p', 'admin', '/api/v1/ai/conversations/:id/clear-context', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (155, 'p', 'admin', '/api/v1/ai/conversations/:id/messages', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (154, 'p', 'admin', '/api/v1/ai/conversations/:id/messages', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (158, 'p', 'admin', '/api/v1/ai/messages/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (148, 'p', 'admin', '/api/v1/ai/models', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (202, 'p', 'admin', '/api/v1/ai/providers/models/fetch', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (174, 'p', 'admin', '/api/v1/ai/test', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (110, 'p', 'admin', '/api/v1/apis', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (112, 'p', 'admin', '/api/v1/apis', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (114, 'p', 'admin', '/api/v1/apis/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (111, 'p', 'admin', '/api/v1/apis/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (113, 'p', 'admin', '/api/v1/apis/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (125, 'p', 'admin', '/api/v1/apis/all', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (126, 'p', 'admin', '/api/v1/apis/groups', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (127, 'p', 'admin', '/api/v1/apis/sync', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (117, 'p', 'admin', '/api/v1/auth/login', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (118, 'p', 'admin', '/api/v1/auth/logout', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (119, 'p', 'admin', '/api/v1/auth/refresh', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (160, 'p', 'admin', '/api/v1/auth/register', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (162, 'p', 'admin', '/api/v1/auth/reset-password', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (166, 'p', 'admin', '/api/v1/auth/reset-password-by-email', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (165, 'p', 'admin', '/api/v1/auth/reset-password-by-username', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (161, 'p', 'admin', '/api/v1/auth/send-email-code', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (120, 'p', 'admin', '/api/v1/auth/userinfo', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (163, 'p', 'admin', '/api/v1/captcha', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (164, 'p', 'admin', '/api/v1/captcha/config', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (191, 'p', 'admin', '/api/v1/captcha/slider', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (192, 'p', 'admin', '/api/v1/captcha/slider/verify', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (128, 'p', 'admin', '/api/v1/configs', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (131, 'p', 'admin', '/api/v1/configs', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (134, 'p', 'admin', '/api/v1/configs/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (132, 'p', 'admin', '/api/v1/configs/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (133, 'p', 'admin', '/api/v1/configs/batch', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (129, 'p', 'admin', '/api/v1/configs/key/:key', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (130, 'p', 'admin', '/api/v1/configs/keys', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (167, 'p', 'admin', '/api/v1/configs/test-email', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (198, 'p', 'admin', '/api/v1/depts', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (200, 'p', 'admin', '/api/v1/depts/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (197, 'p', 'admin', '/api/v1/depts/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (199, 'p', 'admin', '/api/v1/depts/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (201, 'p', 'admin', '/api/v1/depts/manageable-tree', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (196, 'p', 'admin', '/api/v1/depts/tree', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (186, 'p', 'admin', '/api/v1/dict/data', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (188, 'p', 'admin', '/api/v1/dict/data', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (190, 'p', 'admin', '/api/v1/dict/data/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (187, 'p', 'admin', '/api/v1/dict/data/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (189, 'p', 'admin', '/api/v1/dict/data/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (179, 'p', 'admin', '/api/v1/dict/type/:type', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (180, 'p', 'admin', '/api/v1/dict/types', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (183, 'p', 'admin', '/api/v1/dict/types', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (185, 'p', 'admin', '/api/v1/dict/types/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (182, 'p', 'admin', '/api/v1/dict/types/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (184, 'p', 'admin', '/api/v1/dict/types/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (181, 'p', 'admin', '/api/v1/dict/types/all', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (171, 'p', 'admin', '/api/v1/echart/role-status-stats', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (172, 'p', 'admin', '/api/v1/echart/user-register-trend', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (169, 'p', 'admin', '/api/v1/echart/user-role-stats', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (170, 'p', 'admin', '/api/v1/echart/user-status-stats', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (135, 'p', 'admin', '/api/v1/files', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (137, 'p', 'admin', '/api/v1/files/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (136, 'p', 'admin', '/api/v1/files/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (193, 'p', 'admin', '/api/v1/files/batch', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (139, 'p', 'admin', '/api/v1/files/check-md5', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (138, 'p', 'admin', '/api/v1/files/credential', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (144, 'p', 'admin', '/api/v1/files/multipart/abort', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (143, 'p', 'admin', '/api/v1/files/multipart/complete', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (141, 'p', 'admin', '/api/v1/files/multipart/init', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (142, 'p', 'admin', '/api/v1/files/multipart/parts', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (140, 'p', 'admin', '/api/v1/files/save', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (146, 'p', 'admin', '/api/v1/files/upload/chunk', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (145, 'p', 'admin', '/api/v1/files/upload/local', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (116, 'p', 'admin', '/api/v1/logs/login', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (115, 'p', 'admin', '/api/v1/logs/operation', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (168, 'p', 'admin', '/api/v1/logs/route-groups', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (105, 'p', 'admin', '/api/v1/menus', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (107, 'p', 'admin', '/api/v1/menus', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (109, 'p', 'admin', '/api/v1/menus/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (106, 'p', 'admin', '/api/v1/menus/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (108, 'p', 'admin', '/api/v1/menus/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (98, 'p', 'admin', '/api/v1/roles', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (100, 'p', 'admin', '/api/v1/roles', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (102, 'p', 'admin', '/api/v1/roles/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (99, 'p', 'admin', '/api/v1/roles/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (101, 'p', 'admin', '/api/v1/roles/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (104, 'p', 'admin', '/api/v1/roles/:id/apis', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (206, 'p', 'admin', '/api/v1/roles/:id/data-scopes', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (103, 'p', 'admin', '/api/v1/roles/:id/menus', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (207, 'p', 'admin', '/api/v1/roles/:id/permissions', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (208, 'p', 'admin', '/api/v1/roles/data-scope-resources', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (147, 'p', 'admin', '/api/v1/user/avatar', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (124, 'p', 'admin', '/api/v1/user/menus', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (123, 'p', 'admin', '/api/v1/user/password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (121, 'p', 'admin', '/api/v1/user/profile', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (122, 'p', 'admin', '/api/v1/user/profile', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (175, 'p', 'admin', '/api/v1/user/profiles', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (176, 'p', 'admin', '/api/v1/user/profiles/types', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (91, 'p', 'admin', '/api/v1/users', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (93, 'p', 'admin', '/api/v1/users', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (95, 'p', 'admin', '/api/v1/users/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (92, 'p', 'admin', '/api/v1/users/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (94, 'p', 'admin', '/api/v1/users/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (177, 'p', 'admin', '/api/v1/users/:id/offline', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (97, 'p', 'admin', '/api/v1/users/:id/password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (178, 'p', 'admin', '/api/v1/users/:id/profiles', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (96, 'p', 'admin', '/api/v1/users/:id/status', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (194, 'p', 'admin', '/api/v1/users/batch', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (203, 'p', 'admin', '/api/v1/users/batch-password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (195, 'p', 'admin', '/api/v1/users/batch-status', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (293, 'p', 'admin', '/api/v1/users/export', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (289, 'p', 'admin', '/api/v1/users/import', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (286, 'p', 'admin', '/api/v1/users/import-template', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (173, 'p', 'admin', '/api/v1/users/options', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (54, 'p', 'codex_dual_track_1777414199', '/api/v1/users', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (13, 'p', 'doctor', '/api/v1/ai/chat', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (14, 'p', 'doctor', '/api/v1/ai/chat/stream', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (6, 'p', 'doctor', '/api/v1/ai/conversations', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (7, 'p', 'doctor', '/api/v1/ai/conversations', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (10, 'p', 'doctor', '/api/v1/ai/conversations/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (8, 'p', 'doctor', '/api/v1/ai/conversations/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (9, 'p', 'doctor', '/api/v1/ai/conversations/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (16, 'p', 'doctor', '/api/v1/ai/conversations/:id/clear-context', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (12, 'p', 'doctor', '/api/v1/ai/conversations/:id/messages', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (11, 'p', 'doctor', '/api/v1/ai/conversations/:id/messages', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (15, 'p', 'doctor', '/api/v1/ai/messages/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (5, 'p', 'doctor', '/api/v1/ai/models', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (17, 'p', 'doctor', '/api/v1/psy_category', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (21, 'p', 'doctor', '/api/v1/psy_category/options', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (20, 'p', 'doctor', '/api/v1/psy_category/question_stats', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (22, 'p', 'doctor', '/api/v1/psy_paper', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (25, 'p', 'doctor', '/api/v1/psy_paper', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (27, 'p', 'doctor', '/api/v1/psy_paper/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (24, 'p', 'doctor', '/api/v1/psy_paper/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (26, 'p', 'doctor', '/api/v1/psy_paper/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (28, 'p', 'doctor', '/api/v1/psy_paper/batch', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (23, 'p', 'doctor', '/api/v1/psy_paper/options', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (18, 'p', 'doctor', '/api/v1/psy_question', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (19, 'p', 'doctor', '/api/v1/psy_question/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (4, 'p', 'doctor', '/api/v1/user/avatar', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (3, 'p', 'doctor', '/api/v1/user/password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (1, 'p', 'doctor', '/api/v1/user/profile', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (2, 'p', 'doctor', '/api/v1/user/profile', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (46, 'p', 'farmer_certification', '/api/v1/breeding_area', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (50, 'p', 'farmer_certification', '/api/v1/breeding_area', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (52, 'p', 'farmer_certification', '/api/v1/breeding_area/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (49, 'p', 'farmer_certification', '/api/v1/breeding_area/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (51, 'p', 'farmer_certification', '/api/v1/breeding_area/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (53, 'p', 'farmer_certification', '/api/v1/breeding_area/batch', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (48, 'p', 'farmer_certification', '/api/v1/breeding_area/creator/options', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (47, 'p', 'farmer_certification', '/api/v1/breeding_area/options', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (33, 'p', 'farmer_certification', '/api/v1/files/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (32, 'p', 'farmer_certification', '/api/v1/files/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (35, 'p', 'farmer_certification', '/api/v1/files/check-md5', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (34, 'p', 'farmer_certification', '/api/v1/files/credential', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (40, 'p', 'farmer_certification', '/api/v1/files/multipart/abort', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (39, 'p', 'farmer_certification', '/api/v1/files/multipart/complete', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (37, 'p', 'farmer_certification', '/api/v1/files/multipart/init', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (38, 'p', 'farmer_certification', '/api/v1/files/multipart/parts', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (36, 'p', 'farmer_certification', '/api/v1/files/save', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (42, 'p', 'farmer_certification', '/api/v1/files/upload/chunk', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (41, 'p', 'farmer_certification', '/api/v1/files/upload/local', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (43, 'p', 'farmer_certification', '/api/v1/user/avatar', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (31, 'p', 'farmer_certification', '/api/v1/user/password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (29, 'p', 'farmer_certification', '/api/v1/user/profile', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (30, 'p', 'farmer_certification', '/api/v1/user/profile', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (44, 'p', 'farmer_certification', '/api/v1/user/profiles', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (45, 'p', 'farmer_certification', '/api/v1/user/profiles/types', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (252, 'p', 'system_admin', '/api/v1/ai/chat', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (253, 'p', 'system_admin', '/api/v1/ai/chat/stream', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (275, 'p', 'system_admin', '/api/v1/ai/config', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (276, 'p', 'system_admin', '/api/v1/ai/config', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (245, 'p', 'system_admin', '/api/v1/ai/conversations', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (246, 'p', 'system_admin', '/api/v1/ai/conversations', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (249, 'p', 'system_admin', '/api/v1/ai/conversations/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (247, 'p', 'system_admin', '/api/v1/ai/conversations/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (248, 'p', 'system_admin', '/api/v1/ai/conversations/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (255, 'p', 'system_admin', '/api/v1/ai/conversations/:id/clear-context', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (251, 'p', 'system_admin', '/api/v1/ai/conversations/:id/messages', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (250, 'p', 'system_admin', '/api/v1/ai/conversations/:id/messages', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (254, 'p', 'system_admin', '/api/v1/ai/messages/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (244, 'p', 'system_admin', '/api/v1/ai/models', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (273, 'p', 'system_admin', '/api/v1/ai/providers/models/fetch', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (261, 'p', 'system_admin', '/api/v1/ai/test', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (217, 'p', 'system_admin', '/api/v1/auth/login', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (218, 'p', 'system_admin', '/api/v1/auth/logout', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (219, 'p', 'system_admin', '/api/v1/auth/refresh', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (220, 'p', 'system_admin', '/api/v1/auth/userinfo', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (224, 'p', 'system_admin', '/api/v1/configs', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (227, 'p', 'system_admin', '/api/v1/configs', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (230, 'p', 'system_admin', '/api/v1/configs/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (228, 'p', 'system_admin', '/api/v1/configs/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (229, 'p', 'system_admin', '/api/v1/configs/batch', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (225, 'p', 'system_admin', '/api/v1/configs/key/:key', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (226, 'p', 'system_admin', '/api/v1/configs/keys', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (269, 'p', 'system_admin', '/api/v1/depts', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (271, 'p', 'system_admin', '/api/v1/depts/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (268, 'p', 'system_admin', '/api/v1/depts/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (270, 'p', 'system_admin', '/api/v1/depts/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (272, 'p', 'system_admin', '/api/v1/depts/manageable-tree', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (267, 'p', 'system_admin', '/api/v1/depts/tree', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (258, 'p', 'system_admin', '/api/v1/echart/role-status-stats', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (259, 'p', 'system_admin', '/api/v1/echart/user-register-trend', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (256, 'p', 'system_admin', '/api/v1/echart/user-role-stats', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (257, 'p', 'system_admin', '/api/v1/echart/user-status-stats', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (231, 'p', 'system_admin', '/api/v1/files', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (233, 'p', 'system_admin', '/api/v1/files/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (232, 'p', 'system_admin', '/api/v1/files/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (235, 'p', 'system_admin', '/api/v1/files/check-md5', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (234, 'p', 'system_admin', '/api/v1/files/credential', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (240, 'p', 'system_admin', '/api/v1/files/multipart/abort', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (239, 'p', 'system_admin', '/api/v1/files/multipart/complete', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (237, 'p', 'system_admin', '/api/v1/files/multipart/init', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (238, 'p', 'system_admin', '/api/v1/files/multipart/parts', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (236, 'p', 'system_admin', '/api/v1/files/save', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (242, 'p', 'system_admin', '/api/v1/files/upload/chunk', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (241, 'p', 'system_admin', '/api/v1/files/upload/local', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (216, 'p', 'system_admin', '/api/v1/roles', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (243, 'p', 'system_admin', '/api/v1/user/avatar', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (223, 'p', 'system_admin', '/api/v1/user/password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (221, 'p', 'system_admin', '/api/v1/user/profile', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (222, 'p', 'system_admin', '/api/v1/user/profile', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (262, 'p', 'system_admin', '/api/v1/user/profiles', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (263, 'p', 'system_admin', '/api/v1/user/profiles/types', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (209, 'p', 'system_admin', '/api/v1/users', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (211, 'p', 'system_admin', '/api/v1/users', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (213, 'p', 'system_admin', '/api/v1/users/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (210, 'p', 'system_admin', '/api/v1/users/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (212, 'p', 'system_admin', '/api/v1/users/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (264, 'p', 'system_admin', '/api/v1/users/:id/offline', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (215, 'p', 'system_admin', '/api/v1/users/:id/password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (265, 'p', 'system_admin', '/api/v1/users/:id/profiles', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (214, 'p', 'system_admin', '/api/v1/users/:id/status', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (274, 'p', 'system_admin', '/api/v1/users/batch-password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (266, 'p', 'system_admin', '/api/v1/users/batch-status', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (291, 'p', 'system_admin', '/api/v1/users/export', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (288, 'p', 'system_admin', '/api/v1/users/import', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (284, 'p', 'system_admin', '/api/v1/users/import-template', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (260, 'p', 'system_admin', '/api/v1/users/options', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (66, 'p', 'test', '/api/v1/configs', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (67, 'p', 'test', '/api/v1/configs/key/:key', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (86, 'p', 'test', '/api/v1/depts', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (88, 'p', 'test', '/api/v1/depts/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (85, 'p', 'test', '/api/v1/depts/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (87, 'p', 'test', '/api/v1/depts/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (89, 'p', 'test', '/api/v1/depts/manageable-tree', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (84, 'p', 'test', '/api/v1/depts/tree', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (69, 'p', 'test', '/api/v1/files/check-md5', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (68, 'p', 'test', '/api/v1/files/credential', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (74, 'p', 'test', '/api/v1/files/multipart/abort', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (73, 'p', 'test', '/api/v1/files/multipart/complete', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (71, 'p', 'test', '/api/v1/files/multipart/init', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (72, 'p', 'test', '/api/v1/files/multipart/parts', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (70, 'p', 'test', '/api/v1/files/save', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (76, 'p', 'test', '/api/v1/files/upload/chunk', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (75, 'p', 'test', '/api/v1/files/upload/local', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (62, 'p', 'test', '/api/v1/roles', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (77, 'p', 'test', '/api/v1/user/avatar', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (65, 'p', 'test', '/api/v1/user/password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (63, 'p', 'test', '/api/v1/user/profile', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (64, 'p', 'test', '/api/v1/user/profile', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (79, 'p', 'test', '/api/v1/user/profiles', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (80, 'p', 'test', '/api/v1/user/profiles/types', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (55, 'p', 'test', '/api/v1/users', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (57, 'p', 'test', '/api/v1/users', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (59, 'p', 'test', '/api/v1/users/:id', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (56, 'p', 'test', '/api/v1/users/:id', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (58, 'p', 'test', '/api/v1/users/:id', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (61, 'p', 'test', '/api/v1/users/:id/password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (81, 'p', 'test', '/api/v1/users/:id/profiles', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (60, 'p', 'test', '/api/v1/users/:id/status', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (82, 'p', 'test', '/api/v1/users/batch', 'DELETE', '', '', '');
INSERT INTO `casbin_rule` VALUES (90, 'p', 'test', '/api/v1/users/batch-password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (83, 'p', 'test', '/api/v1/users/batch-status', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (290, 'p', 'test', '/api/v1/users/export', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (287, 'p', 'test', '/api/v1/users/import', 'POST', '', '', '');
INSERT INTO `casbin_rule` VALUES (283, 'p', 'test', '/api/v1/users/import-template', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (78, 'p', 'test', '/api/v1/users/options', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (281, 'p', 'user', '/api/v1/user/avatar', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (280, 'p', 'user', '/api/v1/user/password', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (278, 'p', 'user', '/api/v1/user/profile', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (279, 'p', 'user', '/api/v1/user/profile', 'PUT', '', '', '');
INSERT INTO `casbin_rule` VALUES (282, 'p', 'user', '/api/v1/user/profiles', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (277, 'p', 'user', '/api/v1/users', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (292, 'p', 'user', '/api/v1/users/export', 'GET', '', '', '');
INSERT INTO `casbin_rule` VALUES (285, 'p', 'user', '/api/v1/users/import-template', 'GET', '', '', '');

-- ----------------------------
-- Table structure for sys_api
-- ----------------------------
DROP TABLE IF EXISTS `sys_api`;
CREATE TABLE `sys_api`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'API路径',
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '请求方法',
  `group` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'API分组',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '描述',
  `request_params` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '请求参数JSON',
  `response_params` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '响应参数JSON',
  `need_auth` tinyint(1) NULL DEFAULT NULL COMMENT '是否需要认证',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_api_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 247 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_api
-- ----------------------------
INSERT INTO `sys_api` VALUES (1, '2026-01-24 02:39:48.470', '2026-04-29 20:33:50.702', NULL, '/api/v1/users', 'GET', '用户管理', '用户列表', '[{\"name\":\"page\",\"type\":\"integer\",\"description\":\"页码\",\"required\":false,\"in\":\"query\"},{\"name\":\"page_size\",\"type\":\"integer\",\"description\":\"每页数量\",\"required\":false,\"in\":\"query\"},{\"name\":\"username\",\"type\":\"string\",\"description\":\"用户名\",\"required\":false,\"in\":\"query\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态\",\"required\":false,\"in\":\"query\"},{\"name\":\"gender\",\"type\":\"integer\",\"description\":\"性别\",\"required\":false,\"in\":\"query\"},{\"name\":\"role_id\",\"type\":\"integer\",\"description\":\"角色ID\",\"required\":false,\"in\":\"query\"},{\"name\":\"dept_id\",\"type\":\"integer\",\"description\":\"部门ID\",\"required\":false,\"in\":\"query\"},{\"name\":\"unassigned_dept\",\"type\":\"boolean\",\"description\":\"是否筛选未绑定部门用户\",\"required\":false,\"in\":\"query\"}]', '', 1);
INSERT INTO `sys_api` VALUES (2, '2026-01-24 02:39:48.470', '2026-02-01 01:49:32.743', NULL, '/api/v1/users/:id', 'GET', '用户管理', '用户详情', '', '', 1);
INSERT INTO `sys_api` VALUES (3, '2026-01-24 02:39:48.470', '2026-04-29 20:33:50.714', NULL, '/api/v1/users', 'POST', '用户管理', '创建用户', '[{\"name\":\"username\",\"type\":\"string\",\"description\":\"用户名\",\"required\":true,\"in\":\"body\"},{\"name\":\"password\",\"type\":\"string\",\"description\":\"密码\",\"required\":true,\"in\":\"body\"},{\"name\":\"nickname\",\"type\":\"string\",\"description\":\"昵称\",\"required\":false,\"in\":\"body\"},{\"name\":\"gender\",\"type\":\"integer\",\"description\":\"性别(0:未知,1:男,2:女)\",\"required\":false,\"in\":\"body\"},{\"name\":\"email\",\"type\":\"string\",\"description\":\"邮箱\",\"required\":false,\"in\":\"body\"},{\"name\":\"phone\",\"type\":\"string\",\"description\":\"手机号\",\"required\":false,\"in\":\"body\"},{\"name\":\"avatar\",\"type\":\"string\",\"description\":\"头像地址\",\"required\":false,\"in\":\"body\"},{\"name\":\"avatar_file_id\",\"type\":\"integer\",\"description\":\"头像文件ID\",\"required\":false,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(0:禁用,1:启用)\",\"required\":false,\"in\":\"body\"},{\"name\":\"dept_id\",\"type\":\"integer\",\"description\":\"部门ID\",\"required\":false,\"in\":\"body\"},{\"name\":\"role_ids\",\"type\":\"array[integer]\",\"description\":\"角色ID列表\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (4, '2026-01-24 02:39:48.470', '2026-04-29 20:33:50.720', NULL, '/api/v1/users/:id', 'PUT', '用户管理', '更新用户', '[{\"name\":\"nickname\",\"type\":\"string\",\"description\":\"昵称\",\"required\":false,\"in\":\"body\"},{\"name\":\"gender\",\"type\":\"integer\",\"description\":\"性别(0:未知,1:男,2:女)\",\"required\":false,\"in\":\"body\"},{\"name\":\"email\",\"type\":\"string\",\"description\":\"邮箱\",\"required\":false,\"in\":\"body\"},{\"name\":\"phone\",\"type\":\"string\",\"description\":\"手机号\",\"required\":false,\"in\":\"body\"},{\"name\":\"avatar\",\"type\":\"string\",\"description\":\"头像地址\",\"required\":false,\"in\":\"body\"},{\"name\":\"avatar_file_id\",\"type\":\"integer\",\"description\":\"头像文件ID\",\"required\":false,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(0:禁用,1:启用)\",\"required\":false,\"in\":\"body\"},{\"name\":\"dept_id\",\"type\":\"integer\",\"description\":\"部门ID\",\"required\":false,\"in\":\"body\"},{\"name\":\"role_ids\",\"type\":\"array[integer]\",\"description\":\"角色ID列表\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (5, '2026-01-24 02:39:48.470', '2026-02-01 01:49:32.881', NULL, '/api/v1/users/:id', 'DELETE', '用户管理', '删除用户', '', '', 1);
INSERT INTO `sys_api` VALUES (6, '2026-01-24 02:39:48.470', '2026-02-01 01:49:32.925', NULL, '/api/v1/users/:id/status', 'PUT', '用户管理', '修改用户状态', '', '', 1);
INSERT INTO `sys_api` VALUES (7, '2026-01-24 02:39:48.470', '2026-04-29 20:33:50.738', NULL, '/api/v1/users/:id/password', 'PUT', '用户管理', '重置密码', '', '', 1);
INSERT INTO `sys_api` VALUES (8, '2026-01-24 02:39:48.470', '2026-02-01 01:49:30.425', NULL, '/api/v1/roles', 'GET', '角色管理', '角色列表', '', '', 1);
INSERT INTO `sys_api` VALUES (9, '2026-01-24 02:39:48.470', '2026-02-01 01:49:30.539', NULL, '/api/v1/roles/:id', 'GET', '角色管理', '角色详情', '', '', 1);
INSERT INTO `sys_api` VALUES (10, '2026-01-24 02:39:48.470', '2026-05-05 01:15:49.804', NULL, '/api/v1/roles', 'POST', '角色管理', '创建角色', '[{\"name\":\"name\",\"type\":\"string\",\"description\":\"角色名称\",\"required\":true,\"in\":\"body\"},{\"name\":\"code\",\"type\":\"string\",\"description\":\"角色编码\",\"required\":true,\"in\":\"body\"},{\"name\":\"sort\",\"type\":\"integer\",\"description\":\"排序\",\"required\":false,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(0:禁用,1:启用)\",\"required\":false,\"in\":\"body\"},{\"name\":\"is_super_admin\",\"type\":\"boolean\",\"description\":\"是否显式超管\",\"required\":false,\"in\":\"body\"},{\"name\":\"data_scope\",\"type\":\"integer\",\"description\":\"数据范围\",\"required\":false,\"in\":\"body\"},{\"name\":\"remark\",\"type\":\"string\",\"description\":\"备注\",\"required\":false,\"in\":\"body\"},{\"name\":\"dept_ids\",\"type\":\"array[integer]\",\"description\":\"自定义数据范围部门ID列表\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (11, '2026-01-24 02:39:48.470', '2026-05-05 01:15:49.809', NULL, '/api/v1/roles/:id', 'PUT', '角色管理', '更新角色', '[{\"name\":\"name\",\"type\":\"string\",\"description\":\"角色名称\",\"required\":false,\"in\":\"body\"},{\"name\":\"code\",\"type\":\"string\",\"description\":\"角色编码\",\"required\":false,\"in\":\"body\"},{\"name\":\"sort\",\"type\":\"integer\",\"description\":\"排序\",\"required\":false,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(0:禁用,1:启用)\",\"required\":false,\"in\":\"body\"},{\"name\":\"is_super_admin\",\"type\":\"boolean\",\"description\":\"是否显式超管\",\"required\":false,\"in\":\"body\"},{\"name\":\"data_scope\",\"type\":\"integer\",\"description\":\"数据范围\",\"required\":false,\"in\":\"body\"},{\"name\":\"remark\",\"type\":\"string\",\"description\":\"备注\",\"required\":false,\"in\":\"body\"},{\"name\":\"dept_ids\",\"type\":\"array[integer]\",\"description\":\"自定义数据范围部门ID列表\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (12, '2026-01-24 02:39:48.470', '2026-02-01 01:49:30.767', NULL, '/api/v1/roles/:id', 'DELETE', '角色管理', '删除角色', '', '', 1);
INSERT INTO `sys_api` VALUES (13, '2026-01-24 02:39:48.470', '2026-02-01 01:49:30.834', NULL, '/api/v1/roles/:id/menus', 'PUT', '角色管理', '分配菜单', '[{\"name\":\"menu_ids\",\"type\":\"array[integer]\",\"description\":\"菜单ID列表\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (14, '2026-01-24 02:39:48.470', '2026-02-01 01:49:30.908', NULL, '/api/v1/roles/:id/apis', 'PUT', '角色管理', '分配API', '[{\"name\":\"api_ids\",\"type\":\"array[integer]\",\"description\":\"API ID列表\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (15, '2026-01-24 02:39:48.470', '2026-02-01 01:49:30.000', NULL, '/api/v1/menus', 'GET', '菜单管理', '菜单列表', '', '', 1);
INSERT INTO `sys_api` VALUES (16, '2026-01-24 02:39:48.470', '2026-02-01 01:49:30.042', NULL, '/api/v1/menus/:id', 'GET', '菜单管理', '菜单详情', '', '', 1);
INSERT INTO `sys_api` VALUES (17, '2026-01-24 02:39:48.470', '2026-02-01 01:49:30.174', NULL, '/api/v1/menus', 'POST', '菜单管理', '创建菜单', '[{\"name\":\"parent_id\",\"type\":\"integer\",\"description\":\"父菜单ID\",\"required\":false,\"in\":\"body\"},{\"name\":\"name\",\"type\":\"string\",\"description\":\"菜单名称\",\"required\":true,\"in\":\"body\"},{\"name\":\"path\",\"type\":\"string\",\"description\":\"路由路径\",\"required\":false,\"in\":\"body\"},{\"name\":\"component\",\"type\":\"string\",\"description\":\"组件路径\",\"required\":false,\"in\":\"body\"},{\"name\":\"icon\",\"type\":\"string\",\"description\":\"图标\",\"required\":false,\"in\":\"body\"},{\"name\":\"sort\",\"type\":\"integer\",\"description\":\"排序\",\"required\":false,\"in\":\"body\"},{\"name\":\"type\",\"type\":\"integer\",\"description\":\"类型(1:目录,2:菜单,3:按钮)\",\"required\":true,\"in\":\"body\"},{\"name\":\"permission\",\"type\":\"string\",\"description\":\"权限标识\",\"required\":false,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(0:禁用,1:启用)\",\"required\":false,\"in\":\"body\"},{\"name\":\"hidden\",\"type\":\"integer\",\"description\":\"是否隐藏(0:显示,1:隐藏)\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (18, '2026-01-24 02:39:48.470', '2026-02-01 01:49:30.250', NULL, '/api/v1/menus/:id', 'PUT', '菜单管理', '更新菜单', '[{\"name\":\"parent_id\",\"type\":\"integer\",\"description\":\"父菜单ID\",\"required\":false,\"in\":\"body\"},{\"name\":\"name\",\"type\":\"string\",\"description\":\"菜单名称\",\"required\":false,\"in\":\"body\"},{\"name\":\"path\",\"type\":\"string\",\"description\":\"路由路径\",\"required\":false,\"in\":\"body\"},{\"name\":\"component\",\"type\":\"string\",\"description\":\"组件路径\",\"required\":false,\"in\":\"body\"},{\"name\":\"icon\",\"type\":\"string\",\"description\":\"图标\",\"required\":false,\"in\":\"body\"},{\"name\":\"sort\",\"type\":\"integer\",\"description\":\"排序\",\"required\":false,\"in\":\"body\"},{\"name\":\"type\",\"type\":\"integer\",\"description\":\"类型(1:目录,2:菜单,3:按钮)\",\"required\":false,\"in\":\"body\"},{\"name\":\"permission\",\"type\":\"string\",\"description\":\"权限标识\",\"required\":false,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(0:禁用,1:启用)\",\"required\":false,\"in\":\"body\"},{\"name\":\"hidden\",\"type\":\"integer\",\"description\":\"是否隐藏(0:显示,1:隐藏)\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (19, '2026-01-24 02:39:48.470', '2026-02-01 01:49:30.292', NULL, '/api/v1/menus/:id', 'DELETE', '菜单管理', '删除菜单', '', '', 1);
INSERT INTO `sys_api` VALUES (20, '2026-01-24 02:39:48.470', '2026-02-01 01:49:31.456', NULL, '/api/v1/apis', 'GET', 'API管理', 'API列表', '[{\"name\":\"page\",\"type\":\"integer\",\"description\":\"页码\",\"required\":false,\"in\":\"query\"},{\"name\":\"page_size\",\"type\":\"integer\",\"description\":\"每页数量\",\"required\":false,\"in\":\"query\"},{\"name\":\"path\",\"type\":\"string\",\"description\":\"API路径\",\"required\":false,\"in\":\"query\"},{\"name\":\"method\",\"type\":\"string\",\"description\":\"请求方法\",\"required\":false,\"in\":\"query\"},{\"name\":\"group\",\"type\":\"string\",\"description\":\"API分组\",\"required\":false,\"in\":\"query\"}]', '', 1);
INSERT INTO `sys_api` VALUES (21, '2026-01-24 02:39:48.470', '2026-02-01 01:49:31.625', NULL, '/api/v1/apis/:id', 'GET', 'API管理', 'API详情', '', '', 1);
INSERT INTO `sys_api` VALUES (22, '2026-01-24 02:39:48.470', '2026-02-01 01:49:31.690', NULL, '/api/v1/apis', 'POST', 'API管理', '创建API', '[{\"name\":\"path\",\"type\":\"string\",\"description\":\"API路径\",\"required\":true,\"in\":\"body\"},{\"name\":\"method\",\"type\":\"string\",\"description\":\"请求方法\",\"required\":true,\"in\":\"body\"},{\"name\":\"group\",\"type\":\"string\",\"description\":\"API分组\",\"required\":false,\"in\":\"body\"},{\"name\":\"description\",\"type\":\"string\",\"description\":\"描述\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (23, '2026-01-24 02:39:48.470', '2026-02-01 01:49:31.767', NULL, '/api/v1/apis/:id', 'PUT', 'API管理', '更新API', '[{\"name\":\"path\",\"type\":\"string\",\"description\":\"API路径\",\"required\":false,\"in\":\"body\"},{\"name\":\"method\",\"type\":\"string\",\"description\":\"请求方法\",\"required\":false,\"in\":\"body\"},{\"name\":\"group\",\"type\":\"string\",\"description\":\"API分组\",\"required\":false,\"in\":\"body\"},{\"name\":\"description\",\"type\":\"string\",\"description\":\"描述\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (24, '2026-01-24 02:39:48.470', '2026-02-01 01:49:31.808', NULL, '/api/v1/apis/:id', 'DELETE', 'API管理', '删除API', '', '', 1);
INSERT INTO `sys_api` VALUES (25, '2026-01-24 02:39:48.470', '2026-02-01 20:46:48.177', NULL, '/api/v1/logs/operation', 'GET', '日志管理', '操作日志列表', '[{\"name\":\"page\",\"type\":\"integer\",\"description\":\"页码\",\"required\":false,\"in\":\"query\"},{\"name\":\"page_size\",\"type\":\"integer\",\"description\":\"每页数量\",\"required\":false,\"in\":\"query\"},{\"name\":\"username\",\"type\":\"string\",\"description\":\"用户名\",\"required\":false,\"in\":\"query\"},{\"name\":\"method\",\"type\":\"string\",\"description\":\"请求方法\",\"required\":false,\"in\":\"query\"},{\"name\":\"path\",\"type\":\"string\",\"description\":\"请求路径\",\"required\":false,\"in\":\"query\"},{\"name\":\"group\",\"type\":\"string\",\"description\":\"路由分组\",\"required\":false,\"in\":\"query\"},{\"name\":\"summary\",\"type\":\"string\",\"description\":\"路由描述\",\"required\":false,\"in\":\"query\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"HTTP状态码\",\"required\":false,\"in\":\"query\"},{\"name\":\"business_code\",\"type\":\"integer\",\"description\":\"业务状态码\",\"required\":false,\"in\":\"query\"},{\"name\":\"start_time\",\"type\":\"string\",\"description\":\"开始时间\",\"required\":false,\"in\":\"query\"},{\"name\":\"end_time\",\"type\":\"string\",\"description\":\"结束时间\",\"required\":false,\"in\":\"query\"},{\"name\":\"sort_field\",\"type\":\"string\",\"description\":\"排序字段\",\"required\":false,\"in\":\"query\"},{\"name\":\"sort_order\",\"type\":\"string\",\"description\":\"排序方式(ascend/descend)\",\"required\":false,\"in\":\"query\"}]', '', 1);
INSERT INTO `sys_api` VALUES (26, '2026-01-24 02:39:48.470', '2026-02-01 20:46:48.252', NULL, '/api/v1/logs/login', 'GET', '日志管理', '登录日志列表', '[{\"name\":\"page\",\"type\":\"integer\",\"description\":\"页码\",\"required\":false,\"in\":\"query\"},{\"name\":\"page_size\",\"type\":\"integer\",\"description\":\"每页数量\",\"required\":false,\"in\":\"query\"},{\"name\":\"username\",\"type\":\"string\",\"description\":\"用户名\",\"required\":false,\"in\":\"query\"},{\"name\":\"method\",\"type\":\"string\",\"description\":\"请求方法\",\"required\":false,\"in\":\"query\"},{\"name\":\"path\",\"type\":\"string\",\"description\":\"请求路径\",\"required\":false,\"in\":\"query\"},{\"name\":\"group\",\"type\":\"string\",\"description\":\"路由分组\",\"required\":false,\"in\":\"query\"},{\"name\":\"summary\",\"type\":\"string\",\"description\":\"路由描述\",\"required\":false,\"in\":\"query\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"HTTP状态码\",\"required\":false,\"in\":\"query\"},{\"name\":\"business_code\",\"type\":\"integer\",\"description\":\"业务状态码\",\"required\":false,\"in\":\"query\"},{\"name\":\"start_time\",\"type\":\"string\",\"description\":\"开始时间\",\"required\":false,\"in\":\"query\"},{\"name\":\"end_time\",\"type\":\"string\",\"description\":\"结束时间\",\"required\":false,\"in\":\"query\"},{\"name\":\"sort_field\",\"type\":\"string\",\"description\":\"排序字段\",\"required\":false,\"in\":\"query\"},{\"name\":\"sort_order\",\"type\":\"string\",\"description\":\"排序方式(ascend/descend)\",\"required\":false,\"in\":\"query\"}]', '', 1);
INSERT INTO `sys_api` VALUES (27, '2026-01-25 00:44:22.483', '2026-02-01 01:49:27.691', NULL, '/api/v1/auth/login', 'POST', '认证管理', '登录', '[{\"name\":\"username\",\"type\":\"string\",\"description\":\"用户名\",\"required\":true,\"in\":\"body\"},{\"name\":\"password\",\"type\":\"string\",\"description\":\"密码\",\"required\":true,\"in\":\"body\"},{\"name\":\"captcha_id\",\"type\":\"string\",\"description\":\"验证码ID\",\"required\":false,\"in\":\"body\"},{\"name\":\"captcha\",\"type\":\"string\",\"description\":\"验证码\",\"required\":false,\"in\":\"body\"}]', '', 0);
INSERT INTO `sys_api` VALUES (28, '2026-01-25 00:44:22.544', '2026-02-01 01:49:27.986', NULL, '/api/v1/auth/logout', 'POST', '认证管理', '登出', '', '', 1);
INSERT INTO `sys_api` VALUES (29, '2026-01-25 00:44:22.585', '2026-02-07 00:16:47.156', NULL, '/api/v1/auth/refresh', 'POST', '认证管理', '刷新Token', '', '', 0);
INSERT INTO `sys_api` VALUES (30, '2026-01-25 00:44:22.641', '2026-02-01 01:49:28.073', NULL, '/api/v1/auth/userinfo', 'GET', '认证管理', '获取用户信息', '', '', 1);
INSERT INTO `sys_api` VALUES (31, '2026-01-25 00:44:22.684', '2026-02-01 01:49:32.226', NULL, '/api/v1/user/profile', 'GET', '个人中心', '获取个人资料', '', '', 1);
INSERT INTO `sys_api` VALUES (32, '2026-01-25 00:44:22.819', '2026-02-01 01:49:32.268', NULL, '/api/v1/user/profile', 'PUT', '个人中心', '更新个人资料', '[{\"name\":\"nickname\",\"type\":\"string\",\"description\":\"昵称\",\"required\":false,\"in\":\"body\"},{\"name\":\"email\",\"type\":\"string\",\"description\":\"邮箱\",\"required\":false,\"in\":\"body\"},{\"name\":\"phone\",\"type\":\"string\",\"description\":\"手机号\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (33, '2026-01-25 00:44:23.069', '2026-02-01 01:49:32.368', NULL, '/api/v1/user/password', 'PUT', '个人中心', '修改密码', '[{\"name\":\"old_password\",\"type\":\"string\",\"description\":\"旧密码\",\"required\":true,\"in\":\"body\"},{\"name\":\"new_password\",\"type\":\"string\",\"description\":\"新密码\",\"required\":true,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (37, '2026-01-25 00:44:23.309', '2026-02-01 01:49:30.108', NULL, '/api/v1/user/menus', 'GET', '菜单管理', '获取用户菜单', '', '', 1);
INSERT INTO `sys_api` VALUES (38, '2026-01-25 00:44:23.354', '2026-02-01 01:49:31.500', NULL, '/api/v1/apis/all', 'GET', 'API管理', '全部API', '', '', 1);
INSERT INTO `sys_api` VALUES (39, '2026-01-25 00:44:23.409', '2026-02-01 01:49:31.581', NULL, '/api/v1/apis/groups', 'GET', 'API管理', 'API分组', '', '', 1);
INSERT INTO `sys_api` VALUES (40, '2026-01-25 00:44:23.454', '2026-02-01 01:49:31.850', NULL, '/api/v1/apis/sync', 'POST', 'API管理', '同步API', '', '', 1);
INSERT INTO `sys_api` VALUES (42, '2026-01-25 00:44:23.551', '2026-04-29 20:37:13.839', NULL, '/api/v1/configs', 'GET', '系统配置管理', '配置列表', '', '', 1);
INSERT INTO `sys_api` VALUES (43, '2026-01-25 00:44:23.602', '2026-04-29 20:37:13.853', NULL, '/api/v1/configs/key/:key', 'GET', '系统配置管理', '根据key获取配置', '', '', 1);
INSERT INTO `sys_api` VALUES (44, '2026-01-25 00:44:23.661', '2026-05-05 01:15:49.682', NULL, '/api/v1/configs/keys', 'POST', '系统配置管理', '批量获取公开配置', '[{\"name\":\"keys\",\"type\":\"array[string]\",\"description\":\"配置键列表\",\"required\":true,\"in\":\"body\"}]', '', 0);
INSERT INTO `sys_api` VALUES (45, '2026-01-25 00:44:23.702', '2026-04-29 20:37:13.868', NULL, '/api/v1/configs', 'POST', '系统配置管理', '创建配置', '', '', 1);
INSERT INTO `sys_api` VALUES (46, '2026-01-25 00:44:23.744', '2026-04-29 20:37:13.886', NULL, '/api/v1/configs/:id', 'PUT', '系统配置管理', '更新配置', '', '', 1);
INSERT INTO `sys_api` VALUES (47, '2026-01-25 00:44:23.794', '2026-04-29 20:37:13.901', NULL, '/api/v1/configs/batch', 'PUT', '系统配置管理', '批量更新配置', '', '', 1);
INSERT INTO `sys_api` VALUES (49, '2026-01-25 00:44:57.920', '2026-04-29 20:37:13.914', NULL, '/api/v1/configs/:id', 'DELETE', '系统配置管理', '删除配置', '', '', 1);
INSERT INTO `sys_api` VALUES (50, '2026-01-25 01:39:13.641', '2026-02-01 01:49:29.058', '2026-04-29 20:33:50.339', '/api/v1/generator/tables', 'GET', '代码生成', '获取数据库表列表', '', '', 1);
INSERT INTO `sys_api` VALUES (51, '2026-01-25 01:39:13.682', '2026-02-01 01:49:29.099', '2026-04-29 20:33:50.346', '/api/v1/generator/tables/:name/columns', 'GET', '代码生成', '获取表字段信息', '', '', 1);
INSERT INTO `sys_api` VALUES (52, '2026-01-25 01:39:13.724', '2026-02-01 01:49:29.175', '2026-04-29 20:33:50.351', '/api/v1/generator/preview', 'POST', '代码生成', '预览生成代码', '', '', 1);
INSERT INTO `sys_api` VALUES (53, '2026-01-25 01:39:13.782', '2026-02-01 01:49:29.216', '2026-04-29 20:33:50.356', '/api/v1/generator/generate', 'POST', '代码生成', '生成代码', '', '', 1);
INSERT INTO `sys_api` VALUES (54, '2026-01-25 01:39:13.973', '2026-02-01 01:49:29.258', '2026-04-29 20:33:50.360', '/api/v1/generator/modules', 'GET', '代码生成', '获取已生成模块', '', '', 1);
INSERT INTO `sys_api` VALUES (55, '2026-01-25 01:39:14.015', '2026-02-01 01:49:29.316', '2026-04-29 20:33:50.368', '/api/v1/generator/modules/:name', 'DELETE', '代码生成', '删除已生成模块', '', '', 1);
INSERT INTO `sys_api` VALUES (62, '2026-01-25 17:23:32.657', '2026-02-01 01:49:28.240', NULL, '/api/v1/files', 'GET', '文件管理', '文件列表', '', '', 1);
INSERT INTO `sys_api` VALUES (63, '2026-01-25 17:23:32.722', '2026-02-01 01:49:28.299', NULL, '/api/v1/files/:id', 'GET', '文件管理', '文件详情', '', '', 1);
INSERT INTO `sys_api` VALUES (64, '2026-01-25 17:23:32.779', '2026-02-01 01:49:28.357', NULL, '/api/v1/files/:id', 'DELETE', '文件管理', '删除文件', '', '', 1);
INSERT INTO `sys_api` VALUES (65, '2026-01-25 17:23:32.820', '2026-02-01 01:49:28.423', NULL, '/api/v1/files/credential', 'POST', '文件管理', '获取上传凭证', '', '', 1);
INSERT INTO `sys_api` VALUES (66, '2026-01-25 17:23:32.879', '2026-02-01 01:49:28.550', NULL, '/api/v1/files/check-md5', 'POST', '文件管理', 'MD5秒传检查', '', '', 1);
INSERT INTO `sys_api` VALUES (67, '2026-01-25 17:23:32.938', '2026-02-01 01:49:28.608', NULL, '/api/v1/files/save', 'POST', '文件管理', '保存上传文件', '', '', 1);
INSERT INTO `sys_api` VALUES (68, '2026-01-25 17:23:32.997', '2026-02-01 01:49:28.724', NULL, '/api/v1/files/multipart/init', 'POST', '文件管理', '初始化分片上传', '', '', 1);
INSERT INTO `sys_api` VALUES (69, '2026-01-25 17:23:33.072', '2026-02-01 01:49:28.800', NULL, '/api/v1/files/multipart/parts', 'GET', '文件管理', '获取已上传分片', '', '', 1);
INSERT INTO `sys_api` VALUES (70, '2026-01-25 17:23:33.130', '2026-02-01 01:49:28.841', NULL, '/api/v1/files/multipart/complete', 'POST', '文件管理', '完成分片上传', '', '', 1);
INSERT INTO `sys_api` VALUES (71, '2026-01-25 17:23:33.197', '2026-02-01 01:49:28.883', NULL, '/api/v1/files/multipart/abort', 'POST', '文件管理', '取消分片上传', '', '', 1);
INSERT INTO `sys_api` VALUES (72, '2026-01-25 17:23:33.254', '2026-02-01 01:49:28.933', NULL, '/api/v1/files/upload/local', 'POST', '文件管理', '本地文件上传', '', '', 1);
INSERT INTO `sys_api` VALUES (73, '2026-01-25 17:23:33.297', '2026-02-01 01:49:28.983', NULL, '/api/v1/files/upload/chunk', 'POST', '文件管理', '上传分片', '', '', 1);
INSERT INTO `sys_api` VALUES (81, '2026-01-25 23:44:08.234', '2026-02-01 01:49:32.322', NULL, '/api/v1/user/avatar', 'PUT', '个人中心', '更新头像', '', '', 1);
INSERT INTO `sys_api` VALUES (82, '2026-01-27 00:11:03.295', '2026-02-01 01:49:29.392', '2026-04-29 20:33:50.372', '/api/v1/generator/configs', 'POST', '代码生成', '保存配置', '', '', 1);
INSERT INTO `sys_api` VALUES (83, '2026-01-27 00:11:03.336', '2026-02-01 01:49:29.484', '2026-04-29 20:33:50.379', '/api/v1/generator/configs', 'GET', '代码生成', '获取配置列表', '', '', 1);
INSERT INTO `sys_api` VALUES (84, '2026-01-27 00:11:03.431', '2026-02-01 01:49:29.575', '2026-04-29 20:33:50.382', '/api/v1/generator/configs/:id', 'GET', '代码生成', '获取配置详情', '', '', 1);
INSERT INTO `sys_api` VALUES (85, '2026-01-27 00:11:03.470', '2026-02-01 01:49:29.625', '2026-04-29 20:33:50.386', '/api/v1/generator/configs/:id', 'DELETE', '代码生成', '删除配置', '', '', 1);
INSERT INTO `sys_api` VALUES (86, '2026-01-27 00:55:44.510', '2026-02-01 01:49:29.700', '2026-04-29 20:33:50.391', '/api/v1/generator/execute-sql', 'POST', '代码生成', '执行建表SQL', '', '', 1);
INSERT INTO `sys_api` VALUES (97, '2026-01-29 21:31:22.751', '2026-02-01 01:49:26.696', NULL, '/api/v1/ai/models', 'GET', 'AI对话', '获取模型列表', '', '', 0);
INSERT INTO `sys_api` VALUES (98, '2026-01-29 21:31:22.787', '2026-04-28 19:58:58.264', NULL, '/api/v1/ai/conversations', 'GET', 'AI对话', '获取对话列表', '[{\"name\":\"page\",\"type\":\"integer\",\"description\":\"页码\",\"required\":false,\"in\":\"query\"},{\"name\":\"page_size\",\"type\":\"integer\",\"description\":\"每页数量\",\"required\":false,\"in\":\"query\"}]', '', 1);
INSERT INTO `sys_api` VALUES (99, '2026-01-29 21:31:22.869', '2026-04-28 19:58:58.281', NULL, '/api/v1/ai/conversations', 'POST', 'AI对话', '创建对话', '[{\"name\":\"title\",\"type\":\"string\",\"description\":\"对话标题\",\"required\":false,\"in\":\"body\"},{\"name\":\"model\",\"type\":\"string\",\"description\":\"模型名称\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (100, '2026-01-29 21:31:22.937', '2026-04-28 19:58:58.293', NULL, '/api/v1/ai/conversations/:id', 'GET', 'AI对话', '获取对话详情', '', '', 1);
INSERT INTO `sys_api` VALUES (101, '2026-01-29 21:31:23.004', '2026-04-28 19:58:58.303', NULL, '/api/v1/ai/conversations/:id', 'PUT', 'AI对话', '更新对话标题', '', '', 1);
INSERT INTO `sys_api` VALUES (102, '2026-01-29 21:31:23.054', '2026-04-28 19:58:58.316', NULL, '/api/v1/ai/conversations/:id', 'DELETE', 'AI对话', '删除对话', '', '', 1);
INSERT INTO `sys_api` VALUES (103, '2026-01-29 21:31:23.104', '2026-04-28 19:58:58.330', NULL, '/api/v1/ai/conversations/:id/messages', 'GET', 'AI对话', '获取对话消息', '', '', 1);
INSERT INTO `sys_api` VALUES (104, '2026-01-29 21:31:23.154', '2026-04-28 19:58:58.346', NULL, '/api/v1/ai/conversations/:id/messages', 'DELETE', 'AI对话', '清空对话消息', '', '', 1);
INSERT INTO `sys_api` VALUES (105, '2026-01-29 21:31:23.220', '2026-04-28 19:58:58.395', NULL, '/api/v1/ai/chat', 'POST', 'AI对话', 'AI对话', '[{\"name\":\"conversation_id\",\"type\":\"integer\",\"description\":\"对话ID(为0则创建新对话)\",\"required\":false,\"in\":\"body\"},{\"name\":\"model\",\"type\":\"string\",\"description\":\"模型名称\",\"required\":false,\"in\":\"body\"},{\"name\":\"message\",\"type\":\"string\",\"description\":\"用户消息\",\"required\":true,\"in\":\"body\"},{\"name\":\"file_ids\",\"type\":\"array[integer]\",\"description\":\"附件文件ID列表\",\"required\":false,\"in\":\"body\"},{\"name\":\"enable_search\",\"type\":\"boolean\",\"description\":\"是否启用联网搜索\",\"required\":false,\"in\":\"body\"},{\"name\":\"enable_thinking\",\"type\":\"boolean\",\"description\":\"是否启用思考模式\",\"required\":false,\"in\":\"body\"},{\"name\":\"save_conversation\",\"type\":\"boolean\",\"description\":\"是否保存对话记录(默认true)\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (106, '2026-01-29 21:31:23.421', '2026-04-28 19:58:58.410', NULL, '/api/v1/ai/chat/stream', 'POST', 'AI对话', 'AI流式对话', '[{\"name\":\"conversation_id\",\"type\":\"integer\",\"description\":\"对话ID(为0则创建新对话)\",\"required\":false,\"in\":\"body\"},{\"name\":\"model\",\"type\":\"string\",\"description\":\"模型名称\",\"required\":false,\"in\":\"body\"},{\"name\":\"message\",\"type\":\"string\",\"description\":\"用户消息\",\"required\":true,\"in\":\"body\"},{\"name\":\"file_ids\",\"type\":\"array[integer]\",\"description\":\"附件文件ID列表\",\"required\":false,\"in\":\"body\"},{\"name\":\"enable_search\",\"type\":\"boolean\",\"description\":\"是否启用联网搜索\",\"required\":false,\"in\":\"body\"},{\"name\":\"enable_thinking\",\"type\":\"boolean\",\"description\":\"是否启用思考模式\",\"required\":false,\"in\":\"body\"},{\"name\":\"save_conversation\",\"type\":\"boolean\",\"description\":\"是否保存对话记录(默认true)\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (112, '2026-01-29 22:06:41.944', '2026-04-28 19:58:58.378', NULL, '/api/v1/ai/messages/:id', 'DELETE', 'AI对话', '删除单条消息', '', '', 1);
INSERT INTO `sys_api` VALUES (113, '2026-01-30 02:19:07.375', '2026-04-28 19:58:58.361', NULL, '/api/v1/ai/conversations/:id/clear-context', 'POST', 'AI对话', '清空上下文', '', '', 1);
INSERT INTO `sys_api` VALUES (114, '2026-01-31 00:13:15.803', '2026-02-03 00:50:21.106', NULL, '/api/v1/auth/register', 'POST', '认证管理', '注册', '[{\"name\":\"username\",\"type\":\"string\",\"description\":\"用户名\",\"required\":true,\"in\":\"body\"},{\"name\":\"password\",\"type\":\"string\",\"description\":\"密码\",\"required\":true,\"in\":\"body\"},{\"name\":\"email\",\"type\":\"string\",\"description\":\"邮箱\",\"required\":false,\"in\":\"body\"},{\"name\":\"email_code\",\"type\":\"string\",\"description\":\"邮箱验证码\",\"required\":false,\"in\":\"body\"},{\"name\":\"captcha_id\",\"type\":\"string\",\"description\":\"验证码ID\",\"required\":false,\"in\":\"body\"},{\"name\":\"captcha_code\",\"type\":\"string\",\"description\":\"验证码\",\"required\":false,\"in\":\"body\"}]', '', 0);
INSERT INTO `sys_api` VALUES (115, '2026-01-31 00:13:15.842', '2026-02-01 01:49:27.791', NULL, '/api/v1/auth/send-email-code', 'POST', '认证管理', '发送邮箱验证码', '[{\"name\":\"email\",\"type\":\"string\",\"description\":\"邮箱地址\",\"required\":true,\"in\":\"body\"},{\"name\":\"captcha_id\",\"type\":\"string\",\"description\":\"验证码ID\",\"required\":false,\"in\":\"body\"},{\"name\":\"captcha\",\"type\":\"string\",\"description\":\"验证码\",\"required\":false,\"in\":\"body\"}]', '', 0);
INSERT INTO `sys_api` VALUES (117, '2026-01-31 00:13:15.927', '2026-02-01 01:49:27.849', NULL, '/api/v1/auth/reset-password', 'POST', '认证管理', '重置密码(Token)', '[{\"name\":\"token\",\"type\":\"string\",\"description\":\"重置令牌\",\"required\":true,\"in\":\"body\"},{\"name\":\"password\",\"type\":\"string\",\"description\":\"新密码\",\"required\":true,\"in\":\"body\"}]', '', 0);
INSERT INTO `sys_api` VALUES (118, '2026-01-31 00:13:15.961', '2026-02-01 01:49:28.140', NULL, '/api/v1/captcha', 'GET', '验证码管理', '获取图形验证码', '', '', 0);
INSERT INTO `sys_api` VALUES (119, '2026-01-31 00:13:15.992', '2026-02-01 01:49:28.183', NULL, '/api/v1/captcha/config', 'GET', '验证码管理', '获取验证码配置', '', '', 0);
INSERT INTO `sys_api` VALUES (120, '2026-01-31 22:44:30.727', '2026-02-01 01:49:27.899', NULL, '/api/v1/auth/reset-password-by-username', 'POST', '认证管理', '重置密码(用户名)', '[{\"name\":\"username\",\"type\":\"string\",\"description\":\"用户名\",\"required\":true,\"in\":\"body\"},{\"name\":\"new_password\",\"type\":\"string\",\"description\":\"新密码\",\"required\":true,\"in\":\"body\"},{\"name\":\"captcha_id\",\"type\":\"string\",\"description\":\"验证码ID\",\"required\":true,\"in\":\"body\"},{\"name\":\"captcha\",\"type\":\"string\",\"description\":\"验证码\",\"required\":true,\"in\":\"body\"}]', '', 0);
INSERT INTO `sys_api` VALUES (121, '2026-01-31 23:47:32.467', '2026-02-01 01:49:27.932', NULL, '/api/v1/auth/reset-password-by-email', 'POST', '认证管理', '重置密码(邮箱)', '[{\"name\":\"email\",\"type\":\"string\",\"description\":\"邮箱地址\",\"required\":true,\"in\":\"body\"},{\"name\":\"email_code\",\"type\":\"string\",\"description\":\"邮箱验证码\",\"required\":true,\"in\":\"body\"},{\"name\":\"new_password\",\"type\":\"string\",\"description\":\"新密码\",\"required\":true,\"in\":\"body\"}]', '', 0);
INSERT INTO `sys_api` VALUES (122, '2026-02-01 02:42:42.420', '2026-04-29 20:37:13.929', NULL, '/api/v1/configs/test-email', 'POST', '系统配置管理', '发送测试邮件', '[{\"name\":\"email\",\"type\":\"string\",\"description\":\"接收测试邮件的邮箱地址\",\"required\":true,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (123, '2026-02-01 20:46:48.312', '2026-02-01 20:46:48.312', NULL, '/api/v1/logs/route-groups', 'GET', '日志管理', '获取路由分组列表', '', '', 1);
INSERT INTO `sys_api` VALUES (129, '2026-02-03 00:50:21.576', '2026-02-03 00:50:21.576', NULL, '/api/v1/echart/user-role-stats', 'GET', '图表统计', '用户角色占比', '', '', 1);
INSERT INTO `sys_api` VALUES (130, '2026-02-03 00:50:21.841', '2026-02-03 00:50:21.841', NULL, '/api/v1/echart/user-status-stats', 'GET', '图表统计', '用户状态统计', '', '', 1);
INSERT INTO `sys_api` VALUES (131, '2026-02-03 00:50:21.892', '2026-02-03 00:50:21.892', NULL, '/api/v1/echart/role-status-stats', 'GET', '图表统计', '角色状态统计', '', '', 1);
INSERT INTO `sys_api` VALUES (132, '2026-02-03 00:50:21.950', '2026-02-03 00:50:21.950', NULL, '/api/v1/echart/user-register-trend', 'GET', '图表统计', '用户注册趋势', '', '', 1);
INSERT INTO `sys_api` VALUES (159, '2026-02-05 04:04:03.589', '2026-02-05 04:04:03.589', NULL, '/api/v1/users/options', 'GET', '用户管理', '用户选项', '', '', 1);
INSERT INTO `sys_api` VALUES (160, '2026-02-05 07:49:09.965', '2026-04-29 20:33:50.430', NULL, '/api/v1/ai/test', 'POST', 'AI配置', '测试AI配置', '[{\"name\":\"api_key\",\"type\":\"string\",\"description\":\"APIKey\",\"required\":true,\"in\":\"body\"},{\"name\":\"base_url\",\"type\":\"string\",\"description\":\"BaseURL\",\"required\":true,\"in\":\"body\"},{\"name\":\"model\",\"type\":\"string\",\"description\":\"Model\",\"required\":true,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (172, '2026-02-07 00:16:48.036', '2026-02-07 00:16:48.036', NULL, '/api/v1/user/profiles', 'GET', '个人中心', '获取用户所有身份', '', '', 1);
INSERT INTO `sys_api` VALUES (173, '2026-02-07 00:16:48.281', '2026-02-07 00:16:48.281', NULL, '/api/v1/user/profiles/types', 'GET', '个人中心', '获取所有身份类型', '', '', 1);
INSERT INTO `sys_api` VALUES (174, '2026-02-07 00:16:48.373', '2026-02-07 00:16:48.373', NULL, '/api/v1/users/:id/offline', 'POST', '用户管理', '强制下线', '', '', 1);
INSERT INTO `sys_api` VALUES (175, '2026-02-07 00:16:48.446', '2026-02-07 00:16:48.446', NULL, '/api/v1/users/:id/profiles', 'GET', '用户管理', '用户身份', '', '', 1);
INSERT INTO `sys_api` VALUES (204, '2026-02-10 01:15:27.599', '2026-02-10 01:15:27.599', NULL, '/api/v1/dict/type/:type', 'GET', '字典管理', '获取字典数据', '', '', 0);
INSERT INTO `sys_api` VALUES (205, '2026-02-10 01:15:27.656', '2026-02-10 01:15:27.656', NULL, '/api/v1/dict/types', 'GET', '字典管理', '字典类型列表', '', '', 1);
INSERT INTO `sys_api` VALUES (206, '2026-02-10 01:15:27.723', '2026-02-10 01:15:27.723', NULL, '/api/v1/dict/types/all', 'GET', '字典管理', '所有字典类型', '', '', 1);
INSERT INTO `sys_api` VALUES (207, '2026-02-10 01:15:27.763', '2026-02-10 01:15:27.763', NULL, '/api/v1/dict/types/:id', 'GET', '字典管理', '字典类型详情', '', '', 1);
INSERT INTO `sys_api` VALUES (208, '2026-02-10 01:15:27.832', '2026-02-10 01:15:27.832', NULL, '/api/v1/dict/types', 'POST', '字典管理', '创建字典类型', '[{\"name\":\"name\",\"type\":\"string\",\"description\":\"字典名称\",\"required\":true,\"in\":\"body\"},{\"name\":\"type\",\"type\":\"string\",\"description\":\"字典类型\",\"required\":true,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(1:正常,0:停用)\",\"required\":false,\"in\":\"body\"},{\"name\":\"remark\",\"type\":\"string\",\"description\":\"备注\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (209, '2026-02-10 01:15:27.873', '2026-02-10 01:15:27.873', NULL, '/api/v1/dict/types/:id', 'PUT', '字典管理', '更新字典类型', '[{\"name\":\"name\",\"type\":\"string\",\"description\":\"字典名称\",\"required\":false,\"in\":\"body\"},{\"name\":\"type\",\"type\":\"string\",\"description\":\"字典类型\",\"required\":false,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(1:正常,0:停用)\",\"required\":false,\"in\":\"body\"},{\"name\":\"remark\",\"type\":\"string\",\"description\":\"备注\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (210, '2026-02-10 01:15:27.913', '2026-02-10 01:15:27.913', NULL, '/api/v1/dict/types/:id', 'DELETE', '字典管理', '删除字典类型', '', '', 1);
INSERT INTO `sys_api` VALUES (211, '2026-02-10 01:15:27.963', '2026-02-10 01:15:27.963', NULL, '/api/v1/dict/data', 'GET', '字典管理', '字典数据列表', '', '', 1);
INSERT INTO `sys_api` VALUES (212, '2026-02-10 01:15:28.097', '2026-02-10 01:15:28.097', NULL, '/api/v1/dict/data/:id', 'GET', '字典管理', '字典数据详情', '', '', 1);
INSERT INTO `sys_api` VALUES (213, '2026-02-10 01:15:28.140', '2026-02-10 01:15:28.140', NULL, '/api/v1/dict/data', 'POST', '字典管理', '创建字典数据', '[{\"name\":\"dict_type\",\"type\":\"string\",\"description\":\"字典类型\",\"required\":true,\"in\":\"body\"},{\"name\":\"label\",\"type\":\"string\",\"description\":\"字典标签\",\"required\":true,\"in\":\"body\"},{\"name\":\"value\",\"type\":\"string\",\"description\":\"字典键值\",\"required\":true,\"in\":\"body\"},{\"name\":\"sort\",\"type\":\"integer\",\"description\":\"排序\",\"required\":false,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(1:正常,0:停用)\",\"required\":false,\"in\":\"body\"},{\"name\":\"tag_type\",\"type\":\"string\",\"description\":\"标签类型\",\"required\":false,\"in\":\"body\"},{\"name\":\"is_default\",\"type\":\"integer\",\"description\":\"是否默认\",\"required\":false,\"in\":\"body\"},{\"name\":\"remark\",\"type\":\"string\",\"description\":\"备注\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (214, '2026-02-10 01:15:28.215', '2026-02-10 01:15:28.215', NULL, '/api/v1/dict/data/:id', 'PUT', '字典管理', '更新字典数据', '[{\"name\":\"dict_type\",\"type\":\"string\",\"description\":\"字典类型\",\"required\":false,\"in\":\"body\"},{\"name\":\"label\",\"type\":\"string\",\"description\":\"字典标签\",\"required\":false,\"in\":\"body\"},{\"name\":\"value\",\"type\":\"string\",\"description\":\"字典键值\",\"required\":false,\"in\":\"body\"},{\"name\":\"sort\",\"type\":\"integer\",\"description\":\"排序\",\"required\":false,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(1:正常,0:停用)\",\"required\":false,\"in\":\"body\"},{\"name\":\"tag_type\",\"type\":\"string\",\"description\":\"标签类型\",\"required\":false,\"in\":\"body\"},{\"name\":\"is_default\",\"type\":\"integer\",\"description\":\"是否默认\",\"required\":false,\"in\":\"body\"},{\"name\":\"remark\",\"type\":\"string\",\"description\":\"备注\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (215, '2026-02-10 01:15:28.305', '2026-02-10 01:15:28.305', NULL, '/api/v1/dict/data/:id', 'DELETE', '字典管理', '删除字典数据', '', '', 1);
INSERT INTO `sys_api` VALUES (216, '2026-02-10 02:18:36.653', '2026-02-10 02:18:36.653', NULL, '/api/v1/captcha/slider', 'GET', '验证码管理', '获取滑动验证码', '', '', 0);
INSERT INTO `sys_api` VALUES (217, '2026-02-10 02:18:36.715', '2026-02-10 02:18:36.715', NULL, '/api/v1/captcha/slider/verify', 'POST', '验证码管理', '验证滑动验证码', '', '', 0);
INSERT INTO `sys_api` VALUES (218, '2026-02-10 02:18:36.790', '2026-02-10 02:18:36.790', NULL, '/api/v1/files/batch', 'DELETE', '文件管理', '批量删除文件', '', '', 1);
INSERT INTO `sys_api` VALUES (219, '2026-02-10 02:18:36.914', '2026-02-10 02:18:36.914', NULL, '/api/v1/users/batch', 'DELETE', '用户管理', '批量删除用户', '', '', 1);
INSERT INTO `sys_api` VALUES (220, '2026-04-26 23:18:22.046', '2026-04-29 20:33:50.732', NULL, '/api/v1/users/batch-status', 'PUT', '用户管理', '批量修改用户状态', '[{\"name\":\"ids\",\"type\":\"array[integer]\",\"description\":\"用户ID列表\",\"required\":true,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(0:禁用,1:启用)\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (221, '2026-04-27 05:10:18.841', '2026-04-28 19:58:58.184', NULL, '/api/v1/depts/tree', 'GET', '部门管理', '部门树', '', '', 1);
INSERT INTO `sys_api` VALUES (222, '2026-04-27 05:10:18.846', '2026-04-28 19:58:58.211', NULL, '/api/v1/depts/:id', 'GET', '部门管理', '部门详情', '', '', 1);
INSERT INTO `sys_api` VALUES (223, '2026-04-27 05:10:18.850', '2026-04-29 20:33:50.535', NULL, '/api/v1/depts', 'POST', '部门管理', '创建部门', '[{\"name\":\"parent_id\",\"type\":\"integer\",\"description\":\"父部门ID\",\"required\":false,\"in\":\"body\"},{\"name\":\"name\",\"type\":\"string\",\"description\":\"部门名称\",\"required\":true,\"in\":\"body\"},{\"name\":\"sort\",\"type\":\"integer\",\"description\":\"排序\",\"required\":false,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(0:禁用,1:启用)\",\"required\":false,\"in\":\"body\"},{\"name\":\"remark\",\"type\":\"string\",\"description\":\"备注\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (224, '2026-04-27 05:10:18.855', '2026-04-29 20:33:50.541', NULL, '/api/v1/depts/:id', 'PUT', '部门管理', '更新部门', '[{\"name\":\"parent_id\",\"type\":\"integer\",\"description\":\"父部门ID\",\"required\":false,\"in\":\"body\"},{\"name\":\"name\",\"type\":\"string\",\"description\":\"部门名称\",\"required\":true,\"in\":\"body\"},{\"name\":\"sort\",\"type\":\"integer\",\"description\":\"排序\",\"required\":false,\"in\":\"body\"},{\"name\":\"status\",\"type\":\"integer\",\"description\":\"状态(0:禁用,1:启用)\",\"required\":false,\"in\":\"body\"},{\"name\":\"remark\",\"type\":\"string\",\"description\":\"备注\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (225, '2026-04-27 05:10:18.859', '2026-04-28 19:58:58.250', NULL, '/api/v1/depts/:id', 'DELETE', '部门管理', '删除部门', '', '', 1);
INSERT INTO `sys_api` VALUES (226, '2026-04-27 05:52:38.131', '2026-04-28 19:58:58.199', NULL, '/api/v1/depts/manageable-tree', 'GET', '部门管理', '可管理部门树', '', '', 1);
INSERT INTO `sys_api` VALUES (227, '2026-04-28 04:38:07.467', '2026-05-05 01:15:49.648', NULL, '/api/v1/ai/providers/models/fetch', 'POST', 'AI配置', '拉取平台模型列表', '[{\"name\":\"api_key\",\"type\":\"string\",\"description\":\"APIKey\",\"required\":true,\"in\":\"body\"},{\"name\":\"base_url\",\"type\":\"string\",\"description\":\"BaseURL\",\"required\":true,\"in\":\"body\"},{\"name\":\"provider_name\",\"type\":\"string\",\"description\":\"ProviderName\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (228, '2026-04-28 19:58:27.000', '2026-04-29 20:33:50.748', NULL, '/api/v1/users/batch-password', 'PUT', '用户管理', '批量重置密码', '[{\"name\":\"ids\",\"type\":\"array[integer]\",\"description\":\"用户ID列表\",\"required\":true,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (230, '2026-04-29 19:40:11.997', '2026-04-29 19:40:11.997', NULL, '/api/v1/ai/config', 'GET', 'AI配置', '获取AI配置', '', '', 1);
INSERT INTO `sys_api` VALUES (231, '2026-04-29 19:40:12.102', '2026-04-29 19:40:12.102', NULL, '/api/v1/ai/config', 'PUT', 'AI配置', '保存AI配置', '', '', 1);
INSERT INTO `sys_api` VALUES (232, '2026-04-29 20:33:50.402', '2026-04-29 20:33:50.402', NULL, '/api/v1/ai/conversations/batch', 'DELETE', 'AI对话', '批量删除对话', '[{\"name\":\"ids\",\"type\":\"array[integer]\",\"description\":\"对话ID列表\",\"required\":true,\"in\":\"query\"}]', '', 1);
INSERT INTO `sys_api` VALUES (233, '2026-04-29 20:33:50.525', '2026-04-29 20:37:13.947', NULL, '/api/v1/configs/storage/test', 'POST', '系统配置管理', '测试存储配置', '', '', 1);
INSERT INTO `sys_api` VALUES (234, '2026-04-29 20:33:50.598', '2026-04-29 20:33:50.598', NULL, '/api/v1/files/migrate/preview', 'POST', '文件管理', '预览文件迁移', '', '', 1);
INSERT INTO `sys_api` VALUES (235, '2026-04-29 20:33:50.607', '2026-04-29 20:33:50.607', NULL, '/api/v1/files/migrate/execute', 'POST', '文件管理', '执行文件迁移', '', '', 1);
INSERT INTO `sys_api` VALUES (236, '2026-04-29 20:33:50.616', '2026-04-29 20:33:50.616', NULL, '/api/v1/files/migrate/task/current', 'GET', '文件管理', '获取当前文件迁移任务', '', '', 1);
INSERT INTO `sys_api` VALUES (237, '2026-04-30 04:13:13.326', '2026-05-05 01:15:49.819', NULL, '/api/v1/roles/:id/data-scopes', 'PUT', '角色管理', '分配数据权限', '[{\"name\":\"scopes\",\"type\":\"array[object]\",\"description\":\"角色业务功能数据范围列表\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (238, '2026-05-02 01:23:33.450', '2026-05-05 01:15:49.825', NULL, '/api/v1/roles/:id/permissions', 'PUT', '角色管理', '统一保存角色权限', '[{\"name\":\"menu_ids\",\"type\":\"array[integer]\",\"description\":\"菜单ID列表\",\"required\":false,\"in\":\"body\"},{\"name\":\"direct_api_ids\",\"type\":\"array[integer]\",\"description\":\"直接API ID列表\",\"required\":false,\"in\":\"body\"},{\"name\":\"scopes\",\"type\":\"array[object]\",\"description\":\"角色业务功能数据范围列表\",\"required\":false,\"in\":\"body\"}]', '', 1);
INSERT INTO `sys_api` VALUES (239, '2026-05-02 09:55:27.986', '2026-05-02 09:55:27.986', NULL, '/api/v1/roles/data-scope-resources', 'GET', '角色管理', '数据权限资源列表', '', '', 1);
INSERT INTO `sys_api` VALUES (240, '2026-05-04 11:51:09.592', '2026-05-04 11:51:09.592', NULL, '/api/v1/users/import-template', 'GET', '用户管理', '下载导入模板', '', '', 1);
INSERT INTO `sys_api` VALUES (241, '2026-05-04 11:51:09.695', '2026-05-04 11:51:09.695', NULL, '/api/v1/users/import', 'POST', '用户管理', '导入用户', '', '', 1);
INSERT INTO `sys_api` VALUES (242, '2026-05-04 11:51:09.801', '2026-05-04 11:51:09.801', NULL, '/api/v1/users/export', 'GET', '用户管理', '导出用户', '', '', 1);
INSERT INTO `sys_api` VALUES (243, '2026-05-05 01:15:49.769', '2026-05-05 01:15:49.769', NULL, '/api/v1/health', 'GET', '健康检查', '服务健康检查', '', '', 0);
INSERT INTO `sys_api` VALUES (244, '2026-05-05 01:15:49.780', '2026-05-05 01:15:49.780', NULL, '/api/v1/menus/tree-with-apis', 'GET', '菜单管理', '菜单树(带API)', '', '', 1);
INSERT INTO `sys_api` VALUES (245, '2026-05-05 01:15:49.787', '2026-05-05 01:15:49.787', NULL, '/api/v1/menus/:id/apis', 'GET', '菜单管理', '菜单API列表', '', '', 1);
INSERT INTO `sys_api` VALUES (246, '2026-05-05 01:15:49.795', '2026-05-05 01:15:49.795', NULL, '/api/v1/menus/:id/apis', 'PUT', '菜单管理', '更新菜单API', '[{\"name\":\"api_ids\",\"type\":\"array[integer]\",\"description\":\"菜单关联API ID列表\",\"required\":false,\"in\":\"body\"}]', '', 1);

-- ----------------------------
-- Table structure for sys_config
-- ----------------------------
DROP TABLE IF EXISTS `sys_config`;
CREATE TABLE `sys_config`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '配置名称',
  `key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '配置键',
  `value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '配置值',
  `value_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'string' COMMENT '值类型 string/json',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sys_config_key`(`key` ASC) USING BTREE,
  INDEX `idx_sys_config_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 55 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_config
-- ----------------------------
INSERT INTO `sys_config` VALUES (1, '2026-01-25 00:27:27.184', '2026-05-04 03:30:09.275', NULL, '系统名称', 'sys_name', '天之云', 'string', '显示在侧边栏顶部');
INSERT INTO `sys_config` VALUES (2, '2026-01-25 00:27:27.184', '2026-05-04 03:30:09.277', NULL, '系统Logo', 'sys_logo', 'http://115.190.199.115:7000/go-base/2026/05/02/1777664719606847500.png', 'string', '系统Logo图片地址');
INSERT INTO `sys_config` VALUES (13, '2026-01-26 00:44:42.303', '2026-05-05 01:08:08.325', NULL, 'login_title', 'login_title', '欢迎登录', 'string', '');
INSERT INTO `sys_config` VALUES (14, '2026-01-26 00:44:42.304', '2026-05-05 01:08:08.312', NULL, 'login_subtitle', 'login_subtitle', '智慧养殖，生态未来', 'string', '');
INSERT INTO `sys_config` VALUES (15, '2026-01-26 00:44:42.308', '2026-05-05 01:08:08.304', NULL, 'login_bg_image', 'login_bg_image', 'http://115.190.199.115:7000/go-base/2026/05/02/1777664632234976900.png', 'string', '');
INSERT INTO `sys_config` VALUES (16, '2026-01-26 00:44:42.308', '2026-05-05 01:08:08.315', NULL, 'login_bg_color', 'login_bg_color', 'linear-gradient(135deg,#42e695,#3bb2b8)', 'string', '');
INSERT INTO `sys_config` VALUES (17, '2026-01-30 22:01:16.067', '2026-05-05 01:08:08.322', NULL, '', 'login_features', '[{\"icon\":\"CheckCircleOutlined\",\"title\":\"牲畜追溯\",\"desc\":\"\"},{\"icon\":\"SafetyOutlined\",\"title\":\"智能喂养\",\"desc\":\"\"},{\"icon\":\"LineChartOutlined\",\"title\":\"疫苗管理\",\"desc\":\"\"},{\"icon\":\"ThunderboltOutlined\",\"title\":\"疾病预警\",\"desc\":\"\"}]', 'string', '');
INSERT INTO `sys_config` VALUES (18, '2026-01-30 22:01:16.068', '2026-05-05 01:08:08.334', NULL, '', 'login_images', '[]', 'string', '');
INSERT INTO `sys_config` VALUES (19, '2026-01-30 22:01:16.072', '2026-05-05 01:08:08.328', NULL, '', 'login_slogan', '智慧养殖，生态未来', 'string', '');
INSERT INTO `sys_config` VALUES (20, '2026-01-30 22:01:16.072', '2026-05-05 01:08:08.319', NULL, '', 'login_desc', '牧智云是一套专业的生态养殖管理平台，为养殖户提供牲畜管理、疫苗接种、疾病防控、饲料追踪等全流程数字化解决方案', 'string', '');
INSERT INTO `sys_config` VALUES (21, '2026-01-30 22:37:07.412', '2026-05-05 01:08:08.331', NULL, '', 'login_features_max', '4', 'string', '');
INSERT INTO `sys_config` VALUES (22, '2026-01-30 22:37:07.414', '2026-05-05 01:08:08.338', NULL, '', 'login_images_max', '2', 'string', '');
INSERT INTO `sys_config` VALUES (23, '2026-01-31 00:22:04.469', '2026-04-27 20:26:40.932', NULL, '', 'email_password', 'jciaykeercqaibgc', 'string', '');
INSERT INTO `sys_config` VALUES (24, '2026-01-31 00:22:04.471', '2026-04-27 20:26:40.946', NULL, '', 'email_from_name', 'test', 'string', '');
INSERT INTO `sys_config` VALUES (25, '2026-01-31 00:22:04.471', '2026-04-27 20:26:40.954', NULL, '', 'register_email_verify', '1', 'string', '');
INSERT INTO `sys_config` VALUES (26, '2026-01-31 00:22:04.472', '2026-04-27 20:26:40.959', NULL, '', 'email_smtp_host', 'smtp.qq.com', 'string', '');
INSERT INTO `sys_config` VALUES (27, '2026-01-31 00:22:04.472', '2026-04-27 20:26:40.949', NULL, '', 'login_captcha_enabled', '1', 'string', '');
INSERT INTO `sys_config` VALUES (28, '2026-01-31 00:22:04.473', '2026-02-02 09:14:19.672', NULL, '', 'register_captcha_enabled', '0', 'string', '');
INSERT INTO `sys_config` VALUES (29, '2026-01-31 00:22:04.474', '2026-04-27 20:26:40.941', NULL, '', 'frontend_url', 'http://localhost:5173', 'string', '');
INSERT INTO `sys_config` VALUES (30, '2026-01-31 00:22:04.475', '2026-04-27 20:26:40.944', NULL, '', 'email_smtp_port', '587', 'string', '');
INSERT INTO `sys_config` VALUES (31, '2026-01-31 00:22:04.475', '2026-04-27 20:26:40.929', NULL, '', 'email_username', '1440350254@qq.com', 'string', '');
INSERT INTO `sys_config` VALUES (32, '2026-01-31 20:42:31.302', '2026-05-05 01:08:08.295', NULL, '', 'register_logo', 'http://115.190.199.115:7000/go-base/2026/02/15/1771163518773570700.jpg', 'string', '');
INSERT INTO `sys_config` VALUES (33, '2026-02-01 02:14:55.436', '2026-05-05 01:08:08.340', NULL, '', 'enable_register', 'false', 'string', '');
INSERT INTO `sys_config` VALUES (35, '2026-02-07 17:04:08.827', '2026-05-04 03:30:09.281', NULL, '', 'front_mode', 'none', 'string', '');
INSERT INTO `sys_config` VALUES (36, '2026-02-10 01:48:05.727', '2026-04-27 20:26:40.939', NULL, '', 'login_lock_time', '15', 'string', '');
INSERT INTO `sys_config` VALUES (37, '2026-02-10 01:48:05.729', '2026-04-27 20:26:40.935', NULL, '', 'login_captcha_type', 'math', 'string', '');
INSERT INTO `sys_config` VALUES (38, '2026-02-10 01:48:05.729', '2026-04-27 20:26:40.952', NULL, '', 'login_max_retry', '5', 'string', '');
INSERT INTO `sys_config` VALUES (39, '2026-02-10 02:08:32.557', '2026-04-27 20:26:40.958', NULL, '', 'slider_captcha_bg', '/api/v1/upload/2026/02/15/1771163674365522700.jpg', 'string', '');
INSERT INTO `sys_config` VALUES (40, '2026-04-26 23:18:22.035', '2026-05-04 03:30:09.282', NULL, '用户身份按钮显示', 'user_profile_button_visible', 'false', 'string', '后台用户管理列表是否显示身份按钮');
INSERT INTO `sys_config` VALUES (42, '2026-04-27 21:27:30.217', '2026-05-04 10:09:03.350', NULL, '', 'storage_type', 'minio', 'string', '');
INSERT INTO `sys_config` VALUES (43, '2026-04-27 21:27:30.224', '2026-04-27 21:30:32.671', NULL, '', 'storage_config', '{\"base_path\":\"uploads\",\"base_url\":\"/api/v1/upload\"}', 'string', '');
INSERT INTO `sys_config` VALUES (44, '2026-04-28 00:13:29.747', '2026-05-04 10:09:03.344', NULL, '', 'storage_minio_config', '{\"endpoint\":\"115.190.199.115:7000\",\"access_key_id\":\"s6o4KO5Iup1UpIXWKtwc\",\"secret_access_key\":\"3B2P1cjmwODnT5l4yrWoevnGOxC1rJU7uexUjQdt\",\"bucket_name\":\"go-base\",\"use_ssl\":false}', 'string', '');
INSERT INTO `sys_config` VALUES (45, '2026-04-28 00:13:38.258', '2026-05-04 10:09:03.354', NULL, '', 'storage_local_config', '{\"base_path\":\"uploads\",\"base_url\":\"/api/v1/upload\"}', 'string', '');
INSERT INTO `sys_config` VALUES (46, '2026-04-28 00:15:46.192', '2026-05-04 10:09:03.336', NULL, '阿里云 OSS 配置', 'storage_aliyun_config', '{\"endpoint\":\"\",\"access_key_id\":\"YT01\",\"access_key_secret\":\"123456\",\"bucket_name\":\"\",\"region\":\"\"}', 'json', '阿里云 OSS 的已保存配置(JSON)');
INSERT INTO `sys_config` VALUES (47, '2026-04-28 00:15:46.201', '2026-05-04 10:09:03.340', NULL, '腾讯云 COS 配置', 'storage_tencent_config', '{\"region\":\"\",\"secret_id\":\"\",\"secret_key\":\"\",\"bucket\":\"\",\"app_id\":\"\"}', 'json', '腾讯云 COS 的已保存配置(JSON)');
INSERT INTO `sys_config` VALUES (48, '2026-04-28 00:36:13.315', '2026-05-04 10:09:03.347', NULL, '', 'file_delete_mode', 'physical', 'string', '');
INSERT INTO `sys_config` VALUES (49, '2026-04-28 19:58:27.000', '2026-05-04 03:30:09.286', NULL, '用户默认密码', 'user_default_password', '123456', 'string', '后台用户管理单条/批量重置密码默认值');
INSERT INTO `sys_config` VALUES (50, '2026-04-30 21:29:48.644', '2026-04-30 21:51:47.628', NULL, '公开配置键', 'public_config_keys', '[\"sys_name\",\"sys_logo\",\"login_bg_image\",\"login_title\",\"login_subtitle\",\"login_bg_color\",\"login_slogan\",\"login_desc\",\"login_features\",\"login_features_max\",\"login_images\",\"login_images_max\",\"register_logo\",\"enable_register\",\"front_mode\",\"user_profile_button_visible\"]', 'json', '允许匿名批量读取的配置键(JSON数组)，敏感键即使写入也不会公开');
INSERT INTO `sys_config` VALUES (51, '2026-05-01 08:03:59.025', '2026-05-01 08:03:59.025', NULL, 'AI配置', 'ai_config', '{\"default_provider\":\"阿里云百炼\",\"providers\":[{\"name\":\"deepseek\",\"api_key\":\"sk-e6c0478a770a49ad8650884ce456c801\",\"base_url\":\"https://api.deepseek.com\",\"models\":[{\"id\":\"deepseek-v4-flash\",\"name\":\"deepseek-v4-flash\",\"description\":\"\"},{\"id\":\"deepseek-v4-pro\",\"name\":\"deepseek-v4-pro\",\"description\":\"\"}]},{\"name\":\"阿里云百炼\",\"api_key\":\"sk-93a09f74fb7e49dc951d8b30b06e5a04\",\"base_url\":\"https://dashscope.aliyuncs.com/compatible-mode/v1\",\"models\":[{\"id\":\"deepseek-v4-flash\",\"name\":\"deepseek-v4-flash\",\"description\":\"\"},{\"id\":\"kimi-k2.6\",\"name\":\"kimi-k2.6\",\"description\":\"\"},{\"id\":\"deepseek-v4-pro\",\"name\":\"deepseek-v4-pro\",\"description\":\"\"},{\"id\":\"qwen3.6-27b\",\"name\":\"qwen3.6-27b\",\"description\":\"\"}]}]}', 'json', 'AI平台配置');
INSERT INTO `sys_config` VALUES (52, '2026-05-02 03:40:07.109', '2026-05-04 03:30:09.279', NULL, '系统Logo文件ID', 'sys_logo_file_id', '121', 'string', '系统Logo关联文件ID');
INSERT INTO `sys_config` VALUES (53, '2026-05-02 03:40:07.122', '2026-05-05 01:08:08.301', NULL, '注册默认头像文件ID', 'register_logo_file_id', '93', 'string', '注册默认头像关联文件ID');
INSERT INTO `sys_config` VALUES (54, '2026-05-02 03:40:07.133', '2026-05-05 01:08:08.309', NULL, '登录页背景图文件ID', 'login_bg_image_file_id', '120', 'string', '登录页背景图关联文件ID');

-- ----------------------------
-- Table structure for sys_dept
-- ----------------------------
DROP TABLE IF EXISTS `sys_dept`;
CREATE TABLE `sys_dept`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `parent_id` bigint UNSIGNED NOT NULL DEFAULT 0 COMMENT '父部门ID',
  `ancestors` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '祖级列表',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '部门名称',
  `sort` bigint NULL DEFAULT 0 COMMENT '排序',
  `status` bigint NULL DEFAULT 1 COMMENT '状态 1启用 0禁用',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_dept_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_sys_dept_parent_id`(`parent_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_dept
-- ----------------------------
INSERT INTO `sys_dept` VALUES (1, '2026-04-27 05:10:18.800', '2026-04-27 06:16:24.529', '2026-04-27 06:19:30.998', 0, '0', '平台', 1, 1, '系统根部门');
INSERT INTO `sys_dept` VALUES (2, '2026-04-27 05:16:57.192', '2026-04-27 05:16:57.192', '2026-04-27 06:19:29.604', 1, '0,1', '111', 0, 1, '');
INSERT INTO `sys_dept` VALUES (3, '2026-04-27 06:09:50.379', '2026-04-27 06:34:23.493', NULL, 0, '0', '运维部', 0, 1, '');
INSERT INTO `sys_dept` VALUES (4, '2026-04-27 06:19:35.384', '2026-04-27 06:34:35.885', NULL, 0, '0', '开发部', 0, 1, '');
INSERT INTO `sys_dept` VALUES (5, '2026-04-27 06:20:30.215', '2026-04-27 06:34:46.553', NULL, 4, '0,4', '后端组', 0, 1, '');
INSERT INTO `sys_dept` VALUES (6, '2026-04-27 06:34:53.999', '2026-04-27 06:34:53.999', NULL, 4, '0,4', '前端组', 0, 1, '');
INSERT INTO `sys_dept` VALUES (7, '2026-04-27 06:39:46.508', '2026-04-28 19:58:58.168', NULL, 0, '0', '平台', 1, 1, '系统根部门');
INSERT INTO `sys_dept` VALUES (8, '2026-04-30 05:37:13.949', '2026-04-30 05:37:13.949', NULL, 0, '0', 'A公司', 0, 1, '');

-- ----------------------------
-- Table structure for sys_dict_data
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict_data`;
CREATE TABLE `sys_dict_data`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `dict_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '字典类型',
  `label` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '字典标签',
  `value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '字典键值',
  `sort` bigint NULL DEFAULT 0 COMMENT '排序',
  `status` bigint NULL DEFAULT 1 COMMENT '状态(1:正常,0:停用)',
  `tag_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '标签类型(success/info/warning/error)',
  `is_default` bigint NULL DEFAULT 0 COMMENT '是否默认(1:是,0:否)',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_dict_data_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_sys_dict_data_dict_type`(`dict_type` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_dict_data
-- ----------------------------
INSERT INTO `sys_dict_data` VALUES (1, '2026-02-08 00:59:59.132', '2026-04-27 16:53:55.951', NULL, 'common_status', '正常', '1', 0, 1, 'success', 0, '');
INSERT INTO `sys_dict_data` VALUES (2, '2026-02-08 01:00:11.318', '2026-02-08 01:12:09.859', NULL, 'common_status', '禁用', '2', 0, 1, 'default', 0, '');
INSERT INTO `sys_dict_data` VALUES (3, '2026-02-08 01:13:46.981', '2026-02-08 01:13:46.981', '2026-02-10 01:05:00.537', 'area_type', '圈舍', '1', 0, 1, 'processing', 0, '');
INSERT INTO `sys_dict_data` VALUES (4, '2026-02-08 01:14:11.180', '2026-02-08 01:14:11.180', '2026-02-10 01:05:00.537', 'area_type', '鱼塘', '2', 0, 1, 'warning', 0, '');
INSERT INTO `sys_dict_data` VALUES (5, '2026-02-08 01:14:23.826', '2026-02-08 01:14:23.826', '2026-02-10 01:05:00.537', 'area_type', '牧场', '3', 0, 1, 'purple', 0, '');
INSERT INTO `sys_dict_data` VALUES (6, '2026-04-28 04:05:16.680', '2026-04-28 04:05:16.680', '2026-04-28 16:57:55.502', 'sys_gender', '未知', '0', 0, 1, 'default', 1, '默认值');
INSERT INTO `sys_dict_data` VALUES (7, '2026-04-28 04:05:16.685', '2026-04-28 16:58:14.315', NULL, 'sys_gender', '男', '0', 1, 1, 'processing', 0, '');
INSERT INTO `sys_dict_data` VALUES (8, '2026-04-28 04:05:16.689', '2026-04-28 16:58:19.307', NULL, 'sys_gender', '女', '1', 2, 1, 'purple', 0, '');
INSERT INTO `sys_dict_data` VALUES (9, '2026-04-28 17:10:46.761', '2026-04-28 17:10:46.761', '2026-04-29 03:26:42.177', 'sys_gender', '女', '2', 2, 1, 'pink', 0, '');
INSERT INTO `sys_dict_data` VALUES (10, '2026-04-29 03:27:10.236', '2026-04-29 03:27:10.236', '2026-04-29 03:27:32.583', 'sys_gender', '女', '2', 2, 1, 'pink', 0, '');

-- ----------------------------
-- Table structure for sys_dict_type
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict_type`;
CREATE TABLE `sys_dict_type`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '字典名称',
  `type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '字典类型',
  `status` bigint NULL DEFAULT 1 COMMENT '状态(1:正常,0:停用)',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sys_dict_type_type`(`type` ASC) USING BTREE,
  INDEX `idx_sys_dict_type_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_dict_type
-- ----------------------------
INSERT INTO `sys_dict_type` VALUES (1, '2026-02-08 00:59:43.416', '2026-02-08 00:59:43.416', NULL, '通用状态', 'common_status', 1, '');
INSERT INTO `sys_dict_type` VALUES (2, '2026-02-08 01:13:01.337', '2026-02-08 01:13:01.337', '2026-02-10 01:05:00.585', '区域类型', 'area_type', 1, '');
INSERT INTO `sys_dict_type` VALUES (3, '2026-04-28 04:05:16.675', '2026-04-28 04:05:16.675', NULL, '性别', 'sys_gender', 1, '用户性别字典');

-- ----------------------------
-- Table structure for sys_file
-- ----------------------------
DROP TABLE IF EXISTS `sys_file`;
CREATE TABLE `sys_file`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '文件名称',
  `path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '文件路径',
  `url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '访问URL',
  `size` bigint NULL DEFAULT NULL COMMENT '文件大小(字节)',
  `ext` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '文件扩展名',
  `mime_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'MIME类型',
  `md5` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '文件MD5',
  `storage_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '存储类型快照',
  `uploader_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '上传者ID',
  `status` bigint NULL DEFAULT 1 COMMENT '状态 0删除 1正常',
  `storage_bucket` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '存储桶/路径快照',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_file_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_sys_file_md5`(`md5` ASC) USING BTREE,
  INDEX `idx_sys_file_storage_type`(`storage_type` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 126 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_file
-- ----------------------------
INSERT INTO `sys_file` VALUES (93, '2026-02-15 21:51:58.803', '2026-05-04 10:17:32.059', NULL, 'test.jpg', '2026/02/15/1771163518773570700.jpg', 'http://115.190.199.115:7000/go-base/2026/02/15/1771163518773570700.jpg', 34422, 'jpg', 'image/jpeg', '8c5882dcf7600d0ca41487aec2d9a314', 'minio', 1, 1, 'go-base');
INSERT INTO `sys_file` VALUES (105, '2026-04-28 02:35:14.061', '2026-04-28 02:35:27.057', NULL, '20240308210115-0-2344-image-2.jpg', '2026/04/28/1777314913439194500.jpg', 'http://115.190.199.115:7000/go-base/2026/04/28/1777314913439194500.jpg', 248176, 'jpg', 'image/jpeg', 'a897e5d67d4d150b2500578709e14871', 'minio', 1, 0, NULL);
INSERT INTO `sys_file` VALUES (114, '2026-04-30 06:00:11.656', '2026-05-04 10:17:32.406', NULL, 'go.png', '2026/04/30/1777500011289137100.png', 'http://115.190.199.115:7000/go-base/2026/04/30/1777500011289137100.png', 5148, 'png', 'image/png', 'd51f4e40fa45e80f1f2b66d7f3300e78', 'minio', 1, 1, 'go-base');
INSERT INTO `sys_file` VALUES (115, '2026-04-30 06:01:39.114', '2026-05-04 10:17:32.863', NULL, 'sushi.jpg', '2026/04/30/1777500098740563700.jpg', 'http://115.190.199.115:7000/go-base/2026/04/30/1777500098740563700.jpg', 48285, 'jpg', 'image/jpeg', '57b1f6833e8c2275ac136637fbc2c97f', 'minio', 37, 1, 'go-base');
INSERT INTO `sys_file` VALUES (120, '2026-05-02 03:43:57.911', '2026-05-04 10:18:16.830', NULL, '46d3397ee79742c49beb280b17a66f71.jpeg~tplv-a9rns2rl98-image_raw_b.png', '2026/05/02/1777664632234976900.png', 'http://115.190.199.115:7000/go-base/2026/05/02/1777664632234976900.png', 5306844, 'png', 'image/png', '3b48d9c07838cf0fc86650ff93b7c4b9', 'minio', 1, 1, 'go-base');
INSERT INTO `sys_file` VALUES (121, '2026-05-02 03:45:19.959', '2026-05-04 10:18:17.290', NULL, '1776124139529.png', '2026/05/02/1777664719606847500.png', 'http://115.190.199.115:7000/go-base/2026/05/02/1777664719606847500.png', 10019, 'png', 'image/png', '4e241dfe9d3363ae6dbf9541d064d4e1', 'minio', 1, 1, 'go-base');
INSERT INTO `sys_file` VALUES (125, '2026-05-04 04:25:13.575', '2026-05-04 10:18:27.878', NULL, '运维_陈善平_本科_一年.pdf', '2026/05/04/1777839912660346300.pdf', 'http://115.190.199.115:7000/go-base/2026/05/04/1777839912660346300.pdf', 1318348, 'pdf', 'application/pdf', '936c20d1cd1c0fe016aec8ac9a49b053', 'minio', 1, 1, 'go-base');

-- ----------------------------
-- Table structure for sys_file_chunk
-- ----------------------------
DROP TABLE IF EXISTS `sys_file_chunk`;
CREATE TABLE `sys_file_chunk`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `upload_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '上传ID',
  `file_hash` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '文件MD5',
  `chunk_index` bigint NULL DEFAULT NULL COMMENT '分片索引',
  `chunk_size` bigint NULL DEFAULT NULL COMMENT '分片大小',
  `chunk_hash` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '分片MD5',
  `storage_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '存储类型快照',
  `storage_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '分片存储路径',
  `status` bigint NULL DEFAULT 0 COMMENT '状态 0上传中 1已完成',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_file_chunk_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_sys_file_chunk_upload_id`(`upload_id` ASC) USING BTREE,
  INDEX `idx_sys_file_chunk_file_hash`(`file_hash` ASC) USING BTREE,
  INDEX `idx_sys_file_chunk_storage_type`(`storage_type` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_file_chunk
-- ----------------------------

-- ----------------------------
-- Table structure for sys_login_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_login_log`;
CREATE TABLE `sys_login_log`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'IP地址',
  `location` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '登录地点',
  `browser` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '浏览器',
  `os` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作系统',
  `status` bigint NULL DEFAULT NULL COMMENT '状态 1成功 0失败',
  `msg` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '消息',
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_login_log
-- ----------------------------

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `parent_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT '父菜单ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '菜单名称',
  `path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '路由路径',
  `component` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '组件路径',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '图标',
  `sort` bigint NULL DEFAULT 0 COMMENT '排序',
  `type` bigint NULL DEFAULT 1 COMMENT '类型 1目录 2菜单 3按钮',
  `permission` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '权限标识',
  `status` bigint NULL DEFAULT 1 COMMENT '状态 1启用 0禁用',
  `hidden` bigint NULL DEFAULT 0 COMMENT '是否隐藏 0显示 1隐藏',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_menu_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 385 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_menu
-- ----------------------------
INSERT INTO `sys_menu` VALUES (1, '2026-01-24 02:39:47.984', '2026-01-24 05:17:02.075', NULL, 0, '系统管理', '/system', '', 'SettingFilled', 1, 1, '', 1, 0);
INSERT INTO `sys_menu` VALUES (2, '2026-01-24 02:39:48.085', '2026-01-24 05:11:33.852', '2026-04-29 17:20:15.616', 0, '系统监控', '/monitor', '', 'official-MonitorOutlined', 2, 1, '', 1, 0);
INSERT INTO `sys_menu` VALUES (3, '2026-01-24 02:39:48.202', '2026-01-24 19:04:53.047', NULL, 1, '用户管理', '/system/user', 'system/user/index', 'official-UserOutlined', 1, 2, 'system:user:list', 1, 0);
INSERT INTO `sys_menu` VALUES (4, '2026-01-24 02:39:48.202', '2026-01-24 19:05:24.917', NULL, 1, '角色管理', '/system/role', 'system/role/index', 'official-TeamOutlined', 2, 2, 'system:role:list', 1, 0);
INSERT INTO `sys_menu` VALUES (5, '2026-01-24 02:39:48.202', '2026-01-24 19:05:31.229', NULL, 1, '菜单管理', '/system/menu', 'system/menu/index', 'official-MenuOutlined', 3, 2, 'system:menu:list', 1, 0);
INSERT INTO `sys_menu` VALUES (6, '2026-01-24 02:39:48.202', '2026-01-24 19:05:39.252', NULL, 1, 'API管理', '/system/api', 'system/api/index', 'official-ApiFilled', 4, 2, 'system:api:list', 1, 0);
INSERT INTO `sys_menu` VALUES (7, '2026-01-24 02:39:48.202', '2026-05-04 01:25:15.439', NULL, 381, '操作日志', '/monitor/operation-log', 'monitor/operation-log/index', 'FileFilled', 1, 2, 'monitor:operation-log:list', 1, 0);
INSERT INTO `sys_menu` VALUES (8, '2026-01-24 02:39:48.202', '2026-04-29 16:42:43.650', NULL, 381, '登录日志', '/monitor/login-log', 'monitor/login-log/index', 'official-LoginOutlined', 2, 2, 'monitor:login-log:list', 1, 0);
INSERT INTO `sys_menu` VALUES (9, '2026-01-24 19:16:18.962', '2026-01-24 19:16:18.962', '2026-01-24 19:26:39.975', 3, '用户添加', '', '', '', 0, 3, 'system:user:add', 1, 0);
INSERT INTO `sys_menu` VALUES (10, '2026-01-24 19:18:05.943', '2026-01-24 19:18:05.943', '2026-01-24 19:26:37.141', 3, '用户编辑', '', '', '', 0, 3, 'system:user:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (11, '2026-01-25 00:30:18.345', '2026-04-27 20:29:44.415', NULL, 1, '系统配置', '/system/config', 'system/config/index', 'SettingTwoTone', 5, 2, 'system:config:list', 1, 0);
INSERT INTO `sys_menu` VALUES (12, '2026-01-25 00:49:56.151', '2026-01-25 00:55:42.428', '2026-04-29 17:20:13.210', 2, '慢查询日志', '/monitor/slow-log', 'monitor/slow-log/index', 'FileSearchOutlined', 0, 2, 'monitor:slow-log:list', 1, 0);
INSERT INTO `sys_menu` VALUES (13, '2026-01-25 01:40:12.093', '2026-04-26 20:01:05.586', '2026-04-29 17:20:17.756', 0, '系统工具', '/tool', '', 'ToolTwoTone', 3, 1, '', 0, 1);
INSERT INTO `sys_menu` VALUES (14, '2026-01-25 01:41:30.599', '2026-01-25 01:43:26.397', '2026-04-26 18:43:32.438', 13, '代码生成器', '/tool/generator', 'tool/generator/index', 'CodeTwoTone', 0, 2, 'tool:generator:list', 1, 0);
INSERT INTO `sys_menu` VALUES (15, '2026-01-25 02:09:15.803', '2026-01-25 02:12:16.656', '2026-01-25 02:14:37.631', 0, '产品分类', '/product-type', 'product-type/index', '', 0, 2, 'data:productType:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (16, '2026-01-25 14:42:33.918', '2026-01-25 14:42:33.918', NULL, 3, '用户修改', '', '', '', 0, 3, 'system:user:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (17, '2026-01-25 14:59:46.114', '2026-01-25 14:59:46.114', NULL, 6, 'api新增', '', '', '', 0, 3, 'system:api:add', 1, 0);
INSERT INTO `sys_menu` VALUES (18, '2026-01-25 15:00:02.997', '2026-01-25 15:00:02.997', NULL, 6, 'api同步', '', '', '', 0, 3, 'system:api:sync', 1, 0);
INSERT INTO `sys_menu` VALUES (19, '2026-01-25 15:00:29.712', '2026-01-25 15:00:39.380', NULL, 6, 'api编辑', '', '', '', 0, 3, 'system:api:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (20, '2026-01-25 15:04:03.700', '2026-01-25 15:04:03.700', NULL, 6, 'api删除', '', '', '', 0, 3, 'system:api:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (21, '2026-01-25 15:33:05.228', '2026-01-25 15:33:05.228', NULL, 3, '用户重置密码', '', '', '', 0, 3, 'system:user:resetPwd', 1, 0);
INSERT INTO `sys_menu` VALUES (22, '2026-01-25 15:33:24.656', '2026-01-25 15:33:24.656', NULL, 3, '用户删除', '', '', '', 0, 3, 'system:user:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (23, '2026-01-25 15:33:47.058', '2026-01-25 15:33:47.058', NULL, 3, '用户新增', '', '', '', 0, 3, 'system:user:add', 1, 0);
INSERT INTO `sys_menu` VALUES (24, '2026-01-25 17:25:31.612', '2026-01-25 17:34:45.886', NULL, 1, '文件管理', '/system/file', 'system/file/index', 'FileProtectOutlined', 6, 2, 'system:file:list', 1, 0);
INSERT INTO `sys_menu` VALUES (25, '2026-01-25 17:36:00.848', '2026-01-25 17:36:00.848', '2026-04-27 21:27:53.809', 1, '文件存储设置', '/system/storage', 'system/storage', 'AppstoreTwoTone', 7, 2, 'system:storage:list', 1, 0);
INSERT INTO `sys_menu` VALUES (26, '2026-01-27 00:35:23.946', '2026-01-27 00:35:23.946', '2026-01-27 00:50:58.109', 0, '产品类型', '/product_type', 'product_type/index', 'official-ProjectFilled', 0, 2, 'poductType:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (27, '2026-01-27 00:59:59.340', '2026-01-27 00:59:59.340', '2026-01-27 01:37:15.042', 0, '产品类型', '/product_type', 'product_type/index', 'official-ProjectFilled', 0, 2, 'poductType:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (28, '2026-01-27 01:01:39.818', '2026-01-27 01:01:39.818', '2026-01-27 01:02:30.270', 0, '产品类型', '/product_type', 'product_type/index', 'official-ProjectFilled', 0, 2, 'poductType:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (29, '2026-01-27 01:52:02.330', '2026-01-27 01:52:02.330', '2026-01-27 01:52:13.497', 0, '产品类型', '/product_type', 'product_type/index', 'official-ProjectFilled', 0, 2, 'poductType:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (30, '2026-01-27 02:36:49.151', '2026-01-27 02:36:49.151', '2026-01-27 02:44:28.705', 0, '产品类型', '/product_type', 'product_type/index', 'official-ProjectFilled', 0, 2, 'poductType:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (31, '2026-01-27 02:44:43.511', '2026-01-27 02:44:43.511', '2026-01-27 03:00:00.281', 0, '产品类型', '/product_type', 'product_type/index', 'official-ProjectFilled', 0, 2, 'poductType:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (32, '2026-01-27 03:00:18.875', '2026-01-27 03:00:18.875', '2026-01-27 03:06:34.954', 0, '产品类型', '/product_type', 'product_type/index', 'official-ProjectFilled', 0, 2, 'poductType:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (33, '2026-01-27 03:07:20.440', '2026-01-27 03:07:20.440', '2026-01-27 03:17:18.865', 0, '产品类型', '/product_type', 'product_type/index', 'official-ProjectFilled', 0, 2, 'poductType:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (34, '2026-01-27 03:18:16.860', '2026-01-27 03:18:16.860', '2026-02-01 01:27:06.847', 0, '产品类型', '/product_type', 'product_type/index', 'official-ProjectFilled', 0, 2, 'poductType:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (35, '2026-01-27 03:33:35.521', '2026-01-27 03:33:35.521', '2026-01-27 03:38:20.062', 0, '产品信息', '/product', 'product/index', 'official-PropertySafetyFilled', 0, 2, 'product:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (36, '2026-01-27 03:45:42.658', '2026-01-27 03:45:42.658', '2026-01-27 03:45:47.992', 0, '产品信息', '/product', 'product/index', 'official-PropertySafetyFilled', 0, 2, 'product:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (37, '2026-01-27 03:48:00.708', '2026-01-27 03:48:00.708', '2026-01-27 03:48:33.643', 0, '产品信息', '/product', 'product/index', 'official-PropertySafetyFilled', 0, 2, 'product:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (38, '2026-01-27 03:48:36.957', '2026-01-27 03:48:36.957', '2026-01-27 03:55:01.032', 0, '产品信息', '/product', 'product/index', 'official-PropertySafetyFilled', 0, 2, 'product:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (39, '2026-01-27 04:19:03.445', '2026-01-27 04:19:03.445', '2026-01-27 04:36:59.329', 0, '产品信息', '/product', 'product/index', 'official-PropertySafetyFilled', 0, 2, 'product:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (40, '2026-01-27 04:39:57.694', '2026-01-27 04:39:57.694', '2026-01-27 05:06:43.304', 0, '产品信息', '/product', 'product/index', 'official-PropertySafetyFilled', 0, 2, 'product:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (41, '2026-01-27 05:06:49.527', '2026-01-29 20:59:25.876', '2026-01-30 08:20:37.635', 0, '产品信息', '/product', 'product/index', 'custom-phone', 0, 2, 'product:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (42, '2026-01-30 08:20:44.346', '2026-01-30 08:20:44.346', '2026-01-30 08:57:58.413', 0, '产品信息', '/product', 'product/index', 'official-AppstoreTwoTone', 0, 2, 'product:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (43, '2026-01-30 08:58:01.810', '2026-01-30 08:58:01.810', '2026-02-01 01:26:57.600', 0, '产品信息', '/product', 'product/index', 'official-AppstoreTwoTone', 0, 2, 'product:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (44, '2026-02-02 09:50:23.606', '2026-02-02 09:50:23.606', '2026-02-05 06:27:45.842', 0, '信息管理', '/dataManger', '', 'BookTwoTone', 0, 1, '', 1, 0);
INSERT INTO `sys_menu` VALUES (45, '2026-02-02 09:50:50.771', '2026-02-02 09:50:50.771', '2026-02-02 09:53:56.417', 44, '题目分类', '/psy_category', 'psy_category/index', 'official-CalendarTwoTone', 0, 2, 'data:psyCategory:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (46, '2026-02-02 09:54:04.359', '2026-02-02 09:54:04.359', '2026-02-02 09:55:36.477', 44, '题目分类', '/psy_category', 'psy_category/index', 'official-CalendarTwoTone', 0, 2, 'data:psyCategory:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (47, '2026-02-02 09:55:39.449', '2026-02-04 08:36:52.721', '2026-02-05 06:27:44.410', 44, '分类信息', '/dataManger/psy_category', 'dataManger/psy_category/index', 'CalendarTwoTone', 0, 2, 'data:psyCategory:list:list', 1, 0);
INSERT INTO `sys_menu` VALUES (68, '2026-02-03 01:12:20.839', '2026-02-03 01:12:20.839', '2026-02-03 01:12:27.492', 44, '测试', '/demo', 'dataManger/demo/index', 'official-AimOutlined', 0, 2, 'data:test:list', 1, 0);
INSERT INTO `sys_menu` VALUES (69, '2026-02-03 01:12:20.939', '2026-02-03 01:12:20.939', '2026-02-03 01:12:27.446', 68, '查看', '', '', '', 1, 3, 'data:test:list', 1, 0);
INSERT INTO `sys_menu` VALUES (70, '2026-02-03 01:12:20.981', '2026-02-03 01:12:20.981', '2026-02-03 01:12:27.446', 68, '新增', '', '', '', 2, 3, 'data:test:add', 1, 0);
INSERT INTO `sys_menu` VALUES (71, '2026-02-03 01:12:21.030', '2026-02-03 01:12:21.030', '2026-02-03 01:12:27.446', 68, '编辑', '', '', '', 3, 3, 'data:test:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (72, '2026-02-03 01:12:21.157', '2026-02-03 01:12:21.157', '2026-02-03 01:12:27.446', 68, '删除', '', '', '', 4, 3, 'data:test:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (73, '2026-02-03 01:16:37.325', '2026-02-03 01:16:37.325', '2026-02-03 01:16:51.714', 44, '测试', '/demo', 'dataManger/demo/index', 'official-AimOutlined', 0, 2, 'data:test:list', 1, 0);
INSERT INTO `sys_menu` VALUES (74, '2026-02-03 01:16:37.356', '2026-02-03 01:16:37.356', '2026-02-03 01:16:51.653', 73, '查看', '', '', '', 1, 3, 'data:test:list', 1, 0);
INSERT INTO `sys_menu` VALUES (75, '2026-02-03 01:16:37.415', '2026-02-03 01:16:37.415', '2026-02-03 01:16:51.653', 73, '新增', '', '', '', 2, 3, 'data:test:add', 1, 0);
INSERT INTO `sys_menu` VALUES (76, '2026-02-03 01:16:37.456', '2026-02-03 01:16:37.456', '2026-02-03 01:16:51.653', 73, '编辑', '', '', '', 3, 3, 'data:test:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (77, '2026-02-03 01:16:37.498', '2026-02-03 01:16:37.498', '2026-02-03 01:16:51.653', 73, '删除', '', '', '', 4, 3, 'data:test:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (78, '2026-02-03 02:57:56.900', '2026-02-04 08:37:02.333', '2026-02-05 06:27:12.492', 44, '题库信息', '/psy_question', 'dataManger/psy_question/index', 'custom-psy_question', 0, 2, '', 1, 0);
INSERT INTO `sys_menu` VALUES (79, '2026-02-03 02:57:56.976', '2026-02-03 02:57:56.976', '2026-02-05 06:27:12.450', 78, '查看', '', '', '', 1, 3, 'data:psy_question:list', 1, 0);
INSERT INTO `sys_menu` VALUES (80, '2026-02-03 02:57:57.181', '2026-02-03 02:57:57.181', '2026-02-05 06:27:12.450', 78, '新增', '', '', '', 2, 3, 'data:psy_question:add', 1, 0);
INSERT INTO `sys_menu` VALUES (81, '2026-02-03 02:57:57.216', '2026-02-03 02:57:57.216', '2026-02-05 06:27:12.450', 78, '编辑', '', '', '', 3, 3, 'data:psy_question:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (82, '2026-02-03 02:57:57.257', '2026-02-03 02:57:57.257', '2026-02-05 06:27:12.450', 78, '删除', '', '', '', 4, 3, 'data:psy_question:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (83, '2026-02-04 03:08:34.584', '2026-02-04 03:08:34.584', '2026-02-04 03:19:17.069', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (84, '2026-02-04 03:08:34.641', '2026-02-04 03:08:34.641', '2026-02-04 03:19:17.017', 83, '查看', '', '', '', 1, 3, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (85, '2026-02-04 03:08:34.730', '2026-02-04 03:08:34.730', '2026-02-04 03:19:17.017', 83, '新增', '', '', '', 2, 3, 'data:paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (86, '2026-02-04 03:08:34.959', '2026-02-04 03:08:34.959', '2026-02-04 03:19:17.017', 83, '编辑', '', '', '', 3, 3, 'data:paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (87, '2026-02-04 03:08:35.057', '2026-02-04 03:08:35.057', '2026-02-04 03:19:17.017', 83, '删除', '', '', '', 4, 3, 'data:paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (88, '2026-02-04 03:34:15.680', '2026-02-04 03:34:15.680', '2026-02-04 04:00:18.984', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (89, '2026-02-04 03:34:16.075', '2026-02-04 03:34:16.075', '2026-02-04 04:00:18.917', 88, '查看', '', '', '', 1, 3, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (90, '2026-02-04 03:34:16.314', '2026-02-04 03:34:16.314', '2026-02-04 04:00:18.917', 88, '新增', '', '', '', 2, 3, 'data:paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (91, '2026-02-04 03:34:16.398', '2026-02-04 03:34:16.398', '2026-02-04 04:00:18.917', 88, '编辑', '', '', '', 3, 3, 'data:paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (92, '2026-02-04 03:34:16.473', '2026-02-04 03:34:16.473', '2026-02-04 04:00:18.917', 88, '删除', '', '', '', 4, 3, 'data:paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (93, '2026-02-04 04:00:23.278', '2026-02-04 04:00:23.278', '2026-02-04 04:15:15.496', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (94, '2026-02-04 04:00:23.319', '2026-02-04 04:00:23.319', '2026-02-04 04:15:15.328', 93, '查看', '', '', '', 1, 3, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (95, '2026-02-04 04:00:23.435', '2026-02-04 04:00:23.435', '2026-02-04 04:15:15.328', 93, '新增', '', '', '', 2, 3, 'data:paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (96, '2026-02-04 04:00:23.480', '2026-02-04 04:00:23.480', '2026-02-04 04:15:15.328', 93, '编辑', '', '', '', 3, 3, 'data:paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (97, '2026-02-04 04:00:23.533', '2026-02-04 04:00:23.533', '2026-02-04 04:15:15.328', 93, '删除', '', '', '', 4, 3, 'data:paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (98, '2026-02-04 04:19:05.594', '2026-02-04 04:19:05.594', '2026-02-04 04:56:00.062', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (99, '2026-02-04 04:19:05.635', '2026-02-04 04:19:05.635', '2026-02-04 04:56:00.004', 98, '查看', '', '', '', 1, 3, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (100, '2026-02-04 04:19:05.694', '2026-02-04 04:19:05.694', '2026-02-04 04:56:00.004', 98, '新增', '', '', '', 2, 3, 'data:paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (101, '2026-02-04 04:19:05.869', '2026-02-04 04:19:05.869', '2026-02-04 04:56:00.004', 98, '编辑', '', '', '', 3, 3, 'data:paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (102, '2026-02-04 04:19:06.336', '2026-02-04 04:19:06.336', '2026-02-04 04:56:00.004', 98, '删除', '', '', '', 4, 3, 'data:paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (103, '2026-02-04 04:59:34.522', '2026-02-04 04:59:34.522', '2026-02-04 05:34:10.924', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (104, '2026-02-04 04:59:34.604', '2026-02-04 04:59:34.604', '2026-02-04 05:34:10.874', 103, '查看', '', '', '', 1, 3, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (105, '2026-02-04 04:59:34.695', '2026-02-04 04:59:34.695', '2026-02-04 05:34:10.874', 103, '新增', '', '', '', 2, 3, 'data:paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (106, '2026-02-04 04:59:35.163', '2026-02-04 04:59:35.163', '2026-02-04 05:34:10.874', 103, '编辑', '', '', '', 3, 3, 'data:paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (107, '2026-02-04 04:59:35.196', '2026-02-04 04:59:35.196', '2026-02-04 05:34:10.874', 103, '删除', '', '', '', 4, 3, 'data:paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (108, '2026-02-04 05:34:13.828', '2026-02-04 05:34:13.828', '2026-02-04 06:19:26.446', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (109, '2026-02-04 05:34:13.908', '2026-02-04 05:34:13.908', '2026-02-04 06:19:26.397', 108, '查看', '', '', '', 1, 3, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (110, '2026-02-04 05:34:13.983', '2026-02-04 05:34:13.983', '2026-02-04 06:19:26.397', 108, '新增', '', '', '', 2, 3, 'data:paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (111, '2026-02-04 05:34:14.075', '2026-02-04 05:34:14.075', '2026-02-04 06:19:26.397', 108, '编辑', '', '', '', 3, 3, 'data:paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (112, '2026-02-04 05:34:14.148', '2026-02-04 05:34:14.148', '2026-02-04 06:19:26.397', 108, '删除', '', '', '', 4, 3, 'data:paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (113, '2026-02-04 06:19:29.350', '2026-02-04 06:19:29.350', '2026-02-04 07:13:45.901', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (114, '2026-02-04 06:19:29.382', '2026-02-04 06:19:29.382', '2026-02-04 07:13:45.751', 113, '查看', '', '', '', 1, 3, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (115, '2026-02-04 06:19:29.449', '2026-02-04 06:19:29.449', '2026-02-04 07:13:45.751', 113, '新增', '', '', '', 2, 3, 'data:paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (116, '2026-02-04 06:19:29.515', '2026-02-04 06:19:29.515', '2026-02-04 07:13:45.751', 113, '编辑', '', '', '', 3, 3, 'data:paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (117, '2026-02-04 06:19:29.565', '2026-02-04 06:19:29.565', '2026-02-04 07:13:45.751', 113, '删除', '', '', '', 4, 3, 'data:paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (118, '2026-02-04 07:14:34.425', '2026-02-04 07:14:34.425', '2026-02-04 07:27:31.374', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (119, '2026-02-04 07:14:34.589', '2026-02-04 07:14:34.589', '2026-02-04 07:27:31.332', 118, '查看', '', '', '', 1, 3, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (120, '2026-02-04 07:14:34.659', '2026-02-04 07:14:34.659', '2026-02-04 07:27:31.332', 118, '新增', '', '', '', 2, 3, 'data:paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (121, '2026-02-04 07:14:34.734', '2026-02-04 07:14:34.734', '2026-02-04 07:27:31.332', 118, '编辑', '', '', '', 3, 3, 'data:paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (122, '2026-02-04 07:14:34.796', '2026-02-04 07:14:34.796', '2026-02-04 07:27:31.332', 118, '删除', '', '', '', 4, 3, 'data:paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (123, '2026-02-04 07:27:34.017', '2026-02-04 07:27:34.017', '2026-02-04 07:33:12.222', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (124, '2026-02-04 07:27:34.075', '2026-02-04 07:27:34.075', '2026-02-04 07:33:12.181', 123, '查看', '', '', '', 1, 3, 'data:paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (125, '2026-02-04 07:27:34.192', '2026-02-04 07:27:34.192', '2026-02-04 07:33:12.181', 123, '新增', '', '', '', 2, 3, 'data:paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (126, '2026-02-04 07:27:34.291', '2026-02-04 07:27:34.291', '2026-02-04 07:33:12.181', 123, '编辑', '', '', '', 3, 3, 'data:paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (127, '2026-02-04 07:27:34.332', '2026-02-04 07:27:34.332', '2026-02-04 07:33:12.181', 123, '删除', '', '', '', 4, 3, 'data:paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (128, '2026-02-04 19:57:54.418', '2026-02-04 19:57:54.418', '2026-02-04 19:57:59.780', 44, '试卷信息', '/psy_paper1', 'dataManger/psy_paper1/index', 'official-AccountBookTwoTone', 0, 2, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (129, '2026-02-04 19:57:54.604', '2026-02-04 19:57:54.604', '2026-02-04 19:57:59.721', 128, '查看', '', '', '', 1, 3, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (130, '2026-02-04 19:57:54.904', '2026-02-04 19:57:54.904', '2026-02-04 19:57:59.721', 128, '新增', '', '', '', 2, 3, 'dataManger:psy_paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (131, '2026-02-04 19:57:54.944', '2026-02-04 19:57:54.944', '2026-02-04 19:57:59.721', 128, '编辑', '', '', '', 3, 3, 'dataManger:psy_paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (132, '2026-02-04 19:57:54.987', '2026-02-04 19:57:54.987', '2026-02-04 19:57:59.721', 128, '删除', '', '', '', 4, 3, 'dataManger:psy_paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (133, '2026-02-04 20:17:55.364', '2026-02-04 20:17:55.364', '2026-02-04 20:21:09.405', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (134, '2026-02-04 20:17:55.412', '2026-02-04 20:17:55.412', '2026-02-04 20:21:09.363', 133, '查看', '', '', '', 1, 3, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (135, '2026-02-04 20:17:55.510', '2026-02-04 20:17:55.510', '2026-02-04 20:21:09.363', 133, '新增', '', '', '', 2, 3, 'dataManger:psy_paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (136, '2026-02-04 20:17:55.820', '2026-02-04 20:17:55.820', '2026-02-04 20:21:09.363', 133, '编辑', '', '', '', 3, 3, 'dataManger:psy_paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (137, '2026-02-04 20:17:55.930', '2026-02-04 20:17:55.930', '2026-02-04 20:21:09.363', 133, '删除', '', '', '', 4, 3, 'dataManger:psy_paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (138, '2026-02-04 20:32:35.464', '2026-02-04 20:32:35.464', '2026-02-04 20:32:45.185', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (139, '2026-02-04 20:32:35.532', '2026-02-04 20:32:35.532', '2026-02-04 20:32:45.143', 138, '查看', '', '', '', 1, 3, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (140, '2026-02-04 20:32:35.580', '2026-02-04 20:32:35.580', '2026-02-04 20:32:45.143', 138, '新增', '', '', '', 2, 3, 'dataManger:psy_paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (141, '2026-02-04 20:32:35.639', '2026-02-04 20:32:35.639', '2026-02-04 20:32:45.143', 138, '编辑', '', '', '', 3, 3, 'dataManger:psy_paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (142, '2026-02-04 20:32:35.707', '2026-02-04 20:32:35.707', '2026-02-04 20:32:45.143', 138, '删除', '', '', '', 4, 3, 'dataManger:psy_paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (143, '2026-02-04 20:32:47.870', '2026-02-04 20:32:47.870', '2026-02-04 21:00:39.892', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (144, '2026-02-04 20:32:47.928', '2026-02-04 20:32:47.928', '2026-02-04 21:00:39.825', 143, '查看', '', '', '', 1, 3, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (145, '2026-02-04 20:32:47.980', '2026-02-04 20:32:47.980', '2026-02-04 21:00:39.825', 143, '新增', '', '', '', 2, 3, 'dataManger:psy_paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (146, '2026-02-04 20:32:48.036', '2026-02-04 20:32:48.036', '2026-02-04 21:00:39.825', 143, '编辑', '', '', '', 3, 3, 'dataManger:psy_paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (147, '2026-02-04 20:32:48.125', '2026-02-04 20:32:48.125', '2026-02-04 21:00:39.825', 143, '删除', '', '', '', 4, 3, 'dataManger:psy_paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (148, '2026-02-04 21:00:42.274', '2026-02-04 21:00:42.274', '2026-02-04 21:10:10.984', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (149, '2026-02-04 21:00:42.368', '2026-02-04 21:00:42.368', '2026-02-04 21:10:10.917', 148, '查看', '', '', '', 1, 3, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (150, '2026-02-04 21:00:42.419', '2026-02-04 21:00:42.419', '2026-02-04 21:10:10.917', 148, '新增', '', '', '', 2, 3, 'dataManger:psy_paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (151, '2026-02-04 21:00:42.476', '2026-02-04 21:00:42.476', '2026-02-04 21:10:10.917', 148, '编辑', '', '', '', 3, 3, 'dataManger:psy_paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (152, '2026-02-04 21:00:42.545', '2026-02-04 21:00:42.545', '2026-02-04 21:10:10.917', 148, '删除', '', '', '', 4, 3, 'dataManger:psy_paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (153, '2026-02-04 21:10:13.404', '2026-02-04 21:10:13.404', '2026-02-05 00:00:08.484', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (154, '2026-02-04 21:10:13.444', '2026-02-04 21:10:13.444', '2026-02-05 00:00:08.417', 153, '查看', '', '', '', 1, 3, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (155, '2026-02-04 21:10:13.551', '2026-02-04 21:10:13.551', '2026-02-05 00:00:08.417', 153, '新增', '', '', '', 2, 3, 'dataManger:psy_paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (156, '2026-02-04 21:10:13.611', '2026-02-04 21:10:13.611', '2026-02-05 00:00:08.417', 153, '编辑', '', '', '', 3, 3, 'dataManger:psy_paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (157, '2026-02-04 21:10:13.661', '2026-02-04 21:10:13.661', '2026-02-05 00:00:08.417', 153, '删除', '', '', '', 4, 3, 'dataManger:psy_paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (158, '2026-02-04 21:16:33.000', '2026-02-04 21:16:33.000', '2026-02-05 00:00:08.417', 153, '审批', '', '', '', 5, 3, 'dataManger:psy_paper:audit', 1, 0);
INSERT INTO `sys_menu` VALUES (159, '2026-02-05 00:03:12.708', '2026-02-05 00:03:12.708', '2026-02-05 00:11:32.808', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (160, '2026-02-05 00:03:12.802', '2026-02-05 00:03:12.802', '2026-02-05 00:11:32.704', 159, '查看', '', '', '', 1, 3, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (161, '2026-02-05 00:03:12.871', '2026-02-05 00:03:12.871', '2026-02-05 00:11:32.704', 159, '新增', '', '', '', 2, 3, 'dataManger:psy_paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (162, '2026-02-05 00:03:13.173', '2026-02-05 00:03:13.173', '2026-02-05 00:11:32.704', 159, '编辑', '', '', '', 3, 3, 'dataManger:psy_paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (163, '2026-02-05 00:03:13.218', '2026-02-05 00:03:13.218', '2026-02-05 00:11:32.704', 159, '删除', '', '', '', 4, 3, 'dataManger:psy_paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (164, '2026-02-05 00:03:13.286', '2026-02-05 00:03:13.286', '2026-02-05 00:11:32.704', 159, '审批', '', '', '', 5, 3, 'dataManger:psy_paper:audit', 1, 0);
INSERT INTO `sys_menu` VALUES (165, '2026-02-05 00:14:42.352', '2026-02-05 00:14:42.352', '2026-02-05 00:22:10.218', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (166, '2026-02-05 00:14:42.410', '2026-02-05 00:14:42.410', '2026-02-05 00:22:10.170', 165, '查看', '', '', '', 1, 3, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (167, '2026-02-05 00:14:42.452', '2026-02-05 00:14:42.452', '2026-02-05 00:22:10.170', 165, '新增', '', '', '', 2, 3, 'dataManger:psy_paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (168, '2026-02-05 00:14:42.493', '2026-02-05 00:14:42.493', '2026-02-05 00:22:10.170', 165, '编辑', '', '', '', 3, 3, 'dataManger:psy_paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (169, '2026-02-05 00:14:42.562', '2026-02-05 00:14:42.562', '2026-02-05 00:22:10.170', 165, '删除', '', '', '', 4, 3, 'dataManger:psy_paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (170, '2026-02-05 00:14:42.886', '2026-02-05 00:14:42.886', '2026-02-05 00:22:10.170', 165, '审批', '', '', '', 5, 3, 'dataManger:psy_paper:audit', 1, 0);
INSERT INTO `sys_menu` VALUES (171, '2026-02-05 00:22:15.505', '2026-02-05 00:22:15.505', '2026-02-05 00:22:39.221', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (172, '2026-02-05 00:22:15.571', '2026-02-05 00:22:15.571', '2026-02-05 00:22:39.176', 171, '查看', '', '', '', 1, 3, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (173, '2026-02-05 00:22:15.854', '2026-02-05 00:22:15.854', '2026-02-05 00:22:39.176', 171, '新增', '', '', '', 2, 3, 'dataManger:psy_paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (174, '2026-02-05 00:22:15.932', '2026-02-05 00:22:15.932', '2026-02-05 00:22:39.176', 171, '编辑', '', '', '', 3, 3, 'dataManger:psy_paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (175, '2026-02-05 00:22:16.013', '2026-02-05 00:22:16.013', '2026-02-05 00:22:39.176', 171, '删除', '', '', '', 4, 3, 'dataManger:psy_paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (176, '2026-02-05 00:22:16.063', '2026-02-05 00:22:16.063', '2026-02-05 00:22:39.176', 171, '审批', '', '', '', 5, 3, 'dataManger:psy_paper:audit', 1, 0);
INSERT INTO `sys_menu` VALUES (177, '2026-02-05 00:22:51.498', '2026-02-05 00:22:51.498', '2026-02-05 06:25:30.210', 44, '试卷信息', '/psy_paper', 'dataManger/psy_paper/index', 'official-AccountBookTwoTone', 0, 2, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (178, '2026-02-05 00:22:51.707', '2026-02-05 00:22:51.707', '2026-02-05 06:25:30.110', 177, '查看', '', '', '', 1, 3, 'dataManger:psy_paper:list', 1, 0);
INSERT INTO `sys_menu` VALUES (179, '2026-02-05 00:22:51.760', '2026-02-05 00:22:51.760', '2026-02-05 06:25:30.110', 177, '新增', '', '', '', 2, 3, 'dataManger:psy_paper:add', 1, 0);
INSERT INTO `sys_menu` VALUES (180, '2026-02-05 00:22:51.812', '2026-02-05 00:22:51.812', '2026-02-05 06:25:30.110', 177, '编辑', '', '', '', 3, 3, 'dataManger:psy_paper:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (181, '2026-02-05 00:22:51.852', '2026-02-05 00:22:51.852', '2026-02-05 06:25:30.110', 177, '删除', '', '', '', 4, 3, 'dataManger:psy_paper:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (182, '2026-02-05 00:22:51.919', '2026-02-05 00:22:51.919', '2026-02-05 06:25:30.110', 177, '审批', '', '', '', 5, 3, 'dataManger:psy_paper:audit', 1, 0);
INSERT INTO `sys_menu` VALUES (183, '2026-02-07 02:08:58.919', '2026-02-07 02:08:58.919', '2026-02-10 01:16:31.665', 0, '信息管理', '/dataManger', '', 'DatabaseTwoTone', 0, 1, '', 1, 0);
INSERT INTO `sys_menu` VALUES (184, '2026-02-07 02:09:22.941', '2026-02-07 02:09:22.941', '2026-02-10 01:15:02.170', 183, '牲畜分类', '/livestock_category', 'dataManger/livestock_category/index', 'official-AppstoreOutlined', 1, 2, 'dataManger:livestock_category:list', 1, 0);
INSERT INTO `sys_menu` VALUES (185, '2026-02-07 02:09:23.043', '2026-02-07 02:09:23.043', '2026-02-10 01:15:02.102', 184, '查看', '', '', '', 1, 3, 'dataManger:livestock_category:list', 1, 0);
INSERT INTO `sys_menu` VALUES (186, '2026-02-07 02:09:23.098', '2026-02-07 02:09:23.098', '2026-02-10 01:15:02.102', 184, '新增', '', '', '', 2, 3, 'dataManger:livestock_category:add', 1, 0);
INSERT INTO `sys_menu` VALUES (187, '2026-02-07 02:09:23.157', '2026-02-07 02:09:23.157', '2026-02-10 01:15:02.102', 184, '编辑', '', '', '', 3, 3, 'dataManger:livestock_category:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (188, '2026-02-07 02:09:23.268', '2026-02-07 02:09:23.268', '2026-02-10 01:15:02.102', 184, '删除', '', '', '', 4, 3, 'dataManger:livestock_category:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (189, '2026-02-07 02:45:02.603', '2026-02-07 02:45:02.603', '2026-02-07 02:47:04.494', 183, '养殖区域', '/breeding_area', 'dataManger/breeding_area/index', 'official-HomeOutlined', 2, 2, 'dataManger:breeding_area:list', 1, 0);
INSERT INTO `sys_menu` VALUES (190, '2026-02-07 02:45:02.644', '2026-02-07 02:45:02.644', '2026-02-07 02:47:04.412', 189, '查看', '', '', '', 1, 3, 'dataManger:breeding_area:list', 1, 0);
INSERT INTO `sys_menu` VALUES (191, '2026-02-07 02:45:02.735', '2026-02-07 02:45:02.735', '2026-02-07 02:47:04.412', 189, '新增', '', '', '', 2, 3, 'dataManger:breeding_area:add', 1, 0);
INSERT INTO `sys_menu` VALUES (192, '2026-02-07 02:45:02.794', '2026-02-07 02:45:02.794', '2026-02-07 02:47:04.412', 189, '编辑', '', '', '', 3, 3, 'dataManger:breeding_area:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (193, '2026-02-07 02:45:02.853', '2026-02-07 02:45:02.853', '2026-02-07 02:47:04.412', 189, '删除', '', '', '', 4, 3, 'dataManger:breeding_area:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (194, '2026-02-07 02:47:06.965', '2026-02-07 02:47:06.965', '2026-02-08 01:09:55.447', 183, '养殖区域', '/breeding_area', 'dataManger/breeding_area/index', 'official-HomeOutlined', 2, 2, 'dataManger:breeding_area:list', 1, 0);
INSERT INTO `sys_menu` VALUES (195, '2026-02-07 02:47:07.021', '2026-02-07 02:47:07.021', '2026-02-08 01:09:55.414', 194, '查看', '', '', '', 1, 3, 'dataManger:breeding_area:list', 1, 0);
INSERT INTO `sys_menu` VALUES (196, '2026-02-07 02:47:07.243', '2026-02-07 02:47:07.243', '2026-02-08 01:09:55.414', 194, '新增', '', '', '', 2, 3, 'dataManger:breeding_area:add', 1, 0);
INSERT INTO `sys_menu` VALUES (197, '2026-02-07 02:47:07.346', '2026-02-07 02:47:07.346', '2026-02-08 01:09:55.414', 194, '编辑', '', '', '', 3, 3, 'dataManger:breeding_area:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (198, '2026-02-07 02:47:07.471', '2026-02-07 02:47:07.471', '2026-02-08 01:09:55.414', 194, '删除', '', '', '', 4, 3, 'dataManger:breeding_area:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (199, '2026-02-07 17:20:06.181', '2026-02-07 17:20:06.181', '2026-02-07 17:35:17.245', 1, '养殖户认证', '/farmer_certification', 'Layout/farmer_certification/index', 'custom-farmer', 9, 2, 'system:farmer_certification:list', 1, 0);
INSERT INTO `sys_menu` VALUES (200, '2026-02-07 17:20:06.240', '2026-02-07 17:20:06.240', '2026-02-07 17:35:17.187', 199, '查看', '', '', '', 1, 3, 'system:farmer_certification:list', 1, 0);
INSERT INTO `sys_menu` VALUES (201, '2026-02-07 17:20:06.330', '2026-02-07 17:20:06.330', '2026-02-07 17:35:17.187', 199, '新增', '', '', '', 2, 3, 'system:farmer_certification:add', 1, 0);
INSERT INTO `sys_menu` VALUES (202, '2026-02-07 17:20:06.406', '2026-02-07 17:20:06.406', '2026-02-07 17:35:17.187', 199, '编辑', '', '', '', 3, 3, 'system:farmer_certification:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (203, '2026-02-07 17:20:06.481', '2026-02-07 17:20:06.481', '2026-02-07 17:35:17.187', 199, '删除', '', '', '', 4, 3, 'system:farmer_certification:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (204, '2026-02-07 17:20:06.548', '2026-02-07 17:20:06.548', '2026-02-07 17:35:17.187', 199, '审批', '', '', '', 5, 3, 'system:farmer_certification:audit', 1, 0);
INSERT INTO `sys_menu` VALUES (205, '2026-02-07 17:35:21.376', '2026-02-07 17:35:21.376', '2026-02-10 01:14:57.417', 1, '养殖户认证', '/farmer_certification', 'system/farmer_certification/index', 'custom-farmer', 9, 2, 'system:farmer_certification:list', 1, 0);
INSERT INTO `sys_menu` VALUES (206, '2026-02-07 17:35:21.441', '2026-02-07 17:35:21.441', '2026-02-10 01:14:57.359', 205, '查看', '', '', '', 1, 3, 'system:farmer_certification:list', 1, 0);
INSERT INTO `sys_menu` VALUES (207, '2026-02-07 17:35:21.514', '2026-02-07 17:35:21.514', '2026-02-10 01:14:57.359', 205, '新增', '', '', '', 2, 3, 'system:farmer_certification:add', 1, 0);
INSERT INTO `sys_menu` VALUES (208, '2026-02-07 17:35:21.565', '2026-02-07 17:35:21.565', '2026-02-10 01:14:57.359', 205, '编辑', '', '', '', 3, 3, 'system:farmer_certification:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (209, '2026-02-07 17:35:21.661', '2026-02-07 17:35:21.661', '2026-02-10 01:14:57.359', 205, '删除', '', '', '', 4, 3, 'system:farmer_certification:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (210, '2026-02-07 17:35:21.777', '2026-02-07 17:35:21.777', '2026-02-10 01:14:57.359', 205, '审批', '', '', '', 5, 3, 'system:farmer_certification:audit', 1, 0);
INSERT INTO `sys_menu` VALUES (211, '2026-02-08 00:55:58.528', '2026-02-08 00:58:17.997', NULL, 1, '字典管理', '/system/dict', 'system/dict/index', 'AntDesignOutlined', 3, 2, 'system:dict:list', 1, 0);
INSERT INTO `sys_menu` VALUES (212, '2026-02-08 00:57:06.342', '2026-02-08 00:57:06.342', '2026-02-08 00:58:10.037', 3, '字典管理', '', '', '', 3, 2, '', 1, 0);
INSERT INTO `sys_menu` VALUES (213, '2026-02-08 01:10:00.039', '2026-02-08 01:10:00.039', '2026-02-08 02:09:12.246', 183, '养殖区域', '/breeding_area', 'dataManger/breeding_area/index', 'official-HomeOutlined', 2, 2, 'dataManger:breeding_area:list', 1, 0);
INSERT INTO `sys_menu` VALUES (214, '2026-02-08 01:10:00.112', '2026-02-08 01:10:00.112', '2026-02-08 02:09:12.204', 213, '查看', '', '', '', 1, 3, 'dataManger:breeding_area:list', 1, 0);
INSERT INTO `sys_menu` VALUES (215, '2026-02-08 01:10:00.165', '2026-02-08 01:10:00.165', '2026-02-08 02:09:12.204', 213, '新增', '', '', '', 2, 3, 'dataManger:breeding_area:add', 1, 0);
INSERT INTO `sys_menu` VALUES (216, '2026-02-08 01:10:00.206', '2026-02-08 01:10:00.206', '2026-02-08 02:09:12.204', 213, '编辑', '', '', '', 3, 3, 'dataManger:breeding_area:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (217, '2026-02-08 01:10:00.249', '2026-02-08 01:10:00.249', '2026-02-08 02:09:12.204', 213, '删除', '', '', '', 4, 3, 'dataManger:breeding_area:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (218, '2026-02-09 17:14:52.380', '2026-02-09 17:14:52.380', '2026-02-09 17:53:41.350', 183, '养殖区域', '/breeding_area', 'dataManger/breeding_area/index', 'official-HomeOutlined', 2, 2, 'dataManger:breeding_area:list', 1, 0);
INSERT INTO `sys_menu` VALUES (219, '2026-02-09 17:14:52.437', '2026-02-09 17:14:52.437', '2026-02-09 17:53:41.099', 218, '查看', '', '', '', 1, 3, 'dataManger:breeding_area:list', 1, 0);
INSERT INTO `sys_menu` VALUES (220, '2026-02-09 17:14:52.520', '2026-02-09 17:14:52.520', '2026-02-09 17:53:41.099', 218, '新增', '', '', '', 2, 3, 'dataManger:breeding_area:add', 1, 0);
INSERT INTO `sys_menu` VALUES (221, '2026-02-09 17:14:52.578', '2026-02-09 17:14:52.578', '2026-02-09 17:53:41.099', 218, '编辑', '', '', '', 3, 3, 'dataManger:breeding_area:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (222, '2026-02-09 17:14:52.620', '2026-02-09 17:14:52.620', '2026-02-09 17:53:41.099', 218, '删除', '', '', '', 4, 3, 'dataManger:breeding_area:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (223, '2026-02-09 17:53:44.632', '2026-02-09 17:53:44.632', '2026-02-10 01:14:59.652', 183, '养殖区域', '/breeding_area', 'dataManger/breeding_area/index', 'official-HomeOutlined', 2, 2, 'dataManger:breeding_area:list', 1, 0);
INSERT INTO `sys_menu` VALUES (224, '2026-02-09 17:53:44.810', '2026-02-09 17:53:44.810', '2026-02-10 01:14:59.610', 223, '查看', '', '', '', 1, 3, 'dataManger:breeding_area:list', 1, 0);
INSERT INTO `sys_menu` VALUES (225, '2026-02-09 17:53:44.939', '2026-02-09 17:53:44.939', '2026-02-10 01:14:59.610', 223, '新增', '', '', '', 2, 3, 'dataManger:breeding_area:add', 1, 0);
INSERT INTO `sys_menu` VALUES (226, '2026-02-09 17:53:45.265', '2026-02-09 17:53:45.265', '2026-02-10 01:14:59.610', 223, '编辑', '', '', '', 3, 3, 'dataManger:breeding_area:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (227, '2026-02-09 17:53:45.458', '2026-02-09 17:53:45.458', '2026-02-10 01:14:59.610', 223, '删除', '', '', '', 4, 3, 'dataManger:breeding_area:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (228, '2026-02-10 22:21:03.301', '2026-02-10 22:21:03.301', '2026-02-10 22:27:48.940', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (229, '2026-02-10 22:21:03.349', '2026-02-10 22:21:03.349', '2026-02-10 22:27:48.864', 228, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (230, '2026-02-10 22:21:03.397', '2026-02-10 22:21:03.397', '2026-02-10 22:27:48.864', 228, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (231, '2026-02-10 22:21:03.457', '2026-02-10 22:21:03.457', '2026-02-10 22:27:48.864', 228, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (232, '2026-02-10 22:21:03.863', '2026-02-10 22:21:03.863', '2026-02-10 22:27:48.864', 228, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (233, '2026-02-10 22:22:39.341', '2026-02-10 22:22:39.341', '2026-02-10 22:34:45.085', 0, '产品类型', '/productType', 'productType/index', '', 0, 2, 'productType:list', 1, 0);
INSERT INTO `sys_menu` VALUES (234, '2026-02-10 22:22:39.469', '2026-02-10 22:22:39.469', '2026-02-10 22:34:44.999', 233, '查看', '', '', '', 1, 3, 'productType:list', 1, 0);
INSERT INTO `sys_menu` VALUES (235, '2026-02-10 22:22:39.536', '2026-02-10 22:22:39.536', '2026-02-10 22:34:44.999', 233, '新增', '', '', '', 2, 3, 'productType:add', 1, 0);
INSERT INTO `sys_menu` VALUES (236, '2026-02-10 22:22:39.570', '2026-02-10 22:22:39.570', '2026-02-10 22:34:44.999', 233, '编辑', '', '', '', 3, 3, 'productType:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (237, '2026-02-10 22:22:39.646', '2026-02-10 22:22:39.646', '2026-02-10 22:34:44.999', 233, '删除', '', '', '', 4, 3, 'productType:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (238, '2026-02-10 22:28:56.327', '2026-02-10 22:28:56.327', '2026-02-10 22:34:35.226', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (239, '2026-02-10 22:28:56.428', '2026-02-10 22:28:56.428', '2026-02-10 22:34:35.123', 238, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (240, '2026-02-10 22:28:56.850', '2026-02-10 22:28:56.850', '2026-02-10 22:34:35.123', 238, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (241, '2026-02-10 22:28:57.086', '2026-02-10 22:28:57.086', '2026-02-10 22:34:35.123', 238, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (242, '2026-02-10 22:28:57.226', '2026-02-10 22:28:57.226', '2026-02-10 22:34:35.123', 238, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (243, '2026-02-10 22:34:59.358', '2026-02-10 22:34:59.358', '2026-02-10 23:11:01.775', 0, '产品类型', '/product_type', 'product_type/index', '', 0, 2, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (244, '2026-02-10 22:34:59.447', '2026-02-10 22:34:59.447', '2026-02-10 23:10:56.660', 243, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (245, '2026-02-10 22:34:59.541', '2026-02-10 22:34:59.541', '2026-02-10 23:10:57.967', 243, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0);
INSERT INTO `sys_menu` VALUES (246, '2026-02-10 22:34:59.596', '2026-02-10 22:34:59.596', '2026-02-10 23:10:59.136', 243, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (247, '2026-02-10 22:34:59.638', '2026-02-10 22:34:59.638', '2026-02-10 23:11:00.355', 243, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (248, '2026-02-10 22:35:01.542', '2026-02-10 22:35:01.542', '2026-02-10 23:07:43.715', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (249, '2026-02-10 22:35:01.647', '2026-02-10 22:35:01.647', '2026-02-10 23:07:43.616', 248, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (250, '2026-02-10 22:35:01.706', '2026-02-10 22:35:01.706', '2026-02-10 23:07:43.616', 248, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (251, '2026-02-10 22:35:01.813', '2026-02-10 22:35:01.813', '2026-02-10 23:07:43.616', 248, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (252, '2026-02-10 22:35:01.890', '2026-02-10 22:35:01.890', '2026-02-10 23:07:43.616', 248, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (253, '2026-02-10 23:07:50.734', '2026-02-10 23:07:50.734', '2026-02-10 23:10:38.609', 0, '产品类型', '/productType', 'productType/index', '', 0, 2, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (254, '2026-02-10 23:07:50.859', '2026-02-10 23:07:50.859', '2026-02-10 23:10:38.543', 253, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (255, '2026-02-10 23:07:51.076', '2026-02-10 23:07:51.076', '2026-02-10 23:10:38.543', 253, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0);
INSERT INTO `sys_menu` VALUES (256, '2026-02-10 23:07:51.339', '2026-02-10 23:07:51.339', '2026-02-10 23:10:38.543', 253, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (257, '2026-02-10 23:07:51.716', '2026-02-10 23:07:51.716', '2026-02-10 23:10:38.543', 253, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (258, '2026-02-10 23:07:53.011', '2026-02-10 23:07:53.011', '2026-02-10 23:10:35.583', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (259, '2026-02-10 23:07:53.109', '2026-02-10 23:07:53.109', '2026-02-10 23:10:35.525', 258, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (260, '2026-02-10 23:07:53.282', '2026-02-10 23:07:53.282', '2026-02-10 23:10:35.525', 258, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (261, '2026-02-10 23:07:53.412', '2026-02-10 23:07:53.412', '2026-02-10 23:10:35.525', 258, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (262, '2026-02-10 23:07:53.484', '2026-02-10 23:07:53.484', '2026-02-10 23:10:35.525', 258, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (263, '2026-02-10 23:11:14.198', '2026-02-10 23:11:14.198', '2026-02-11 00:41:55.191', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (264, '2026-02-10 23:11:14.346', '2026-02-10 23:11:14.346', '2026-02-11 00:41:55.124', 263, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (265, '2026-02-10 23:11:14.412', '2026-02-10 23:11:14.412', '2026-02-11 00:41:55.124', 263, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (266, '2026-02-10 23:11:14.475', '2026-02-10 23:11:14.475', '2026-02-11 00:41:55.124', 263, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (267, '2026-02-10 23:11:15.062', '2026-02-10 23:11:15.062', '2026-02-11 00:41:55.124', 263, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (268, '2026-02-10 23:11:16.057', '2026-02-10 23:11:16.057', '2026-02-12 22:41:05.319', 0, '产品类型', '/productType', 'productType/index', '', 0, 2, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (269, '2026-02-10 23:11:16.158', '2026-02-10 23:11:16.158', '2026-02-12 22:41:05.252', 268, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (270, '2026-02-10 23:11:16.229', '2026-02-10 23:11:16.229', '2026-02-12 22:41:05.252', 268, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0);
INSERT INTO `sys_menu` VALUES (271, '2026-02-10 23:11:16.299', '2026-02-10 23:11:16.299', '2026-02-12 22:41:05.252', 268, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (272, '2026-02-10 23:11:16.372', '2026-02-10 23:11:16.372', '2026-02-12 22:41:05.252', 268, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (273, '2026-02-11 00:41:58.047', '2026-02-11 00:41:58.047', '2026-02-11 00:44:35.268', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (274, '2026-02-11 00:41:58.110', '2026-02-11 00:41:58.110', '2026-02-11 00:44:35.217', 273, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (275, '2026-02-11 00:41:58.210', '2026-02-11 00:41:58.210', '2026-02-11 00:44:35.217', 273, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (276, '2026-02-11 00:41:58.328', '2026-02-11 00:41:58.328', '2026-02-11 00:44:35.217', 273, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (277, '2026-02-11 00:41:58.402', '2026-02-11 00:41:58.402', '2026-02-11 00:44:35.217', 273, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (278, '2026-02-11 00:44:37.662', '2026-02-11 00:44:37.662', '2026-02-11 00:57:03.186', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (279, '2026-02-11 00:44:37.769', '2026-02-11 00:44:37.769', '2026-02-11 00:57:03.127', 278, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (280, '2026-02-11 00:44:37.877', '2026-02-11 00:44:37.877', '2026-02-11 00:57:03.127', 278, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (281, '2026-02-11 00:44:37.911', '2026-02-11 00:44:37.911', '2026-02-11 00:57:03.127', 278, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (282, '2026-02-11 00:44:37.977', '2026-02-11 00:44:37.977', '2026-02-11 00:57:03.127', 278, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (283, '2026-02-11 00:57:07.255', '2026-02-11 00:57:07.255', '2026-02-11 01:38:04.155', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (284, '2026-02-11 00:57:07.287', '2026-02-11 00:57:07.287', '2026-02-11 01:38:04.099', 283, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (285, '2026-02-11 00:57:07.329', '2026-02-11 00:57:07.329', '2026-02-11 01:38:04.099', 283, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (286, '2026-02-11 00:57:07.418', '2026-02-11 00:57:07.418', '2026-02-11 01:38:04.099', 283, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (287, '2026-02-11 00:57:07.463', '2026-02-11 00:57:07.463', '2026-02-11 01:38:04.099', 283, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (288, '2026-02-11 01:38:06.419', '2026-02-11 01:38:06.419', '2026-02-11 01:49:27.087', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (289, '2026-02-11 01:38:06.464', '2026-02-11 01:38:06.464', '2026-02-11 01:49:27.038', 288, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (290, '2026-02-11 01:38:06.514', '2026-02-11 01:38:06.514', '2026-02-11 01:49:27.038', 288, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (291, '2026-02-11 01:38:06.583', '2026-02-11 01:38:06.583', '2026-02-11 01:49:27.038', 288, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (292, '2026-02-11 01:38:06.648', '2026-02-11 01:38:06.648', '2026-02-11 01:49:27.038', 288, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (293, '2026-02-11 01:49:41.635', '2026-02-11 01:49:41.635', '2026-02-12 22:41:01.347', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (294, '2026-02-11 01:49:41.811', '2026-02-11 01:49:41.811', '2026-02-12 22:41:01.101', 293, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (295, '2026-02-11 01:49:41.844', '2026-02-11 01:49:41.844', '2026-02-12 22:41:01.101', 293, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (296, '2026-02-11 01:49:41.878', '2026-02-11 01:49:41.878', '2026-02-12 22:41:01.101', 293, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (297, '2026-02-11 01:49:41.911', '2026-02-11 01:49:41.911', '2026-02-12 22:41:01.101', 293, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (298, '2026-02-11 02:16:31.159', '2026-02-11 02:16:31.159', '2026-02-11 02:16:46.030', 0, '33', '333', '', 'AccountBookOutlined', 0, 1, '', 1, 0);
INSERT INTO `sys_menu` VALUES (299, '2026-02-14 20:57:27.509', '2026-02-14 20:57:27.509', '2026-02-14 21:02:25.411', 0, '产品类型', '/productType', 'productType/index', '', 0, 2, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (300, '2026-02-14 20:57:27.606', '2026-02-14 20:57:27.606', '2026-02-14 21:02:25.370', 299, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (301, '2026-02-14 20:57:27.681', '2026-02-14 20:57:27.681', '2026-02-14 21:02:25.370', 299, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0);
INSERT INTO `sys_menu` VALUES (302, '2026-02-14 20:57:27.721', '2026-02-14 20:57:27.721', '2026-02-14 21:02:25.370', 299, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (303, '2026-02-14 20:57:27.775', '2026-02-14 20:57:27.775', '2026-02-14 21:02:25.370', 299, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (304, '2026-02-14 20:57:29.771', '2026-02-14 20:57:29.771', '2026-02-14 21:02:19.989', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (305, '2026-02-14 20:57:29.896', '2026-02-14 20:57:29.896', '2026-02-14 21:02:19.884', 304, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (306, '2026-02-14 20:57:29.995', '2026-02-14 20:57:29.995', '2026-02-14 21:02:19.884', 304, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (307, '2026-02-14 20:57:30.509', '2026-02-14 20:57:30.509', '2026-02-14 21:02:19.884', 304, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (308, '2026-02-14 20:57:30.615', '2026-02-14 20:57:30.615', '2026-02-14 21:02:19.884', 304, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (309, '2026-02-14 21:02:27.214', '2026-02-14 21:02:27.214', '2026-02-14 22:03:32.147', 0, '产品类型', '/productType', 'productType/index', '', 0, 2, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (310, '2026-02-14 21:02:27.281', '2026-02-14 21:02:27.281', '2026-02-14 22:03:32.096', 309, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (311, '2026-02-14 21:02:27.339', '2026-02-14 21:02:27.339', '2026-02-14 22:03:32.096', 309, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0);
INSERT INTO `sys_menu` VALUES (312, '2026-02-14 21:02:27.399', '2026-02-14 21:02:27.399', '2026-02-14 22:03:32.096', 309, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (313, '2026-02-14 21:02:27.449', '2026-02-14 21:02:27.449', '2026-02-14 22:03:32.096', 309, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (314, '2026-02-14 22:04:53.237', '2026-02-14 22:04:53.237', '2026-02-14 22:04:55.253', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (315, '2026-02-14 22:04:53.395', '2026-02-14 22:04:53.395', '2026-02-14 22:04:55.215', 314, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (316, '2026-02-14 22:04:53.478', '2026-02-14 22:04:53.478', '2026-02-14 22:04:55.215', 314, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (317, '2026-02-14 22:04:53.528', '2026-02-14 22:04:53.528', '2026-02-14 22:04:55.215', 314, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (318, '2026-02-14 22:04:53.561', '2026-02-14 22:04:53.561', '2026-02-14 22:04:55.215', 314, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (319, '2026-02-14 22:04:57.032', '2026-02-14 22:04:57.032', '2026-02-14 22:28:18.509', 0, '产品类型', '/productType', 'productType/index', '', 0, 2, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (320, '2026-02-14 22:04:57.129', '2026-02-14 22:04:57.129', '2026-02-14 22:28:18.397', 319, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (321, '2026-02-14 22:04:57.214', '2026-02-14 22:04:57.214', '2026-02-14 22:28:18.397', 319, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0);
INSERT INTO `sys_menu` VALUES (322, '2026-02-14 22:04:57.387', '2026-02-14 22:04:57.387', '2026-02-14 22:28:18.397', 319, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (323, '2026-02-14 22:04:57.454', '2026-02-14 22:04:57.454', '2026-02-14 22:28:18.397', 319, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (324, '2026-02-14 22:28:20.915', '2026-02-14 22:28:20.915', '2026-02-14 22:53:01.013', 0, '产品类型', '/productType', 'productType/index', '', 0, 2, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (325, '2026-02-14 22:28:20.980', '2026-02-14 22:28:20.980', '2026-02-14 22:53:00.972', 324, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (326, '2026-02-14 22:28:21.039', '2026-02-14 22:28:21.039', '2026-02-14 22:53:00.972', 324, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0);
INSERT INTO `sys_menu` VALUES (327, '2026-02-14 22:28:21.140', '2026-02-14 22:28:21.140', '2026-02-14 22:53:00.972', 324, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (328, '2026-02-14 22:28:21.239', '2026-02-14 22:28:21.239', '2026-02-14 22:53:00.972', 324, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (329, '2026-02-14 22:53:02.999', '2026-02-14 22:53:02.999', '2026-02-14 23:04:22.747', 0, '产品类型', '/productType', 'productType/index', '', 0, 2, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (330, '2026-02-14 22:53:03.140', '2026-02-14 22:53:03.140', '2026-02-14 23:04:21.385', 329, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (331, '2026-02-14 22:53:03.181', '2026-02-14 22:53:03.181', '2026-02-14 23:04:20.041', 329, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0);
INSERT INTO `sys_menu` VALUES (332, '2026-02-14 22:53:03.311', '2026-02-14 22:53:03.311', '2026-02-14 23:04:18.641', 329, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (333, '2026-02-14 22:53:03.373', '2026-02-14 22:53:03.373', '2026-02-14 23:04:17.195', 329, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (334, '2026-02-14 23:04:49.000', '2026-02-14 23:04:49.000', '2026-02-14 23:05:31.072', 0, '产品类型', '/productType', 'productType/index', '', 0, 2, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (335, '2026-02-14 23:04:49.000', '2026-02-14 23:04:49.000', '2026-02-14 23:05:31.020', 334, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (336, '2026-02-14 23:04:49.000', '2026-02-14 23:04:49.000', '2026-02-14 23:05:31.020', 334, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0);
INSERT INTO `sys_menu` VALUES (337, '2026-02-14 23:04:49.000', '2026-02-14 23:04:49.000', '2026-02-14 23:05:31.020', 334, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (338, '2026-02-14 23:04:49.000', '2026-02-14 23:04:49.000', '2026-02-14 23:05:31.020', 334, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (339, '2026-02-14 23:04:49.000', '2026-02-14 23:04:49.000', '2026-02-14 23:05:31.020', 334, '导出', '', '', '', 5, 3, 'product_type:export', 1, 0);
INSERT INTO `sys_menu` VALUES (340, '2026-02-14 23:04:49.000', '2026-02-14 23:04:49.000', '2026-02-14 23:05:31.020', 334, '导入', '', '', '', 6, 3, 'product_type:import', 1, 0);
INSERT INTO `sys_menu` VALUES (341, '2026-02-14 23:25:10.793', '2026-02-14 23:25:10.793', '2026-02-14 23:27:52.951', 0, '产品类型', '/productType', 'productType/index', '', 0, 2, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (342, '2026-02-14 23:25:10.981', '2026-02-14 23:25:10.981', '2026-02-14 23:27:52.902', 341, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (343, '2026-02-14 23:25:11.028', '2026-02-14 23:25:11.028', '2026-02-14 23:27:52.902', 341, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0);
INSERT INTO `sys_menu` VALUES (344, '2026-02-14 23:25:11.078', '2026-02-14 23:25:11.078', '2026-02-14 23:27:52.902', 341, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (345, '2026-02-14 23:25:11.144', '2026-02-14 23:25:11.144', '2026-02-14 23:27:52.902', 341, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (346, '2026-02-14 23:25:11.195', '2026-02-14 23:25:11.195', '2026-02-14 23:27:52.902', 341, '导出', '', '', '', 5, 3, 'product_type:export', 1, 0);
INSERT INTO `sys_menu` VALUES (347, '2026-02-14 23:25:11.262', '2026-02-14 23:25:11.262', '2026-02-14 23:27:52.902', 341, '导入', '', '', '', 6, 3, 'product_type:import', 1, 0);
INSERT INTO `sys_menu` VALUES (348, '2026-02-14 23:27:55.068', '2026-02-14 23:27:55.068', '2026-02-14 23:29:13.715', 0, '产品类型', '/productType', 'productType/index', '', 0, 2, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (349, '2026-02-14 23:27:55.112', '2026-02-14 23:27:55.112', '2026-02-14 23:29:13.672', 348, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (350, '2026-02-14 23:27:55.212', '2026-02-14 23:27:55.212', '2026-02-14 23:29:13.672', 348, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0);
INSERT INTO `sys_menu` VALUES (351, '2026-02-14 23:27:55.427', '2026-02-14 23:27:55.427', '2026-02-14 23:29:13.672', 348, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (352, '2026-02-14 23:27:55.535', '2026-02-14 23:27:55.535', '2026-02-14 23:29:13.672', 348, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (353, '2026-02-14 23:27:55.602', '2026-02-14 23:27:55.602', '2026-02-14 23:29:13.672', 348, '导出', '', '', '', 5, 3, 'product_type:export', 1, 0);
INSERT INTO `sys_menu` VALUES (354, '2026-02-14 23:27:55.651', '2026-02-14 23:27:55.651', '2026-02-14 23:29:13.672', 348, '导入', '', '', '', 6, 3, 'product_type:import', 1, 0);
INSERT INTO `sys_menu` VALUES (355, '2026-02-14 23:29:16.312', '2026-02-14 23:29:16.312', '2026-02-17 17:45:01.404', 0, '产品类型', '/productType', 'productType/index', '', 0, 2, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (356, '2026-02-14 23:29:16.402', '2026-02-14 23:29:16.402', '2026-02-17 17:45:01.401', 355, '查看', '', '', '', 1, 3, 'product_type:list', 1, 0);
INSERT INTO `sys_menu` VALUES (357, '2026-02-14 23:29:16.478', '2026-02-14 23:29:16.478', '2026-02-17 17:45:01.401', 355, '新增', '', '', '', 2, 3, 'product_type:add', 1, 0);
INSERT INTO `sys_menu` VALUES (358, '2026-02-14 23:29:16.638', '2026-02-14 23:29:16.638', '2026-02-17 17:45:01.401', 355, '编辑', '', '', '', 3, 3, 'product_type:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (359, '2026-02-14 23:29:16.685', '2026-02-14 23:29:16.685', '2026-02-17 17:45:01.401', 355, '删除', '', '', '', 4, 3, 'product_type:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (360, '2026-02-14 23:29:16.718', '2026-02-14 23:29:16.718', '2026-02-17 17:45:01.401', 355, '导出', '', '', '', 5, 3, 'product_type:export', 1, 0);
INSERT INTO `sys_menu` VALUES (361, '2026-02-14 23:29:16.768', '2026-02-14 23:29:16.768', '2026-02-17 17:45:01.401', 355, '导入', '', '', '', 6, 3, 'product_type:import', 1, 0);
INSERT INTO `sys_menu` VALUES (362, '2026-02-14 23:32:38.612', '2026-02-14 23:32:38.612', '2026-02-17 17:44:58.547', 0, '产品信息', '/product', 'product/index', 'official-AccountBookFilled', 0, 2, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (363, '2026-02-14 23:32:38.705', '2026-02-14 23:32:38.705', '2026-02-17 17:44:58.542', 362, '查看', '', '', '', 1, 3, 'product:list', 1, 0);
INSERT INTO `sys_menu` VALUES (364, '2026-02-14 23:32:38.772', '2026-02-14 23:32:38.772', '2026-02-17 17:44:58.542', 362, '新增', '', '', '', 2, 3, 'product:add', 1, 0);
INSERT INTO `sys_menu` VALUES (365, '2026-02-14 23:32:38.820', '2026-02-14 23:32:38.820', '2026-02-17 17:44:58.542', 362, '编辑', '', '', '', 3, 3, 'product:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (366, '2026-02-14 23:32:38.921', '2026-02-14 23:32:38.921', '2026-02-17 17:44:58.542', 362, '删除', '', '', '', 4, 3, 'product:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (367, '2026-02-14 23:32:39.175', '2026-02-14 23:32:39.175', '2026-02-17 17:44:58.542', 362, '导出', '', '', '', 5, 3, 'product:export', 1, 0);
INSERT INTO `sys_menu` VALUES (368, '2026-02-14 23:32:39.419', '2026-02-14 23:32:39.419', '2026-02-17 17:44:58.542', 362, '导入', '', '', '', 6, 3, 'product:import', 1, 0);
INSERT INTO `sys_menu` VALUES (369, '2026-04-26 23:18:22.198', '2026-04-27 06:16:24.635', NULL, 3, '批量启用', '', '', '', 5, 3, 'system:user:batchEnable', 1, 0);
INSERT INTO `sys_menu` VALUES (370, '2026-04-26 23:18:22.220', '2026-04-27 06:16:24.647', NULL, 3, '批量禁用', '', '', '', 6, 3, 'system:user:batchDisable', 1, 0);
INSERT INTO `sys_menu` VALUES (371, '2026-04-27 05:10:18.823', '2026-04-27 20:17:27.131', NULL, 1, '部门管理', '/system/dept', 'system/dept/index', 'AppstoreTwoTone', 3, 2, 'system:dept:list', 1, 0);
INSERT INTO `sys_menu` VALUES (372, '2026-04-27 05:10:18.829', '2026-04-27 06:16:24.590', NULL, 371, '新增', '', '', '', 1, 3, 'system:dept:add', 1, 0);
INSERT INTO `sys_menu` VALUES (373, '2026-04-27 05:10:18.833', '2026-04-27 06:16:24.595', NULL, 371, '编辑', '', '', '', 2, 3, 'system:dept:edit', 1, 0);
INSERT INTO `sys_menu` VALUES (374, '2026-04-27 05:10:18.838', '2026-04-27 06:16:24.604', NULL, 371, '删除', '', '', '', 3, 3, 'system:dept:delete', 1, 0);
INSERT INTO `sys_menu` VALUES (375, '2026-04-28 03:24:48.528', '2026-04-28 04:21:46.825', NULL, 0, 'AI工具', '/ai-tools', 'Layout', 'ToolTwoTone', 3, 1, 'ai:tools', 1, 0);
INSERT INTO `sys_menu` VALUES (376, '2026-04-28 03:24:48.536', '2026-04-28 04:22:11.778', NULL, 375, 'AI对话', '/ai/chat', 'ai', 'custom-aiChat', 1, 2, 'ai:chat:list', 1, 0);
INSERT INTO `sys_menu` VALUES (377, '2026-04-28 03:24:48.544', '2026-04-28 04:21:59.652', NULL, 375, 'AI配置', '/ai/config', 'ai/config/index', 'SettingOutlined', 2, 2, 'ai:config:list', 1, 0);
INSERT INTO `sys_menu` VALUES (378, '2026-04-28 03:30:55.629', '2026-04-28 03:43:55.987', '2026-04-28 03:44:10.009', 0, '111', '/dd', 'system/dict/index', '', 0, 2, '', 1, 1);
INSERT INTO `sys_menu` VALUES (379, '2026-04-28 03:31:44.027', '2026-04-28 03:34:49.001', '2026-04-28 03:35:02.334', 378, 'test', '/test', 'system/role/index', '', 0, 2, '', 1, 0);
INSERT INTO `sys_menu` VALUES (380, '2026-04-28 19:58:27.000', '2026-04-28 19:58:27.000', NULL, 3, '批量重置密码', '', '', '', 7, 3, 'system:user:batchResetPwd', 1, 0);
INSERT INTO `sys_menu` VALUES (381, '2026-04-29 16:42:43.638', '2026-04-29 20:59:20.926', NULL, 0, '操作审计', '/system/operation-audit', 'Layout', 'custom-paper', 8, 1, '', 1, 0);
INSERT INTO `sys_menu` VALUES (382, '2026-05-02 01:27:44.994', '2026-05-02 01:27:45.040', NULL, 3, '强制下线', '', '', '', 0, 3, 'system:user:forceOffline', 1, 0);
INSERT INTO `sys_menu` VALUES (383, '2026-05-04 11:51:09.550', '2026-05-04 11:51:09.550', NULL, 3, '导入用户', '', '', '', 8, 3, 'system:user:import', 1, 0);
INSERT INTO `sys_menu` VALUES (384, '2026-05-04 11:51:09.572', '2026-05-04 11:51:09.572', NULL, 3, '导出用户', '', '', '', 9, 3, 'system:user:export', 1, 0);

-- ----------------------------
-- Table structure for sys_menu_api
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu_api`;
CREATE TABLE `sys_menu_api`  (
  `sys_menu_id` bigint UNSIGNED NOT NULL,
  `sys_api_id` bigint UNSIGNED NOT NULL,
  PRIMARY KEY (`sys_menu_id`, `sys_api_id`) USING BTREE,
  INDEX `fk_sys_menu_api_sys_api`(`sys_api_id` ASC) USING BTREE,
  CONSTRAINT `fk_sys_menu_api_sys_api` FOREIGN KEY (`sys_api_id`) REFERENCES `sys_api` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_sys_menu_api_sys_menu` FOREIGN KEY (`sys_menu_id`) REFERENCES `sys_menu` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_menu_api
-- ----------------------------
INSERT INTO `sys_menu_api` VALUES (382, 174);

-- ----------------------------
-- Table structure for sys_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_operation_log`;
CREATE TABLE `sys_operation_log`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'IP地址',
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '请求方法',
  `path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '请求路径',
  `request` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '请求参数',
  `response` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '响应结果',
  `status` bigint NULL DEFAULT NULL COMMENT 'HTTP状态码',
  `latency` bigint NULL DEFAULT NULL COMMENT '耗时(ms)',
  `user_agent` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'User-Agent',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `group` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '路由分组',
  `summary` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '路由描述',
  `business_code` bigint NULL DEFAULT NULL COMMENT '业务状态码',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 537 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_operation_log
-- ----------------------------

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '角色名称',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '角色编码',
  `sort` bigint NULL DEFAULT 0 COMMENT '排序',
  `status` bigint NULL DEFAULT 1 COMMENT '状态 1启用 0禁用',
  `data_scope` bigint NULL DEFAULT 1 COMMENT '数据范围 1全部 2自定义 3本部门 4本部门及下级 5仅本人',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
  `user_count` bigint NULL DEFAULT NULL,
  `is_super_admin` tinyint(1) NULL DEFAULT 0 COMMENT '是否显式超管',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sys_role_code`(`code` ASC) USING BTREE,
  INDEX `idx_sys_role_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO `sys_role` VALUES (1, '2026-01-24 02:39:47.708', '2026-05-04 11:51:09.837', NULL, '超级管理员', 'admin', 1, 1, 1, '拥有所有权限', NULL, 1);
INSERT INTO `sys_role` VALUES (2, '2026-01-24 02:57:36.645', '2026-05-04 11:51:09.829', NULL, '普通用户', 'user', 0, 1, 5, '普通注册用户', NULL, 0);
INSERT INTO `sys_role` VALUES (3, '2026-01-25 15:21:40.670', '2026-05-04 11:51:09.820', NULL, '系统管理员', 'system_admin', 0, 1, 1, '系统的最高权限用户，负责整个平台的运营和管理', NULL, 1);
INSERT INTO `sys_role` VALUES (6, '2026-04-29 02:37:34.212', '2026-05-04 11:51:09.806', NULL, 'test', 'test', 0, 1, 1, '', NULL, 0);
INSERT INTO `sys_role` VALUES (7, '2026-04-29 06:09:21.462', '2026-04-29 06:09:21.462', '2026-04-29 06:09:21.826', 'CodexDualTrack1777414160', 'codex_dual_track_1777414160', 999, 1, 1, 'codex dual track verification', NULL, 0);
INSERT INTO `sys_role` VALUES (8, '2026-04-29 06:09:59.805', '2026-04-29 06:10:00.865', '2026-04-29 06:10:01.295', 'CodexDualTrack1777414199', 'codex_dual_track_1777414199', 999, 1, 1, 'codex dual track verification', NULL, 0);
INSERT INTO `sys_role` VALUES (9, '2026-04-30 05:00:21.127', '2026-04-30 05:00:21.127', '2026-04-30 06:22:15.848', 'test1', 'test1', 0, 1, 1, '', NULL, 0);

-- ----------------------------
-- Table structure for sys_role_api
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_api`;
CREATE TABLE `sys_role_api`  (
  `sys_role_id` bigint UNSIGNED NOT NULL,
  `sys_api_id` bigint UNSIGNED NOT NULL,
  PRIMARY KEY (`sys_role_id`, `sys_api_id`) USING BTREE,
  INDEX `fk_sys_role_api_sys_api`(`sys_api_id` ASC) USING BTREE,
  CONSTRAINT `fk_sys_role_api_sys_api` FOREIGN KEY (`sys_api_id`) REFERENCES `sys_api` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_sys_role_api_sys_role` FOREIGN KEY (`sys_role_id`) REFERENCES `sys_role` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_role_api
-- ----------------------------
INSERT INTO `sys_role_api` VALUES (1, 1);
INSERT INTO `sys_role_api` VALUES (2, 1);
INSERT INTO `sys_role_api` VALUES (3, 1);
INSERT INTO `sys_role_api` VALUES (6, 1);
INSERT INTO `sys_role_api` VALUES (8, 1);
INSERT INTO `sys_role_api` VALUES (1, 2);
INSERT INTO `sys_role_api` VALUES (3, 2);
INSERT INTO `sys_role_api` VALUES (6, 2);
INSERT INTO `sys_role_api` VALUES (1, 3);
INSERT INTO `sys_role_api` VALUES (3, 3);
INSERT INTO `sys_role_api` VALUES (6, 3);
INSERT INTO `sys_role_api` VALUES (1, 4);
INSERT INTO `sys_role_api` VALUES (3, 4);
INSERT INTO `sys_role_api` VALUES (6, 4);
INSERT INTO `sys_role_api` VALUES (1, 5);
INSERT INTO `sys_role_api` VALUES (3, 5);
INSERT INTO `sys_role_api` VALUES (6, 5);
INSERT INTO `sys_role_api` VALUES (1, 6);
INSERT INTO `sys_role_api` VALUES (3, 6);
INSERT INTO `sys_role_api` VALUES (6, 6);
INSERT INTO `sys_role_api` VALUES (1, 7);
INSERT INTO `sys_role_api` VALUES (3, 7);
INSERT INTO `sys_role_api` VALUES (6, 7);
INSERT INTO `sys_role_api` VALUES (1, 8);
INSERT INTO `sys_role_api` VALUES (3, 8);
INSERT INTO `sys_role_api` VALUES (6, 8);
INSERT INTO `sys_role_api` VALUES (1, 9);
INSERT INTO `sys_role_api` VALUES (1, 10);
INSERT INTO `sys_role_api` VALUES (1, 11);
INSERT INTO `sys_role_api` VALUES (1, 12);
INSERT INTO `sys_role_api` VALUES (1, 13);
INSERT INTO `sys_role_api` VALUES (1, 14);
INSERT INTO `sys_role_api` VALUES (1, 15);
INSERT INTO `sys_role_api` VALUES (1, 16);
INSERT INTO `sys_role_api` VALUES (1, 17);
INSERT INTO `sys_role_api` VALUES (1, 18);
INSERT INTO `sys_role_api` VALUES (1, 19);
INSERT INTO `sys_role_api` VALUES (1, 20);
INSERT INTO `sys_role_api` VALUES (1, 21);
INSERT INTO `sys_role_api` VALUES (1, 22);
INSERT INTO `sys_role_api` VALUES (1, 23);
INSERT INTO `sys_role_api` VALUES (1, 24);
INSERT INTO `sys_role_api` VALUES (1, 25);
INSERT INTO `sys_role_api` VALUES (1, 26);
INSERT INTO `sys_role_api` VALUES (1, 27);
INSERT INTO `sys_role_api` VALUES (3, 27);
INSERT INTO `sys_role_api` VALUES (1, 28);
INSERT INTO `sys_role_api` VALUES (3, 28);
INSERT INTO `sys_role_api` VALUES (1, 29);
INSERT INTO `sys_role_api` VALUES (3, 29);
INSERT INTO `sys_role_api` VALUES (1, 30);
INSERT INTO `sys_role_api` VALUES (3, 30);
INSERT INTO `sys_role_api` VALUES (1, 31);
INSERT INTO `sys_role_api` VALUES (2, 31);
INSERT INTO `sys_role_api` VALUES (3, 31);
INSERT INTO `sys_role_api` VALUES (6, 31);
INSERT INTO `sys_role_api` VALUES (1, 32);
INSERT INTO `sys_role_api` VALUES (2, 32);
INSERT INTO `sys_role_api` VALUES (3, 32);
INSERT INTO `sys_role_api` VALUES (6, 32);
INSERT INTO `sys_role_api` VALUES (1, 33);
INSERT INTO `sys_role_api` VALUES (2, 33);
INSERT INTO `sys_role_api` VALUES (3, 33);
INSERT INTO `sys_role_api` VALUES (6, 33);
INSERT INTO `sys_role_api` VALUES (1, 37);
INSERT INTO `sys_role_api` VALUES (1, 38);
INSERT INTO `sys_role_api` VALUES (1, 39);
INSERT INTO `sys_role_api` VALUES (1, 40);
INSERT INTO `sys_role_api` VALUES (1, 42);
INSERT INTO `sys_role_api` VALUES (3, 42);
INSERT INTO `sys_role_api` VALUES (6, 42);
INSERT INTO `sys_role_api` VALUES (1, 43);
INSERT INTO `sys_role_api` VALUES (3, 43);
INSERT INTO `sys_role_api` VALUES (6, 43);
INSERT INTO `sys_role_api` VALUES (1, 44);
INSERT INTO `sys_role_api` VALUES (3, 44);
INSERT INTO `sys_role_api` VALUES (1, 45);
INSERT INTO `sys_role_api` VALUES (3, 45);
INSERT INTO `sys_role_api` VALUES (1, 46);
INSERT INTO `sys_role_api` VALUES (3, 46);
INSERT INTO `sys_role_api` VALUES (1, 47);
INSERT INTO `sys_role_api` VALUES (3, 47);
INSERT INTO `sys_role_api` VALUES (1, 49);
INSERT INTO `sys_role_api` VALUES (3, 49);
INSERT INTO `sys_role_api` VALUES (1, 62);
INSERT INTO `sys_role_api` VALUES (3, 62);
INSERT INTO `sys_role_api` VALUES (1, 63);
INSERT INTO `sys_role_api` VALUES (3, 63);
INSERT INTO `sys_role_api` VALUES (1, 64);
INSERT INTO `sys_role_api` VALUES (3, 64);
INSERT INTO `sys_role_api` VALUES (1, 65);
INSERT INTO `sys_role_api` VALUES (3, 65);
INSERT INTO `sys_role_api` VALUES (6, 65);
INSERT INTO `sys_role_api` VALUES (1, 66);
INSERT INTO `sys_role_api` VALUES (3, 66);
INSERT INTO `sys_role_api` VALUES (6, 66);
INSERT INTO `sys_role_api` VALUES (1, 67);
INSERT INTO `sys_role_api` VALUES (3, 67);
INSERT INTO `sys_role_api` VALUES (6, 67);
INSERT INTO `sys_role_api` VALUES (1, 68);
INSERT INTO `sys_role_api` VALUES (3, 68);
INSERT INTO `sys_role_api` VALUES (6, 68);
INSERT INTO `sys_role_api` VALUES (1, 69);
INSERT INTO `sys_role_api` VALUES (3, 69);
INSERT INTO `sys_role_api` VALUES (6, 69);
INSERT INTO `sys_role_api` VALUES (1, 70);
INSERT INTO `sys_role_api` VALUES (3, 70);
INSERT INTO `sys_role_api` VALUES (6, 70);
INSERT INTO `sys_role_api` VALUES (1, 71);
INSERT INTO `sys_role_api` VALUES (3, 71);
INSERT INTO `sys_role_api` VALUES (6, 71);
INSERT INTO `sys_role_api` VALUES (1, 72);
INSERT INTO `sys_role_api` VALUES (3, 72);
INSERT INTO `sys_role_api` VALUES (6, 72);
INSERT INTO `sys_role_api` VALUES (1, 73);
INSERT INTO `sys_role_api` VALUES (3, 73);
INSERT INTO `sys_role_api` VALUES (6, 73);
INSERT INTO `sys_role_api` VALUES (1, 81);
INSERT INTO `sys_role_api` VALUES (2, 81);
INSERT INTO `sys_role_api` VALUES (3, 81);
INSERT INTO `sys_role_api` VALUES (6, 81);
INSERT INTO `sys_role_api` VALUES (1, 97);
INSERT INTO `sys_role_api` VALUES (3, 97);
INSERT INTO `sys_role_api` VALUES (1, 98);
INSERT INTO `sys_role_api` VALUES (3, 98);
INSERT INTO `sys_role_api` VALUES (1, 99);
INSERT INTO `sys_role_api` VALUES (3, 99);
INSERT INTO `sys_role_api` VALUES (1, 100);
INSERT INTO `sys_role_api` VALUES (3, 100);
INSERT INTO `sys_role_api` VALUES (1, 101);
INSERT INTO `sys_role_api` VALUES (3, 101);
INSERT INTO `sys_role_api` VALUES (1, 102);
INSERT INTO `sys_role_api` VALUES (3, 102);
INSERT INTO `sys_role_api` VALUES (1, 103);
INSERT INTO `sys_role_api` VALUES (3, 103);
INSERT INTO `sys_role_api` VALUES (1, 104);
INSERT INTO `sys_role_api` VALUES (3, 104);
INSERT INTO `sys_role_api` VALUES (1, 105);
INSERT INTO `sys_role_api` VALUES (3, 105);
INSERT INTO `sys_role_api` VALUES (1, 106);
INSERT INTO `sys_role_api` VALUES (3, 106);
INSERT INTO `sys_role_api` VALUES (1, 112);
INSERT INTO `sys_role_api` VALUES (3, 112);
INSERT INTO `sys_role_api` VALUES (1, 113);
INSERT INTO `sys_role_api` VALUES (3, 113);
INSERT INTO `sys_role_api` VALUES (1, 114);
INSERT INTO `sys_role_api` VALUES (1, 115);
INSERT INTO `sys_role_api` VALUES (1, 117);
INSERT INTO `sys_role_api` VALUES (1, 118);
INSERT INTO `sys_role_api` VALUES (1, 119);
INSERT INTO `sys_role_api` VALUES (1, 120);
INSERT INTO `sys_role_api` VALUES (1, 121);
INSERT INTO `sys_role_api` VALUES (1, 122);
INSERT INTO `sys_role_api` VALUES (1, 123);
INSERT INTO `sys_role_api` VALUES (1, 129);
INSERT INTO `sys_role_api` VALUES (3, 129);
INSERT INTO `sys_role_api` VALUES (1, 130);
INSERT INTO `sys_role_api` VALUES (3, 130);
INSERT INTO `sys_role_api` VALUES (1, 131);
INSERT INTO `sys_role_api` VALUES (3, 131);
INSERT INTO `sys_role_api` VALUES (1, 132);
INSERT INTO `sys_role_api` VALUES (3, 132);
INSERT INTO `sys_role_api` VALUES (1, 159);
INSERT INTO `sys_role_api` VALUES (3, 159);
INSERT INTO `sys_role_api` VALUES (6, 159);
INSERT INTO `sys_role_api` VALUES (1, 160);
INSERT INTO `sys_role_api` VALUES (3, 160);
INSERT INTO `sys_role_api` VALUES (1, 172);
INSERT INTO `sys_role_api` VALUES (2, 172);
INSERT INTO `sys_role_api` VALUES (3, 172);
INSERT INTO `sys_role_api` VALUES (6, 172);
INSERT INTO `sys_role_api` VALUES (1, 173);
INSERT INTO `sys_role_api` VALUES (3, 173);
INSERT INTO `sys_role_api` VALUES (6, 173);
INSERT INTO `sys_role_api` VALUES (1, 174);
INSERT INTO `sys_role_api` VALUES (3, 174);
INSERT INTO `sys_role_api` VALUES (1, 175);
INSERT INTO `sys_role_api` VALUES (3, 175);
INSERT INTO `sys_role_api` VALUES (6, 175);
INSERT INTO `sys_role_api` VALUES (1, 204);
INSERT INTO `sys_role_api` VALUES (1, 205);
INSERT INTO `sys_role_api` VALUES (1, 206);
INSERT INTO `sys_role_api` VALUES (1, 207);
INSERT INTO `sys_role_api` VALUES (1, 208);
INSERT INTO `sys_role_api` VALUES (1, 209);
INSERT INTO `sys_role_api` VALUES (1, 210);
INSERT INTO `sys_role_api` VALUES (1, 211);
INSERT INTO `sys_role_api` VALUES (1, 212);
INSERT INTO `sys_role_api` VALUES (1, 213);
INSERT INTO `sys_role_api` VALUES (1, 214);
INSERT INTO `sys_role_api` VALUES (1, 215);
INSERT INTO `sys_role_api` VALUES (1, 216);
INSERT INTO `sys_role_api` VALUES (1, 217);
INSERT INTO `sys_role_api` VALUES (1, 218);
INSERT INTO `sys_role_api` VALUES (1, 219);
INSERT INTO `sys_role_api` VALUES (6, 219);
INSERT INTO `sys_role_api` VALUES (1, 220);
INSERT INTO `sys_role_api` VALUES (3, 220);
INSERT INTO `sys_role_api` VALUES (6, 220);
INSERT INTO `sys_role_api` VALUES (1, 221);
INSERT INTO `sys_role_api` VALUES (3, 221);
INSERT INTO `sys_role_api` VALUES (6, 221);
INSERT INTO `sys_role_api` VALUES (1, 222);
INSERT INTO `sys_role_api` VALUES (3, 222);
INSERT INTO `sys_role_api` VALUES (6, 222);
INSERT INTO `sys_role_api` VALUES (1, 223);
INSERT INTO `sys_role_api` VALUES (3, 223);
INSERT INTO `sys_role_api` VALUES (6, 223);
INSERT INTO `sys_role_api` VALUES (1, 224);
INSERT INTO `sys_role_api` VALUES (3, 224);
INSERT INTO `sys_role_api` VALUES (6, 224);
INSERT INTO `sys_role_api` VALUES (1, 225);
INSERT INTO `sys_role_api` VALUES (3, 225);
INSERT INTO `sys_role_api` VALUES (6, 225);
INSERT INTO `sys_role_api` VALUES (1, 226);
INSERT INTO `sys_role_api` VALUES (3, 226);
INSERT INTO `sys_role_api` VALUES (6, 226);
INSERT INTO `sys_role_api` VALUES (1, 227);
INSERT INTO `sys_role_api` VALUES (3, 227);
INSERT INTO `sys_role_api` VALUES (1, 228);
INSERT INTO `sys_role_api` VALUES (3, 228);
INSERT INTO `sys_role_api` VALUES (6, 228);
INSERT INTO `sys_role_api` VALUES (1, 230);
INSERT INTO `sys_role_api` VALUES (3, 230);
INSERT INTO `sys_role_api` VALUES (1, 231);
INSERT INTO `sys_role_api` VALUES (3, 231);
INSERT INTO `sys_role_api` VALUES (1, 237);
INSERT INTO `sys_role_api` VALUES (1, 238);
INSERT INTO `sys_role_api` VALUES (1, 239);
INSERT INTO `sys_role_api` VALUES (1, 240);
INSERT INTO `sys_role_api` VALUES (2, 240);
INSERT INTO `sys_role_api` VALUES (3, 240);
INSERT INTO `sys_role_api` VALUES (6, 240);
INSERT INTO `sys_role_api` VALUES (1, 241);
INSERT INTO `sys_role_api` VALUES (3, 241);
INSERT INTO `sys_role_api` VALUES (6, 241);
INSERT INTO `sys_role_api` VALUES (1, 242);
INSERT INTO `sys_role_api` VALUES (2, 242);
INSERT INTO `sys_role_api` VALUES (3, 242);
INSERT INTO `sys_role_api` VALUES (6, 242);

-- ----------------------------
-- Table structure for sys_role_data_scope
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_data_scope`;
CREATE TABLE `sys_role_data_scope`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `role_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '角色ID',
  `resource_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '业务功能资源码',
  `data_scope` bigint NULL DEFAULT 1 COMMENT '数据范围 1全部 2自定义 3本部门 4本部门及下级 5仅本人',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_role_resource`(`role_id` ASC, `resource_code` ASC) USING BTREE,
  INDEX `idx_sys_role_data_scope_deleted_at`(`deleted_at` ASC) USING BTREE,
  CONSTRAINT `fk_sys_role_feature_data_scopes` FOREIGN KEY (`role_id`) REFERENCES `sys_role` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 54 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role_data_scope
-- ----------------------------
INSERT INTO `sys_role_data_scope` VALUES (52, '2026-05-02 02:08:00.309', '2026-05-02 02:08:00.312', NULL, 6, 'system:user-management', 4);
INSERT INTO `sys_role_data_scope` VALUES (53, '2026-05-02 02:08:00.315', '2026-05-02 02:08:00.327', NULL, 6, 'system:dept-management', 4);

-- ----------------------------
-- Table structure for sys_role_data_scope_dept
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_data_scope_dept`;
CREATE TABLE `sys_role_data_scope_dept`  (
  `sys_role_data_scope_id` bigint UNSIGNED NOT NULL,
  `sys_dept_id` bigint UNSIGNED NOT NULL,
  PRIMARY KEY (`sys_role_data_scope_id`, `sys_dept_id`) USING BTREE,
  INDEX `fk_sys_role_data_scope_dept_sys_dept`(`sys_dept_id` ASC) USING BTREE,
  CONSTRAINT `fk_sys_role_data_scope_dept_sys_dept` FOREIGN KEY (`sys_dept_id`) REFERENCES `sys_dept` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_sys_role_data_scope_dept_sys_role_data_scope` FOREIGN KEY (`sys_role_data_scope_id`) REFERENCES `sys_role_data_scope` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role_data_scope_dept
-- ----------------------------

-- ----------------------------
-- Table structure for sys_role_dept
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_dept`;
CREATE TABLE `sys_role_dept`  (
  `sys_role_id` bigint UNSIGNED NOT NULL,
  `sys_dept_id` bigint UNSIGNED NOT NULL,
  PRIMARY KEY (`sys_role_id`, `sys_dept_id`) USING BTREE,
  INDEX `fk_sys_role_dept_sys_dept`(`sys_dept_id` ASC) USING BTREE,
  CONSTRAINT `fk_sys_role_dept_sys_dept` FOREIGN KEY (`sys_dept_id`) REFERENCES `sys_dept` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_sys_role_dept_sys_role` FOREIGN KEY (`sys_role_id`) REFERENCES `sys_role` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role_dept
-- ----------------------------

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu`  (
  `sys_role_id` bigint UNSIGNED NOT NULL,
  `sys_menu_id` bigint UNSIGNED NOT NULL,
  PRIMARY KEY (`sys_role_id`, `sys_menu_id`) USING BTREE,
  INDEX `fk_sys_role_menu_sys_menu`(`sys_menu_id` ASC) USING BTREE,
  CONSTRAINT `fk_sys_role_menu_sys_menu` FOREIGN KEY (`sys_menu_id`) REFERENCES `sys_menu` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_sys_role_menu_sys_role` FOREIGN KEY (`sys_role_id`) REFERENCES `sys_role` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_role_menu
-- ----------------------------
INSERT INTO `sys_role_menu` VALUES (1, 1);
INSERT INTO `sys_role_menu` VALUES (6, 1);
INSERT INTO `sys_role_menu` VALUES (8, 1);
INSERT INTO `sys_role_menu` VALUES (3, 3);
INSERT INTO `sys_role_menu` VALUES (6, 3);
INSERT INTO `sys_role_menu` VALUES (8, 3);
INSERT INTO `sys_role_menu` VALUES (1, 4);
INSERT INTO `sys_role_menu` VALUES (1, 5);
INSERT INTO `sys_role_menu` VALUES (1, 6);
INSERT INTO `sys_role_menu` VALUES (1, 7);
INSERT INTO `sys_role_menu` VALUES (3, 7);
INSERT INTO `sys_role_menu` VALUES (1, 8);
INSERT INTO `sys_role_menu` VALUES (3, 8);
INSERT INTO `sys_role_menu` VALUES (1, 11);
INSERT INTO `sys_role_menu` VALUES (3, 11);
INSERT INTO `sys_role_menu` VALUES (3, 16);
INSERT INTO `sys_role_menu` VALUES (6, 16);
INSERT INTO `sys_role_menu` VALUES (1, 17);
INSERT INTO `sys_role_menu` VALUES (1, 18);
INSERT INTO `sys_role_menu` VALUES (1, 19);
INSERT INTO `sys_role_menu` VALUES (1, 20);
INSERT INTO `sys_role_menu` VALUES (3, 21);
INSERT INTO `sys_role_menu` VALUES (6, 21);
INSERT INTO `sys_role_menu` VALUES (3, 22);
INSERT INTO `sys_role_menu` VALUES (6, 22);
INSERT INTO `sys_role_menu` VALUES (3, 23);
INSERT INTO `sys_role_menu` VALUES (8, 23);
INSERT INTO `sys_role_menu` VALUES (1, 24);
INSERT INTO `sys_role_menu` VALUES (3, 24);
INSERT INTO `sys_role_menu` VALUES (1, 211);
INSERT INTO `sys_role_menu` VALUES (3, 369);
INSERT INTO `sys_role_menu` VALUES (6, 369);
INSERT INTO `sys_role_menu` VALUES (3, 370);
INSERT INTO `sys_role_menu` VALUES (6, 370);
INSERT INTO `sys_role_menu` VALUES (1, 371);
INSERT INTO `sys_role_menu` VALUES (3, 371);
INSERT INTO `sys_role_menu` VALUES (6, 371);
INSERT INTO `sys_role_menu` VALUES (1, 372);
INSERT INTO `sys_role_menu` VALUES (3, 372);
INSERT INTO `sys_role_menu` VALUES (1, 373);
INSERT INTO `sys_role_menu` VALUES (3, 373);
INSERT INTO `sys_role_menu` VALUES (1, 374);
INSERT INTO `sys_role_menu` VALUES (3, 374);
INSERT INTO `sys_role_menu` VALUES (1, 375);
INSERT INTO `sys_role_menu` VALUES (3, 375);
INSERT INTO `sys_role_menu` VALUES (1, 376);
INSERT INTO `sys_role_menu` VALUES (3, 376);
INSERT INTO `sys_role_menu` VALUES (1, 377);
INSERT INTO `sys_role_menu` VALUES (3, 377);
INSERT INTO `sys_role_menu` VALUES (3, 380);
INSERT INTO `sys_role_menu` VALUES (6, 380);
INSERT INTO `sys_role_menu` VALUES (1, 381);
INSERT INTO `sys_role_menu` VALUES (3, 381);
INSERT INTO `sys_role_menu` VALUES (3, 383);
INSERT INTO `sys_role_menu` VALUES (6, 383);
INSERT INTO `sys_role_menu` VALUES (8, 383);
INSERT INTO `sys_role_menu` VALUES (3, 384);
INSERT INTO `sys_role_menu` VALUES (6, 384);
INSERT INTO `sys_role_menu` VALUES (8, 384);

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名',
  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '密码',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '昵称',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '手机号',
  `status` bigint NULL DEFAULT 1 COMMENT '状态 1启用 0禁用',
  `dept_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT '部门ID',
  `avatar_file_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '头像文件ID',
  `gender` bigint NULL DEFAULT 0 COMMENT '性别 0未知 1男 2女',
  `created_by` bigint UNSIGNED NULL DEFAULT 0 COMMENT '创建人ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sys_user_username`(`username` ASC) USING BTREE,
  INDEX `idx_sys_user_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `fk_sys_user_avatar_file`(`avatar_file_id` ASC) USING BTREE,
  INDEX `idx_sys_user_dept_id`(`dept_id` ASC) USING BTREE,
  INDEX `idx_sys_user_created_by`(`created_by` ASC) USING BTREE,
  CONSTRAINT `fk_sys_user_avatar_file` FOREIGN KEY (`avatar_file_id`) REFERENCES `sys_file` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_sys_user_dept` FOREIGN KEY (`dept_id`) REFERENCES `sys_dept` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 41 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_user
-- ----------------------------
INSERT INTO `sys_user` VALUES (1, '2026-01-24 02:39:47.881', '2026-04-30 05:35:41.260', NULL, 'admin', '$2a$10$nW3w6LIPh.29XBIczfzrVOTjFJMBqbxty48vdanFbzmyZxB6yk0mW', '管理员', '1440350254@qq.com', '17688744031', 1, 7, 93, 0, 0);
INSERT INTO `sys_user` VALUES (18, '2026-01-25 15:26:36.448', '2026-04-28 16:58:42.222', NULL, 'A1001', '$2a$10$Hg7bxA3CSbxLve3wCeEfuu3HJmtaO3LytEXkeKo81qc0MFkK0hfMu', '龙舞', 'eee33@qq.com', '1768874223', 1, 6, 93, 1, 0);
INSERT INTO `sys_user` VALUES (24, '2026-01-31 22:04:53.330', '2026-04-28 16:58:38.110', NULL, 'lisan', '$2a$10$Gijr8jyphPbEKMpNmjsRje1Ku.1qXm7deGz0VqqNH2CwU83RwVHv2', '李三', 'dsdsd444@qilincsp.cn', '', 1, 6, 93, 0, 0);
INSERT INTO `sys_user` VALUES (30, '2026-02-05 03:35:05.611', '2026-04-28 20:11:00.938', NULL, 'csp', '$2a$10$Lst0s1FcBF/Fy86v.tbGXu5u17UTz7NJOOYIhOUOhFI9Eu3LX7hfq', '', '', '', 1, 3, 93, 1, 0);
INSERT INTO `sys_user` VALUES (33, '2026-04-27 05:55:08.272', '2026-05-02 04:04:07.375', '2026-04-30 05:46:18.866', 'test_deleted_33_1777499178', '$2a$10$Lst0s1FcBF/Fy86v.tbGXu5u17UTz7NJOOYIhOUOhFI9Eu3LX7hfq', '', '', '', 1, 8, NULL, 0, 0);
INSERT INTO `sys_user` VALUES (35, '2026-04-29 06:09:59.899', '2026-04-29 06:10:01.272', '2026-04-29 06:10:01.275', 'codex_dual_1777414199_deleted_35_1777414201', '$2a$10$ewxm7xXmb8IAeD7aSyL8PeFidYsYo8JIHQ50acpfo3SsPYKpQgJBq', 'codex dual track', 'codex_dual_1777414199@example.com', '', 1, 5, 93, 0, 0);
INSERT INTO `sys_user` VALUES (36, '2026-04-30 04:42:35.436', '2026-05-02 04:04:04.767', '2026-04-30 05:10:13.517', 'test1111_deleted_36_1777497013', '$2a$10$dAz0NIGkXz0tv7Aa6ui1beIKnNf0ZOaDFNrvgjW/0O1N.A2U1z4Zq', '', '', '', 1, 5, NULL, 0, 33);
INSERT INTO `sys_user` VALUES (37, '2026-04-30 06:00:13.647', '2026-04-30 06:14:56.901', NULL, 'test', '$2a$10$/Bhgu47AwsGKeM0R5k.MIOhNCSeIJo4FqImdC7pq6At3z5wZp2z46', '11', 'test@qq.com', '17688744031', 1, 8, 114, 0, 1);
INSERT INTO `sys_user` VALUES (38, '2026-04-30 06:01:40.456', '2026-04-30 06:02:10.174', NULL, '111', '$2a$10$9eDqAjhMyQJ8Wlxu5Fl0KuRy30HP7/EkigxkZlHNKFYLki3MuDuOy', '', '11', '11', 1, 8, 115, 0, 37);
INSERT INTO `sys_user` VALUES (39, '2026-05-05 01:07:18.601', '2026-05-05 02:08:10.113', '2026-05-05 02:08:10.115', 'zhangsan_deleted_39_1777918090', '$2a$10$UXtHWZLAMMq1ST.D2saYTOw7.pkhlyER5L6QsPogzGIMzHr3X6xb6', '', 'zhangsan@example.com', '13800138000', 1, 3, NULL, 0, 1);
INSERT INTO `sys_user` VALUES (40, '2026-05-05 02:08:19.222', '2026-05-05 02:08:19.226', NULL, 'zhangsan', '$2a$10$37lfTaPkT3JXBjoOOMwLLeS3qVVLuhERE4aHnZ/Xc3DTpE0uTXdmO', '', 'zhangsan@example.com', '13800138000', 1, 3, 93, 0, 1);

-- ----------------------------
-- Table structure for sys_user_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_user_role`;
CREATE TABLE `sys_user_role`  (
  `sys_user_id` bigint UNSIGNED NOT NULL,
  `sys_role_id` bigint UNSIGNED NOT NULL,
  PRIMARY KEY (`sys_user_id`, `sys_role_id`) USING BTREE,
  INDEX `fk_sys_user_role_sys_role`(`sys_role_id` ASC) USING BTREE,
  CONSTRAINT `fk_sys_user_role_sys_role` FOREIGN KEY (`sys_role_id`) REFERENCES `sys_role` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_sys_user_role_sys_user` FOREIGN KEY (`sys_user_id`) REFERENCES `sys_user` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_user_role
-- ----------------------------
INSERT INTO `sys_user_role` VALUES (1, 1);
INSERT INTO `sys_user_role` VALUES (24, 2);
INSERT INTO `sys_user_role` VALUES (30, 2);
INSERT INTO `sys_user_role` VALUES (38, 2);
INSERT INTO `sys_user_role` VALUES (40, 2);
INSERT INTO `sys_user_role` VALUES (18, 3);
INSERT INTO `sys_user_role` VALUES (37, 6);

SET FOREIGN_KEY_CHECKS = 1;
