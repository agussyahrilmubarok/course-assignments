const CollaborationHandler = require('./handler');
const routes = require('./route');

module.exports = {
    name: 'collaborations',
    version: '1.0.0',
    register: async (server, {collaborationService, playlistService, validator}) => {
        const collaborationHandler = new CollaborationHandler(collaborationService, playlistService, validator);
        server.route(routes(collaborationHandler));
    },
};
