/*ok*/
ALTER TABLE "public"."qr_code"
    ADD COLUMN "image_path" varchar(500) NOT NULL DEFAULT '',
    ADD COLUMN "key" varchar(255) NOT NULL DEFAULT '',
    ADD COLUMN "ext" varchar(255) NOT NULL DEFAULT '',
    ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT '2022-07-28 14:00:00+08',
    ADD COLUMN "qr_code_scene" varchar(255) NOT NULL DEFAULT '',
    ADD COLUMN "qr_code_source" varchar(255) NOT NULL DEFAULT '',
    ADD COLUMN "content" varchar(255) NOT NULL DEFAULT '',

    ALTER COLUMN "qr_code_id" TYPE varchar(255) COLLATE "pg_catalog"."default",
    ALTER COLUMN "openid" TYPE varchar(255) COLLATE "pg_catalog"."default",
    ALTER COLUMN "openid" SET DEFAULT '',
    ALTER COLUMN "openid" SET NOT NULL ,
    ALTER COLUMN "description" TYPE varchar(255) COLLATE "pg_catalog"."default",
    ALTER COLUMN "description" SET DEFAULT '',
    ALTER COLUMN "description" SET NOT NULL ,
    ALTER COLUMN "image_url" SET DEFAULT '废弃',
    ALTER COLUMN "qr_code_type" SET DEFAULT '废弃';

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

