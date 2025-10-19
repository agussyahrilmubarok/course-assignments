class DeleteCommentUseCase {
    constructor({ commentRepository }) {
        this._commentRepository = commentRepository;
    }

    async execute(useCasePayload) {
        const {
            threadId, commentId, userId,
        } = useCasePayload;

        await this._commentRepository.verifyCommentIsExist(commentId, threadId);
        await this._commentRepository.verifyCommentAccess(commentId, userId);
        await this._commentRepository.deleteCommentById(commentId);
    }
}

module.exports = DeleteCommentUseCase;
