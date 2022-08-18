/*ok*/
ALTER TABLE "public"."badge"
    ADD COLUMN "image_url" varchar(1000) NOT NULL DEFAULT '',
    ADD COLUMN "is_new" bool NOT NULL DEFAULT false,
    ALTER COLUMN "code" TYPE varchar(100) COLLATE "pg_catalog"."default",
    ALTER COLUMN "code" SET NOT NULL,
    ALTER COLUMN "code" SET DEFAULT '';

COMMENT ON COLUMN "public"."badge"."id" IS '用户证书记录';
COMMENT ON COLUMN "public"."badge"."code" IS '证书编号';
COMMENT ON COLUMN "public"."badge"."openid" IS '用户openid';
COMMENT ON COLUMN "public"."badge"."certificate_id" IS '证书id';
COMMENT ON COLUMN "public"."badge"."product_item_id" IS '对应商品id';
COMMENT ON COLUMN "public"."badge"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."badge"."partnership" IS '合作伙伴';
COMMENT ON COLUMN "public"."badge"."order_id" IS '订单编号';
COMMENT ON COLUMN "public"."badge"."image_url" IS '证书图片';
COMMENT ON COLUMN "public"."badge"."is_new" IS '是否是新获得';