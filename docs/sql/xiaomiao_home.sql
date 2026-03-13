# CREATE DATABASE xiaomiao_home CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE xiaomiao_home;

-- 用户表
CREATE TABLE
    IF NOT EXISTS t_user
(
    id           bigint PRIMARY KEY COMMENT '用户id',
    username     varchar(64) NOT NULL COMMENT '用户名称',
    nickname     varchar(64) NOT NULL COMMENT '显示名称',
    password     varchar(255) COMMENT '密码',
    salt         varchar(32) COMMENT '密码加盐字段',
    telephone    varchar(32) COMMENT '手机号',
    email        varchar(30) COMMENT '邮件',
    status       tinyint(1)  default 0 COMMENT '账户状态, 0-禁用, 1-启用',
    signature    varchar(64) COMMENT '个人签名',
    avatar       varchar(255) COMMENT '头像地址',
    position     varchar(32) COMMENT '职位',
    bio          varchar(255) COMMENT '个人简介',
    mfa_status   tinyint(1)  default 1 COMMENT 'mfa状态, 0-禁用, 1-启用',
    mfa_secret   varchar(255) COMMENT 'mfa加密密钥',
    remark  longtext COMMENT '描述',
    deleted_flag tinyint(1)                   DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_user bigint NOT NULL COMMENT '创建用户',
    updated_user bigint NOT NULL COMMENT '更新用户',
    deleted_user bigint DEFAULT NULL COMMENT '删除用户',
    created_time datetime    NOT NULL         DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL         DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_username (username) USING BTREE,
    UNIQUE KEY uk_username (username, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户表';

-- 角色表
CREATE TABLE
    IF NOT EXISTS t_role
(
    id           bigint PRIMARY KEY COMMENT '角色id',
    name         varchar(64) NOT NULL COMMENT '角色名',
    status       tinyint(1)  default 0 COMMENT '角色状态, 0-禁用, 1-启用',
    remark  longtext COMMENT '描述',
    deleted_flag tinyint(1)                   DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_user bigint NOT NULL COMMENT '创建用户',
    updated_user bigint NOT NULL COMMENT '更新用户',
    deleted_user bigint DEFAULT NULL COMMENT '删除用户',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime             DEFAULT '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_role_name (name) USING BTREE,
    UNIQUE KEY uk_role_name (name, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='角色表';

-- 用户组表
CREATE TABLE
    IF NOT EXISTS t_user_group
(
    id           bigint PRIMARY KEY COMMENT '主键',
    group_name   VARCHAR(50) NOT NULL COMMENT '组名',
    remark  longtext COMMENT '描述',
    status       tinyint(1)    NOT NULL COMMENT '状态; 0 禁用 1 正常',
    created_user bigint NOT NULL COMMENT '创建用户',
    updated_user bigint NOT NULL COMMENT '更新用户',
    deleted_user bigint DEFAULT NULL COMMENT '删除用户',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记; 0 未删除；1 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime             DEFAULT '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_group_name (group_name) USING BTREE,
    UNIQUE KEY uk_group_name (group_name, deleted_flag, deleted_time)
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户组表';

-- 用户组角色关系表
CREATE TABLE
    IF NOT EXISTS t_user_group_role_relation
(
    id           bigint PRIMARY KEY COMMENT '主键',
    group_id   bigint NOT NULL COMMENT '用户组ID',
    role_id   bigint NOT NULL COMMENT '角色ID',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记; 0 未删除；1 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime             DEFAULT '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_role_id (role_id) USING BTREE,
    KEY idx_group_id (group_id) USING BTREE,
    UNIQUE KEY uk_group_id_role_id (group_id, role_id, deleted_flag, deleted_time)
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户组角色关系表';

-- 用户组用户关系表
CREATE TABLE
    IF NOT EXISTS t_user_group_user_relation
(
    id           bigint PRIMARY KEY COMMENT '主键',
    user_id   bigint NOT NULL COMMENT '用户ID',
    group_id   bigint NOT NULL COMMENT '用户组ID',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记; 0 未删除；1 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime             DEFAULT '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_user_id (user_id) USING BTREE,
    KEY idx_group_id (group_id) USING BTREE,
    UNIQUE KEY uk_user_id_group_id (user_id, group_id, deleted_flag, deleted_time)
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户组用户关系表';

-- 角色权限表
CREATE TABLE
    IF NOT EXISTS t_role_permission
(
    id           bigint PRIMARY KEY COMMENT '主键',
    role_id   bigint NOT NULL COMMENT '角色ID',
    permission_code VARCHAR(255) NOT NULL COMMENT '权限编码',
    created_user bigint NOT NULL COMMENT '创建用户',
    updated_user bigint NOT NULL COMMENT '更新用户',
    deleted_user bigint DEFAULT NULL COMMENT '删除用户',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记; 0 未删除；1 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime             DEFAULT '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_role_id (role_id) USING BTREE,
    KEY idx_permission_code (permission_code) USING BTREE,
    UNIQUE KEY uk_role_id_permission_code (role_id, permission_code, deleted_flag, deleted_time)
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='角色权限表';

-- 消息通知表
CREATE TABLE
    IF NOT EXISTS t_user_notification
(
    id           bigint PRIMARY KEY COMMENT '主键',
    user_id bigint NOT NULL COMMENT '用户ID',
    name varchar(100) NOT NULL COMMENT '通知名称',
    type   VARCHAR(100) NOT NULL COMMENT '通知类型',
    content longtext NOT NULL COMMENT '消息内容',
    status tinyint(1) NOT NULL COMMENT '状态; 0: 未读 1: 已读',
    read_time datetime DEFAULT NULL COMMENT '阅读时间',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记; 0 未删除；1 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime             DEFAULT '1970-01-01 08:00:00' COMMENT '删除时间'
) ENGINE = innodb
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci
    ROW_FORMAT = DYNAMIC COMMENT ='消息通知表';

-- 系统设置表
CREATE TABLE
    IF NOT EXISTS t_system_setting
(
    id           bigint PRIMARY KEY COMMENT '设置id',
    name         varchar(64) NOT NULL COMMENT '设置名称',
    value        varchar(255) DEFAULT NULL COMMENT '设置值',
    remark  longtext COMMENT '描述',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记; 0 未删除；1 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime             DEFAULT '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_setting_name (name) USING BTREE,
    UNIQUE KEY uk_setting_name (name, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='系统设置表';

-- 用户设置表
CREATE TABLE
    IF NOT EXISTS t_user_setting
(
    id           varchar(50) PRIMARY KEY COMMENT '设置id',
    user_id      varchar(50) NOT NULL COMMENT '用户id',
    name         varchar(64) NOT NULL COMMENT '设置名称',
    value        varchar(255) DEFAULT NULL COMMENT '设置值',
    remark  longtext COMMENT '描述',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_setting_name (name) USING BTREE,
    KEY idx_setting_user (user_id) USING BTREE,
    UNIQUE KEY uk_setting_name (name, user_id, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户设置表';
