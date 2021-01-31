DROP DATABASE IF EXISTS sandbox_db;
CREATE DATABASE sandbox_db;
USE sandbox_db;

DROP TABLE IF EXISTS `issues`;

CREATE TABLE `issues` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` text NOT NULL,
  `description` text NOT NULL,
  PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `issues` (`id`, `title`, `description`) VALUES
    (1, 'title1', 'description1'),
    (2, 'title2', 'description2');

DROP TABLE IF EXISTS `labels`;

CREATE TABLE `labels` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` text NOT NULL,
  PRIMARY KEY(`id`)
);

INSERT INTO `labels` (`id`, `name`) VALUES
    (1, 'name1'),
    (2, 'name2');

DROP TABLE IF EXISTS `issue_labels`;

CREATE TABLE `issue_labels` (
  `issue_id` int(11),
  `label_id` int(11),
  PRIMARY KEY(`issue_id`, `label_id`),
  FOREIGN KEY (`issue_id`) REFERENCES `issues` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`label_id`) REFERENCES `labels` (`id`) ON DELETE CASCADE
);

INSERT INTO `issue_labels` (`issue_id`, `label_id`) VALUES
    (1, 1),
    (1, 2);
