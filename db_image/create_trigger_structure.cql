CREATE KEYSPACE userdb WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}  AND durable_writes = true;

create table userdb.buyTriggers (
    tid uuid, 
    userId varchar, 
    stock varchar, 
    pendingCash int, 
    triggerValue int, 
    stockAmount int, 
    pendingStocks int, 
    pending boolean, 
    primary key (userId, stock));

create table userdb.sellTriggers (
    tid uuid, 
    userId varchar, 
    stock varchar, 
    pendingCash int, 
    triggerValue int, 
    stockAmount int, 
    pendingStocks int, 
    pending boolean, 
    primary key (userId, stock));
