const products = require("./products");

class ProductsHandler {
  constructor(validator) {
    this._validator = validator;

    this.postProductHandler = this.postProductHandler.bind(this);
    this.getProductsHandler = this.getProductsHandler.bind(this);
  }

  async postProductHandler(request, h) {
    const { name, price, category } = request.payload;

    products.push({ name, price, category });

    const response = h.response({
      status: "success",
      message: "Produk berhasil dimasukkan",
      data: {
        name,
      },
    });
    response.code(201);
    return response;
  }

  async getProductsHandler(request, h) {
    const { name = "" } = request.query;

    await this._validator.validateQuery({ name });

    if (name !== "") {
      const product = products.filter((product) => product.name === name);
      const response = h.response({
        status: "success",
        message: "Produk berhasil ditampilkan",
        data: {
          product,
        },
      });
      return response;
    }

    const response = h.response({
      status: "success",
      message: "Produk berhasil ditampilkan",
      data: {
        products,
      },
    });
    return response;
  }
}

module.exports = ProductsHandler;
