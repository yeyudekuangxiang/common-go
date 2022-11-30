--
-- alter table "user"
--     add vip integer default 1 not null;
--
-- comment on column "user".vip is '1普通用户 2黄v用户 3蓝v用户';


alter table "user"
    add partners integer default 0 not null;

comment on column "user".partners is '0非乐活家 1乐活家 ';

alter table "user"
    add status integer default 1 not null;

comment on column "user".status is '0全部 1正常 2禁言 3封号';

alter table "user"
    add auth integer default 2 not null;

comment on column "user".auth is '权限字段：0浏览 1评论 2无权限';