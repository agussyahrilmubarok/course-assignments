const ClientError = require('./ClientError');

// this error will be thrown when requested resource is not found
class NotFoundError extends ClientError {
  constructor(message) {
    super(message, 404);
    this.name = 'NotFoundError';
  }
}

module.exports = NotFoundError;
