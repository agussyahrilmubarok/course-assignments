// this handler will be used to handle requests related to festivals feature
class FestivalsHandler {
  constructor(festivalsService) {
    this._festivalsService = festivalsService;
  }

  // this handler will be used to handle GET /festivals
  async getFestivalsHandler() {
    const festivals = await this._festivalsService.getFestivals();

    return {
      status: 'success',
      message: 'list of festivals',
      data: {
        festivals,
      },
    };
  }

  // this handler will be used to handle GET /festivals/{id}
  async getFestivalByIdHandler(request) {
    const { id } = request.params;

    const festival = await this._festivalsService.getFestival(id);

    return {
      status: 'success',
      message: `festival of ${id}`,
      data: {
        festival,
      },
    };
  }
}

module.exports = FestivalsHandler;
