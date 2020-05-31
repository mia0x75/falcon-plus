/**
 * utf8mb4_unicode_ci is based on the official Unicode rules
 *                    for universal sorting and comparison,
 *                    which sorts accurately in a wide range
 *                    of languages.
 *
 * utf8mb4_general_ci is a simplified set of sorting rules which
 *                    aims to do as well as it can while taking
 *                    many short-cuts designed to improve speed.
 *                    It does not follow the Unicode rules and
 *                    will result in undesirable sorting or comparison
 *                    in some situations, such as when using
 *                    particular languages or characters.
 */

CREATE DATABASE `dashboard`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;
USE `dashboard`;
SET NAMES utf8mb4;


DROP TABLE IF EXISTS `teams`;
CREATE TABLE `teams`
(
   `id`        INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '自增主键',
   `name`      VARCHAR(50)
               NOT NULL
               COMMENT '分组名称',
   `resume`    VARCHAR(200)
               NOT NULL
               COMMENT ''
               DEFAULT '',
   `creator`   INT UNSIGNED
               NOT NULL
               COMMENT '创建者标识'
               DEFAULT '0',
   `create_at` INT UNSIGNED
               NOT NULL
               COMMENT '创建时间'
               DEFAULT UNIX_TIMESTAMP(),
   PRIMARY KEY (`id`),
   UNIQUE KEY `unique_1` (`name`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '分组';


/**
 * role: -1:blocked 0:normal 1:admin 2:root
 */
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
   `id`        INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '自增主键',
   `name`      VARCHAR(50)
               NOT NULL
               COMMENT '登录名称',
   `passwd`    CHAR(64)
               NOT NULL
               COMMENT '哈希后的密码'
               DEFAULT '',
   `cnname`    VARCHAR(15)
               NOT NULL
               COMMENT '中文名称'
               DEFAULT '',
   `email`     VARCHAR(75)
               NOT NULL
               COMMENT '邮箱'
               DEFAULT '',
   `phone`     VARCHAR(16)
               NOT NULL
               COMMENT '电话'
               DEFAULT '',
   `im`        VARCHAR(50)
               NOT NULL
               COMMENT '即时通讯账户'
               DEFAULT '',
   `role`      TINYINT
               NOT NULL
               COMMENT '角色'
               DEFAULT 0,
   `creator`   INT UNSIGNED
               NOT NULL
               COMMENT '创建者标识'
               DEFAULT 0,
   `create_at` INT UNSIGNED
               NOT NULL
               COMMENT '创建时间'
               DEFAULT UNIX_TIMESTAMP(),
   PRIMARY KEY (`id`),
   UNIQUE KEY `unique_1` (`name`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '用户';


DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions`
(
   `id`        INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '自增主键',
   `user_id`   INT UNSIGNED
               NOT NULL
               COMMENT '用户标识',
   `sign`      CHAR(32)
               NOT NULL
               COMMENT '会话标识',
   `expire`    INT UNSIGNED
               NOT NULL
               COMMENT '过期时间',
   `create_at` INT UNSIGNED
               NOT NULL
               COMMENT '创建时间'
               DEFAULT UNIX_TIMESTAMP(),
   PRIMARY KEY (`id`),
   KEY `unique_1` (`user_id`),
   KEY `index_1` (`sign`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '会话';


/**
 * 这里的机器是从机器管理系统中同步过来的
 * 系统拿出来单独部署需要为hbs增加功能，心跳上来的机器写入host表
 */
DROP TABLE IF EXISTS `hosts`;
CREATE TABLE `hosts`
(
   `id`             INT UNSIGNED
                    NOT NULL
                    AUTO_INCREMENT
                    COMMENT '',
   `hostname`       VARCHAR(200)
                    NOT NULL
                    COMMENT ''
                    DEFAULT '',
   `ip`             VARCHAR(15)
                    NOT NULL
                    COMMENT ''
                    DEFAULT '',
   `agent_version`  VARCHAR(20)
                    NOT NULL
                    COMMENT ''
                    DEFAULT '',
   `plugin_version` VARCHAR(20)
                    NOT NULL
                    COMMENT ''
                    DEFAULT '',
   `maintain_begin` INT UNSIGNED
                    NOT NULL
                    COMMENT ''
                    DEFAULT 0,
   `maintain_end`   INT UNSIGNED
                    NOT NULL
                    COMMENT ''
                    DEFAULT 0,
   `update_at`      INT UNSIGNED
                    COMMENT '',
   `create_at`      INT UNSIGNED
                    NOT NULL
                    COMMENT ''
                    DEFAULT UNIX_TIMESTAMP(),
   PRIMARY KEY (`id`),
   UNIQUE KEY `unique_1` (`hostname`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


/**
 * 机器分组信息
 * come_from 0: 从机器管理同步过来的；1: 从页面创建的
 */
DROP TABLE IF EXISTS `groups`;
CREATE TABLE `groups`
(
   `id`          INT UNSIGNED
                 NOT NULL
                 AUTO_INCREMENT
                 COMMENT '',
   `name`        VARCHAR(75)
                 NOT NULL
                 COMMENT ''
                 DEFAULT '',
   `creator`     INT UNSIGNED
                 NOT NULL
                 COMMENT '',
   `create_at`   INT UNSIGNED
                 NOT NULL
                 COMMENT ''
                 DEFAULT UNIX_TIMESTAMP(),
   `come_from`   TINYINT
                 NOT NULL
                 COMMENT ''
                 DEFAULT '0',
   PRIMARY KEY (`id`),
   UNIQUE KEY `unique_1` (`name`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


/**
 * 监控策略模板
 * name全局唯一，命名的时候可以适当带上一些前缀，比如：sa.falcon.base
 */
DROP TABLE IF EXISTS `templates`;
CREATE TABLE `templates`
(
   `id`        INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '自增主键',
   `name`      VARCHAR(75)
               NOT NULL
               COMMENT ''
               DEFAULT '',
   `parent_id` INT UNSIGNED
               NOT NULL
               COMMENT ''
               DEFAULT 0,
   `action_id` INT UNSIGNED
               NOT NULL
               COMMENT ''
               DEFAULT 0,
   `creator`   INT UNSIGNED
               NOT NULL
               COMMENT '',
   `create_at` INT UNSIGNED
               NOT NULL
               COMMENT ''
               DEFAULT UNIX_TIMESTAMP(),
   PRIMARY KEY (`id`),
   UNIQUE KEY `unique_1` (`name`),
   KEY `index_1` (`creator`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


DROP TABLE IF EXISTS `strategies`;
CREATE TABLE `strategies`
(
   `id`          INT UNSIGNED
                 NOT NULL
                 AUTO_INCREMENT
                 COMMENT '',
   `metric`      VARCHAR(100)
                 NOT NULL
                 COMMENT ''
                 DEFAULT '',
   `tags`        VARCHAR(200)
                 NOT NULL
                 COMMENT ''
                 DEFAULT '',
   `max_step`    INT
                 NOT NULL
                 COMMENT ''
                 DEFAULT '1',
   `priority`    TINYINT
                 NOT NULL
                 COMMENT ''
                 DEFAULT '0',
   `func`        VARCHAR(16)
                 NOT NULL
                 COMMENT ''
                 DEFAULT 'all(#1)',
   `op`          VARCHAR(8)
                 NOT NULL
                 COMMENT ''
                 DEFAULT '',
   `right_value` VARCHAR(50)
                 NOT NULL
                 COMMENT '',
   `note`        VARCHAR(200)
                 NOT NULL
                 COMMENT ''
                 DEFAULT '',
   `run_begin`   VARCHAR(16)
                 NOT NULL
                 COMMENT ''
                 DEFAULT '',
   `run_end`     VARCHAR(16)
                 NOT NULL
                 COMMENT ''
                 DEFAULT '',
   `template_id` INT UNSIGNED
                 NOT NULL
                 COMMENT ''
                 DEFAULT '0',
   PRIMARY KEY (`id`),
   KEY `index_1` (`template_id`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


DROP TABLE IF EXISTS `expressions`;
CREATE TABLE `expressions`
(
   `id`          INT UNSIGNED
                 NOT NULL
                 AUTO_INCREMENT
                 COMMENT '',
   `expression`  VARCHAR(500)
                 NOT NULL
                 COMMENT '',
   `func`        VARCHAR(50)
                 NOT NULL
                 COMMENT ''
                 DEFAULT 'all(#1)',
   `op`          VARCHAR(8)
                 NOT NULL
                 COMMENT ''
                 DEFAULT '',
   `right_value` VARCHAR(20)
                 NOT NULL
                 COMMENT ''
                 DEFAULT '',
   `max_step`    INT
                 NOT NULL
                 COMMENT ''
                 DEFAULT '1',
   `priority`    TINYINT
                 NOT NULL
                 COMMENT ''
                 DEFAULT '0',
   `note`        VARCHAR(200)
                 NOT NULL
                 COMMENT ''
                 DEFAULT '',
   `action_id`   INT UNSIGNED
                 NOT NULL
                 COMMENT ''
                 DEFAULT '0',
   `creator`     INT UNSIGNED
                 NOT NULL
                 COMMENT '',
   `pause`       TINYINT
                 NOT NULL
                 COMMENT ''
                 DEFAULT '0',
   PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


DROP TABLE IF EXISTS `plugin_dir`;
CREATE TABLE `plugin_dir`
(
   `id`        INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '自增主键',
   `group_id`  INT UNSIGNED
               NOT NULL
               COMMENT '',
   `dir`       VARCHAR(255)
               NOT NULL
               COMMENT '',
   `creator`   INT UNSIGNED
               NOT NULL
               COMMENT '',
   `create_at` INT UNSIGNED
               NOT NULL
               COMMENT '创建时间'
               DEFAULT UNIX_TIMESTAMP(),
   PRIMARY KEY (`id`),
   KEY `index_1` (`group_id`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


DROP TABLE IF EXISTS `actions`;
CREATE TABLE `actions`
(
   `id`                   INT UNSIGNED
                          NOT NULL
                          AUTO_INCREMENT
                          COMMENT '',
   `uic`                  VARCHAR(255)
                          NOT NULL
                          COMMENT ''
                          DEFAULT '',
   `url`                  VARCHAR(255)
                          NOT NULL
                          COMMENT ''
                          DEFAULT '',
   `callback`             TINYINT
                          NOT NULL
                          COMMENT ''
                          DEFAULT '0',
   `before_callback_sms`  TINYINT
                          NOT NULL
                          COMMENT ''
                          DEFAULT '0',
   `before_callback_mail` TINYINT
                          NOT NULL
                          COMMENT ''
                          DEFAULT '0',
   `after_callback_sms`   TINYINT
                          NOT NULL
                          COMMENT ''
                          DEFAULT '0',
   `after_callback_mail`  TINYINT
                          NOT NULL
                          COMMENT ''
                          DEFAULT '0',
   PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


/**
 * nodata mock config
 */
DROP TABLE IF EXISTS `mockcfg`;
CREATE TABLE `mockcfg`
(
   `id`        BIGINT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '',
   `name`      VARCHAR(200)
               NOT NULL
               DEFAULT ''
               COMMENT 'name of mockcfg, used for uuid',
   `obj`       TEXT
               NOT NULL
               DEFAULT ''
               COMMENT 'desc of object',
   `obj_type`  VARCHAR(200)
               NOT NULL
               DEFAULT ''
               COMMENT 'type of object, host or group or other',
   `metric`    VARCHAR(128)
               NOT NULL
               COMMENT ''
               DEFAULT '',
   `tags`      VARCHAR(500)
               NOT NULL
               COMMENT ''
               DEFAULT '',
   `dstype`    VARCHAR(32)
               NOT NULL
               COMMENT ''
               DEFAULT 'GAUGE',
   `step`      INT UNSIGNED
               NOT NULL
               COMMENT ''
               DEFAULT 60,
   `mock`      DOUBLE
               NOT NULL
               DEFAULT 0
               COMMENT 'mocked value when nodata occurs',
   `creator`   INT UNSIGNED
               NOT NULL
               COMMENT '',
   `create_at` INT UNSIGNED
               NOT NULL
               DEFAULT UNIX_TIMESTAMP()
               COMMENT '创建时间',
   `update_at` INT UNSIGNED
               COMMENT '修改时间',
   PRIMARY KEY (`id`),
   UNIQUE KEY `unique_1` (`name`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


/**
 *  aggregator cluster metric config table
 */
DROP TABLE IF EXISTS `clusters`;
CREATE TABLE `clusters`
(
   `id`          INT UNSIGNED
                 NOT NULL
                 AUTO_INCREMENT
                 COMMENT '',
   `group_id`    INT
                 NOT NULL
                 COMMENT '',
   `numerator`   TEXT # 原来是 VARCHAR(10240) 这是一个不好的设计
                 NOT NULL
                 COMMENT '',
   `denominator` TEXT # 原来是 VARCHAR(10240) 这是一个不好的设计
                 NOT NULL
                 COMMENT '',
   `endpoint`    VARCHAR(200)
                 NOT NULL
                 COMMENT '',
   `metric`      VARCHAR(200)
                 NOT NULL
                 COMMENT '',
   `tags`        VARCHAR(200)
                 NOT NULL
                 COMMENT '',
   `ds_type`     VARCHAR(200)
                 NOT NULL
                 COMMENT '',
   `step`        INT
                 NOT NULL
                 COMMENT '',
   `update_at`   INT UNSIGNED
                 COMMENT '',
   `creator`     INT UNSIGNED
                 NOT NULL
                 COMMENT '',
   `create_at`   INT UNSIGNED
                 NOT NULL
                 COMMENT ''
                 DEFAULT UNIX_TIMESTAMP(),
   PRIMARY KEY (`id`)
)
ENGINE = InnoDB
DEFAULT CHARSET= utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


/**
 * alert links
 */
DROP TABLE IF EXISTS `alert_link`;
CREATE TABLE `alert_link`
(
   `id`        INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '',
   `path`      VARCHAR(16)
               NOT NULL
               COMMENT ''
               DEFAULT '',
   `content`   TEXT
               NOT NULL
               COMMENT '',
   `create_at` INT UNSIGNED
               NOT NULL
               COMMENT ''
               DEFAULT UNIX_TIMESTAMP(),
   PRIMARY KEY (id),
   UNIQUE KEY `unique_1` (path)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


DROP TABLE IF EXISTS `graphs`;
CREATE TABLE `graphs`
(
   `id`        INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '',
   `title`     VARCHAR(100) # 数据类型是否要调整
               NOT NULL
               COMMENT '',
   `hosts`     TEXT # 原来使用 VARCHAR(10240)
               NOT NULL
               COMMENT ''
               DEFAULT '',
   `counters`  TEXT # 原来使用 VARCHAR(10240)
               NOT NULL
               COMMENT ''
               DEFAULT '',
   `screen_id` INT UNSIGNED
               NOT NULL
               COMMENT '',
   `timespan`  INT UNSIGNED
               NOT NULL
               COMMENT ''
               DEFAULT '3600',
   `type`      CHAR(2)
               NOT NULL
               COMMENT ''
               DEFAULT 'h',
   `method`    CHAR(8) # 数据类型是否要调整
               COMMENT ''
               DEFAULT '',
   `position`  INT UNSIGNED
               NOT NULL
               COMMENT ''
               DEFAULT '0',
   `tags`      VARCHAR(500)
               NOT NULL
               COMMENT ''
               DEFAULT '',
   `create_at` INT UNSIGNED
               NOT NULL
               COMMENT ''
               DEFAULT UNIX_TIMESTAMP(),
   PRIMARY KEY (`id`),
   KEY `index_1` (`screen_id`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


DROP TABLE IF EXISTS `screens`;
CREATE TABLE `screens`
(
   `id`        INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '',
   `pid`       INT UNSIGNED
               NOT NULL
               COMMENT ''
               DEFAULT '0',
   `name`      CHAR(200) # 数据类型是否要调整
               NOT NULL
               COMMENT '',
   `create_at` INT UNSIGNED
               NOT NULL
               COMMENT ''
               DEFAULT UNIX_TIMESTAMP(),
   `update_at` INT UNSIGNED
               COMMENT '',
   PRIMARY KEY (`id`),
   KEY `index_1` (`pid`),
   UNIQUE KEY `unique_1` (`pid`,`name`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


DROP TABLE IF EXISTS `drafts`;
CREATE TABLE `drafts`
(
   `id`        INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '',
   `endpoints` TEXT # 原来使用 VARCHAR(10240)
               NOT NULL
               COMMENT ''
               DEFAULT '',
   `counters`  TEXT # 原来使用 VARCHAR(10240)
               NOT NULL
               COMMENT ''
               DEFAULT '',
   `sign`      CHAR(32)
               NOT NULL
               COMMENT '',
   `create_at` INT UNSIGNED
               NOT NULL
               COMMENT ''
               DEFAULT UNIX_TIMESTAMP(),
   PRIMARY KEY (`id`),
   UNIQUE KEY `unique_1` (`sign`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


DROP TABLE IF EXISTS `endpoints`;
CREATE TABLE `endpoints`
(
   `id`        INT UNSIGNED
               NOT NULL
               AUTO_INCREMENT
               COMMENT '',
   `endpoint`  VARCHAR(200)
               NOT NULL
               COMMENT ''
               DEFAULT '',
   `ts`        INT
               COMMENT '',
   `create_at` INT UNSIGNED
               NOT NULL
               DEFAULT UNIX_TIMESTAMP()
               COMMENT '创建时间',
   `update_at` INT UNSIGNED
               COMMENT '修改时间',
   PRIMARY KEY (`id`),
   UNIQUE KEY `unique_1` (`endpoint`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


DROP TABLE IF EXISTS `counters`;
CREATE TABLE `counters`
(
   `id`          INT UNSIGNED
                 NOT NULL
                 AUTO_INCREMENT
                 COMMENT '',
   `endpoint_id` INT UNSIGNED
                 NOT NULL
                 COMMENT '',
   `counter`     VARCHAR(200)
                 NOT NULL
                 COMMENT ''
                 DEFAULT '',
   `step`        INT
                 NOT NULL
                 COMMENT '单位秒'
                 DEFAULT 60,
   `type`        VARCHAR(16)
                 NOT NULL
                 COMMENT 'GAUGE|COUNTER|DERIVE',
   `ts`          INT
                 COMMENT '',
   `create_at`   INT UNSIGNED
                 NOT NULL
                 DEFAULT UNIX_TIMESTAMP()
                 COMMENT '创建时间',
   `update_at`   INT UNSIGNED
                 COMMENT '修改时间',
   PRIMARY KEY (`id`),
   UNIQUE KEY `unique_1` (`endpoint_id`, `counter`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


DROP TABLE IF EXISTS `tags`;
CREATE TABLE `tags`
(
   `id`          INT UNSIGNED
                 NOT NULL
                 AUTO_INCREMENT
                 COMMENT '',
   `tag`         VARCHAR(200)
                 NOT NULL
                 COMMENT 'srv=tv'
                 DEFAULT '',
   `endpoint_id` INT UNSIGNED
                 NOT NULL
                 COMMENT '',
   `ts`          INT
                 COMMENT '',
   `create_at`   INT UNSIGNED
                 NOT NULL
                 DEFAULT UNIX_TIMESTAMP()
                 COMMENT '创建时间',
   `update_at`   INT UNSIGNED
                 COMMENT '修改时间',
   PRIMARY KEY (`id`),
   UNIQUE KEY `unique_1` (`tag`, `endpoint_id`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '';


/*
* 建立告警归档资料表, 主要存储各个告警的最后触发状况
*/
DROP TABLE IF EXISTS `cases`;
CREATE TABLE `cases`
(
   `id`             VARCHAR(50)
                    NOT NULL
                    COMMENT '字符主键',
   `endpoint`       VARCHAR(100)
                    NOT NULL
                    COMMENT '监控主机',
   `metric`         VARCHAR(200)
                    NOT NULL
                    COMMENT '监控指标',
   `func`           VARCHAR(50)
                    COMMENT '计算函数',
   `cond`           VARCHAR(200)
                    NOT NULL
                    COMMENT '告警条件',
   `note`           VARCHAR(200)
                    COMMENT '备注',
   `max_step`       INT UNSIGNED
                    COMMENT '告警次数阈值',
   `current_step`   INT UNSIGNED
                    COMMENT '已触发告警次数',
   `priority`       INT
                    NOT NULL
                    COMMENT '告警级别',
   `status`         VARCHAR(20)
                    NOT NULL
                    COMMENT '告警状态',
   `create_at`      INT UNSIGNED
                    NOT NULL
                    DEFAULT UNIX_TIMESTAMP()
                    COMMENT '创建时间',
   `update_at`      INT UNSIGNED
                    COMMENT '更新时间',
   `closed_at`      INT UNSIGNED
                    COMMENT '关闭时间',
   `closed_note`    VARCHAR(200)
                    COMMENT '关闭说明',
   `user_modified`  INT UNSIGNED
                    COMMENT '操作用户',
   `expression_id`  INT UNSIGNED
                    COMMENT '表达式 portal.expression.id',
   `strategy_id`    INT UNSIGNED
                    COMMENT '策略 portal.strategy.id',
   `template_id`    INT UNSIGNED
                    COMMENT '模版 portal.template.id',
   `process_note`   MEDIUMINT
                    COMMENT '处理意见',
   `process_status` VARCHAR(20)
                    COMMENT '处理状态'
                    DEFAULT 'unresolved',
   PRIMARY KEY (`id`),
   INDEX `index_1` (`endpoint`, `strategy_id`, `template_id`)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '告警事件';


/*
* 建立告警归档资料表, 存储各个告警触发状况的历史状态
*/
DROP TABLE IF EXISTS `events`;
CREATE TABLE `events`
(
   `id`           INT UNSIGNED
                  NOT NULL
                  AUTO_INCREMENT
                  COMMENT '自增主键',
   `case_id`      VARCHAR(50)
                  COMMENT '事件源头 cases.id',
   `step`         INT UNSIGNED
                  COMMENT '告警次数',
   `cond`         VARCHAR(200)
                  NOT NULL
                  COMMENT '告警条件',
   `status`       INT UNSIGNED
                  COMMENT ''
                  DEFAULT 0,
   `create_at`    INT UNSIGNED
                  NOT NULL
                  DEFAULT UNIX_TIMESTAMP()
                  COMMENT '创建时间',
   PRIMARY KEY (`id`),
   INDEX `index_1` (`case_id`),
   FOREIGN KEY (`case_id`) REFERENCES `cases`(`id`)
      ON DELETE CASCADE
      ON UPDATE CASCADE
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '告警事件归档';


/*
* 告警留言表
*/
DROP TABLE IF EXISTS `notes`;
CREATE TABLE `notes`
(
   `id`           INT UNSIGNED
                  NOT NULL
                  AUTO_INCREMENT
                  COMMENT '自增主键',
   `event_caseId` VARCHAR(50)
                  COMMENT '',
   `note`         VARCHAR(200)
                  COMMENT '',
   `case_id`      VARCHAR(20)
                  COMMENT '',
   `status`       VARCHAR(15)
                  COMMENT '',
   `create_at`    INT UNSIGNED
                  DEFAULT UNIX_TIMESTAMP()
                  COMMENT '',
   `creator`      INT UNSIGNED
                  COMMENT '',
   PRIMARY KEY (`id`),
   INDEX `index_1` (`event_caseId`),
   FOREIGN KEY (`event_caseId`) REFERENCES `cases`(`id`)
      ON DELETE CASCADE
      ON UPDATE CASCADE,
   FOREIGN KEY (`creator`) REFERENCES `users`(`id`)
      ON DELETE CASCADE
      ON UPDATE CASCADE
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '告警处置信息';


DROP TABLE IF EXISTS `edges`;
CREATE TABLE `edges`
(
   `id`             INT UNSIGNED
                    NOT NULL
                    AUTO_INCREMENT
                    COMMENT '自增主键',
   `type`           TINYINT UNSIGNED
                    NOT NULL
                    COMMENT '类别',
   `ancestor_id`    INT UNSIGNED
                    NOT NULL
                    COMMENT '先代',
   `descendant_id`  INT UNSIGNED
                    NOT NULL
                    COMMENT '后代',
   `create_at`      INT UNSIGNED
                    DEFAULT UNIX_TIMESTAMP()
                    COMMENT '创建事件',
   `creator`        INT UNSIGNED
                    NOT NULL
                    COMMENT '创建者',
   PRIMARY KEY (`id`),
   UNIQUE `unique_1` (`ancestor_id`, `descendant_id`, `type`),
   INDEX `index_1` (`ancestor_id`, `type`),
   INDEX `index_2` (`descendant_id`, `type`),
   FOREIGN KEY (`creator`) REFERENCES `users`(`id`)
      ON DELETE CASCADE
      ON UPDATE CASCADE
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = '关联关系';


/*初始数据 - 密码: 123456*/
INSERT INTO `users` (`name`, `passwd`, `role`) VALUES ('root', '$2a$10$lhDzHQ2MnXi66OR44MUWZerrP3hzUK1TUe.j6NX0YpbvUkii35bxG', 2);

