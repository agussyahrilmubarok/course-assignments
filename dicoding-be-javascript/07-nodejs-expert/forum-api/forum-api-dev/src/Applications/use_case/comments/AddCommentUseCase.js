const NewComment = require('../../../Domains/comments/entities/NewComment');

class AddCommentUseCase {
    constructor({ commentRepository, threadRepository }) {
        this._commentRepository = commentRepository;
        this._threadRepository = threadRepository;
    }

    async execute(threadId, payload, owner) {
        const newComment = new NewComment(payload);

        await this._threadRepository.verifyThread(threadId);

        return this._commentRepository.addCommentToThread(threadId, newComment, owner);
    }
}

module.exports = AddCommentUseCase;
