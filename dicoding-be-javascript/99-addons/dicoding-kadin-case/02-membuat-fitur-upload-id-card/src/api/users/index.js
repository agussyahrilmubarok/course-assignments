const UsersHandler = require('./handler');
const UsersValidator = require('./validator');
const routes = require('./routes');

// this plugin will be used to create users feature
const users = {
  name: 'users',
  version: '1.0.0',
  register: async (server, { usersService, storageService }) => {
    const usersHandler = new UsersHandler(usersService, storageService, UsersValidator);
    server.route(routes(usersHandler));
  },
};

module.exports = users;
