const AuthenticationsValidator = require('./validator');
const AuthenticationsHandler = require('./handler');
const routes = require('./routes');

// this plugin will be used to create authentication feature
const authentications = {
  name: 'authentications',
  version: '1.0.0',
  register: async (server, { authenticationsService, usersService, jwtTokenize }) => {
    const authenticationsHandler = new AuthenticationsHandler(
      authenticationsService,
      usersService,
      jwtTokenize,
      AuthenticationsValidator,
    );

    server.route(routes(authenticationsHandler));
  },
};

module.exports = authentications;
