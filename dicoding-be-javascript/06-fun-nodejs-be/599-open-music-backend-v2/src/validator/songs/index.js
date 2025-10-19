const {songPayloadSchema} = require('./schema');
const InvariantError = require('../../exceptions/InvariantError');

class SongValidator {
    validateSongPayload = (payload) => {
        const validationResult = songPayloadSchema.validate(payload);

        if (validationResult.error){
            throw new InvariantError(validationResult.error.message);
        }

        return validationResult.value;
    };
};

module.exports = {SongValidator};
