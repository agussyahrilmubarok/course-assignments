class ModelReply {
    constructor(payload) {
        this._verifyPayload(payload);

        const {
            id, username, date, content, isDeleted,
        } = payload;

        this.id = id;
        this.username = username;
        this.date = date;
        this.content = (isDeleted) ? '**balasan telah dihapus**' : content;
        this.isDeleted = isDeleted;
    }

    _verifyPayload({
        id, username, date, content, isDeleted,
    }) {
        if (!id || !username || !date || !content || typeof isDeleted === 'undefined') {
            throw new Error('COMMENT.NOT_CONTAIN_NEEDED_PROPERTY');
        }

        if (
            typeof id !== 'string' || typeof username !== 'string'
            || !(date instanceof Date) || typeof content !== 'string'
            || typeof isDeleted !== 'boolean'
        ) {
            throw new Error('COMMENT.NOT_MEET_DATA_TYPE_SPECIFICATION');
        }
    }
}

module.exports = ModelReply;
