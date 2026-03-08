const ClientError = require('../../exceptions/ClientError');

class AlbumLikesHandler {
  constructor(service, albumsService) {
    this._service = service;
    this._albumsService = albumsService;

    this.postAlbumLikeHandler = this.postAlbumLikeHandler.bind(this);
    this.getAlbumLikeHandler = this.getAlbumLikeHandler.bind(this);
    this.deleteAlbumLikeHandler = this.deleteAlbumLikeHandler.bind(this);
  }

  async postAlbumLikeHandler(request, h) {
    try {
      const {id: albumId} = request.params;
      const {id: userId} = request.auth.credentials;

      await this._albumsService.getAlbumById(albumId);
      await this._service.likeAlbum(userId, albumId);

      const response = h.response({
        status: 'success',
        message: 'Menyukai album',
      });
      response.code(201);
      return response;
    } catch (error) {
      if (error instanceof ClientError) {
        const response = h.response({
          status: 'fail',
          message: error.message,
        });
        response.code(error.statusCode);
        return response;
      }

      const response = h.response({
        status: 'error',
        message: 'Maaf, terjadi kegagalan pada server kami.',
      });
      response.code(500);
      console.error(error);
      return response;
    }
  }

  async getAlbumLikeHandler(request, h) {
    try {
      const {id: albumId} = request.params;
      const {likes, isCache = 0} = await this._service.getLikeAlbum(albumId);

      const response = h.response({
        status: 'success',
        message: 'Melihat jumlah yang menyukai album',
        data: {
          likes: likes.length,
        },
      });
      response.code(200);

      if (isCache) response.header('X-Data-Source', 'cache');

      return response;
    } catch (error) {
      if (error instanceof ClientError) {
        const response = h.response({
          status: 'fail',
          message: error.message,
        });
        response.code(error.statusCode);
        return response;
      }

      const response = h.response({
        status: 'error',
        message: 'Maaf, terjadi kegagalan pada server kami.',
      });
      response.code(500);
      console.error(error);
      return response;
    }
  }

  async deleteAlbumLikeHandler(request, h) {
    try {
      const {id: albumId} = request.params;
      const {id: userId} = request.auth.credentials;

      await this._albumsService.getAlbumById(albumId);
      await this._service.unlikeAlbum(userId, albumId);

      const response = h.response({
        status: 'success',
        message: 'Batal menyukai album.',
      });

      response.code(200);
      return response;
    } catch (error) {
      if (error instanceof ClientError) {
        const response = h.response({
          status: 'fail',
          message: error.message,
        });
        response.code(error.statusCode);
        return response;
      }

      const response = h.response({
        status: 'error',
        message: 'Maaf, terjadi kegagalan pada server kami.',
      });
      response.code(500);
      console.error(error);
      return response;
    }
  }
}

module.exports = AlbumLikesHandler;
