//用户表

CREATE TABLE "public"."question_user"
(
    "user_id"       int8                                        NOT NULL,
    "third_id"      varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
    "invited_by_id" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
    "phone"         varchar(100) COLLATE "pg_catalog"."default",
    "channel"       varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
    "ip"            varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
    "city"          varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
    "created_at"    timestamptz(6) NOT NULL DEFAULT now(),
    "updated_at"    timestamptz(6) NOT NULL DEFAULT now(),
    CONSTRAINT "question_user_pkey" PRIMARY KEY ("user_id")
);

ALTER TABLE "public"."question_user" OWNER TO "miniprogram";

COMMENT
ON COLUMN "public"."question_user"."user_id" IS '问卷答案表';
COMMENT
ON COLUMN "public"."question_user"."third_id" IS '答题人openid';
COMMENT
ON COLUMN "public"."question_user"."invited_by_id" IS '邀请人openid';
COMMENT
ON COLUMN "public"."question_user"."phone" IS '电话';
COMMENT
ON COLUMN "public"."question_user"."channel" IS '渠道';
COMMENT
ON COLUMN "public"."question_user"."ip" IS '提交ip';
COMMENT
ON COLUMN "public"."question_user"."city" IS '提交城市';
