const SongHandler = require('./handler');
const routes = require('./route');

module.exports = {
    name: 'songs',
    version: '1.0.0',
    register: async (server, {service, validator}) => {
        const songHandler = new SongHandler(service, validator);
        server.route(routes(songHandler));
    },
};
