require('dotenv').config();

const Hapi = require('@hapi/hapi');
const {Pool} = require('pg');
const albums = require('./api/albums');
const AlbumService = require('./services/albums/AlbumService');
const {AlbumValidator} = require('./validator/albums');
const songs = require('./api/songs');
const SongService = require('./services/songs/SongService');
const {SongValidator} = require('./validator/songs');
const ClientError = require('./exceptions/ClientError');

const init = async () => {
  const pool = new Pool({
    user: process.env.PGUSER,
    host: process.env.PGHOST,
    database: process.env.PGDATABASE,
    password: process.env.PGPASSWORD,
    port: process.env.PGPORT,
  });

  pool.connect((err) => {
    if (err) throw err;
    console.log('Connected to PostgreSQL Server!');
  });

  const albumService = new AlbumService();
  const albumValidator = new AlbumValidator();
  const songService = new SongService();
  const songValidator = new SongValidator();
  const server = Hapi.server({
    host: process.env.HOST,
    port: process.env.PORT,
    routes: {
      cors: {
        origin: ['*'],
      },
    },
  });

  await server.register([
    {
      plugin: albums,
      options: {
        service: albumService,
        validator: albumValidator,
      },
    },
    {
      plugin: songs,
      options: {
        service: songService,
        validator: songValidator,
      },
    },
  ]);

  server.ext('onPreResponse', (request, h) => {
     // mendapatkan konteks response dari request
     const {response} = request;
     console.log(response);
     if (response instanceof Error) {
         // penanganan client error secara internal.
         if (response instanceof ClientError) {
             const newResponse = h.response({
                 status: 'fail',
                 message: response.message,
             });
             newResponse.code(response.statusCode);
             return newResponse;
         }

         // mempertahankan penanganan client error oleh hapi secara native, seperti 404, etc.
         if (!response.isServer) {
             return h.continue;
         }

         // penanganan server error sesuai kebutuhan
         const newResponse = h.response({
             status: 'error',
             message: 'terjadi kegagalan pada server kami',
         });
         newResponse.code(500);
         return newResponse;
     }

     // jika bukan error, lanjutkan dengan response sebelumnya (tanpa terintervensi)
     return h.continue;
  });

  await server.start();
  console.log('Server running on %s', server.info.uri);
};

init();
