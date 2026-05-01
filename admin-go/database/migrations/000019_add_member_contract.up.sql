-- ================================================================
-- 000019: 电子签合同（模板 + 签署记录）
-- 业务背景：
--   注册成功后引导会员签《会员协议》一份；不签不阻断业务，但前端显示"未签"提示。
--   未来可扩展为升级/下单等场景再签。
--   签名图片用前端 canvas 输出 base64 PNG 上传，PDF 异步生成（gopdf）。
-- ================================================================

SET NAMES utf8mb4;

-- 合同模板（管理员后台维护，可有多版本）
CREATE TABLE IF NOT EXISTS `member_contract_template` (
  `id` bigint unsigned NOT NULL COMMENT '模板ID（Snowflake）',
  `template_name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模板名称|search:like|keyword:on|priority:100',
  `template_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'register' COMMENT '模板类型|search:select|enum:register=注册协议,upgrade=升级协议,custom=自定义',
  `content`       mediumtext  COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模板正文（HTML，支持{{nickname}}{{phone}}{{date}}等占位符）|search:off',
  `is_default`    tinyint NOT NULL DEFAULT '0' COMMENT '是否默认模板:0=否,1=是|search:select',
  `remark`        varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注|search:off',
  `sort`          int NOT NULL DEFAULT '0' COMMENT '排序（升序）',
  `status`        tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启|search:select',
  `tenant_id`     bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id`   bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by`    bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id`       bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at`    datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at`    datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at`    datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_template_type` (`template_type`),
  KEY `idx_is_default` (`is_default`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员合同模板';

-- 合同签署记录
CREATE TABLE IF NOT EXISTS `member_contract` (
  `id`               bigint unsigned NOT NULL COMMENT '合同ID（Snowflake）',
  `user_id`          bigint unsigned NOT NULL DEFAULT '0' COMMENT '会员|ref:member_user.nickname|search:select',
  `contract_no`      varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '合同编号|search:eq|keyword:on|priority:100',
  `contract_type`    varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'register' COMMENT '合同类型|search:select|enum:register=注册协议,upgrade=升级协议,custom=自定义',
  `template_id`      bigint unsigned NOT NULL DEFAULT '0' COMMENT '模板|ref:member_contract_template.template_name',
  `related_id`       bigint unsigned NOT NULL DEFAULT '0' COMMENT '关联业务ID（订单/升级记录等）',
  `signed_html`      mediumtext  COLLATE utf8mb4_unicode_ci COMMENT '签署时实际渲染的 HTML（已替换占位符，含签名图）|search:off',
  `signature_image`  mediumtext  COLLATE utf8mb4_unicode_ci COMMENT '手写签名 base64 PNG（data:image/png;base64,...）|search:off',
  `signed_at`        datetime DEFAULT NULL COMMENT '签署时间|search:date',
  `signed_ip`        varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '签署IP',
  `signed_user_agent` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT 'UA',
  `pdf_path`         varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT 'PDF存储路径（OSS或本地）',
  `pdf_status`       tinyint NOT NULL DEFAULT '0' COMMENT 'PDF生成状态:0=未生成,1=生成中,2=已生成,3=失败|search:select',
  `pdf_error`        varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT 'PDF生成错误信息',
  `remark`           varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注|search:off',
  `sort`             int NOT NULL DEFAULT '0' COMMENT '排序',
  `status`           tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=作废,1=正常|search:select',
  `tenant_id`        bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id`      bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by`       bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id`          bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at`       datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at`       datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at`       datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_contract_no` (`contract_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_contract_type` (`contract_type`),
  KEY `idx_pdf_status` (`pdf_status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员合同记录';

-- 默认注册协议模板
INSERT INTO `member_contract_template`
  (`id`, `template_name`, `template_type`, `content`, `is_default`, `status`, `created_at`, `updated_at`)
VALUES
  (
    1,
    '会员注册协议',
    'register',
    CONCAT(
      '<h1 style="text-align:center">会员注册协议</h1>',
      '<p>甲方：基金平台运营方</p>',
      '<p>乙方：{{nickname}}（手机号：{{phone}}）</p>',
      '<h2>一、协议概述</h2>',
      '<p>乙方自愿注册成为本平台会员，并承诺遵守平台规则。</p>',
      '<h2>二、会员权利</h2>',
      '<p>1. 享受平台进货、寄售、奖励金、推广奖等功能；</p>',
      '<p>2. 享受平台等级体系带来的限购、奖励等差异化服务。</p>',
      '<h2>三、会员义务</h2>',
      '<p>1. 不得违规进行交易；</p>',
      '<p>2. 不得使用本平台从事违法行为；</p>',
      '<p>3. 妥善保管账号密码，账号下行为视为本人行为。</p>',
      '<h2>四、其他</h2>',
      '<p>本协议自乙方电子签名时生效。</p>',
      '<p style="margin-top:40px">签署日期：{{date}}</p>'
    ),
    1, 1, NOW(), NOW()
  )
ON DUPLICATE KEY UPDATE `updated_at` = NOW();
