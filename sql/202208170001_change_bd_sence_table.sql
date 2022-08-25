-- alter table bd_scene
--     rename column key to "secret";

alter table bd_scene
    add app_id varchar;

comment on column bd_scene.app_id is 'app_id';

alter table bd_scene
    add domain varchar default 'e'::bpchar not null;