UPDATE `system_menu`
SET `icon` = 'ControlOutlined', `updated_at` = NOW()
WHERE `type` = 2
  AND `path` = '/upload/config'
  AND `deleted_at` IS NULL
  AND (`icon` IS NULL OR `icon` = '');

UPDATE `system_menu`
SET `icon` = 'FolderOpenOutlined', `updated_at` = NOW()
WHERE `type` = 2
  AND `path` = '/upload/dir'
  AND `deleted_at` IS NULL
  AND (`icon` IS NULL OR `icon` = '');

UPDATE `system_menu`
SET `icon` = 'PartitionOutlined', `updated_at` = NOW()
WHERE `type` = 2
  AND `path` = '/upload/dir-rule'
  AND `deleted_at` IS NULL
  AND (`icon` IS NULL OR `icon` = '');

UPDATE `system_menu`
SET `icon` = 'FileTextOutlined', `updated_at` = NOW()
WHERE `type` = 2
  AND `path` = '/upload/file'
  AND `deleted_at` IS NULL
  AND (`icon` IS NULL OR `icon` = '');
