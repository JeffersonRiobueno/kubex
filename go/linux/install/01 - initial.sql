CREATE SCHEMA `kubex` ;

create table kubex.runner(
    id int not null auto_increment primary key,
    proyect varchar(500),
    rama varchar(500),
    folder varchar(500),
    pApp varchar(500),
    pBd varchar(500),
    status int

);


create table kubex.portApp(
    id int not null auto_increment primary key,
    portApp int
);

create table kubex.portBD(
    id int not null auto_increment primary key,
    portBD int
);