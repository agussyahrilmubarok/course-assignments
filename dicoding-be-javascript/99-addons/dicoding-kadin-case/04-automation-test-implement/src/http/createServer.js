const Hapi = require('@hapi/hapi');
const Jwt = require('@hapi/jwt');
const Inert = require('@hapi/inert');
const config = require('../utils/config');
const ClientError = require('../exceptions/ClientError');
const users = require('../api/users');
const UsersService = require('../services/postgres/UsersService');
const authentications = require('../api/authentications');
const AuthenticationsService = require('../services/postgres/AuthenticationsService');
const JwtTokenManager = require('../tokenize/JwtTokenManager');
const festivals = require('../api/festivals');
const FestivalsService = require('../services/postgres/FestivalsService');
const BookingsService = require('../services/postgres/BookingsService');
const QueueService = require('../services/rabbitmq/QueueService');
const bookings = require('../api/bookings');
const LocalStorageService = require('../services/storages/LocalStorageService');
const CacheService = require('../services/redis/CacheService');

async function createServer() {
  // create services that will be used by the plugin
  const cacheService = new CacheService();
  const usersService = new UsersService();
  const authenticationsService = new AuthenticationsService();
  const festivalsService = new FestivalsService();
  const bookingsService = new BookingsService(cacheService);
  const queueService = new QueueService();
  const storageService = new LocalStorageService();

  // create HTTP server using hapi
  const server = Hapi.server({
    host: config.application.host,
    port: config.application.port,
  });

  // register external plugin
  await server.register([
    {
      plugin: Jwt.plugin,
    },
    {
      plugin: Inert,
    },
  ]);

  // define strategy for jwt token
  server.auth.strategy(config.application.authenticationName, 'jwt', {
    keys: config.jwtTokenize.accessTokenKey,
    // verify the token
    // the false value for aud, iss, and sub is to disable the verification on those properties
    verify: {
      aud: false,
      iss: false,
      sub: false,
      maxAgeSec: config.jwtTokenize.accessTokenAges,
    },
    // function that allows additional validation based on the decoded token
    // we used this function to get the user id from the decoded token,
    // but you can also use this function to validate other things
    validate: (artifacts) => ({
      isValid: true,
      credentials: {
        id: artifacts.decoded.payload.id,
      },
    }),
  });

  // register internal plugin
  await server.register([
    {
      plugin: users,
      options: {
        usersService,
        storageService,
      },
    },
    {
      plugin: authentications,
      options: {
        authenticationsService,
        usersService,
        jwtTokenize: JwtTokenManager,
      },
    },
    {
      plugin: festivals,
      options: {
        festivalsService,
      },
    },
    {
      plugin: bookings,
      options: {
        bookingsService,
        festivalsService,
        queueService,
      },
    },
  ]);

  // intercept the response before it is sent to the client
  // this is used to handle error response
  // if the response is an instance of ClientError, then we will send the error response
  // why we define this interceptor here?
  // because we don't want to catch the error one by one in each route handler
  server.ext('onPreResponse', (request, h) => {
    const { response } = request;

    // if the response is an instance of ClientError, then we will send the error response
    if (response instanceof ClientError) {
      return h.response({
        status: 'fail',
        message: response.message,
        data: {},
      }).code(response.statusCode);
    }

    // otherwise, return the response without any change
    return h.continue;
  });

  return server;
}

module.exports = createServer;
