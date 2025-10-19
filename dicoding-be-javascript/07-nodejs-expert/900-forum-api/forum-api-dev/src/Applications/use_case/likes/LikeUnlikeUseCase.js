class LikeUnlikeUseCase {
    constructor({ threadRepository, commentRepository, likeRepository }) {
        this._threadRepository = threadRepository;
        this._commentRepository = commentRepository;
        this._likeRepository = likeRepository;
    }

    async execute(useCasePayload) {
        await this._threadRepository.verifyThread(useCasePayload.threadId);
        await this._commentRepository.verifyCommentIsExist(
            useCasePayload.commentId,
            useCasePayload.threadId,
        );

        const isLike = await this._likeRepository.verifyLikeIsExists(useCasePayload);

        if (isLike > 0) {
            return this._likeRepository.unLikeComment(useCasePayload);
        }

        return this._likeRepository.addLikeToComment(useCasePayload);
    }
}

module.exports = LikeUnlikeUseCase;
