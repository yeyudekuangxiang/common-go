alter table point_collect_history
    add point varchar(50);
    add additional_order varchar(200);

alter table point_collect_history
alter column type type varchar(50) using type::varchar(50);