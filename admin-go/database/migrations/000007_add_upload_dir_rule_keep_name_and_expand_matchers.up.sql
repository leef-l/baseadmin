ALTER TABLE `upload_dir_rule`
  MODIFY COLUMN `file_type` text COMMENT '匹配条件，多个可换行',
  ADD COLUMN `keep_name` tinyint(1) NOT NULL DEFAULT '0' COMMENT '保留原文件名:0=否,1=是' AFTER `save_path`;
