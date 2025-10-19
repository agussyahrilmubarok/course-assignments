const routes = (handler) => [
  {
    method: "POST",
    path: "/products",
    handler: handler.postProductHandler,
  },
  {
    method: "GET",
    path: "/products",
    handler: handler.getProductsHandler,
  },
];

module.exports = routes;
