/*
 Navicat Premium Data Transfer

 Source Server         : wechat-bot
 Source Server Type    : MySQL
 Source Server Version : 80024
 Source Host           : 175.178.74.2:3306
 Source Schema         : wechat-bot

 Target Server Type    : MySQL
 Target Server Version : 80024
 File Encoding         : 65001

 Date: 06/09/2023 11:45:28
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for w_user
-- ----------------------------
DROP TABLE IF EXISTS `w_user`;
CREATE TABLE `w_user` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'id',
  `wechat_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '微信ID',
  `nick_name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '昵称',
  `balance` int NOT NULL DEFAULT '0' COMMENT '余额',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `is_admin` char(1) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '2' COMMENT '是否管理员 1是 2否',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_wechat_id` (`wechat_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;
