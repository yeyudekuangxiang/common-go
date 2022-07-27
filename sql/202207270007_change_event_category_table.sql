ALTER TABLE "public"."event_category"
    ADD COLUMN "sort" int4 NOT NULL DEFAULT 0,
    ADD COLUMN "icon" varchar(500) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."event_category"."id" IS '公益活动分类表';
COMMENT ON COLUMN "public"."event_category"."event_category_id" IS '分类标识';
COMMENT ON COLUMN "public"."event_category"."active" IS '是否上线';
COMMENT ON COLUMN "public"."event_category"."image_url" IS '分类主图';
COMMENT ON COLUMN "public"."event_category"."title" IS '分类名称';
COMMENT ON COLUMN "public"."event_category"."icon" IS '分类图标';