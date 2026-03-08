const pool = require('../../../database/postgres/pool');

const AuthorizationError = require('../../../../Commons/exceptions/AuthorizationError');
const NotFoundError = require('../../../../Commons/exceptions/NotFoundError');

const RepliesRepositoryPostgres = require('../RepliesRepositoryPostgres');

const AddedReply = require('../../../../Domains/replies/entities/AddedReply');
const NewReply = require('../../../../Domains/replies/entities/NewReply');

const UsersTableTestHelper = require('../../../../../tests/UsersTableTestHelper');
const ThreadsTableTestHelper = require('../../../../../tests/ThreadsTableTestHelper');
const CommentsTableTestHelper = require('../../../../../tests/CommentsTableTestHelper');
const RepliesTableTestHelper = require('../../../../../tests/RepliesTableTestHelper');

describe('RepliesRepositoryPostgres', () => {
    const dummyUser = {
        id: 'user-123',
        username: 'dicoding',
    };

    const dummyUser2 = {
        id: 'user-xyz',
        username: 'jhon doe',
    };

    const dummyThread = {
        id: 'thread-123',
        owner: 'user-123',
    };

    const dummyComment = {
        id: 'comment-123',
        owner: 'user-123',
    };

    beforeAll(async () => {
        await UsersTableTestHelper.addUser({ ...dummyUser });
        await UsersTableTestHelper.addUser({ ...dummyUser2 });
        await ThreadsTableTestHelper.addThread({ ...dummyThread });
        await CommentsTableTestHelper.addCommentToThread({
            commentId: dummyComment.id,
            threadId: dummyThread.id,
            owner: dummyComment.owner,
        });
    });

    afterEach(async () => {
        await RepliesTableTestHelper.cleanTable();
    });

    afterAll(async () => {
        await UsersTableTestHelper.cleanTable();
        await pool.end();
    });

    describe('addReplyToComment method', () => {
        it('should presist comment and return added comment correctly', async () => {
            // Arrange
            const newComment = new NewReply({ content: 'A reply' });
            const fakeIdGenerator = () => '123';

            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, fakeIdGenerator);

            // Action
            await repliesRepositoryPostgres
                .addReplyToComment(dummyComment.id, newComment, dummyUser.id);

            // Assert
            const comments = await RepliesTableTestHelper.findReplyById('reply-123');

            expect(comments).toHaveLength(1);
        });

        it('should return AddedReply correctly', async () => {
            // Arrange
            const newReply = new NewReply({ content: 'A reply' });
            const fakeIdGenerator = () => '123';

            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, fakeIdGenerator);

            // Action
            const addedReply = await repliesRepositoryPostgres
                .addReplyToComment(dummyComment.id, newReply, dummyUser.id);

            // Assert
            expect(addedReply).toStrictEqual(new AddedReply({
                id: 'reply-123',
                content: newReply.content,
                owner: dummyUser.id,
            }));
        });
    });

    describe('verifyReplyAccess method', () => {
        it('should throw AuthorizationError if user is not the comment, reply owner', async () => {
            // Arrange
            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, {});

            await RepliesTableTestHelper.addReplyToComment({ id: 'reply-123', owner: dummyUser.id });

            // Action & Assert
            await expect(repliesRepositoryPostgres.verifyReplyAccess('reply-123', dummyUser2.id))
                .rejects.toThrowError(AuthorizationError);
        });

        it('should not throw Authorization when user is the comment reply owner', async () => {
            // Arrange
            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, {});

            await RepliesTableTestHelper.addReplyToComment({});

            // Action & Assert
            await expect(repliesRepositoryPostgres.verifyReplyAccess('reply-123', dummyUser.id))
                .resolves.not.toThrowError(AuthorizationError);
        });
    });

    describe('deleteReply method', () => {
        it('should thow NotFoundError if reply is not found', async () => {
            // Arrange
            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, {});

            // Action & Assert
            await expect(repliesRepositoryPostgres.deleteReply('reply-123'))
                .rejects.toThrowError(NotFoundError);
        });

        it('should not thow NotFoundError if reply is valid', async () => {
            // Arrange
            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, {});
            await RepliesTableTestHelper.addReplyToComment({
                replyId: 'reply-123',
                commentId: dummyComment.id,
                owner: dummyUser.id,
            });

            // Action & Assert
            await expect(repliesRepositoryPostgres.deleteReply('reply-123'))
                .resolves.not.toThrowError(NotFoundError);
        });

        it('should update the reply delete status', async () => {
            // Arrange
            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, {});

            await RepliesTableTestHelper.addReplyToComment({ id: 'reply-123', owner: dummyUser.id });

            // Action
            await repliesRepositoryPostgres.deleteReply('reply-123', dummyUser.id);

            // Assert
            const [reply] = await RepliesTableTestHelper.findReplyById('reply-123');

            expect(reply.is_deleted).toEqual(true);
        });
    });

    describe('getRepliesFromComment method', () => {
        it('should return empty array if no replies are found', async () => {
            // Arrange
            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, {});

            // Action
            const replies = await repliesRepositoryPostgres.repliesFromComment([dummyComment.id]);

            // Assert
            expect(replies).toEqual([]);
        });

        it('should return array of replies with expected comment value', async () => {
            // Arrange
            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, {});

            await RepliesTableTestHelper.addReplyToComment({
                replyId: 'reply-123',
                commentId: dummyComment.id,
                content: 'A reply',
                owner: dummyUser.id,
            });

            // Action
            const [reply] = await repliesRepositoryPostgres.repliesFromComment('comment-123');

            // Assert
            expect(reply.id).toStrictEqual('reply-123');
            expect(reply.username).toStrictEqual(dummyUser.username);
            expect(reply.content).toStrictEqual('A reply');
            expect(reply.date.getDate()).toStrictEqual(new Date().getDate());
            expect(reply.isDeleted).toEqual(false);
        });

        it('should return comment with custom content if comment is deleted', async () => {
            // Arrange
            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, {});

            await RepliesTableTestHelper.addReplyToComment({
                replyId: 'reply-123',
                content: 'sebuah balasan',
                owner: dummyUser.id,
            });

            await RepliesTableTestHelper.addReplyToComment({
                replyId: 'reply-abc',
                content: 'sebuah balasan',
                owner: dummyUser2.id,
            });

            await RepliesTableTestHelper.deleteReply('reply-abc');

            // Action
            const [reply, deletedReply] = await repliesRepositoryPostgres
                .repliesFromComment(dummyComment.id);

            // Assert
            expect(reply.id).toStrictEqual('reply-123');
            expect(reply.username).toStrictEqual(dummyUser.username);
            expect(reply.content).toStrictEqual('sebuah balasan');
            expect(reply.date.getDate()).toStrictEqual(new Date().getDate());
            expect(reply.isDeleted).toEqual(false);

            expect(deletedReply.id).toStrictEqual('reply-abc');
            expect(deletedReply.username).toStrictEqual(dummyUser2.username);
            expect(deletedReply.content).toStrictEqual('**balasan telah dihapus**');
            expect(deletedReply.date.getDate()).toStrictEqual(new Date().getDate());
            expect(deletedReply.isDeleted).toEqual(true);
        });
    });

    describe('verifyReplyIsExist', () => {
        it('should throw NotFoundError if the reply is not found', async () => {
            // Arrange
            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, {});

            // Action & Assert
            expect(repliesRepositoryPostgres.verifyReplyIsExist('reply-123'))
                .rejects.toThrowError(NotFoundError);
        });

        it('should throw NotFoundError if the reply is invalid', async () => {
            // Arrange
            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, {});

            await RepliesTableTestHelper.addReplyToComment({
                replyId: 'reply-123',
                commentId: 'comment-123',
                content: 'sebuah balasan',
                owner: dummyUser.id,
            });

            // Action & Assert
            expect(repliesRepositoryPostgres.verifyReplyIsExist('reply-xyz'))
                .rejects.toThrowError(NotFoundError);
        });

        it('should not throw NotFoundError if the reply is valid', async () => {
            // Arrange
            const repliesRepositoryPostgres = new RepliesRepositoryPostgres(pool, {});

            await RepliesTableTestHelper.addReplyToComment({
                replyId: 'reply-123',
                commentId: 'comment-123',
                content: 'sebuah balasan',
                owner: dummyUser.id,
            });

            // Action & Assert
            expect(repliesRepositoryPostgres.verifyReplyIsExist('reply-123'))
                .resolves.not.toThrowError(NotFoundError);
        });
    });
});
