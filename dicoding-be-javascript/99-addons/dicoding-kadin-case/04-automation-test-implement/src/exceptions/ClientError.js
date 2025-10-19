// this error indicates that error is caused by client
// we treat this as an abstract class, so we won't instantiate this class directly
// why should treat as abstract class?
// because ClientError is too general, we will use this class to create more specific error
class ClientError extends Error {
  constructor(message, statusCode = 400) {
    super(message);
    this.statusCode = statusCode;
    this.name = 'ClientError';
  }
}

module.exports = ClientError;
