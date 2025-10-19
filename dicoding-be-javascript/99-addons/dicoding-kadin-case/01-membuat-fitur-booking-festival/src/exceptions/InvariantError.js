const ClientError = require('./ClientError');

// this error will be thrown when business logic is not met
class InvariantError extends ClientError {
  constructor(message, statusCode = 400) {
    super(message, statusCode);
    this.name = 'InvariantError';
  }
}

module.exports = InvariantError;
