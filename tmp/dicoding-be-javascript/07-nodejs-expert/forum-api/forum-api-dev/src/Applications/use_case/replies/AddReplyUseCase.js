const NewReply = require('../../../Domains/replies/entities/NewReply');

class AddReplyUseCase {
    constructor({ repliesRepository, commentRepository, threadRepository }) {
        this._repliesRepository = repliesRepository;
        this._commentRepository = commentRepository;
        this._threadRepository = threadRepository;
    }

    async execute(threadId, commentId, payload, owner) {
        const newReply = new NewReply(payload);

        await this._threadRepository.verifyThread(threadId);

        await this._commentRepository.verifyCommentIsExist(commentId, threadId);

        return this._repliesRepository.addReplyToComment(commentId, newReply, owner);
    }
}

module.exports = AddReplyUseCase;
