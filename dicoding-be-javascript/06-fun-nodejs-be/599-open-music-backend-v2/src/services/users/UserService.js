const {Pool} = require('pg');
const bcrypt = require('bcrypt');
const InvariantError = require('../../exceptions/InvariantError');
const AuthenticationError = require('../../exceptions/AuthenticationError');
const {nanoid} = require('nanoid');
const NotFoundError = require('../../exceptions/NotFoundError');

class UserService {
    constructor() {
        this._pool = new Pool();
    }

    async verifyUserCredential({username, password}) {
        const query = {
            text: 'SELECT id, username, fullname, password FROM users WHERE username LIKE $1',
            values: [`%${username}`],
        };

        const {rows, rowCount} = await this._pool.query(query);

        if (!rowCount) {
            throw new AuthenticationError('Kredensial yang Anda berikan salah');
        }

        const {id, password: hashedPassword} = rows[0];

        const match = await bcrypt.compare(password, hashedPassword);

        if (!match) {
            throw new AuthenticationError('Kredensial yang Anda berikan salah');
        }

        return id;
    }
    
    async verifyUsername(username) {
        const query = {
            text: 'SELECT username from users WHERE username = $1',
            values: [username],
        };

        const {rowCount} = await this._pool.query(query);

        if (rowCount > 0) {
            throw new InvariantError('Gagal menambahkan user. Username sudah digunakan.');
        }
    }

    async addUser({username, password, fullname}) {

        await this.verifyUsername(username);

        const id = `user-${nanoid(16)}`;
        const hashedPassword = await bcrypt.hash(password, 10);
        const createdAt = new Date().toISOString();

        const query = {
            text: 'INSERT INTO users VALUES($1, $2, $3, $4, $5, $6) RETURNING id',
            values: [id, username, fullname, hashedPassword, createdAt, createdAt],
        };

        const {rows} = await this._pool.query(query);

        if (!rows[0].id) {
            throw new InvariantError('data user tidak dapat ditambahkan');
        }

        return rows[0].id;
    }

    async verifyUserInDatabase(userId) {
        const query = {
            text: 'SELECT id FROM users WHERE id = $1',
            values: [userId],
        };

        const {rowCount} = await this._pool.query(query);

        if (!rowCount) {
            throw new NotFoundError('User tidak ditemukan');
        }
    }
}

module.exports = UserService;
