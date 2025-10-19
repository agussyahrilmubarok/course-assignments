const { nanoid } = require('nanoid');
const { createPool } = require('./pool');

// this class will be used to handle all the database operations related to bookings table
class BookingsService {
  constructor() {
    this._pool = createPool();
  }

  /**
   *  * this method will be used to add booking to database */
  async persistBooking({
    userId, festivalId, bookingDate, quantity,
  }) {
    const bookingsId = `booking-${nanoid(16)}`;

    const query = {
      text: 'INSERT INTO bookings VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)',
      values: [bookingsId, userId, festivalId, bookingDate, quantity, null, 0, new Date(), null],
    };

    await this._pool.query(query);

    return bookingsId;
  }

  /**
   * this method will be used to confirm bookings
   */
  async confirmBooking(confirmationCode) {
    const query = {
      text: 'UPDATE bookings SET status = 1 WHERE confirmation_code = $1',
      values: [confirmationCode],
    };

    await this._pool.query(query);
  }

  /**
   * this method will be used to get user id by confirmation code
   */
  async getUserIdByConfirmationCode({ confirmationCode }) {
    const query = {
      text: 'SELECT user_id FROM bookings WHERE confirmation_code = $1',
      values: [confirmationCode],
    };

    const { rows } = await this._pool.query(query);

    if (!rows.length) {
      return null;
    }

    return rows[0].user_id;
  }

  /**
   * this method will be used to get user id by booking id
   */
  async getUserIdByBookingId(bookingId) {
    const query = {
      text: 'SELECT user_id FROM bookings WHERE id = $1',
      values: [bookingId],
    };

    const { rows } = await this._pool.query(query);

    if (!rows.length) {
      return null;
    }

    return rows[0].user_id;
  }

  /**
   * this method will be used to get booking by id
   */
  async getBookingById(bookingId) {
    const query = {
      text: 'SELECT * FROM bookings WHERE id = $1',
      values: [bookingId],
    };

    const { rows } = await this._pool.query(query);

    const [booking] = rows;

    // delete the confirmation code from the booking object
    delete booking.confirmation_code;

    return booking;
  }

  /**
   * this method will be used to softly delete booking
   */
  async softDeleteBooking(bookingId) {
    const query = {
      text: 'UPDATE bookings SET status = -1 WHERE id = $1',
      values: [bookingId],
    };

    await this._pool.query(query);
  }
}

module.exports = BookingsService;
