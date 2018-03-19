# triggerserver
Server which handles all the triggers after they have been setup


Data Model Change: ( added extra column "pending" to keep track of processed triggers)

create table buyTriggers (tid uuid, userId varchar, stock varchar, pendingCash int, triggerValue int, stockAmount int, pendingStocks int, pending boolean, primary key (userId, stock));

create table sellTriggers (tid uuid, userId varchar, stock varchar, pendingCash int, triggerValue int, stockAmount int, pendingStocks int, pending boolean, primary key (userId, stock));
