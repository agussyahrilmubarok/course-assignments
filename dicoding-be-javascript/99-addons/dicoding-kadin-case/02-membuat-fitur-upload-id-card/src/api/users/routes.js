const config = require('../../utils/config');

// this function will be used to create users feature routes
const routes = (handler) => [
  {
    method: 'POST',
    path: '/users',
    // the reason we used anonymous function here is because ...
    // we want to keep the context of this on the class handler
    handler: (request, h) => handler.postUserHandler(request, h),
  },
  {
    method: 'PATCH',
    path: '/users/me/id-cards',
    handler: (request, h) => handler.patchUsersMeIdCardHandler(request, h),
    options: {
      auth: config.application.authenticationName,
      // set the payload configuration to allow multipart
      payload: {
        allow: 'multipart/form-data',
        multipart: true,
        // we used stream output because more efficient in memory usage
        output: 'stream',
      },
    },
  },
  {
    method: 'GET',
    path: '/users/me/id-cards',
    handler: (request, h) => handler.getUsersMeIdCardHandler(request, h),
    options: {
      auth: config.application.authenticationName,
    },
  },
];

module.exports = routes;
