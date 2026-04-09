DELETE FROM `system_role_menu`
WHERE `role_id` = 1000000000000000002
  AND `menu_id` IN (
    1000000000000000024,
    1000000000000000036,
    1000000000000000044,
    1000000000000000054,
    314253730302136320,
    314253730436354048,
    314253730520240128,
    314253730637680640
  );

DELETE FROM `system_menu`
WHERE `id` IN (
    1000000000000000024,
    1000000000000000036,
    1000000000000000044,
    1000000000000000054,
    314253730302136320,
    314253730436354048,
    314253730520240128,
    314253730637680640
  )
  AND `permission` IN (
    'system:dept:batch-delete',
    'system:role:batch-delete',
    'system:menu:batch-delete',
    'system:user:batch-delete',
    'upload:config:batch-delete',
    'upload:dir:batch-delete',
    'upload:dir_rule:batch-delete',
    'upload:file:batch-delete'
  );
