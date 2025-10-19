const dotenv = require('dotenv');
dotenv.config();

// this is the config object that will be used throughout the application
// the values are taken from the .env file
const config = {
  application: {
    port: process.env.APPLICATION_PORT,
  },
  postgres: {
    host: process.env.PGHOST,
    database: process.env.PGDATABASE,
    port: process.env.PGPORT,
    user: process.env.PGUSER,
    password: process.env.PGPASSWORD,
  },
  rabbitMq: {
    server: process.env.RABBITMQ_SERVER,
  },
  mail: {
    host: process.env.SMTP_HOST,
    port: Number(process.env.SMTP_PORT),
    user: process.env.SMTP_USER,
    password: process.env.SMTP_PASSWORD,
  }
}

module.exports = config;
