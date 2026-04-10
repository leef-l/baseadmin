-- ============================================================
-- codegen 离线验证用表（不需要实际建表，仅用于 verify_codegen.go 模拟）
-- 覆盖当前协作约定中的 codegen 验收场景
-- ============================================================

-- 1. demo_category: 树形表 + 简单字段（验证树形、排序、Tooltip）
-- 2. demo_article: 复杂表 + 外键 + 所有组件类型（验证外键、枚举、金额、搜索等）
-- 3. demo_tag: 最简表（验证无外键、无树形、无特殊字段的基础 CRUD）
-- 4. demo_user_review: 多段模块名（user_review）+ 跨应用外键（验证 moduleName 保留下划线）

-- ========== 树形分类表 ==========
CREATE TABLE IF NOT EXISTS `demo_category` (
  `id` BIGINT UNSIGNED NOT NULL COMMENT 'ID',
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父分类',
  `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
  `icon` VARCHAR(100) DEFAULT '' COMMENT '图标',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序（升序）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0=禁用,1=启用',
  `created_by` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `dept_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME,
  `updated_at` DATETIME,
  `deleted_at` DATETIME,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分类';

-- ========== 复杂文章表（覆盖所有字段类型）==========
CREATE TABLE IF NOT EXISTS `demo_article` (
  `id` BIGINT UNSIGNED NOT NULL COMMENT 'ID',
  -- 外键：指向同应用的树形表（TreeSelect）
  `category_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '分类',
  -- 外键：指向其他应用的普通表（Select，跨应用，显式声明 ref）
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '作者|ref:system_users.username',
  -- 基础文本
  `title` VARCHAR(200) NOT NULL COMMENT '文章标题',
  `order_no` VARCHAR(50) NOT NULL COMMENT '文章编号',
  -- 图片上传
  `cover` VARCHAR(500) DEFAULT '' COMMENT '封面',
  -- 文件上传
  `attachment_file` VARCHAR(500) DEFAULT '' COMMENT '附件',
  -- 富文本
  `body_content` TEXT COMMENT '正文内容',
  -- JSON 编辑器
  `extra_json` TEXT COMMENT 'JSON扩展',
  -- URL
  `link_url` VARCHAR(500) DEFAULT '' COMMENT '外部链接',
  -- 枚举：3 值 → Radio
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0=草稿,1=已发布,2=已下架',
  -- 枚举：多值 → Select
  `type` TINYINT NOT NULL DEFAULT 1 COMMENT '类型:1=普通,2=置顶,3=推荐,4=热门',
  -- 枚举：2 值 → Switch
  `is_top` TINYINT NOT NULL DEFAULT 0 COMMENT '是否置顶:0=否,1=是',
  -- 金额（分）
  `price` INT NOT NULL DEFAULT 0 COMMENT '价格（分）',
  -- 密码
  `pay_password` VARCHAR(100) DEFAULT '' COMMENT '支付密码',
  -- 排序
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序（升序）',
  -- 图标
  `icon` VARCHAR(100) DEFAULT '' COMMENT '图标',
  -- 邮箱（验证规则）
  `email` VARCHAR(100) DEFAULT '' COMMENT '联系邮箱',
  -- 手机号（验证规则）
  `phone` VARCHAR(20) DEFAULT '' COMMENT '联系电话',
  -- 备注（Textarea + 搜索）
  `remark` TEXT COMMENT '备注',
  -- 字典字段
  `level` VARCHAR(20) DEFAULT '' COMMENT '等级:dict:article_level',
  -- 无 comment 的字段（测试回退）
  `extra_field` VARCHAR(100) DEFAULT '',
  -- 自定义时间字段
  `publish_at` DATETIME COMMENT '发布时间',
  `expire_at` DATETIME COMMENT '过期时间',
  -- 公共字段
  `created_by` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `dept_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME,
  `updated_at` DATETIME,
  `deleted_at` DATETIME,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章';

-- ========== 最简标签表 ==========
CREATE TABLE IF NOT EXISTS `demo_tag` (
  `id` BIGINT UNSIGNED NOT NULL COMMENT 'ID',
  `name` VARCHAR(50) NOT NULL COMMENT '标签名称',
  `color` VARCHAR(20) DEFAULT '' COMMENT '颜色',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0=禁用,1=启用',
  `created_by` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `dept_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME,
  `updated_at` DATETIME,
  `deleted_at` DATETIME,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标签';

-- ========== 多段模块名表（user_review → 包名 user_review）==========
-- 唯一核心验证点：表名里第一个下划线之后保留下划线，moduleName=user_review，
-- ModelName=UserReview，包路径是 app/demo/internal/logic/user_review/。
CREATE TABLE IF NOT EXISTS `demo_user_review` (
  `id` BIGINT UNSIGNED NOT NULL COMMENT 'ID',
  -- 跨应用外键：验证多段表名下的 ref 解析
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '审核人|ref:system_users.username',
  -- 枚举 3 值 → Radio
  `review_type` TINYINT NOT NULL DEFAULT 1 COMMENT '类型:1=内容,2=行为,3=申诉',
  -- 搜索字段
  `content` VARCHAR(500) DEFAULT '' COMMENT '审核内容',
  -- 数值
  `score` INT NOT NULL DEFAULT 0 COMMENT '评分',
  -- 枚举 2 值 → Switch
  `is_passed` TINYINT NOT NULL DEFAULT 0 COMMENT '是否通过:0=否,1=是',
  -- 排序 + Tooltip
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序（升序）',
  -- 枚举 Switch
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0=禁用,1=启用',
  -- 公共字段
  `created_by` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `dept_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME,
  `updated_at` DATETIME,
  `deleted_at` DATETIME,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户审核';

-- ========== 被引用的跨应用表（system_users 已存在）==========
-- demo_article.user_id → system_users（跨应用外键）
-- demo_article.category_id → demo_category（同应用树形外键）
-- demo_user_review.user_id → system_users（多段模块名 + 跨应用外键）
