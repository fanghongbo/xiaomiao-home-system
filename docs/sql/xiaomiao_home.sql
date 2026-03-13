USE xiaomiao_home;

-- 用户表
CREATE TABLE
    IF NOT EXISTS t_user
(
    id           bigint PRIMARY KEY COMMENT '主键',
    username     VARCHAR(50) NOT NULL COMMENT '用户名',
    nickname     VARCHAR(50) NOT NULL COMMENT '姓名',
    remark       VARCHAR(255) default '' COMMENT '备注',
    status       tinyint(1)    NOT NULL COMMENT '状态; 0 禁用 1 正常',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记; 0 未删除；1 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime             DEFAULT '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_username (username) USING BTREE,
    KEY idx_nickname (nickname) USING BTREE,
    UNIQUE KEY uk_username (username, deleted_flag, deleted_time)
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户表';

-- 用户角色表
CREATE TABLE
    IF NOT EXISTS t_user_role
(
    id           bigint PRIMARY KEY COMMENT '主键',
    role_id   bigint NOT NULL COMMENT '角色ID',
    user_id   bigint NOT NULL COMMENT '用户ID',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记; 0 未删除；1 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime             DEFAULT '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_role_id (role_id) USING BTREE,
    KEY idx_user_id (user_id) USING BTREE,
    UNIQUE KEY uk_role_id_user_id (role_id, user_id, deleted_flag, deleted_time)
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户角色表';

-- 用户组表
CREATE TABLE
    IF NOT EXISTS t_user_group
(
    id           bigint PRIMARY KEY COMMENT '主键',
    group_name   VARCHAR(50) NOT NULL COMMENT '组名',
    remark       VARCHAR(255) default '' COMMENT '备注',
    status       tinyint(1)    NOT NULL COMMENT '状态; 0 禁用 1 正常',
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

-- 用户组角色表
CREATE TABLE
    IF NOT EXISTS t_user_group_role
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
  ROW_FORMAT = DYNAMIC COMMENT ='用户组角色表';

  -- 用户组用户关系表
CREATE TABLE
    IF NOT EXISTS t_user_user_group
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

-- 角色表
CREATE TABLE
    IF NOT EXISTS t_role
(
    id           bigint PRIMARY KEY COMMENT '主键',
    role_name   VARCHAR(50) NOT NULL COMMENT '角色名',
    remark       VARCHAR(255) default '' COMMENT '备注',
    status       tinyint(1)    NOT NULL COMMENT '状态; 0 禁用 1 正常',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记; 0 未删除；1 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime             DEFAULT '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_role_name (role_name) USING BTREE,
    UNIQUE KEY uk_role_name (role_name, deleted_flag, deleted_time)
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
    notification_name varchar(100) NOT NULL COMMENT '通知名称',
    notification_type   VARCHAR(100) NOT NULL COMMENT '通知类型',
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
