class Listener {
  constructor(PlaylistService, MailSender) {
    this._playlistService = PlaylistService;
    this._mailSender = MailSender;

    this.listen = this.listen.bind(this);
  }

  async listen(message) {
    try {
      const {playlistId, targetEmail} = JSON.parse(
        message.content.toString(),
      );

      const data = await this._playlistService.getPlaylistSongs(
        playlistId,
      );

      const result = await this._mailSender.sendEmail(
        targetEmail,
        JSON.stringify(data),
      );

      console.log(result);
    } catch (error) {
      console.error(error);
    }
  }
}

module.exports = Listener;
