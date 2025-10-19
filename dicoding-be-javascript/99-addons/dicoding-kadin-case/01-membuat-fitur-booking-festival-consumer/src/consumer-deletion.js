const amqp = require("amqplib");
const DatabaseService = require("./services/DatabaseService");
const MailService = require("./services/MailService");
const config = require("./utils/config");

// this is the consumer that will send deletion notification through email
async function run() {
  // create instances of the services
  const databaseService = new DatabaseService();
  const mailService = new MailService();

  // connect to RabbitMQ
  const connection = await amqp.connect(config.rabbitMq.server);
  // create channel for consuming messages
  const channel = await connection.createChannel();
  // before consuming messages, we need to make sure that the queue exists
  await channel.assertQueue("booking:delete", { durable: true });

  // consume messages from the queue
  await channel.consume("booking:delete", async (message) => {
    // decode message to get the bookingId
    const { bookingId } = JSON.parse(message.content.toString());

    // get user information for sending the email
    const user = await databaseService.getUserInformationByBookingId(bookingId);

    // set booking as deleted
    await databaseService.setBookingAsDeleted(bookingId);

    // send the email
    await mailService.sendDeletedBookingNotificationEmail(user, { bookingId });

    // acknowledge the message
    // this will remove the message from the queue
    channel.ack(message);
  });

}

run();
