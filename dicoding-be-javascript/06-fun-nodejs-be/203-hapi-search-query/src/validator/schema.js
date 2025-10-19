const Joi = require("joi");

const ProductQuerySchema = Joi.object({
  name: Joi.string().empty('')
});

module.exports = ProductQuerySchema;