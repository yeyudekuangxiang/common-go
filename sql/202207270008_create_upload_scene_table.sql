/*ok*/
create table upload_scene
(
    id         serial
        primary key,
    max_count  integer           not null,
    max_size   bigint            not null,
    must_login boolean           not null,
    oss_dir    varchar(100)      not null,
    mime_types varchar(500)      not null,
    scene      varchar(50)       not null,
    scene_name varchar(50)       not null,
    created_at timestamp         not null,
    updated_at timestamp         not null,
    max_age    integer default 0 not null
);

comment on column upload_scene.id is '上传文件记录表';

comment on column upload_scene.max_count is '每天最多上传多少次文件';

comment on column upload_scene.max_size is '上传文件大小限制 单位B 1MB=1024KB=1048576B';

comment on column upload_scene.must_login is '用户是否必须登录';

comment on column upload_scene.oss_dir is '对象存储路径';

comment on column upload_scene.mime_types is '可上传的文件mime类型多个用英文逗号隔开 image/png,image/jpg';

comment on column upload_scene.scene is '上传场景标识 必须是英文字母 例如 userAvatar';

comment on column upload_scene.scene_name is '上传场景标识名称 例如 用户头像';

comment on column upload_scene.max_age is '缓存时长 单位秒 0表示不缓存';

alter table upload_scene
    owner to miniprogram;

