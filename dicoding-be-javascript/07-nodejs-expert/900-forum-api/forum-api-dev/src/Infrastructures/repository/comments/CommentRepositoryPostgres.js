const CommentRepository = require('../../../Domains/comments/CommentRepository');
const AddedComment = require('../../../Domains/comments/entities/AddedComment');
const ModelComment = require('../../../Domains/comments/entities/ModelComment');

const AuthorizationError = require('../../../Commons/exceptions/AuthorizationError');
const NotFoundError = require('../../../Commons/exceptions/NotFoundError');

class CommentRepositoryPostgres extends CommentRepository {
    constructor(pool, idGenerator) {
        super();
        this._pool = pool;
        this._idGenerator = idGenerator;
    }

    async addCommentToThread(threadId, { content }, owner) {
        const id = `comment-${this._idGenerator()}`;

        const query = {
            text: `INSERT INTO comments VALUES ($1, $2, $3, $4)
            RETURNING id, content, owner`,
            values: [id, threadId, content, owner],
        };

        const { rows } = await this._pool.query(query);

        return new AddedComment({ ...rows[0] });
    }

    async verifyCommentAccess(commentId, userId) {
        const query = {
            text: 'SELECT owner FROM comments WHERE id = $1 AND owner = $2',
            values: [commentId, userId],
        };

        const result = await this._pool.query(query);

        if (!result.rowCount) {
            throw new AuthorizationError('anda tidak dapat mengakses resource ini');
        }
    }

    async verifyCommentIsExist(commentId, threadId) {
        const query = {
            text: 'SELECT 1 FROM comments WHERE id = $1 AND thread_id = $2',
            values: [commentId, threadId],
        };

        const result = await this._pool.query(query);

        if (!result.rowCount) {
            throw new NotFoundError('komentar tidak ditemukan');
        }
    }

    async deleteCommentById(commentId) {
        const query = {
            text: `
                UPDATE comments SET is_deleted = TRUE
                WHERE id = $1 RETURNING id
                `,
            values: [commentId],
        };

        await this._pool.query(query);
    }

    async commentsFromThread(threadId) {
        const query = {
            text: `
            SELECT 
                comments.id, users.username, 
                comments.date, comments.content, 
                comments.is_deleted AS "isDeleted"
            FROM comments 
            LEFT JOIN users ON comments.owner = users.id
            WHERE thread_id = $1
            GROUP BY comments.id, users.username
            ORDER BY date
            `,
            values: [threadId],
        };

        const { rows, rowCount } = await this._pool.query(query);

        if (!rowCount) return [];

        return rows.map((val) => new ModelComment(val));
    }
}

module.exports = CommentRepositoryPostgres;
