const config = require('../../utils/config');

const routes = (handler) => [
  {
    method: 'POST',
    path: '/bookings',
    // the reason we used anonymous function here is because ...
    // we want to keep the context of this on the class handler
    handler: (request, h) => handler.postBookingHandler(request, h),
    options: {
      // set the authentication strategy to 'festival-ticket__api'
      auth: config.application.authenticationName,
    },
  },
  {
    method: 'GET',
    path: '/bookings/confirms/{confirmationCode}',
    handler: (request, h) => handler.getBookingConfirmHandler(request, h),
    options: {
      auth: config.application.authenticationName,
    },
  },
  {
    method: 'GET',
    path: '/bookings/{id}',
    handler: (request, h) => handler.getBookingByIdHandler(request, h),
    options: {
      auth: config.application.authenticationName,
    },
  },
  {
    method: 'DELETE',
    path: '/bookings/{id}',
    handler: (request, h) => handler.deleteBookingByIdHandler(request, h),
    options: {
      auth: config.application.authenticationName,
    },
  },
];

module.exports = routes;
