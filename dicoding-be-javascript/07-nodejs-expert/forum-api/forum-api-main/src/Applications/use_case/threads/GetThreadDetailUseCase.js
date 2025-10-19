class GetThreadDetailUseCase {
    constructor({
        threadRepository, commentRepository, repliesRepository, likeRepository,
    }) {
        this._threadRepository = threadRepository;
        this._commentRepository = commentRepository;
        this._repliesRepository = repliesRepository;
        this._likeRepository = likeRepository;
    }

    async execute(threadId) {
        const threads = await this._threadRepository.getThreadById(threadId);
        const rowComment = await this._commentRepository.commentsFromThread(threadId);

        const comments = await Promise.all(rowComment.map(async (val) => {
            const likeCount = await this._likeRepository.getLikeCountComment(val.id);
            const replies = await this._repliesRepository.repliesFromComment(val.id);

            return { ...val, likeCount, replies };
        }));

        return { ...threads, comments };
    }
}

module.exports = GetThreadDetailUseCase;
