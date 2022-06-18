```sh
docker run -it -d --name postgresql --env-file db.env -p 5432:5432 bitnami/postgresql:14

sudo apt install postgresql-client-common postgresql-client-14

psql -h 127.0.0.1 -U admin -p 5432 -d conduit
```

```sql
conduit=> \l
                                  List of databases
   Name    |  Owner   | Encoding |   Collate   |    Ctype    |   Access privileges   
-----------+----------+----------+-------------+-------------+-----------------------
 conduit   | admin    | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =Tc/admin            +
           |          |          |             |             | admin=CTc/admin
 postgres  | postgres | UTF8     | en_US.UTF-8 | en_US.UTF-8 | 
 template0 | postgres | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =c/postgres          +
           |          |          |             |             | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =c/postgres          +
           |          |          |             |             | postgres=CTc/postgres

conduit=> \c
You are now connected to database "conduit" as user "admin".
conduit=> \d
 public | article_tags      | table    | admin
 public | articles          | table    | admin
 public | articles_id_seq   | sequence | admin
 public | comments          | table    | admin
 public | comments_id_seq   | sequence | admin
 public | favorites         | table    | admin
 public | followings        | table    | admin
 public | schema_migrations | table    | admin
 public | tags              | table    | admin
 public | tags_id_seq       | sequence | admin
 public | users             | table    | admin
 public | users_id_seq      | sequence | admin
```

refs: https://techviewleo.com/how-to-install-postgresql-database-on-ubuntu/
