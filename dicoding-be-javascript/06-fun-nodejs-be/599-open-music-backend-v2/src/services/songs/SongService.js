const {nanoid} = require('nanoid');
const {Pool} = require('pg');
const NotFoundError = require('../../exceptions/NotFoundError');
const InvariantError = require('../../exceptions/InvariantError');
const {mapDBToModel} = require('../../utils');

class SongService {
    constructor() {
        this._pool = new Pool();
    }

    async verifySongInDatabase(songId) {
        const query = {
            text: 'SELECT id FROM songs WHERE id = $1',
            values: [songId],
        };

        const {rowCount} = await this._pool.query(query);

        return rowCount;
    }

    async addSong({title, year, performer, genre, duration, albumId}){
        const id = `song-${nanoid(16)}`;
        const createdAt = new Date().toISOString();

        const query = {
            text: 'INSERT INTO songs VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id',
            values: [id, title, year, performer, genre, duration, albumId, createdAt, createdAt],
        };

        const {rows} = await this._pool.query(query);

        if (!rows[0].id){
            throw new InvariantError('data songs tidak dapat ditambahkan');
        }
        return rows[0].id;    
    }

    async getSongs(requestParam){

        const {title, performer} = requestParam;

        if ((title != undefined) && (performer != undefined)){
            const query = {
                text: 'SELECT id, title, performer FROM songs WHERE title ILIKE $1 AND performer ILIKE $2',
                values: [`%${title}%`, `%${performer}%`],
            };

            const {rows} = await this._pool.query(query);
            return rows;
        }

        if (title) {
            const query = {
                text: 'SELECT id, title, performer FROM songs WHERE title ILIKE $1',
                values: [`%${title}%`],
            };

            const {rows} = await this._pool.query(query);
            return rows;
        }

        if (performer) {
            const query = {
                text: 'SELECT id, title, performer FROM songs WHERE performer ILIKE $1',
                values: [`%${performer}%`],
            };

            const {rows} = await this._pool.query(query);
            return rows;
        }

        const query = {
            text: 'SELECT id, title, performer FROM songs',
        };

        const {rows} = await this._pool.query(query);
        return rows;
    }

    async getSongById(id) {
        const query = {
            text: 'SELECT id, title, year, performer, genre, duration, album_id FROM songs WHERE id = $1',
            values: [id],
        };

        const {rows, rowCount} = await this._pool.query(query);

        if (!rowCount){
            throw new NotFoundError('Song tidak ditemukan');
        }

        return rows.map(mapDBToModel)[0];
    }

    async editSongById(id, {title, year, performer, genre, duration, albumId}) {
        const updatedAt = new Date().toISOString();

        const query = {
            text: 'UPDATE songs SET title = $1, year = $2, performer = $3, genre = $4, duration = $5, album_id = $6, updated_at = $7 WHERE id = $8 RETURNING id',
            values: [title, year, performer, genre, duration, albumId, updatedAt, id],
        };

        const {rowCount} = await this._pool.query(query);

        if (!rowCount){
            throw new NotFoundError('Gagal memperbaharui songs. Id tidak ditemukan');
        }
    }

    async deleteSongById(id) {
        const query = {
            text: 'DELETE FROM songs WHERE id = $1 RETURNING id',
            values: [id],
        };

        const {rowCount} = await this._pool.query(query);

        if (!rowCount){
            throw new NotFoundError('Song gagal dihapus. Id tidak ditemukan');
        }
    }
};

module.exports = SongService;
