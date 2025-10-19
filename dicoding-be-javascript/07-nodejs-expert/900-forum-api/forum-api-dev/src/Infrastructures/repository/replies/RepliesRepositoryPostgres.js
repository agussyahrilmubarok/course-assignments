const RepliesRepository = require('../../../Domains/replies/RepliesRepository');
const AddedReply = require('../../../Domains/replies/entities/AddedReply');
const ModelReply = require('../../../Domains/replies/entities/ModelReply');

const AuthorizationError = require('../../../Commons/exceptions/AuthorizationError');
const NotFoundError = require('../../../Commons/exceptions/NotFoundError');

class RepliesRepositoryPostgres extends RepliesRepository {
    constructor(pool, idGenerator) {
        super();
        this._pool = pool;
        this._idGenerator = idGenerator;
    }

    async addReplyToComment(commentId, { content }, owner) {
        const id = `reply-${this._idGenerator()}`;

        const query = {
            text: 'INSERT INTO replies VALUES($1, $2, $3, $4) RETURNING id, content, owner',
            values: [id, commentId, content, owner],
        };

        const { rows } = await this._pool.query(query);

        return new AddedReply({ ...rows[0] });
    }

    async verifyReplyAccess(replyId, userId) {
        const query = {
            text: 'SELECT 1 FROM replies WHERE id = $1 AND owner = $2',
            values: [replyId, userId],
        };

        const { rowCount } = await this._pool.query(query);

        if (!rowCount) {
            throw new AuthorizationError('anda tidak dapat mengakses resource ini');
        }
    }

    async verifyReplyIsExist(replyId) {
        const query = {
            text: 'SELECT id FROM replies WHERE id = $1',
            values: [replyId],
        };

        const result = await this._pool.query(query);

        if (!result.rowCount) {
            throw new NotFoundError('balasan tidak ditemukan');
        }
    }

    async deleteReply(replyId) {
        const query = {
            text: `
                UPDATE replies SET is_deleted = TRUE
                WHERE id = $1 RETURNING id
                `,
            values: [replyId],
        };

        const result = await this._pool.query(query);

        if (!result.rowCount) {
            throw new NotFoundError('balasan tidak ditemukan');
        }
    }

    async repliesFromComment(commentId) {
        const query = {
            text: `
            SELECT
                replies.id,
                users.username,
                replies.date,
                replies.content,
                replies.is_deleted AS "isDeleted"
            FROM replies
            LEFT JOIN users
                ON replies.owner = users.id
            WHERE comment_id = $1
            GROUP BY
                replies.id, users.username
            ORDER BY date
            `,
            // Hint Reviewer Notes
            // text: `
            //         SELECT replies.*, users.username
            //         FROM replies
            //         INNER JOIN users ON users.id = replies.owner
            //         WHERE replies.comment_id = ANY($1::text[])
            //     `,
            values: [commentId],
        };

        const result = await this._pool.query(query);

        if (!result.rowCount) {
            return [];
        }
        return result.rows.map((val) => new ModelReply(val));
    }
}

module.exports = RepliesRepositoryPostgres;
