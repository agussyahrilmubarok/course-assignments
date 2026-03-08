const ThreadRepository = require('../../../Domains/threads/ThreadRepository');
const ModelThread = require('../../../Domains/threads/entities/ModelThread');
const AddedThread = require('../../../Domains/threads/entities/AddedThread');
const NotFoundError = require('../../../Commons/exceptions/NotFoundError');

class ThreadRepositoryPostgres extends ThreadRepository {
    constructor(pool, idGenerator) {
        super();
        this._pool = pool;
        this._idGenerator = idGenerator;
    }

    async addThread({ title, body }, owner) {
        const id = `thread-${this._idGenerator()}`;

        const query = {
            text: 'INSERT INTO threads VALUES ($1, $2, $3, $4) RETURNING id, title, owner',
            values: [id, title, body, owner],
        };

        const { rows } = await this._pool.query(query);
        // const result = await this._pool.query(query);

        return new AddedThread(rows[0]);
    }

    async getThreadById(threadId) {
        const query = {
            text: `
            SELECT 
                threads.id, threads.title, threads.body, 
                threads.date, users.username
            FROM threads 
            LEFT JOIN users ON threads.owner = users.id
            WHERE threads.id = $1
            GROUP BY threads.id, users.username
            `,
            values: [threadId],
        };

        const { rows, rowCount } = await this._pool.query(query);

        if (!rowCount) throw new NotFoundError('thread tidak ditemukan');

        return new ModelThread({ ...rows[0] });
    }

    async verifyThread(threadId) {
        const query = {
            text: 'SELECT id FROM threads WHERE id = $1',
            values: [threadId],
        };

        const { rowCount } = await this._pool.query(query);

        if (!rowCount) throw new NotFoundError('thread tidak ditemukan');
    }
}

module.exports = ThreadRepositoryPostgres;
