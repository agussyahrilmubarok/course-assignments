const { nanoid } = require('nanoid');
const bcrypt = require('bcrypt');
const { createPool } = require('./pool');
const AuthenticationError = require('../../exceptions/AuthenticationError');

// this class will be used to handle all the database operations related to users table
class UsersService {
  constructor() {
    // create database connection pool
    this._pool = createPool();
  }

  /**
   * this method will be used to add user to database
   */
  async persistUsers(user) {
    const id = `user-${nanoid(16)}`;
    const { name, email, password } = user;
    const hashedPassword = await bcrypt.hash(password, 10);

    const query = {
      text: 'INSERT INTO users VALUES($1, $2, $3, $4)',
      values: [id, name, email, hashedPassword],
    };

    await this._pool.query(query);

    return {
      id,
      name,
      email,
    };
  }

  /**
   * this method will be used to get user by id from database
   */
  async isEmailAvailable(email) {
    const query = {
      text: 'SELECT id FROM users WHERE email = $1',
      values: [email],
    };
    const result = await this._pool.query(query);

    return !result.rowCount > 0;
  }

  /**
   * this method will be used to verify user credential from database
   */
  async verifyUserCredential(email, password) {
    const query = {
      text: 'SELECT id, password FROM users WHERE email = $1',
      values: [email],
    };
    const result = await this._pool.query(query);

    // if user with email not found, throw error
    if (!result.rowCount) {
      throw new AuthenticationError('kredensial yang Anda berikan salah');
    }

    const { id, password: hashedPassword } = result.rows[0];
    // compare password with hashed password
    const match = await bcrypt.compare(password, hashedPassword);

    // if password not match, throw error
    // why we used same error message for both email and password?
    // to prevent attacker to know which one is wrong, email or password
    if (!match) {
      throw new AuthenticationError('kredensial yang Anda berikan salah');
    }

    return id;
  }
}

module.exports = UsersService;
