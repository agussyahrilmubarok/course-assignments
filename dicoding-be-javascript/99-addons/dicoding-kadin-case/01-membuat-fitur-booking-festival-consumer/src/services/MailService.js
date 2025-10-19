const nodemailer = require("nodemailer");
const config = require("../utils/config");

// this is the service that will handle all mail operations
class MailService {
  constructor() {
    // create a new transporter
    this._transporter = nodemailer.createTransport({
      host: config.mail.host,
      port: config.mail.port,
      auth: {
        user: config.mail.user,
        pass: config.mail.password
      }
    });
  }

  // this method will be used to send the booking verification email
  async sendBookingVerificationEmail(userInformation, bookingInformation) {
    const { email, name } = userInformation;
    const { bookingId, totalPrice, confirmationCode } = bookingInformation;

    const confirmationUrl = `http://localhost:${config.application.port}/bookings/confirms/${confirmationCode}`;

    const message = {
      from: 'Festival Ticket App',
      to: email,
      subject: `Your ${bookingId} confirmation`,
      text: `Hi ${name},\n\nYour booking has been confirmed. Please pay ${totalPrice}€ to the following bank account: 1234-5678-9012-3456. Your confirmation code is ${confirmationCode}. You can confirm your booking by clicking on the following link: ${confirmationUrl}`,
    }

    return this._transporter.sendMail(message);
  }

  // this method will be used to send the booking confirmation email
  async sendDeletedBookingNotificationEmail(userInformation, bookingInformation) {
    const { email, name } = userInformation;
    const { bookingId } = bookingInformation;

    const message = {
      from: 'Festival Ticket App',
      to: email,
      subject: `Your ${bookingId} has been deleted`,
      text: `Hi ${name},\n\nYour booking has been deleted.`
    }

    return this._transporter.sendMail(message);
  }
}

module.exports = MailService;
