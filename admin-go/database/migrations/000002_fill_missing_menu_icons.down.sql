UPDATE `system_menu`
SET `icon` = '', `updated_at` = NOW()
WHERE `type` = 2
  AND `path` = '/upload/config'
  AND `deleted_at` IS NULL
  AND `icon` = 'ControlOutlined';

UPDATE `system_menu`
SET `icon` = '', `updated_at` = NOW()
WHERE `type` = 2
  AND `path` = '/upload/dir'
  AND `deleted_at` IS NULL
  AND `icon` = 'FolderOpenOutlined';

UPDATE `system_menu`
SET `icon` = '', `updated_at` = NOW()
WHERE `type` = 2
  AND `path` = '/upload/dir-rule'
  AND `deleted_at` IS NULL
  AND `icon` = 'PartitionOutlined';

UPDATE `system_menu`
SET `icon` = '', `updated_at` = NOW()
WHERE `type` = 2
  AND `path` = '/upload/file'
  AND `deleted_at` IS NULL
  AND `icon` = 'FileTextOutlined';
