/* eslint-disable key-spacing */
/* eslint-disable camelcase */
const mapDBToModel = ({
  id,
  title,
  year,
  performer,
  genre,
  duration,
  album_id,
}) => ({
  id,
  title,
  year,
  performer,
  genre,
  duration,
  albumId: album_id,
});

const mapDBToAlbumSongService = ({id, name, year}, song) => ({
  id,
  name,
  year,
  songs: song,
});

module.exports = {mapDBToModel, mapDBToAlbumSongService};
