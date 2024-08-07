create table token
(
    name           varchar not null,
    token          varchar not null,
    owner_id       bigint  not null,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    last_login_at  timestamp with time zone,
    binding_bot_id bigint,
    constraint token_pk
        primary key (token),
    constraint token_pk_2
        unique (name)
);

create index token_owner_id_index
    on token (owner_id);

