alter table "user"
drop column position;

alter table "user"
drop column position_icon;


alter table "user"
    add vip integer default 1 not null;

comment on column "user".vip is '1普通用户 2黄v用户 3蓝v用户';


alter table "user"
    add lohoja integer default 2 not null;

comment on column "user".lohoja is '1乐活家用户 2非乐活家用户';