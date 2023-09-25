create table balances(
  id uuid not null primary key,
  user_id uuid not null unique,
  amount int not null default 0,
  updated_at timestamp not null default current_timestamp
);
insert into balances (id, user_id, amount) values 
('ab596652-e526-4838-ab50-c0caa3d7488b', 'f2f0e0d1-e37e-45c3-ad06-e6c2a66544fc', 1000), 
('5231d5de-3157-41af-b1f6-950d8e12f0ec', '089557bc-ddf2-4ec5-8077-d8bf09fe3ddc', 1000);