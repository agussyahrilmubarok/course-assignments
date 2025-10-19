const Joi = require('joi');
const InvariantError = require('../../exceptions/InvariantError');

// Joi schema for POST /authentications
const PostAuthenticationSchema = Joi.object({
  email: Joi.string().email().required(),
  password: Joi.string().required(),
});

// Joi schema for PUT /authentications
const PutAuthenticationSchema = Joi.object({
  refreshToken: Joi.string().required(),
});

// Joi schema for DELETE /authentications
const DeleteAuthenticationSchema = Joi.object({
  refreshToken: Joi.string().required(),
});

// this validator will be used to validate the payload from client related to authentication feature
const AuthenticationsValidator = {
  validatePostAuthenticationPayload(payload) {
    const validationResult = PostAuthenticationSchema.validate(payload);

    // if the payload is not valid, then the validator will throw an error
    if (validationResult.error) {
      throw new InvariantError(validationResult.error.message);
    }

    // if the payload is valid, then the validator will return the payload
    return validationResult.value;
  },

  validatePutAuthenticationPayload(payload) {
    const validationResult = PutAuthenticationSchema.validate(payload);

    if (validationResult.error) {
      throw new InvariantError(validationResult.error.message);
    }

    return validationResult.value;
  },

  validateDeleteAuthenticationPayload(payload) {
    const validationResult = DeleteAuthenticationSchema.validate(payload);

    if (validationResult.error) {
      throw new InvariantError(validationResult.error.message);
    }

    return validationResult.value;
  },
};

module.exports = AuthenticationsValidator;
