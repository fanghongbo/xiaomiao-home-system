# CREATE DATABASE xiaomiao_home CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE xiaomiao_home;

-- 用户表
CREATE TABLE
    IF NOT EXISTS t_user
(
    id           bigint PRIMARY KEY COMMENT '用户id',
    nickname     varchar(64) NOT NULL COMMENT '昵称',
    status       tinyint(1)  default 0 COMMENT '账户状态, 0: 禁用, 1: 启用',
    avatar       varchar(255) COMMENT '头像地址',
    signature    varchar(255) COMMENT '个人签名',
    remark  longtext COMMENT '描述',
    deleted_flag tinyint(1)                   DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL         DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL         DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    UNIQUE KEY uk_nickname (nickname, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户表';

-- 用户密码凭据（仅开通账号密码登录的用户有记录，与 t_user 一对一）
CREATE TABLE
    IF NOT EXISTS t_user_password
(
    id           bigint PRIMARY KEY COMMENT '主键',
    user_id      bigint NOT NULL COMMENT '用户id',
    password     varchar(255) NOT NULL COMMENT '密码哈希',
    salt         varchar(32) NOT NULL COMMENT '密码加盐',
    deleted_flag tinyint(1) DEFAULT 0 COMMENT '删除标记, 0: 未删除, 1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime    DEFAULT '1970-01-01 08:00:00' COMMENT '删除时间',
    UNIQUE KEY uk_user_id (user_id, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户密码凭据';

-- 用户身份表
CREATE TABLE
    IF NOT EXISTS t_user_identity
(
    id           bigint PRIMARY KEY COMMENT '用户身份id',
    user_id      bigint NOT NULL COMMENT '用户id',
    identity_type   varchar(32) NOT NULL COMMENT '身份类型 password:密码登录 sms:短信登录 wechat:微信登录 qq:QQ登录 email:邮箱登录',
    identity_id    varchar(191) NOT NULL COMMENT '身份id 账号:用户名 手机号:手机号 微信:unionid QQ:openid 邮箱:email',
    verified_flag  tinyint(1) DEFAULT 0 COMMENT '是否完成验证 0:未验证 1:已验证',
    remark  longtext COMMENT '描述',
    deleted_flag tinyint(1)                   DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL         DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL         DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_user_id (user_id) USING BTREE,
    UNIQUE KEY uk_identity_type_identity_id (identity_type, identity_id, deleted_flag, deleted_time) USING BTREE,
    UNIQUE KEY uk_user_id_identity_type_identity_id (user_id, identity_type, identity_id, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户身份表';

-- 角色表
CREATE TABLE
    IF NOT EXISTS t_role
(
    id           bigint PRIMARY KEY COMMENT '角色id',
    name         varchar(64) NOT NULL COMMENT '角色名',
    status       tinyint(1)  default 0 COMMENT '角色状态, 0-禁用, 1-启用',
    remark  longtext COMMENT '描述',
    deleted_flag tinyint(1)                   DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime             DEFAULT '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_role_name (name) USING BTREE,
    UNIQUE KEY uk_role_name (name, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='角色表';

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
    id           bigint PRIMARY KEY COMMENT '设置id',
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
