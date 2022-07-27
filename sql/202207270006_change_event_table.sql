ALTER TABLE "public"."event"
    ADD COLUMN "sort" int4 NOT NULL DEFAULT 0,
    ADD COLUMN "event_template_type" varchar(50) NOT NULL DEFAULT '',
    ADD COLUMN "tag" varchar(255) NOT NULL DEFAULT '',
    ADD COLUMN "template_setting" varchar(2000) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."event"."id" IS '公益活动表';
COMMENT ON COLUMN "public"."event"."event_category_id" IS '公益活动所属分类标识';
COMMENT ON COLUMN "public"."event"."event_id" IS '公益活动标识';
COMMENT ON COLUMN "public"."event"."title" IS '公益活动标题';
COMMENT ON COLUMN "public"."event"."subtitle" IS '公益活动副标题';
COMMENT ON COLUMN "public"."event"."active" IS '是否上线';
COMMENT ON COLUMN "public"."event"."cover_image_url" IS '项目主图 343 × 200';
COMMENT ON COLUMN "public"."event"."start_time" IS '公益活动开始时间';
COMMENT ON COLUMN "public"."event"."end_time" IS '公益活动结束时间';
COMMENT ON COLUMN "public"."event"."product_item_id" IS '关联的商品编号';
COMMENT ON COLUMN "public"."event"."participation_count" IS '已参与次数';
COMMENT ON COLUMN "public"."event"."participation_subtitle" IS '用于展示支持次数或者co2';
COMMENT ON COLUMN "public"."event"."sort" IS '排序 从小到大排序';
COMMENT ON COLUMN "public"."event"."event_template_type" IS '证书模版类型';
COMMENT ON COLUMN "public"."event"."tag" IS '标签,多个标签用英文逗号隔开';
COMMENT ON COLUMN "public"."event"."template_setting" IS '公益活动模版配置';
