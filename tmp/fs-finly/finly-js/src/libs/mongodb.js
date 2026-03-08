const mongoose = require('mongoose');
const logger = require('./logger');
const { config } = require('./config');

const mongoUri = config.MONGO_URI;
const mongoTimeout = config.MONGO_TIMEOUT;

const mongooseConnect = async () => {
  try {
    logger.debug('Attempting MongoDB connection', {
      uri: mongoUri,
      timeout: mongoTimeout,
    });

    await mongoose.connect(mongoUri, {
      bufferCommands: false,
      serverSelectionTimeoutMS: mongoTimeout,
    });

    logger.info('MongoDB connected successfully', { uri: mongoUri });
  } catch (err) {
    logger.error('MongoDB connection failed', {
      message: err.message,
      stack: err.stack,
      uri: mongoUri,
    });
    throw err;
  }
};

module.exports = { mongooseConnect };
