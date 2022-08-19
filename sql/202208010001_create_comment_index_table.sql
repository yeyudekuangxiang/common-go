create table comment_index
(
    id              bigserial
        constraint comment_index_pk
            primary key,
    message         text              not null,
    member_id       bigint            not null,
    root_comment_id bigint  default 0 not null,
    to_comment_id   bigint  default 0,
    floor           integer default 0 not null,
    count           integer default 0 not null,
    root_count      integer default 0 not null,
    like_count      integer default 0 not null,
    hate_count      integer default 0 not null,
    state           integer default 0 not null,
    attrs           integer,
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
    obj_id          bigint            not null,
    obj_type        smallint          not null,
    version         bigint,
    del_reason      varchar
);

comment on table comment_index is '评论表';

comment on column comment_index.message is '评论内容';

comment on column comment_index.member_id is '发表者id';

comment on column comment_index.root_comment_id is '根评论id 不为0是回复评论';

comment on column comment_index.to_comment_id is '父评论ID，为0是root评论';

comment on column comment_index.floor is '评论楼层';

comment on column comment_index.count is '评论总数';

comment on column comment_index.root_count is '根评论总数';

comment on column comment_index.like_count is '点赞数';

comment on column comment_index.hate_count is '点踩总数，备用';

comment on column comment_index.state is '状态 （0-正常，1-隐藏）';

comment on column comment_index.attrs is '属性（bit 1-运营置顶 2-up置顶）
举例：  01 up置顶
举例2： 10 运营置顶';

comment on column comment_index.obj_id is '对象id';

comment on column comment_index.obj_type is '对象类型';

comment on column comment_index.version is '版本号 保留字段';

alter table comment_index
    owner to miniprogram;




create table comment_content
(
    comment_id    bigint not null
        constraint comment_content_pk
            primary key,
    message       varchar,
    at_member_ids varchar,
    device        varchar,
    platform      smallint,
    ip            varchar,
    created_at    timestamp with time zone,
    updated_at    timestamp with time zone
);

comment on column comment_content.comment_id is '和comment_index的id 1对1';

comment on column comment_content.message is '评论内容';

alter table comment_content
    owner to miniprogram;