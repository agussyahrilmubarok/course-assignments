const ProductQuerySchema = require("./schema")

const ProductValidator = {
  validateQuery: (payload) => {
    const validationResult = ProductQuerySchema.validate(payload);
    if (validationResult.error) {
      throw new Error(validationResult.error.message);
    }
    return validationResult;
  },
};

module.exports = ProductValidator;