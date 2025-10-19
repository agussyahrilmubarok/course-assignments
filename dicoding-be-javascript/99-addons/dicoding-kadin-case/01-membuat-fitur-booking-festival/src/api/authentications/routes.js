const routes = (handler) => [
  {
    method: 'POST',
    path: '/authentications',
    // the reason we used anonymous function here is because ...
    // we want to keep the context of this on the class handler
    handler: (request, h) => handler.postAuthenticationHandler(request, h),
  },
  {
    method: 'PUT',
    path: '/authentications',
    handler: (request, h) => handler.putAuthenticationHandler(request, h),
  },
  {
    method: 'DELETE',
    path: '/authentications',
    handler: (request, h) => handler.deleteAuthenticationHandler(request, h),
  },
];

module.exports = routes;
