const LikeRepository = require('../../../Domains/likes/LikeRepository');

class LikeRepositoryPostgres extends LikeRepository {
    constructor(pool, idGenerator) {
        super();
        this._pool = pool;
        this._idGenerator = idGenerator;
    }

    async addLikeToComment(payload) {
        const { commentId, userId } = payload;
        const id = `like-${this._idGenerator()}`;

        const query = {
            text: 'INSERT INTO comment_likes VALUES($1, $2, $3)',
            values: [id, commentId, userId],
        };

        const { rowCount } = await this._pool.query(query);

        return rowCount;
    }

    async verifyLikeIsExists(payload) {
        const { commentId, userId } = payload;

        const query = {
            text: 'SELECT 1 FROM comment_likes WHERE comment_id = $1 AND owner = $2',
            values: [commentId, userId],
        };

        const { rowCount } = await this._pool.query(query);

        return rowCount;
    }

    async unLikeComment(payload) {
        const { commentId, userId } = payload;

        const query = {
            text: 'DELETE FROM comment_likes WHERE comment_id = $1 AND owner = $2',
            values: [commentId, userId],
        };

        const { rowCount } = await this._pool.query(query);

        return rowCount;
    }

    async getLikeCountComment(commentId) {
        const query = {
            text: 'SELECT id FROM comment_likes WHERE comment_id = $1',
            values: [commentId],
        };

        const result = await this._pool.query(query);
        return result.rowCount;
    }
}

module.exports = LikeRepositoryPostgres;
