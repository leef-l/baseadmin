SET NAMES utf8mb4;

ALTER TABLE `member_user`
  DROP COLUMN `total_purchase_count`,
  DROP COLUMN `last_purchase_date`,
  DROP COLUMN `today_purchase_count`,
  DROP COLUMN `daily_purchase_limit`;

ALTER TABLE `member_level`
  DROP COLUMN `daily_purchase_limit`;

DROP TABLE IF EXISTS `member_business_config`;
