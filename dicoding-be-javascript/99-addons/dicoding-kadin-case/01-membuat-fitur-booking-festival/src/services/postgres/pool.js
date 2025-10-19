const { Pool } = require('pg');
const config = require('../../utils/config');

// this variable will be used to store the pool instance
let pool = null;

// this function will be used to create a pool instance
function createPool() {
  // if pool instance is already created, return the instance
  if (pool) {
    return pool;
  }

  // create a new pool instance when it is not created yet and store it in pool variable
  // this technique is called singleton
  pool = new Pool({
    host: config.postgres.host,
    port: config.postgres.port,
    user: config.postgres.user,
    database: config.postgres.database,
    password: config.postgres.password,
  });

  return pool;
}

module.exports = { createPool };
