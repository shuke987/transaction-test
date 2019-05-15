begin;
update bank set money=money-7 where id=1;
update bank set money=money+7 where id=2;
commit;