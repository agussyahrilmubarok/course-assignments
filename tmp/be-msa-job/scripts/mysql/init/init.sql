CREATE DATABASE IF NOT EXISTS ebank_accounts_db;
CREATE DATABASE IF NOT EXISTS ebank_cards_db;
CREATE DATABASE IF NOT EXISTS ebank_loans_db;

CREATE USER IF NOT EXISTS 'accounts_user'@'%' IDENTIFIED BY 'accounts_pass';
CREATE USER IF NOT EXISTS 'cards_user'@'%' IDENTIFIED BY 'cards_pass';
CREATE USER IF NOT EXISTS 'loans_user'@'%' IDENTIFIED BY 'loans_pass';

GRANT ALL PRIVILEGES ON ebank_accounts_db.* TO 'accounts_user'@'%';
GRANT ALL PRIVILEGES ON ebank_cards_db.* TO 'cards_user'@'%';
GRANT ALL PRIVILEGES ON ebank_loans_db.* TO 'loans_user'@'%';

FLUSH PRIVILEGES;
