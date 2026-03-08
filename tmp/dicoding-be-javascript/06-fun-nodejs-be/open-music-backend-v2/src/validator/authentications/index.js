const {
    PostAuthenticationPayloadSchema,
    PutAuthenticationPayloadSchema,
    DeleteAuthenticationPayloadSchema,
} = require('./schema');
const InvariantError = require('../../exceptions/InvariantError');

class AuthenticationValidator {
    validatePostAuthenticationPayload = (payload) => {
        const validationResult = PostAuthenticationPayloadSchema.validate(payload);

        if (validationResult.error) {
            throw new InvariantError(validationResult.error.message);
        }

        return validationResult.value;
    };

    validatePutAuthenticationPayload = (payload) => {
        const validationResult = PutAuthenticationPayloadSchema.validate(payload);

        if (validationResult.error) {
            throw new InvariantError(validationResult.error.message);
        }

        return validationResult.value;
    };

    validateDeleteAuthenticationPayload = (payload) => {
        const validationResult = DeleteAuthenticationPayloadSchema.validate(payload);

        if (validationResult.error) {
            throw new InvariantError(validationResult.error.message);
        }

        return validationResult.value;
    };
};

module.exports = {AuthenticationValidator};
