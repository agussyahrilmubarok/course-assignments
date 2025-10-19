const amqp = require('amqplib');
const config = require('../../utils/config');

// this class contains all the methods needed to handle rabbitmq server
class QueueService {
  /**
   * This method is used to send message to a queue in rabbitmq server
   */
  async sendMessage(queue, message) {
    // connect to rabbitmq server
    const connection = await amqp.connect(config.rabbitMq.server);
    // create channel that will be used to send message
    const channel = await connection.createChannel();

    // assert the queue to make sure the queue is exist
    // if the queue is not exist, it will be created automatically
    await channel.assertQueue(queue, {
      durable: true,
    });

    // send message to the queue
    await channel.sendToQueue(queue, Buffer.from(message));

    // close the connection after the message is sent
    // we give it 1-second delay to make sure the message is sent before the connection is closed
    setTimeout(() => {
      connection.close();
    }, 1000);
  }
}

module.exports = QueueService;
