// this function will be used to create users feature routes
const routes = (handler) => [
  {
    method: 'POST',
    path: '/users',
    // the reason we used anonymous function here is because ...
    // we want to keep the context of this on the class handler
    handler: (request, h) => handler.postUserHandler(request, h),
  },
];

module.exports = routes;
