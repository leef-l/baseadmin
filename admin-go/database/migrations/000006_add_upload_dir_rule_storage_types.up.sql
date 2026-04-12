ALTER TABLE `upload_dir_rule`
  ADD COLUMN `storage_types` varchar(20) DEFAULT '1,2,3' COMMENT '适用存储类型，多个用逗号分隔' AFTER `file_type`;

UPDATE `upload_dir_rule`
SET `storage_types` = '1,2,3'
WHERE `storage_types` IS NULL OR `storage_types` = '';
