const Joi = require('joi');

const PlaylistPayloadScheme = Joi.object({
  name: Joi.string().required(),
});

module.exports = {PlaylistPayloadScheme};
