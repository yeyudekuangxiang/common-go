create table comment_like
(
    id         bigint            not null
        constraint id
            primary key,
    comment_id bigint            not null,
    user_id    bigint            not null,
    status     integer default 1 not null,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

comment on table comment_like is '评论点赞记录表';

comment on column comment_like.comment_id is '评论id';

comment on column comment_like.user_id is '点赞用户id';

comment on column comment_like.status is '0取消赞 1点赞 ';

alter table comment_like
    owner to miniprogram;

