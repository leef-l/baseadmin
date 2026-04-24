SET @column_exists := (
  SELECT COUNT(*)
  FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'upload_dir_rule'
    AND COLUMN_NAME = 'storage_types'
);

SET @repair_sql := IF(
  @column_exists = 0,
  'ALTER TABLE `upload_dir_rule` ADD COLUMN `storage_types` varchar(20) DEFAULT ''1,2,3'' COMMENT ''适用存储类型，多个用逗号分隔'' AFTER `file_type`',
  'SELECT 1'
);

PREPARE repair_stmt FROM @repair_sql;
EXECUTE repair_stmt;
DEALLOCATE PREPARE repair_stmt;

UPDATE `upload_dir_rule`
SET `storage_types` = '1,2,3'
WHERE `storage_types` IS NULL OR `storage_types` = '';
