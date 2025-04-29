CREATE TABLE `students` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `info` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `status` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `age` int(10) unsigned not null default 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of students
-- ----------------------------
INSERT INTO `students` VALUES ('1', 'tom', 'a top student', '2022-06-02 15:28:55', '2022-06-02 15:27:01', '1', 0);
INSERT INTO `students` VALUES ('3', 'jimmy', 'a good student', null, null, '0', 0);
INSERT INTO `students` VALUES ('4', 'you', 'fea tea', null, null, '1', 0);
INSERT INTO `students` VALUES ('6', 'ju', '', null, null, '1', 0);
