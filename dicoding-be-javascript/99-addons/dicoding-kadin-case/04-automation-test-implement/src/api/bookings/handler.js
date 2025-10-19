const ForbiddenError = require('../../exceptions/ForbiddenError');
const InvariantError = require('../../exceptions/InvariantError');

// this handler will be used to handle requests related to bookings feature
class BookingsHandler {
  constructor(bookingsService, festivalsService, queueService, validator) {
    this._bookingsService = bookingsService;
    this._festivalsService = festivalsService;
    this._queueService = queueService;
    this._validator = validator;
  }

  // this handler will be used to handle POST /bookings
  async postBookingHandler(request, h) {
    // get user id from request.auth.credentials
    const { id: userId } = request.auth.credentials;

    // validate the payload
    // if the payload is not valid, then the validator will throw an error
    // if the payload is valid, then the validator will return the payload
    const {
      festivalId,
      quantity, bookingDate,
    } = this._validator.validatePostBookingPayload(request.payload);

    // check if the festival is available
    const isFestivalAvailable = await this._festivalsService.isFestivalAvailable(festivalId);

    // if the festival is not available, then throw an error
    if (!isFestivalAvailable) {
      throw new InvariantError('festival tidak valid');
    }

    // persist booking to database and get booking id for sending confirmation email to queue
    const bookingId = await this._bookingsService.persistBooking({
      userId,
      festivalId,
      quantity,
      bookingDate,
    });

    // payload for sending confirmation email to queue
    const queueMessage = JSON.stringify({ bookingId });

    // send confirmation email to queue
    await this._queueService.sendMessage('booking:send_confirmation', queueMessage);

    // return response with booking id
    return h.response({
      status: 'success',
      message: 'konfirmasi booking akan dikirimkan melalui email',
      data: {
        bookingId,
      },
    }).code(201);
  }

  // this handler will be used to handle GET /bookings/confirm/{confirmationCode}
  async getBookingConfirmHandler(request) {
    const { id: userId } = request.auth.credentials;
    const { confirmationCode } = request.params;

    // get booking owner by confirmation code
    const owner = await this._bookingsService.getUserIdByConfirmationCode(confirmationCode);

    // if the booking owner is not the same as the user id, then throw an error
    if (owner !== userId) {
      throw new ForbiddenError('anda tidak berhak mengakses resource ini');
    }

    // confirm booking to database
    const booking = await this._bookingsService.confirmBooking(confirmationCode);

    // return response with booking data
    return {
      status: 'success',
      data: {
        booking,
      },
    };
  }

  // this handler will be used to handle GET /bookings
  async getBookingByIdHandler(request, h) {
    const { id: userId } = request.auth.credentials;
    const { id: bookingId } = request.params;

    // get booking owner by booking id
    const owner = await this._bookingsService.getUserIdByBookingId(bookingId);

    // if the booking owner is not the same as the user id, then throw an error
    if (owner !== userId) {
      throw new ForbiddenError('anda tidak berhak mengakses resource ini');
    }

    // get booking data by booking id
    // value is the booking data
    // from is the data source, either 'cache' or 'database'
    const { from, value } = await this._bookingsService.getBookingById(bookingId);

    // create response with booking data
    const response = h.response({
      status: 'success',
      message: 'booking berhasil ditemukan',
      data: {
        booking: value,
      },
    });

    // set header X-Data-Source to indicate the data source
    response.header('X-Data-Source', from);

    return response;
  }

  // this handler will be used to handle DELETE /bookings/{id}
  async deleteBookingByIdHandler(request) {
    const { id: userId } = request.auth.credentials;
    const { id: bookingId } = request.params;

    // get booking owner by booking id
    const owner = await this._bookingsService.getUserIdByBookingId(bookingId);

    // if the booking owner is not the same as the user id, then throw an error
    if (owner !== userId) {
      throw new ForbiddenError('anda tidak berhak mengakses resource ini');
    }

    // soft delete booking by booking id
    await this._bookingsService.softDeleteBooking(bookingId);

    // send message to queue to delete booking
    await this._queueService.sendMessage('booking:delete', JSON.stringify({ bookingId }));

    return {
      status: 'success',
      message: 'booking berhasil dihapus',
      data: {},
    };
  }
}

module.exports = BookingsHandler;
