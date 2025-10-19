const dotenv = require('dotenv');

dotenv.config();

/**
 * this object contains all the configuration needed by the application
 * the values are taken from the .env file that gives the flexibility to
 * change the configuration without having to change the code
 */
const config = {
  application: {
    host: process.env.APPLICATION_HOST,
    port: process.env.APPLICATION_PORT,
    publicUrl: process.env.APPLICATION_PUBLIC_URL,
    authenticationName: 'festival-ticket__api',
  },
  postgres: {
    host: process.env.PGHOST,
    database: process.env.PGDATABASE,
    port: process.env.PGPORT,
    user: process.env.PGUSER,
    password: process.env.PGPASSWORD,
  },
  jwtTokenize: {
    accessTokenKey: process.env.ACCESS_TOKEN_KEY,
    refreshTokenKey: process.env.REFRESH_TOKEN_KEY,
    accessTokenAges: process.env.ACCESS_TOKEN_AGES,
  },
};

module.exports = config;
