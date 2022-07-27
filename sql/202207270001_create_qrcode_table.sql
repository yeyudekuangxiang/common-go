create table qr_code
(
    id             bigserial
        constraint qr_code_pk
            primary key,
    qr_code_id     varchar(255)                               not null
        unique,
    openid         varchar(255) default ''::character varying not null,
    description    varchar(255) default ''::character varying not null,
    image_path     varchar(500)                               not null,
    key            varchar(255) default ''::character varying not null,
    ext            varchar(255)                               not null,
    created_at     timestamp with time zone                   not null,
    qr_code_scene  varchar(255)                               not null,
    qr_code_source varchar(255)                               not null,
    content        varchar(255)                               not null
);

comment on column qr_code.id is '自增编号';

comment on column qr_code.qr_code_id is '二维码id';

comment on column qr_code.openid is 'openid';

comment on column qr_code.description is '二维码描述';

comment on column qr_code.image_path is '二维码路径';

comment on column qr_code.key is '二维码唯一key';

comment on column qr_code.ext is '额外参数';

comment on column qr_code.created_at is '创建时间';

comment on column qr_code.qr_code_scene is '场景';

comment on column qr_code.qr_code_source is '来源 ';

comment on column qr_code.content is '内容';

alter table qr_code
    owner to miniprogram;

create unique index qrcode_key_union_index
    on qr_code (key);

