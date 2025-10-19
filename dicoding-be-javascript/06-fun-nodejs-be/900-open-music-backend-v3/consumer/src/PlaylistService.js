require('dotenv').config();
const {Pool} = require('pg');

class PlaylistService {
  constructor() {
    this._pool = new Pool({
      host: process.env.PGHOST,
      port: process.env.PGPORT,
      database: process.env.PGDATABASE,
      user: process.env.PGUSER,
      password: process.env.PGPASSWORD,
    });
  }

  async getPlaylistSongs(playlistId) {
    const queryPlaylist = {
      text: `SELECT playlists.id, playlists.name 
      FROM playlists
      WHERE playlists.id = $1`,
      values: [playlistId],
    };
    const resultPlaylist = await this._pool.query(queryPlaylist);
    const playlistById = resultPlaylist.rows;

    const querySongs = {
      text: `SELECT songs.id, songs.title, songs.performer
      FROM playlists
      INNER JOIN playlist_songs ON playlist_songs.playlist_id = playlists.id
      INNER JOIN songs ON songs.id = playlist_songs.song_id
      WHERE playlists.id = $1`,
      values: [playlistId],
    };
    const resultSongs = await this._pool.query(querySongs);
    const songs = resultSongs.rows;

    const playlist = {
      playlist: {
        id: playlistById[0].id,
        name: playlistById[0].name,
        songs,
      },
    };

    console.log(playlist);

    return playlist;
  }
}

module.exports = PlaylistService;
