const {playlistPayloadSchema, playlistSongPayloadSchema} = require('./schema');
const InvariantError = require('../../exceptions/InvariantError');

class PlaylistValidator {
    validatePlaylistPayload = (payload) => {
        const validationResult = playlistPayloadSchema.validate(payload);

        if (validationResult.error){
            throw new InvariantError(validationResult.error.message);
        }

        return validationResult.value;
    };

    validatePlaylistSongPayload = (payload) => {
        const validationResult = playlistSongPayloadSchema.validate(payload);

        if (validationResult.error){
            throw new InvariantError(validationResult.error.message);
        }

        return validationResult.value;
    };
};

module.exports = {PlaylistValidator};
