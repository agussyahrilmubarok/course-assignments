const pool = require('../../database/postgres/pool');

const AuthenticationsTableTestHelper = require('../../../../tests/AuthenticationsTableTestHelper');
const UsersTableTestHelper = require('../../../../tests/UsersTableTestHelper');
const ThreadsTableTestHelper = require('../../../../tests/ThreadsTableTestHelper');
const CommentsTableTestHelper = require('../../../../tests/CommentsTableTestHelper');
const RepliesTableTestHelper = require('../../../../tests/RepliesTableTestHelper');
const CommentLikesTableTestHelper = require('../../../../tests/CommentLikesTableTestHelper');

const container = require('../../container');
const createServer = require('../createServer');

describe('/threads endpoint', () => {
    afterAll(async () => {
        await pool.end();
    });

    afterEach(async () => {
        await UsersTableTestHelper.cleanTable();
        await AuthenticationsTableTestHelper.cleanTable();
        await ThreadsTableTestHelper.cleanTable();
        await CommentsTableTestHelper.cleanTable();
        await RepliesTableTestHelper.cleanTable();
        await CommentLikesTableTestHelper.cleanTable();
    });

    describe('when POST /threads', () => {
        it('should response 201 and persisted thread', async () => {
            // Arrange
            const requestAddUser = {
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
            };

            const requestPayload = {
                title: 'A thread',
                body: 'A thread body',
            };

            const server = await createServer(container);

            // Action
            /* adding user */
            await server.inject({
                method: 'POST',
                url: '/users',
                payload: {
                    username: requestAddUser.username,
                    password: requestAddUser.password,
                    fullname: requestAddUser.fullname,
                },
            });

            /* login user */
            const responseAtuh = await server.inject({
                method: 'POST',
                url: '/authentications',
                payload: {
                    username: requestAddUser.username,
                    password: requestAddUser.password,
                },
            });

            const responseAtuhJson = JSON.parse(responseAtuh.payload);
            const accessToken = responseAtuhJson.data.accessToken;

            /* adding thread */
            const response = await server.inject({
                method: 'POST',
                url: '/threads',
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: requestPayload,
            });

            // Assert
            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toEqual(201);
            expect(responseJson.status).toEqual('success');
            expect(responseJson.data.addedThread).toBeDefined();
            expect(responseJson.data.addedThread.id).toBeDefined();
            expect(responseJson.data.addedThread.title).toBeDefined();
            expect(responseJson.data.addedThread.owner).toBeDefined();
        });

        it('should response 400 status when request a bad payload', async () => {
            // Arrange
            const requestAddUser = {
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
            };

            const requestPayload = {
                title: 'A thread',
                body: 123,
            };

            const server = await createServer(container);

            // Action
            /* adding user */
            await server.inject({
                method: 'POST',
                url: '/users',
                payload: {
                    username: requestAddUser.username,
                    password: requestAddUser.password,
                    fullname: requestAddUser.fullname,
                },
            });

            /* login user */
            const responseAtuh = await server.inject({
                method: 'POST',
                url: '/authentications',
                payload: {
                    username: requestAddUser.username,
                    password: requestAddUser.password,
                },
            });

            const responseAtuhJson = JSON.parse(responseAtuh.payload);
            const accessToken = responseAtuhJson.data.accessToken;

            /* adding thread */
            const response = await server.inject({
                method: 'POST',
                url: '/threads',
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: requestPayload,
            });

            // Assert
            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toEqual(400);
            expect(responseJson.status).toEqual('fail');
            expect(responseJson.message).toEqual('tipe data pada thread tidak valid');
        });

        it('should response 401 status when the request does not have an authentication', async () => {
            // Arrange
            const server = await createServer(container);

            // Action
            const response = await server.inject({
                method: 'POST',
                url: '/threads',
            });

            // Assert
            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toEqual(401);
            expect(responseJson.error).toBeDefined();
            expect(responseJson.message).toEqual('Missing authentication');
        });
    });

    describe('when GET /threads/{threadId}', () => {
        it('should response 200 and show thread by id', async () => {
            // Arrange
            const threadId = 'thread-123';
            await UsersTableTestHelper.addUser({
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
            });

            await ThreadsTableTestHelper.addThread({
                id: 'thread-123',
                title: 'sebuah thread',
                body: 'sebuah body thread',
                owner: 'user-123',
            });

            const server = await createServer(container);

            // Action
            const response = await server.inject({
                method: 'GET',
                url: `/threads/${threadId}`,
            });

            // Assert
            const responseJson = JSON.parse(response.payload);
            expect(response.statusCode).toEqual(200);
            expect(responseJson.status).toEqual('success');
            expect(responseJson.data.thread).toBeDefined();
        });

        it('should respond with 200 with thread details and comments', async () => {
            const requestAddUser1 = {
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
            };
            const requestAddUser2 = {
                id: 'user-xyz',
                username: 'dicoding2',
                password: 'supersecret',
                fullname: 'Dicoding Indonesia 2',
            };

            const server = await createServer(container);

            await UsersTableTestHelper.addUser(requestAddUser1);
            await UsersTableTestHelper.addUser(requestAddUser2);

            await ThreadsTableTestHelper.addThread({
                id: 'thread-123',
                owner: requestAddUser1.id,
            });

            await CommentsTableTestHelper.addCommentToThread({
                commentId: 'comment-123',
                threadId: 'thread-123',
                content: `komentar dari ${requestAddUser1.username}`,
                owner: requestAddUser1.id,
            });

            await CommentsTableTestHelper.addCommentToThread({
                commentId: 'comment-xyz',
                threadId: 'thread-123',
                content: `komentar dari ${requestAddUser2.username}`,
                owner: requestAddUser2.id,
            });

            await CommentsTableTestHelper.deleteComment('comment-xyz');

            await RepliesTableTestHelper.addReplyToComment({
                replyId: 'reply-123',
                commentId: 'comment-123',
                content: `balasan dari ${requestAddUser1.username}`,
                owner: requestAddUser1.id,
            });

            await RepliesTableTestHelper.addReplyToComment({
                replyId: 'reply-xyz',
                commentId: 'comment-123',
                content: `balasan dari ${requestAddUser2.username}`,
                owner: requestAddUser2.id,
            });

            await RepliesTableTestHelper.deleteReply('reply-xyz');

            await CommentLikesTableTestHelper.addLikeToComment({
                id: 'like-123',
                commentId: 'comment-123',
                owner: requestAddUser1.id,
            });

            // Action
            const response = await server.inject({
                method: 'GET',
                url: '/threads/thread-123',
            });

            // Assert
            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toStrictEqual(200);
            expect(responseJson.status).toStrictEqual('success');

            const threads = responseJson.data.thread;

            expect(threads.id).toStrictEqual('thread-123');
            expect(threads.title).toStrictEqual('sebuah thread');
            expect(threads.body).toStrictEqual('sebuah body thread');
            expect(new Date(threads.date).getDate()).toStrictEqual(new Date().getDate());
            expect(threads.username).toStrictEqual(requestAddUser1.username);

            const [requestAddUser1Comment, requestAddUser2Comment] = threads.comments;

            expect(requestAddUser1Comment.id).toStrictEqual('comment-123');
            expect(requestAddUser1Comment.username).toStrictEqual(requestAddUser1.username);
            expect(new Date(requestAddUser1Comment.date).getDate())
                .toStrictEqual(new Date().getDate());
            expect(requestAddUser1Comment.content)
                .toStrictEqual(`komentar dari ${requestAddUser1.username}`);
            expect(requestAddUser1Comment.likeCount).toEqual(1);

            expect(requestAddUser2Comment.id).toStrictEqual('comment-xyz');
            expect(requestAddUser2Comment.username).toStrictEqual(requestAddUser2.username);
            expect(new Date(requestAddUser2Comment.date).getDate())
                .toStrictEqual(new Date().getDate());
            expect(requestAddUser2Comment.content).toStrictEqual('**komentar telah dihapus**');
            expect(requestAddUser2Comment.likeCount).toEqual(0);

            const [requestAddUser1Reply, requestAddUser2Reply] = requestAddUser1Comment.replies;

            expect(requestAddUser1Reply.id).toStrictEqual('reply-123');
            expect(requestAddUser1Reply.username).toStrictEqual(requestAddUser1.username);
            expect(new Date(requestAddUser1Reply.date).getDate())
                .toStrictEqual(new Date().getDate());
            expect(requestAddUser1Reply.content).toStrictEqual(`balasan dari ${requestAddUser1.username}`);

            expect(requestAddUser2Reply.id).toStrictEqual('reply-xyz');
            expect(requestAddUser2Reply.username).toStrictEqual(requestAddUser2.username);
            expect(new Date(requestAddUser2Reply.date).getDate())
                .toStrictEqual(new Date().getDate());
            expect(requestAddUser2Reply.content).toStrictEqual('**balasan telah dihapus**');
        });

        it('should response 404 status when the thread is not found', async () => {
            // Arrange
            const requestAddUser = {
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
            };

            const server = await createServer(container);

            await UsersTableTestHelper.addUser(requestAddUser);

            // Action
            const response = await server.inject({
                method: 'GET',
                url: '/threads/thread-123',
            });

            // Assert
            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toEqual(404);
            expect(responseJson.status).toEqual('fail');
            expect(responseJson.message).toEqual('thread tidak ditemukan');
        });
    });
});
