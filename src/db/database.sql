

## 规则表
CREATE TABLE `rule` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '规则主键 id，自增',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '规则名称',
  `default_value` text COMMENT '命中规则时返回值',
  `flag` tinyint(11) NOT NULL DEFAULT '0' COMMENT '是否有效: 0有效 1无效',
  `comment` varchar(500) DEFAULT NULL COMMENT '备注',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `last_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_rule` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

## 决策项
CREATE TABLE `rule_decision_tree` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 id，自增',
  `rule_id` int(11) DEFAULT NULL COMMENT '关联的规则 id',
  `tree` text COMMENT '决策树(由规则项id组成的逻辑表达式)',
  `priority` int(11) NOT NULL DEFAULT '1' COMMENT '优先级(数字越小 优先级越高)',
  `value` text COMMENT '决策值',
  `flag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否有效: 0有效 1无效',
  `create_time_utc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_modified_utc` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


## 规则项表
CREATE TABLE `rule_item` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 id，自增',
  `operator` varchar(20) NOT NULL COMMENT '规则运算符',
  `factor_type` varchar(50) DEFAULT NULL COMMENT '规则因子类型',
  `factor_val` text COMMENT '规则因子值',
  `flag` tinyint(11) NOT NULL DEFAULT '0' COMMENT '是否有效: 0有效 1无效',
  `comment` varchar(500) DEFAULT NULL COMMENT '规则项备注',
  `create_time_utc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_modified_utc` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


