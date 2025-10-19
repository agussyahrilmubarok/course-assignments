const ClientError = require('./ClientError');

// this error will be thrown when authentication failed
class AuthenticationError extends ClientError {
  constructor(message) {
    super(message, 401);
    this.name = 'AuthenticationError';
  }
}

module.exports = AuthenticationError;
