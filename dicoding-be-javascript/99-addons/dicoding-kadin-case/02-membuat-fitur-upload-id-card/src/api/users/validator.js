const Joi = require('joi');
const InvariantError = require('../../exceptions/InvariantError');

// Joi schema for POST /users
const PostUserPayloadSchema = Joi.object({
  name: Joi.string().required(),
  email: Joi.string().email().required(),
  password: Joi.string().required(),
});

const PatchUserPayloadSchema = Joi.object({
  images: Joi.object({
    pipe: Joi.func().required(),
    hapi: Joi.object({
      filename: Joi.string().required(),
      headers: Joi.object({
        'content-type': Joi.string().valid('image/jpeg', 'image/png', 'image/webp', 'image/jpg').required(),
      }).unknown(),
    }),
  }).unknown(),
});

const UsersValidator = {
  validatePostUserPayload(payload) {
    const validationResult = PostUserPayloadSchema.validate(payload);

    if (validationResult.error) {
      throw new InvariantError(validationResult.error.message);
    }

    return validationResult.value;
  },

  validatePatchUserPayload(payload) {
    const validationResult = PatchUserPayloadSchema.validate(payload);

    if (validationResult.error) {
      throw new InvariantError(validationResult.error.message);
    }

    return validationResult.value;
  },
};

module.exports = UsersValidator;
