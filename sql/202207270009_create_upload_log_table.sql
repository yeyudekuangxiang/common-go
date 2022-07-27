create table upload_log
(
    id            bigserial
        primary key,
    log_id        varchar(100)                                not null,
    oss_path      varchar(1000)                               not null,
    size          bigint                                      not null,
    url           varchar(1000) default ''::character varying not null,
    user_id       bigint        default 0                     not null,
    scene_id      integer                                     not null,
    created_at    timestamp                                   not null,
    updated_at    timestamp                                   not null,
    operator_id   bigint        default 0                     not null,
    operator_type smallint      default 1                     not null
);

comment on column upload_log.id is '上传文件表';

comment on column upload_log.log_id is '文件编号';

comment on column upload_log.oss_path is '阿里云oss路径';

comment on column upload_log.size is '文件大小 单位B';

comment on column upload_log.url is '文件链接';

comment on column upload_log.user_id is '用户编号';

comment on column upload_log.scene_id is '上传场景编号';

comment on column upload_log.operator_id is '操作员';

comment on column upload_log.operator_type is '操作员类型 1用户 2管理员';

alter table upload_log
    owner to miniprogram;

