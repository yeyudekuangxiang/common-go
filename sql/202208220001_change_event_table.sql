ALTER TABLE "public"."event"
    ADD COLUMN "limit" varchar(50) NOT NULL DEFAULT '';

COMMENT ON COLUMN "public"."event"."limit" IS '兑换次数限制
空表示不限制
1-D 按天限制次数
1-W 按周限制次数
1-M 按月限制次数
1-Y 按年限制次数
';