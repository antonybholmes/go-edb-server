rm data/users.db
cat users.sql | sqlite3 data/users.db
