const FestivalsHandler = require('./handler');
const routes = require('./routes');

// this plugin will be used to create festivals feature
const festivals = {
  name: 'festivals',
  version: '1.0.0',
  register: async (server, { festivalsService }) => {
    const handler = new FestivalsHandler(festivalsService);
    server.route(routes(handler));
  },
};

module.exports = festivals;
