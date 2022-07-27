ALTER TABLE "public"."banner"
    ADD COLUMN "sort" int4 NOT NULL DEFAULT 0,
    ADD COLUMN "scene" varchar(20) NOT NULL DEFAULT 'home',
    ADD COLUMN "app_id" varchar(100) NOT NULL DEFAULT '',
    ADD COLUMN "ext" varchar(1000) NOT NULL DEFAULT '',
    ADD COLUMN "type" varchar(20) NOT NULL DEFAULT 'path',
    ADD COLUMN "status" int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "public"."banner"."id" IS '轮播图';
COMMENT ON COLUMN "public"."banner"."name" IS '轮播图名称';
COMMENT ON COLUMN "public"."banner"."image_url" IS '轮播图图片';
COMMENT ON COLUMN "public"."banner"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."banner"."update_time" IS '更新时间';
COMMENT ON COLUMN "public"."banner"."redirect" IS '跳转路径';
COMMENT ON COLUMN "public"."banner"."sort" IS '排序';
COMMENT ON COLUMN "public"."banner"."scene" IS '证书场景 home首页 event携手 topic社区';
COMMENT ON COLUMN "public"."banner"."app_id" IS '跳转到三方小程序时小程序appid';
COMMENT ON COLUMN "public"."banner"."ext" IS '额外参数';
COMMENT ON COLUMN "public"."banner"."type" IS '跳转类型 mini第三方小程序 path内部小程序路径';
COMMENT ON COLUMN "public"."banner"."status" IS '状态 1上线 2下线';