const ModelThread = require('../../../../Domains/threads/entities/ModelThread');
const ThreadRepository = require('../../../../Domains/threads/ThreadRepository');
const ModelComment = require('../../../../Domains/comments/entities/ModelComment');
const CommentRepository = require('../../../../Domains/comments/CommentRepository');
const ModelReply = require('../../../../Domains/replies/entities/ModelReply');
const RepliesRepository = require('../../../../Domains/replies/RepliesRepository');
const GetThreadDetailUseCase = require('../GetThreadDetailUseCase');
const LikeRepository = require('../../../../Domains/likes/LikeRepository');

describe('GetThreadDetailUseCase', () => {
    it('should orchestracting the get thread details by thread id action correctly', async () => {
        // Arrange
        const mockThread = new ModelThread({
            id: 'thread-123',
            title: 'sebuah thread',
            body: 'sebuah body thread',
            date: new Date(),
            username: 'dicoding',
        });

        const mockComments = [
            new ModelComment({
                id: 'comment-123',
                username: 'dicoding',
                date: new Date(),
                content: 'sebuah comment',
                isDeleted: false,
            }),
        ];

        const mockReplies = [
            new ModelReply({
                id: 'reply-123',
                username: 'dicoding',
                date: new Date(),
                content: 'sebuah balasan',
                isDeleted: false,
            }),

            new ModelReply({
                id: 'reply-abc',
                username: 'dicoding',
                date: new Date(),
                content: 'sebuah balasan',
                isDeleted: true,
            }),
        ];

        /** creating dependency of use case */
        const mockThreadRepository = new ThreadRepository();
        const mockCommentRepository = new CommentRepository();
        const mockRepliesRepository = new RepliesRepository();
        const mockLikeRepository = new LikeRepository();

        /** mocking needed function */
        mockThreadRepository.getThreadById = jest.fn()
            .mockImplementation(() => Promise.resolve(mockThread));
        mockCommentRepository.commentsFromThread = jest.fn()
            .mockImplementation(() => Promise.resolve(mockComments));
        mockRepliesRepository.repliesFromComment = jest.fn()
            .mockImplementation(() => Promise.resolve(mockReplies));
        mockLikeRepository.getLikeCountComment = jest.fn()
            .mockImplementation(() => Promise.resolve(1));

        /** creating use case instance */
        const usecase = new GetThreadDetailUseCase({
            threadRepository: mockThreadRepository,
            commentRepository: mockCommentRepository,
            repliesRepository: mockRepliesRepository,
            likeRepository: mockLikeRepository,
        });

        // Action
        const threadDetails = await usecase.execute('thread-123');

        // Assert
        const expectedThreadDetails = {
            id: 'thread-123',
            title: 'sebuah thread',
            body: 'sebuah body thread',
            date: new Date(),
            username: 'dicoding',
            comments: [
                {
                    id: 'comment-123',
                    username: 'dicoding',
                    date: new Date(),
                    content: 'sebuah comment',
                    likeCount: 1,
                    replies: [
                        {
                            id: 'reply-123',
                            username: 'dicoding',
                            date: new Date(),
                            content: 'sebuah balasan',
                        },
                        {
                            id: 'reply-abc',
                            username: 'dicoding',
                            date: new Date(),
                            content: '**balasan telah dihapus**',
                        },
                    ],
                },
            ],
        };

        const comments = threadDetails.comments;
        const expectedComments = expectedThreadDetails.comments;

        const replies = comments[0].replies;
        const expectedReplies = expectedComments[0].replies;

        expect(threadDetails.id).toStrictEqual(expectedThreadDetails.id);
        expect(threadDetails.title).toStrictEqual(expectedThreadDetails.title);
        expect(threadDetails.body).toStrictEqual(expectedThreadDetails.body);
        expect(threadDetails.date.getDate()).toStrictEqual(expectedThreadDetails.date.getDate());
        expect(threadDetails.username).toStrictEqual(expectedThreadDetails.username);

        expect(comments[0].id).toStrictEqual(expectedComments[0].id);
        expect(comments[0].username).toStrictEqual(expectedComments[0].username);
        expect(comments[0].date.getDate()).toStrictEqual(expectedComments[0].date.getDate());
        expect(comments[0].content).toStrictEqual(expectedComments[0].content);
        expect(comments[0].isDeleted).toEqual(false);
        expect(comments[0].likeCount).toEqual(1);

        expect(replies[0].id).toStrictEqual(expectedReplies[0].id);
        expect(replies[0].username).toStrictEqual(expectedReplies[0].username);
        expect(replies[0].date.getDate()).toStrictEqual(expectedReplies[0].date.getDate());
        expect(replies[0].content).toStrictEqual(expectedReplies[0].content);
        expect(replies[0].isDeleted).toEqual(false);

        expect(replies[1].id).toStrictEqual(expectedReplies[1].id);
        expect(replies[1].username).toStrictEqual(expectedReplies[1].username);
        expect(replies[1].date.getDate()).toStrictEqual(expectedReplies[1].date.getDate());
        expect(replies[1].content).toStrictEqual(expectedReplies[1].content);
        expect(replies[1].isDeleted).toEqual(true);

        expect(mockThreadRepository.getThreadById).toBeCalledWith('thread-123');
        expect(mockCommentRepository.commentsFromThread).toBeCalledWith('thread-123');
        expect(mockRepliesRepository.repliesFromComment).toBeCalledWith('comment-123');
        expect(mockLikeRepository.getLikeCountComment).toBeCalledWith('comment-123');
    });
});
