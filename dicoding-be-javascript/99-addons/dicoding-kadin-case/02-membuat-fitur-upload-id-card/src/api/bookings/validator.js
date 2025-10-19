const Joi = require('joi');
const InvariantError = require('../../exceptions/InvariantError');

// Joi schema for POST /bookings
const PostBookingPayloadSchema = Joi.object({
  festivalId: Joi.string().required(),
  quantity: Joi.number().min(1).required(),
  bookingDate: Joi.date().iso().required(),
});

// this validator will be used to validate the payload from client related to bookings feature
const BookingsValidator = {
  validatePostBookingPayload(payload) {
    const validationResult = PostBookingPayloadSchema.validate(payload);

    // if the payload is not valid, then the validator will throw an error
    if (validationResult.error) {
      throw new InvariantError(validationResult.error.message);
    }

    return validationResult.value;
  },
};

module.exports = BookingsValidator;
