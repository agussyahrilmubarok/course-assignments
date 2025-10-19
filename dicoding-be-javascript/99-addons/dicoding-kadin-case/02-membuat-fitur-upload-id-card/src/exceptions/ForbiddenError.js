const ClientError = require('./ClientError');

// this error will be thrown when user is not allowed to access a resource
class ForbiddenError extends ClientError {
  constructor(message) {
    super(message, 403);
    this.name = 'ForbiddenError';
  }
}

module.exports = ForbiddenError;
