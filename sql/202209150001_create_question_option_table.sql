CREATE TABLE "public"."question_option"
(
    "id"              serial8                                     NOT NULL,
    "title"           varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
    "sort"            int4  DEFAULT 0,
    "subject_id"      int8  DEFAULT 0,
    "carbon"          float DEFAULT 0,
    "remind"          text COLLATE "pg_catalog"."default",
    "jump_subject"    int8  DEFAULT 0,
    "related_subject" varchar(100) COLLATE "pg_catalog"."default",
    "created_at"      timestamptz(6) NOT NULL DEFAULT now(),
    "updated_at"      timestamptz(6) NOT NULL DEFAULT now(),
    CONSTRAINT "question_option_pkey" PRIMARY KEY ("id")
);

ALTER TABLE "public"."question_option" OWNER TO "miniprogram";

COMMENT
ON COLUMN "public"."question_option"."id" IS '问卷答案表';
COMMENT
ON COLUMN "public"."question_option"."title" IS '选项';
COMMENT
ON COLUMN "public"."question_option"."sort" IS '排序';
COMMENT
ON COLUMN "public"."question_option"."subject_id" IS '题目id';
COMMENT
ON COLUMN "public"."question_option"."remind" IS '提醒';
COMMENT
ON COLUMN "public"."question_option"."jump_subject" IS '跳转题目';
COMMENT
ON COLUMN "public"."question_option"."related_subject" IS '关联题目';
