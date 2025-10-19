const { createClient } = require('redis');
const config = require('../../utils/config');

// this class contains all the methods needed to handle redis server
class CacheService {
  constructor() {
    // create redis client
    this._client = createClient({
      socket: {
        host: config.redis.server,
      },
    });

    // handle error when connecting to redis server
    this._client.on('error', (error) => {
      console.error(error);
    });
  }

  /**
   * This method is used to connect to redis server, execute the action, and disconnect the client
   * We create this method to avoid repeating the same code when
   * connecting and disconnecting to redis server
   */
  async _doAction(action) {
    await this._client.connect();
    const result = await action();
    await this._client.disconnect();
    return result;
  }

  /**
   * This method is used to set the value of a key in redis server
   * This method use _doAction method to connect and disconnect to redis server
   * to avoid repeating the same code
   */
  async set(key, value, expirationInSecond = 3600) {
    await this._doAction(async () => {
      await this._client.set(key, value, {
        EX: expirationInSecond,
      });
    });
  }

  /**
   * This method is used to get the value of a key in redis server
   * This method use _doAction method to connect and disconnect to redis server
   * to avoid repeating the same code
   */
  async get(key) {
    return this._doAction(() => this._client.get(key));
  }

  /**
   * This method is used to delete the value of a key in redis server
   * This method use _doAction method to connect and disconnect to redis server
   * to avoid repeating the same code
   */
  async delete(key) {
    return this._doAction(() => this._client.del(key));
  }
}

module.exports = CacheService;
