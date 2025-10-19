require('dotenv').config();
const amqp = require('amqplib');
const PlaylistService = require('./PlaylistService');
const MailSender = require('./MailSender');
const Listener = require('./Listener');

const init = async () => {
  const mailSender = new MailSender();
  const playlistService = new PlaylistService();
  const listener = new Listener(playlistService, mailSender);

  const connection = await amqp.connect(process.env.RABBITMQ_SERVER);
  console.log('Terhubung ke server RabbitMQ');

  const channel = await connection.createChannel();

  await channel.assertQueue('export:playlists', {
    durable: true,
  });

  channel.consume('export:playlists', listener.listen, {noAck: true});
};

init();
