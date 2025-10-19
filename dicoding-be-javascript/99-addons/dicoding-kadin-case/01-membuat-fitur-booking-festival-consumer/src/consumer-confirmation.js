const amqp = require("amqplib");
const DatabaseService = require("./services/DatabaseService");
const config = require("./utils/config");
const MailService = require("./services/MailService");
const { nanoid } = require("nanoid");

// this is the consumer that will send the confirmation email
async function run() {
  // create instances of the services
  const databaseService = new DatabaseService();
  const mailService = new MailService();

  // create connection to RabbitMQ
  const connection = await amqp.connect(config.rabbitMq.server);
  // create channel for consuming messages
  const channel = await connection.createChannel();
  // before consuming messages, we need to make sure that the queue exists
  // the queue is idempotent, so it will only be created if it doesn't exist already
  await channel.assertQueue("booking:send_confirmation", { durable: true });

  // consume messages from the queue
  await channel.consume("booking:send_confirmation", async (message) => {
    // decode message to get the bookingId
    const { bookingId } = JSON.parse(message.content.toString());

    // calculating total price and putting it to the database
    const { quantity, festival_id } = await databaseService.getBookingInformationByBookingId(bookingId);
    const price = await databaseService.getTicketPriceByFestivalId(festival_id);
    const totalPrice = quantity * price;
    await databaseService.putTotalPriceToBooking(bookingId, totalPrice);

    // generate and put confirmation code to the database
    const confirmationCode = nanoid(16);
    await databaseService.putConfirmationCode(bookingId, confirmationCode);

    // get user information for sending the email
    const user = await databaseService.getUserInformationByBookingId(bookingId);

    // send the email
    await mailService.sendBookingVerificationEmail(
      user,
      { bookingId, totalPrice, confirmationCode },);

    // acknowledge the message
    // this will remove the message from the queue
    channel.ack(message);
  });
}

run();
