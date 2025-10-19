const Joi = require('joi');
const InvariantError = require('../../exceptions/InvariantError');

// Joi schema for POST /users
const PostUserPayloadSchema = Joi.object({
  name: Joi.string().required(),
  email: Joi.string().email().required(),
  password: Joi.string().required(),
});

const UsersValidator = {
  validatePostUserPayload(payload) {
    const validationResult = PostUserPayloadSchema.validate(payload);

    if (validationResult.error) {
      throw new InvariantError(validationResult.error.message);
    }

    return validationResult.value;
  },
};

module.exports = UsersValidator;
