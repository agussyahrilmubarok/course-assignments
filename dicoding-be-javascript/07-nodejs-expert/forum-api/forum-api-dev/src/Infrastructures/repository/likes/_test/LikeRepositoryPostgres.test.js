const pool = require('../../../database/postgres/pool');

const UsersTableTestHelper = require('../../../../../tests/UsersTableTestHelper');
const ThreadsTableTestHelper = require('../../../../../tests/ThreadsTableTestHelper');
const CommentsTableTestHelper = require('../../../../../tests/CommentsTableTestHelper');
const CommentLikesTableTestHelper = require('../../../../../tests/CommentLikesTableTestHelper');

const LikeRepositoryPostgres = require('../LikeRepositoryPostgres');

describe('LikeRepositoryPostgres', () => {
    const dummyUser = {
        id: 'user-123',
        username: 'dicoding',
    };

    const dummyThread = {
        id: 'thread-123',
        owner: 'user-123',
    };

    const dummyComment = {
        id: 'comment-123',
        owner: 'user-123',
    };
    it('should be instance of LikeUnlikeRepository domain', () => {
        // Arrange
        const likeRepositoryPostgres = new LikeRepositoryPostgres({}, {});

        // Action & Assert
        expect(likeRepositoryPostgres).toBeInstanceOf(LikeRepositoryPostgres);
    });

    describe('behavior test', () => {
        beforeAll(async () => {
            await UsersTableTestHelper.addUser({ ...dummyUser });
            await ThreadsTableTestHelper.addThread({ ...dummyThread });
            await CommentsTableTestHelper.addCommentToThread({ ...dummyComment });
        });

        afterEach(async () => {
            await CommentLikesTableTestHelper.cleanTable();
        });

        afterAll(async () => {
            await UsersTableTestHelper.cleanTable();
            await ThreadsTableTestHelper.cleanTable();
            await CommentLikesTableTestHelper.cleanTable();
            await pool.end();
        });

        describe('addLikeToComment method', () => {
            it('should add like to comment', async () => {
                // Arrange
                const payload = {
                    commentId: dummyComment.id,
                    userId: dummyUser.id,
                };

                const fakeIdGenerator = () => '123';
                const likeRepositoryPostgres = new LikeRepositoryPostgres(pool, fakeIdGenerator);

                // Action
                const isLike = await likeRepositoryPostgres.addLikeToComment(payload);

                // Assert
                expect(isLike).toBeTruthy();
            });
        });

        describe('verifyLikeIsExists method', () => {
            it('should verify if comment is liked', async () => {
                // Arrange
                const payload = {
                    commentId: dummyComment.id,
                    userId: dummyUser.id,
                };

                await CommentLikesTableTestHelper.addLikeToComment({
                    id: 'like-123',
                    commentId: dummyComment.id,
                    owner: dummyUser.id,
                });

                const likeRepositoryPostgres = new LikeRepositoryPostgres(pool, {});

                // Action
                const isLike = await likeRepositoryPostgres.verifyLikeIsExists(payload);

                // Assert
                expect(isLike).toBeTruthy();
            });

            it('should return false if comment is not liked', async () => {
                // Arrange
                const likeRepositoryPostgres = new LikeRepositoryPostgres(pool, {});

                // Action
                const isLike = await likeRepositoryPostgres.verifyLikeIsExists({
                    commentId: 'comment-123',
                    userId: 'user',
                });

                // Assert
                expect(isLike).toBeFalsy();
            });
        });

        describe('unLikeComment method', () => {
            it('should unlike comment', async () => {
                // Arrange
                const payload = {
                    commentId: dummyComment.id,
                    userId: dummyUser.id,
                };

                await CommentLikesTableTestHelper.addLikeToComment({
                    id: 'like-123',
                    commentId: dummyComment.id,
                    owner: dummyUser.id,
                });

                const likeRepositoryPostgres = new LikeRepositoryPostgres(pool, {});

                // Action
                const isLike = await likeRepositoryPostgres.unLikeComment(payload);

                // Assert
                expect(isLike).toBeTruthy();
            });
        });

        describe('getLikeCountComment method', () => {
            it('should return like count of comment', async () => {
                // Arrange
                const payload = {
                    threadId: 'thread-123',
                    commentId: 'comment-123',
                    userId: 'user-123',
                };

                await UsersTableTestHelper.addUser({
                    id: 'user-321',
                    username: 'jhon_doe',
                });

                await CommentLikesTableTestHelper.addLikeToComment({
                    id: 'like-123',
                    commentId: dummyComment.id,
                    owner: dummyUser.id,
                });

                const likeRepositoryPostgres = new LikeRepositoryPostgres(pool, {});

                // Action
                const like = await likeRepositoryPostgres
                    .getLikeCountComment(payload.commentId);

                // Assert
                expect(like).toStrictEqual(1);
            });
        });
    });
});
