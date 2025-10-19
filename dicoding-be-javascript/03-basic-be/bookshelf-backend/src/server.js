const hapi = require('@hapi/hapi');
const routes = require('./book.routes');

const init = async () => {
  const server = hapi.server({
    port: 9000,
    host: 'localhost',
    routes: {
      cors: {
        origin: ['*'],
      },
    },
  });
  server.route(routes);
  await server.start();
  console.log(`Server run at ${server.info.uri}\n`);
};
init();
