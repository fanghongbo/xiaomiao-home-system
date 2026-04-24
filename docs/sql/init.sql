# CREATE DATABASE xiaomiao_home CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE xiaomiao_home;

-- 用户表
CREATE TABLE
    IF NOT EXISTS t_user
(
    id           bigint PRIMARY KEY COMMENT '用户id',
    version      int NOT NULL COMMENT '版本号',
    status       tinyint(1)  default 0 COMMENT '账户状态, 0: 禁用, 1: 启用',
    deleted_flag tinyint(1)                   DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL         DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL         DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_id_version (id, version) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户表';

  -- 用户信息版本表
CREATE TABLE
    IF NOT EXISTS t_user_version
(
    id           bigint PRIMARY KEY COMMENT '用户信息版本id',
    version      int NOT NULL COMMENT '版本号',
    user_id      bigint NOT NULL COMMENT '用户id',
    nickname     varchar(64) NOT NULL COMMENT '昵称',
    gender       tinyint(1)  default 0 COMMENT '性别, 0: 保密, 1: 男, 2: 女',
    birthday     date DEFAULT NULL COMMENT '生日',
    province_id  bigint DEFAULT NULL COMMENT '省份id',
    city_id      bigint DEFAULT NULL COMMENT '城市id',
    address      varchar(255) DEFAULT NULL COMMENT '详细地址',
    avatar       varchar(255) COMMENT '头像地址',
    signature    varchar(100) COMMENT '个人签名',
    remark  longtext COMMENT '描述',
    deleted_flag tinyint(1)                   DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL         DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL         DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_user_id_version (user_id, version) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户信息版本表';

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
    identity_type   varchar(32) NOT NULL COMMENT '身份类型 password:密码校验 sms:短信校验 wechat:微信校验 qq:QQ校验 email:邮箱校验',
    identity_id    varchar(191) NOT NULL COMMENT '身份id 手机号:手机号 微信:unionid QQ:openid email:邮箱地址 password:用户名',
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
    user_id      bigint NOT NULL COMMENT '用户id',
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

-- 小猫信息表
CREATE TABLE
    IF NOT EXISTS t_cat
(
    id           bigint PRIMARY KEY COMMENT '小猫id',
    version      int NOT NULL COMMENT '版本号',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_id_version (id, version) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='小猫信息表';

  -- 小猫信息版本表
CREATE TABLE
    IF NOT EXISTS t_cat_version
(
    id           bigint PRIMARY KEY COMMENT '小猫信息版本id',
    version      int NOT NULL COMMENT '版本号',
    cat_id       bigint NOT NULL COMMENT '小猫id',
    name         varchar(64) NOT NULL COMMENT '小猫名称',
    gender       tinyint(1)  DEFAULT 0 COMMENT '性别, 0: 未知, 1: 弟弟, 2: 妹妹',
    cat_type     tinyint(1) DEFAULT 0 COMMENT '类型, 0: 未知, 1: 流浪猫, 2: 我的小猫',
    breed_type   int(11) DEFAULT 0 COMMENT '品种类型 0: 未知',
    weight       DECIMAL(5,2) DEFAULT 0 COMMENT '小猫体重, 单位kg, 例如:2.50',
    birthday     date DEFAULT NULL COMMENT '生日',
    neuter_status tinyint(1) DEFAULT 0 COMMENT '绝育状态, 0: 未知, 1: 已绝育, 2: 未绝育',
    health_status tinyint(1) DEFAULT 0 COMMENT '健康状态, 0: 未知, 1: 健康, 2: 生病, 3: 残疾, 4: 其他',
    health_desc varchar(500) DEFAULT '' COMMENT '疾病/缺陷说明',
    dewormed_status tinyint(1) DEFAULT 0 COMMENT '驱虫状态, 0: 未知, 1: 已驱虫, 2: 未驱虫',
    vaccine_status tinyint(1) DEFAULT 0 COMMENT '疫苗状态, 0: 未知, 1: 全程接种, 2: 部分接种, 3: 未接种',
    vaccine_types varchar(100) DEFAULT '' COMMENT '疫苗类型,逗号分隔: 0: 未知, 1: 猫三联, 2: 狂犬疫苗, 3: 猫白血病, 4: 其他',
    vaccine_last_date date DEFAULT NULL COMMENT '最后接种日期',
    vaccine_cert_image varchar(500) DEFAULT '' COMMENT '疫苗本凭证图片地址',
    remark  longtext COMMENT '描述',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_name (name) USING BTREE,
    KEY idx_cat_type (cat_type) USING BTREE,
    KEY idx_breed_type (breed_type) USING BTREE,
    KEY idx_weight (weight) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='小猫信息表';

-- 小猫图片关联表
CREATE TABLE
    IF NOT EXISTS t_cat_image
(
    id           bigint PRIMARY KEY COMMENT '关联id',
    cat_id   bigint NOT NULL COMMENT '小猫版本id',
    image    varchar(500) NOT NULL COMMENT '图片url',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_cat_id (cat_id) USING BTREE,
    KEY idx_image_url (image_url) USING BTREE,
    UNIQUE KEY uk_cat_id_image (cat_id, image, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='小猫图片关联表';

-- 发布内容表
CREATE TABLE
    IF NOT EXISTS t_post
(
    id           bigint PRIMARY KEY COMMENT '发布信息id',
    version      int NOT NULL COMMENT '版本号',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_id_version (id, version) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='发布内容表';

  -- 发布内容版本表
CREATE TABLE
    IF NOT EXISTS t_post_version
(
    id           bigint PRIMARY KEY COMMENT '发布信息版本id',
    version      int NOT NULL COMMENT '版本号',
    post_id      bigint NOT NULL COMMENT '发布id',
    title         varchar(64) NOT NULL COMMENT '标题',
    post_type  tinyint(1) NOT NULL COMMENT '发布类型, 1: 领养, 2: 寻猫, 3: 日常, 4: 求助',
    province_id  bigint DEFAULT NULL COMMENT '省份id',
    city_id      bigint DEFAULT NULL COMMENT '城市id',
    lost_time    datetime DEFAULT NULL COMMENT '丢失时间',
    address      varchar(255) DEFAULT NULL COMMENT '详细地址/丢失地址',
    audit_status tinyint(1) NOT NULL COMMENT '审核状态, 0: 待审核, 1: 审核通过, 2: 审核不通过',
    audit_remark longtext DEFAULT NULL COMMENT '审核备注',
    audit_time datetime DEFAULT NULL COMMENT '审核时间',
    post_status tinyint(1) NOT NULL COMMENT '发布状态, 0: 待发布, 1: 已发布, 2: 已下架',
    remark  longtext COMMENT '描述',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_post_type (post_type) USING BTREE,
    KEY idx_audit_status (audit_status) USING BTREE,
    KEY idx_post_status (post_status) USING BTREE,
    KEY idx_province_id (province_id) USING BTREE,
    KEY idx_city_id (city_id) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='发布内容表';

-- 发布小猫关联表
CREATE TABLE
    IF NOT EXISTS t_user_post
(
    id           bigint PRIMARY KEY COMMENT '关联id',
    user_id       bigint NOT NULL COMMENT '用户id',
    post_id   bigint NOT NULL COMMENT '发布id',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_post_id (post_id) USING BTREE,
    KEY idx_user_id (user_id) USING BTREE,
    UNIQUE KEY uk_post_id_user_id (post_id, user_id, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户发布内容关联表';

-- 用户喜欢关联表
CREATE TABLE
    IF NOT EXISTS t_user_like
(
    id           bigint PRIMARY KEY COMMENT '关联id',
    user_id       bigint NOT NULL COMMENT '用户id',
    post_id   bigint NOT NULL COMMENT '发布id',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_post_id (post_id) USING BTREE,
    KEY idx_user_id (user_id) USING BTREE,
    UNIQUE KEY uk_user_id_post_id (user_id, post_id, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户喜欢关联表';

-- 用户收藏表
CREATE TABLE
    IF NOT EXISTS t_user_collect
(
    id           bigint PRIMARY KEY COMMENT '收藏id',
    user_id      bigint NOT NULL COMMENT '用户id',
    post_id   bigint NOT NULL COMMENT '发布id',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_user_id (user_id) USING BTREE,
    KEY idx_post_id (post_id) USING BTREE,
    UNIQUE KEY uk_user_id_post_id (user_id, post_id, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户收藏表';

-- 发布小猫关联表
CREATE TABLE
    IF NOT EXISTS t_post_cat
(
    id           bigint PRIMARY KEY COMMENT '关联id',
    post_id   bigint NOT NULL COMMENT '发布id',
    cat_id       bigint NOT NULL COMMENT '小猫版本id',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_post_id (post_id) USING BTREE,
    KEY idx_cat_id (cat_id) USING BTREE,
    UNIQUE KEY uk_post_id_cat_id (post_id, cat_id, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='发布小猫关联表';

-- 发布图片关联表
CREATE TABLE
    IF NOT EXISTS t_post_image
(
    id           bigint PRIMARY KEY COMMENT '关联id',
    post_id   bigint NOT NULL COMMENT '发布id',
    image    varchar(255) NOT NULL COMMENT '图片url',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_post_id (post_id) USING BTREE,
    KEY idx_image_url (image) USING BTREE,
    UNIQUE KEY uk_post_id_image (post_id, image, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='发布小猫关联表';

-- 用户小猫关联表
CREATE TABLE
    IF NOT EXISTS t_user_cat
(
    id           bigint PRIMARY KEY COMMENT '关联id',
    user_id      bigint NOT NULL COMMENT '用户id',
    cat_id       bigint NOT NULL COMMENT '小猫id',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_user_id (user_id) USING BTREE,
    KEY idx_cat_id (cat_id) USING BTREE,
    UNIQUE KEY uk_user_id_cat_id (user_id, cat_id, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='用户小猫关联表';

-- 领养表
CREATE TABLE
    IF NOT EXISTS t_adopt
(
    id           bigint PRIMARY KEY COMMENT '关联id',
    user_id      bigint NOT NULL COMMENT '用户id',
    cat_id       bigint NOT NULL COMMENT '小猫id',
    adopt_status tinyint(1) NOT NULL COMMENT '领养状态, 0: 待领养, 1: 已领养, 2: 已取消',
    remark longtext DEFAULT NULL COMMENT '备注',
    deleted_flag tinyint(1)           DEFAULT 0 COMMENT '删除标记, 0: 未删除,  1: 已删除',
    created_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_time datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_time datetime                     DEFAULT  '1970-01-01 08:00:00' COMMENT '删除时间',
    KEY idx_user_id (user_id) USING BTREE,
    KEY idx_cat_id (cat_id) USING BTREE,
    UNIQUE KEY uk_user_id_cat_id (user_id, cat_id, deleted_flag, deleted_time) USING BTREE
) ENGINE = innodb
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC COMMENT ='领养表';
