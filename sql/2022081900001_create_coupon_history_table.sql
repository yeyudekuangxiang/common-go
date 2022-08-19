create table coupon_history
(
    id               bigserial
        constraint event_coupon_history_pk
            primary key,
    open_id          varchar                                   not null,
    coupon_type      varchar                                   not null,
    code varchar                                   ,
    create_time      timestamp with time zone               not null,
    update_time      timestamp with time zone default now() not null,
    constraint event_coupon_history_openid_coupontype_code_key
        unique (open_id, coupon_type, code)
);

alter table coupon_history
    owner to miniprogram;