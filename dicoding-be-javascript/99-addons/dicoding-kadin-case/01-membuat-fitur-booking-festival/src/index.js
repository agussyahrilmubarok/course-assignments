const createServer = require('./http/createServer');

(async () => {
  // create Hapi server instance
  const server = await createServer();

  // start Hapi server
  await server.start();

  console.log(`Server running at: ${server.info.uri}`);
})();
