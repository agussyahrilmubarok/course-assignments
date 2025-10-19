const Joi = require('joi');

const ExportPlaylistsPayloadSchema = Joi.object({
  targetEmail: Joi.string()
    .email({tlds: true})
    .required()
    .custom((value, helpers) => {
      const domain = value.split('@')[1];
      if (domain === 'test.com') {
        return helpers.error('invalid email domain');
      }
      return value;
    }),
});

module.exports = ExportPlaylistsPayloadSchema;
