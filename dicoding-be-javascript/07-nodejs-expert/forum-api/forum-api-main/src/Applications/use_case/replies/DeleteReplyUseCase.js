class DeleteReplyUseCase {
    constructor({ repliesRepository }) {
        this._repliesRepository = repliesRepository;
    }

    async execute(useCaseDeletePayload) {
        const {
            replyId, userId,
        } = useCaseDeletePayload;

        await this._repliesRepository.verifyReplyIsExist(replyId);

        await this._repliesRepository.verifyReplyAccess(replyId, userId);

        return this._repliesRepository.deleteReply(replyId);
    }
}

module.exports = DeleteReplyUseCase;
