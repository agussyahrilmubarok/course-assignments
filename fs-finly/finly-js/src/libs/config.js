const getEnv = (key, defaultValue) => process.env[key] || defaultValue;

const config = {
  APP_NAME: getEnv('APP_NAME', 'finly'),
  APP_PORT: getEnv('APP_PORT', '8080'),
  APP_LEVEL: getEnv('APP_LEVEL', 'development'),
  MONGO_URI: getEnv('MONGODB_URI', 'mongodb://mongodb:27017/finly'),
  MONGO_TIMEOUT: parseInt(getEnv('MONGODB_TIMEOUT', '20000'), 10),
  AUTH_SECRET: getEnv('AUTH_SECRET', 'authenticate-secret-key'),
};

module.exports = { config };
