const Joi = require('joi');
const currentYear = new Date().getFullYear();

const albumPayloadSchema = Joi.object({
    name: Joi.string().required(),
    year: Joi.number().integer().min(1900).max(currentYear).required(),
});

module.exports = {albumPayloadSchema};
