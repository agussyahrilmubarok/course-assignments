const { nanoid } = require('nanoid');
const { createPool } = require('./pool');

// this class will be used to handle all the database operations related to bookings table
class BookingsService {
  constructor(cacheService) {
    this._pool = createPool();
    // inject cache service
    this._cacheService = cacheService;
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

    // delete the cache
    const bookingId = await this.getBookingIdByConfirmationCode(confirmationCode);

    // we need to delete the cache because the booking status has changed
    await this._cacheService.del(`booking:${bookingId}`);
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
    // check the cache first
    const cacheKey = `booking:${bookingId}`;

    // get the cache
    const cache = await this._cacheService.get(cacheKey);

    // if cache is exist, then return the cache
    if (cache) {
      return {
        value: JSON.parse(cache),
        from: 'cache', // this will be used to check where the data is coming from
      };
    }

    const query = {
      text: 'SELECT * FROM bookings WHERE id = $1',
      values: [bookingId],
    };

    const { rows } = await this._pool.query(query);

    const [booking] = rows;

    // delete the confirmation code from the booking object
    delete booking.confirmation_code;

    // store the booking object in cache
    await this._cacheService.set(cacheKey, JSON.stringify(booking));

    return {
      value: booking,
      from: 'database', // this will be used to check where the data is coming from
    };
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

    // we need to delete the cache because the booking status has changed
    await this._cacheService.delete(`booking:${bookingId}`);
  }

  /**
   * this method will be used to get booking id by confirmation code
   */
  async getBookingIdByConfirmationCode(confirmationCode) {
    const query = {
      text: 'SELECT id FROM bookings WHERE confirmation_code = $1',
      values: [confirmationCode],
    };

    const { rows } = await this._pool.query(query);

    if (!rows.length) {
      return null;
    }

    return rows[0].id;
  }
}

module.exports = BookingsService;
