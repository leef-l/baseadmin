ALTER TABLE `upload_dir_rule`
  MODIFY COLUMN `file_type` varchar(255) DEFAULT '' COMMENT '文件类型，多个用逗号分隔',
  DROP COLUMN `keep_name`;
