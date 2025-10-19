const routes = (handler) => [
  {
    method: 'GET',
    path: '/festivals',
    // the reason we used anonymous function here is because ...
    // we want to keep the context of this on the class handler
    handler: (request, h) => handler.getFestivalsHandler(request, h),
  },
  {
    method: 'GET',
    path: '/festivals/{id}',
    handler: (request, h) => handler.getFestivalByIdHandler(request, h),
  },
];

module.exports = routes;
