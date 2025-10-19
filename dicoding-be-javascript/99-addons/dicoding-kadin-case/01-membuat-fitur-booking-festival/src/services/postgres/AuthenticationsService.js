const { createPool } = require('./pool');
const AuthenticationError = require('../../exceptions/AuthenticationError');

// this class will be used to handle all the database operations related to authentications table
class AuthenticationsService {
  constructor() {
    // create database connection pool
    this._pool = createPool();
  }

  /**
   * this method will be used to add refresh token to database
   * this method will be used when user login to the application
   */
  async persistRefreshToken(refreshToken) {
    const query = {
      text: 'INSERT INTO authentications VALUES($1)',
      values: [refreshToken],
    };

    await this._pool.query(query);
  }

  /**
   * this method will be used to delete refresh token from database
   * this method will be used when user logout from the application
   */
  async deleteRefreshToken(refreshToken) {
    const query = {
      text: 'DELETE FROM authentications WHERE refresh_token = $1',
      values: [refreshToken],
    };

    await this._pool.query(query);
  }

  /**
   * this method will be used to verify refresh token from database
   * this method will be used when user want to refresh access token
   */
  async verifyRefreshToken(refreshToken) {
    const query = {
      text: 'SELECT * FROM authentications WHERE refresh_token = $1',
      values: [refreshToken],
    };

    const result = await this._pool.query(query);

    if (!result.rowCount) {
      throw new AuthenticationError('refresh token tidak valid');
    }
  }
}

module.exports = AuthenticationsService;
