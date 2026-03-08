const Joi = require('joi');

const playlistPayloadSchema = Joi.object({
    name: Joi.string().required(),
});

const playlistSongPayloadSchema = Joi.object({
    songId: Joi.string().required(),
});

module.exports = {playlistPayloadSchema, playlistSongPayloadSchema};
