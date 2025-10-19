const {nanoid} = require('nanoid');
const {Pool} = require('pg');
const NotFoundError = require('../../exceptions/NotFoundError');
const InvariantError = require('../../exceptions/InvariantError');

class AlbumService {
    constructor() {
        this._pool = new Pool();
    }

    async addAlbum({name, year}){
        const id = 'album-' + nanoid(16);
        const createdAt = new Date().toISOString();
        const updatedAt = createdAt;

        const query = {
            text: 'INSERT INTO albums VALUES($1, $2, $3, $4, $5) RETURNING id',
            values: [id, name, year, createdAt, updatedAt],
        };

        const result = await this._pool.query(query);

        if (!result.rows[0].id){
            throw new InvariantError('data album tidak dapat ditambahkan');
        }
        return result.rows[0].id;    
    }

    async getAlbumById(id) {
        const queryAlbum = {
            text: 'SELECT id, name, year FROM albums WHERE id = $1',
            values: [id],
        };

        const resultAlbum = await this._pool.query(queryAlbum);

        if (!resultAlbum.rows.length) {
            throw new NotFoundError('Album tidak ditemukan');
        }

        const querySong = {
            text: 'SELECT songs.id, songs.title, songs.performer FROM albums JOIN songs ON albums.id = songs.album_id WHERE albums.id = $1',
            values: [id],
        };

        const resultSong = await this._pool.query(querySong);
        
        return {album: resultAlbum.rows[0], songs: resultSong.rows};
    }

    async editAlbumById(id, {name, year}) {
        const updatedAt = new Date().toISOString();

        const query = {
            text: 'UPDATE albums SET name = $1, year = $2, updated_at = $3 WHERE id = $4 RETURNING id',
            values: [name, year, updatedAt, id],
        };

        const result = await this._pool.query(query);

        if (!result.rows.length){
            throw new NotFoundError('Gagal memperbaharui album. Id tidak ditemukan');
        }
    }

    async deleteAlbumById(id) {
        const query = {
            text: 'DELETE FROM albums WHERE id = $1 RETURNING id',
            values: [id],
        };

        const result = await this._pool.query(query);

        if (!result.rows.length){
            throw new NotFoundError('Album gagal dihapus. Id tidak ditemukan');
        }
    }
};

module.exports = AlbumService;
