create or replace function fill_updated_at()
returns trigger as $$
begin
   new.updated_at = now();
   return new;
end;
$$ language 'plpgsql'
;

create table rarity (
    id          numeric(10) not null,
    name        varchar(30) not null,
    alias       char(1) not null,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_rarity primary key(id)
)
;

create trigger update_rarity_updated_at before update on rarity for each row execute procedure fill_updated_at();

create index ix_rarity_name on rarity (name asc);
create index ix_rarity_alias on rarity (alias asc);

insert into rarity values (1, 'Common', 'C');
insert into rarity values (2, 'Uncommon', 'U');
insert into rarity values (3, 'Rare', 'R');
insert into rarity values (4, 'Mythic Rare', 'M');
insert into rarity values (5, 'Special', 'S');

create table type (
    id numeric(10) not null,
    name varchar(50) not null,
    number numeric(3) not null,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_type primary key (id)
)
;

create trigger update_type_updated_at before update on type for each row execute procedure fill_updated_at();

create index ix_type_name on type (name asc);
create index ix_type_number on type (number asc);

insert into type values (1, 'Creature', 1);
insert into type values (2, 'Artifact', 2);
insert into type values (3, 'EnchantmENT', 3);
insert into type values (4, 'Instant', 4);
insert into type values (5, 'Sorcery', 5);
insert into type values (6, 'PlaneswaLKER', 6);
insert into type values (7, 'Land', 7);

create table expansion (
    id numeric(10) not null,
    name varchar(100) not null,
    alias varchar(10) not null,
    metadata jsonb,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_expansion primary key (id)
)
;

create trigger update_expansion_updated_at before update on expansion for each row execute procedure fill_updated_at();

create index ix_expansion_name on expansion (name asc);
create index ix_expansion_alias on expansion (alias asc);

create table expansion_asset (
    id_expansion numeric(10) not null,
    id_rarity numeric(10) not null,
    id_asset     numeric(10) not null,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_expansion_asset primary key (id_expansion, id_rarity),
    constraint fk_expansion_asset_expansion foreign key (id_expansion) references expansion (id),
    constraint fk_expansion_asset_rarity foreign key (id_rarity) references rarity (id)
);

create trigger update_expansion_asset_updated_at before update on expansion_asset for each row execute procedure fill_updated_at();

create table card (
    id                     numeric(10) not null,
    name                   varchar(255) not null,
    id_expansion           numeric(10) not null,
    id_asset               numeric(10) not null,
    id_rarity              numeric(10),
    id_type                numeric(10),
    type                   varchar(255),
    multiverseid           varchar(50),
    multiverse_number      varchar(15),
    manacost               varchar(255),
    converted_manacost     numeric(12,2),
    color                  varchar(50),
    text                   varchar(4000),
    combat_power           varchar(255),
    flavor                 varchar(4000),
    artist                 varchar(255),
    rate                   numeric(14,4) not null default 0,
    rate_votes             numeric(10) not null default 0,
    power                  numeric(12,2),
    toughness              numeric(12,2),
    loyalty                numeric(12,2),
    metadata               jsonb,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_card  primary key (id),
    constraint fk_card_expansion foreign key (id_expansion) references expansion (id),
    constraint fk_card_rarity foreign key (id_rarity) references rarity (id),
    constraint fk_card_type foreign key (id_type) references type (id)
)
;

create trigger update_card_updated_at before update on card for each row execute procedure fill_updated_at();

create index ix_card_name on card (name asc);
create index ix_card_multiverseid on card (multiverseid asc);
create index ix_card_manacost on card (manacost asc);
create index ix_card_converted_manacost on card (converted_manacost asc);
create index ix_card_type on card (type asc);
