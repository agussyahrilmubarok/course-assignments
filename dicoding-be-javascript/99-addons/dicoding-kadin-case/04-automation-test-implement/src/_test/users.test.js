const createServer = require('../http/createServer');
const { createPool } = require('../services/postgres/pool');

describe('resources /users', () => {
  let server;
  let pool;

  beforeAll(async () => {
    server = await createServer();
    pool = await createPool();
  });

  afterAll(() => {
    // close database connection pool after all test done
    pool.end();
  });

  describe('when POST /users', () => {
    it('[Fail] Register user with Bad Payload', async () => {
      const payload = {};

      const response = await server.inject({
        method: 'POST',
        url: '/users',
        payload,
      });

      expect(response.statusCode).toEqual(400);
      expect(response.headers['content-type']).toEqual('application/json; charset=utf-8');
    });

    it('[Success] Register Email', async () => {
      const payload = {
        name: 'John Doe',
        email: `john_${+new Date()}@example.com`,
        password: 'secret',
      };

      const response = await server.inject({
        method: 'POST',
        url: '/users',
        payload,
      });

      expect(response.statusCode).toEqual(201);
    });

    it('[Fail] Register Email that already registered', async () => {
      // arrange
      const payload = {
        name: 'John Doe',
        email: 'john@example.com',
        password: 'secret',
      };

      await server.inject({
        method: 'POST',
        url: '/users',
        payload,
      });

      // action
      const response = await server.inject({
        method: 'POST',
        url: '/users',
        payload,
      });

      // assert
      expect(response.statusCode).toEqual(400);
      const responseJson = JSON.parse(response.payload);
      expect(responseJson.message).toEqual('email sudah digunakan');
    });
  });
});
