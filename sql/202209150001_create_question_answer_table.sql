//答案表
CREATE TABLE "public"."question_answer"
(
    "id"          serial8 NOT NULL,
    "question_id" int8    NOT NULL,
    "subject_id"  int8    NOT NULL,
    "user_id"     int8    NOT NULL,
    "carbon"      float DEFAULT 0,
    "answer"      text COLLATE "pg_catalog"."default",
    "created_at"  timestamptz(6) NOT NULL DEFAULT now(),
    "updated_at"  timestamptz(6) NOT NULL DEFAULT now(),
    CONSTRAINT "question_answer_pkey" PRIMARY KEY ("id")
);

ALTER TABLE "public"."question_answer" OWNER TO "miniprogram";

COMMENT
ON COLUMN "public"."question_answer"."id" IS '问卷答案表';
COMMENT
ON COLUMN "public"."question_answer"."question_id" IS '问卷id';
COMMENT
ON COLUMN "public"."question_answer"."subject_id" IS '题目id';
COMMENT
ON COLUMN "public"."question_answer"."user_id" IS '用户id';
COMMENT
ON COLUMN "public"."question_answer"."answer" IS '答案';
