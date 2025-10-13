CREATE KEYSPACE IF NOT EXISTS catalog
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
