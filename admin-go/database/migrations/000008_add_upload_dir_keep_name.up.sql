ALTER TABLE `upload_dir`
  ADD COLUMN `keep_name` tinyint(1) NOT NULL DEFAULT '0' COMMENT '保留原文件名:0=否,1=是' AFTER `path`;
