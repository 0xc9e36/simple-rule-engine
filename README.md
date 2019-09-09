## 规则引擎


适用于大量 if else 的场景，将业务规则与具体应用代码分离, 抽象成一个简单的规则引擎系统，可以极大降低编码的复杂度.


### 1. 实现原理

![](http://assets.processon.com/chart_image/5d74e7d1e4b017f7e03263f6.png)

首先将业务规则录入数据库中，采用唯一规则名标识，然后输入要执行的规则名， 和一组规则因子(相当于原来 if 条件变量),  规则引擎会返回两个值，其中第一个是规则是否匹配(true 或者 false)，另一个是匹配结果；
如果规则匹配结果是 false，返回的是默认值;
如果规则匹配结果是 true，返回的值是命中 tree 对应的 value;


规则表，用于保存所有的业务规则，一条记录代表一条规则, default_value 代表系统没有匹配到 name 规则时返回值.
```
CREATE TABLE `rule` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '规则主键 id，自增',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '规则名称',
  `default_value` text COMMENT '未命中规则时返回值',
  `flag` tinyint(11) NOT NULL DEFAULT '0' COMMENT '是否有效: 0有效 1无效',
  `comment` varchar(500) DEFAULT NULL COMMENT '备注',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `last_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_rule` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```


决策表，和具体某项规则进行关联，如果有高优先级的 tree 命中，直接返回.
```
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
```


规则项表, 用来保存每一个规则因子比较条件(比较规则引擎的输入参数与数据库中规则因子的值), operator 需要自己实现.
```
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
```

### 2.如何使用

举个例子，我们有这么一个需求， 给定一个分数等级区间如下:    
A:   分数 >= 90    
B:   分数 >= 80    
C:   分数 >= 70    
D:   分数 >= 60    
E:   分数 60 以下

然后给出一个具体的分数，要求输出该分数在哪个区间.


首先我们在规则表中新加入一条记录, 规则名为 score_level, 输出值为等级区间， 如果规则没有匹配到，默认输出 E.
```
INSERT INTO `rule` (`id`, `name`, `default_value`, `flag`, `comment`, `create_time`, `last_modified`) VALUES(1, 'score_level', 'E', 0, '分数等级规则，默认为 E，60 分以下', '2019-09-08 18:30:07', '2019-09-08 18:51:52');
```


然后在规则项表里加入如下记录:
```
INSERT INTO `rule_item` (`id`, `operator`, `factor_type`, `factor_val`, `flag`, `comment`, `create_time_utc`, `last_modified_utc`)
VALUES
(1, '>=', 'score', '90', 0, '分数 >= 90规则项', '2019-09-08 18:30:14', '2019-09-08 18:46:45'),
(2, '>=', 'score', '80', 0, '分数 >= 80规则项', '2019-09-08 18:32:50', '2019-09-08 18:33:12'),
(3, '>=', 'score', '70', 0, '分数 >= 70规则项', '2019-09-08 18:33:26', '2019-09-08 18:33:31'),
(4, '>=', 'score', '60', 0, '分数 >= 60规则项', '2019-09-08 18:51:09', '2019-09-08 18:51:24');
```


最后决策表里加入记录:
```
INSERT INTO `rule_decision_tree` (`id`, `rule_id`, `tree`, `priority`, `value`, `flag`, `create_time_utc`, `last_modified_utc`)
VALUES
(1, 1, '1', 1, 'A', 0, '2019-09-08 18:39:06', '2019-09-08 18:39:06'),
(2, 1, '2', 2, 'B', 0, '2019-09-08 18:39:25', '2019-09-08 18:39:25'),
(3, 1, '3', 3, 'C', 0, '2019-09-08 18:40:10', '2019-09-08 18:40:10'),
(4, 1, '4', 4, 'D', 0, '2019-09-08 18:51:42', '2019-09-08 18:51:42');

```

代码实现:
```
factors := []rule.Factor{
    {
        FactorType: "score",
        FactorValue: "59",
    },

    //这里还可以有其他规则因子...
}

val, matched := rule.MatchRule("score_level", factors...)
```

val 代表规则匹配返回值，如果没有则返回默认值， matched 代表是否命中规则决策树.
