const winston = require('winston');
const DailyRotateFile = require('winston-daily-rotate-file');
const { config } = require('./config');

const serviceName = config.APP_NAME;
const isProduction = config.APP_LEVEL == 'production' ? true : false;

const dailyRotateTransport = new DailyRotateFile({
  dirname: './logs',
  filename: 'finly-%DATE%.log',
  datePattern: 'YYYY-MM-DD',
  zippedArchive: true,
  maxSize: '20m',
  maxFiles: '14d',
});

const log = winston.createLogger({
  level: isProduction ? 'warn' : 'debug',
  defaultMeta: {
    service: serviceName,
  },
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.json()
  ),
  transports: [
    new winston.transports.Console({
      format: isProduction
        ? winston.format.combine(
            winston.format.timestamp(),
            winston.format.simple()
          )
        : winston.format.combine(
            winston.format.colorize(),
            winston.format.simple()
          ),
    }),
    dailyRotateTransport,
  ],
});

module.exports = log;
