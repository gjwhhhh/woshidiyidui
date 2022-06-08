-- 删除无效点赞 让点赞表主键重新有序递增
START TRANSACTION;
DELETE FROM `dy_favorite` WHERE `dy_favorite`.is_deleted = 1;
ALTER TABLE `dy_favorite` DROP `id`;
ALTER TABLE `dy_favorite` ADD `id` int NOT NULL FIRST;
ALTER TABLE `dy_favorite` MODIFY COLUMN `id` int NOT NULL AUTO_INCREMENT,ADD PRIMARY KEY(id);
COMMIT;

-- 删除无效评论，让评论表主键重新有序递增
START TRANSACTION;
DELETE FROM `dy_comment` WHERE `dy_comment`.is_deleted = 1;
ALTER TABLE `dy_comment` DROP `id`;
ALTER TABLE `dy_comment` ADD `id` int NOT NULL FIRST;
ALTER TABLE `dy_comment` MODIFY COLUMN `id` int NOT NULL AUTO_INCREMENT,ADD PRIMARY KEY(id);
COMMIT;