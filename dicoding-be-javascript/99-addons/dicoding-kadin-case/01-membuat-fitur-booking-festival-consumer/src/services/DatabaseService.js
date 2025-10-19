const { Pool } = require('pg');
const config = require("../utils/config");

// this is the service that will handle all database operations
class DatabaseService {
  constructor() {
    // create a new pool of connections
    this._pool = new Pool({
      host: config.postgres.host,
      database: config.postgres.database,
      port: config.postgres.port,
      user: config.postgres.user,
      password: config.postgres.password,
    });
  }

  // this method will be used to get the user information
  async getTicketPriceByFestivalId(festivalId) {
    const query = {
      text: `SELECT price FROM festivals WHERE id = $1`,
      values: [festivalId],
    };

    const result = await this._pool.query(query);

    return result.rows[0].price;
  }

  // this method will be used to get the booking information
  async getBookingInformationByBookingId(bookingId) {
    const query = {
      text: `SELECT bookings.quantity, bookings.festival_id
			FROM bookings
			WHERE bookings.id = $1`,
      values: [bookingId],
    };

    const result = await this._pool.query(query);

    return result.rows[0];
  }

  // this method will be used to put the total price to the database
  async putTotalPriceToBooking(bookingId, price) {
    const query = {
      text: `UPDATE bookings SET total_price = $1 WHERE id = $2`,
      values: [price, bookingId],
    };

    await this._pool.query(query);
  }

  // this method will be used to put the confirmation code to the database
  async putConfirmationCode(bookingId, confirmationCode) {
    const query = {
      text: `UPDATE bookings SET confirmation_code = $1 WHERE id = $2`,
      values: [confirmationCode, bookingId],
    };

    await this._pool.query(query);
  }

  // this method will be used to get the user information
  async getUserInformationByBookingId(bookingId) {
    const query = {
      text: `SELECT users.email, users.name
			FROM bookings
			LEFT JOIN users ON users.id = bookings.user_id
			WHERE bookings.id = $1`,
      values: [bookingId],
    };

    const result = await this._pool.query(query);

    return result.rows[0];
  }

  // this method will be used to set the booking as deleted
  async setBookingAsDeleted(bookingId) {
    const query = {
      text: `UPDATE bookings SET status = -1 WHERE id = $1`,
      values: [bookingId],
    };

    await this._pool.query(query);
  }
}

module.exports = DatabaseService;
