CREATE TABLE "public"."scan_log" (
                                "id" serial8 NOT NULL,
                                "image_url" varchar(1000) NOT NULL,
                                "hash" varchar(255) NOT NULL,
                                "count" int4 NOT NULL DEFAULT 0,
                                "scan_result" varchar(1000) NOT NULL DEFAULT '',
                                "created_at" timestamptz NOT NULL,
                                "updated_at" timestamptz NOT NULL,
                                PRIMARY KEY ("id"),
                                CONSTRAINT "scan_hash_index" UNIQUE ("hash")
)
;

COMMENT ON COLUMN "public"."scan_log"."image_url" IS '图片地址';

COMMENT ON COLUMN "public"."scan_log"."hash" IS '图片hash sha256';

COMMENT ON COLUMN "public"."scan_log"."count" IS '扫描次数';

COMMENT ON COLUMN "public"."scan_log"."created_at" IS '创建时间';

COMMENT ON COLUMN "public"."scan_log"."updated_at" IS '更新时间';

COMMENT ON COLUMN "public"."scan_log"."scan_result" IS '扫描结果';