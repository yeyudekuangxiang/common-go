ALTER TABLE "public"."event"
    ADD COLUMN "is_show" int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "public"."event"."is_show" IS '是否显示 1显示 2不显示';