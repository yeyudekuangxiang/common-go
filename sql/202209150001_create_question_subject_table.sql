//题目表

CREATE TABLE "public"."question_subject"
(
    "id"          serial8                                      NOT NULL,
    "question_id" int8 DEFAULT 0,
    "category_id" int8 DEFAULT 0,
    "subject_id"  int8 DEFAULT 0,
    "title"       varchar(1000) COLLATE "pg_catalog"."default" NOT NULL,
    "type"        int4 DEFAULT 1,
    "is_hide"     int4 DEFAULT 1,
    "remind"      text COLLATE "pg_catalog"."default",
    "sort"        int4 DEFAULT 0,
    "created_at"  timestamptz(6) NOT NULL DEFAULT now(),
    "updated_at"  timestamptz(6) NOT NULL DEFAULT now(),
    CONSTRAINT "question_subject_pkey" PRIMARY KEY ("id")
);

ALTER TABLE "public"."question_subject" OWNER TO "miniprogram";

COMMENT
ON COLUMN "public"."question_subject"."id" IS '问卷题表';
COMMENT
ON COLUMN "public"."question_subject"."title" IS '标题';
COMMENT
ON COLUMN "public"."question_subject"."type" IS '问卷类型：1正常 2 填空';
COMMENT
ON COLUMN "public"."question_subject"."is_hide" IS '1 隐藏 0 正常';
COMMENT
ON COLUMN "public"."question_subject"."remind" IS '提醒';
COMMENT
ON COLUMN "public"."question_subject"."sort" IS '排序';
COMMENT
ON COLUMN "public"."question_subject"."question_id" IS '父级id';
COMMENT
ON COLUMN "public"."question_subject"."category_id" IS '分类id';
