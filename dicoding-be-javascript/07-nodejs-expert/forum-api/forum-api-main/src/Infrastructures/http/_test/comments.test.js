const pool = require('../../database/postgres/pool');

const AuthenticationsTableTestHelper = require('../../../../tests/AuthenticationsTableTestHelper');
const UsersTableTestHelper = require('../../../../tests/UsersTableTestHelper');
const ThreadsTableTestHelper = require('../../../../tests/ThreadsTableTestHelper');
const CommentsTableTestHelper = require('../../../../tests/CommentsTableTestHelper');
const RepliesTableTestHelper = require('../../../../tests/RepliesTableTestHelper');

const container = require('../../container');
const createServer = require('../createServer');

describe('/comments endpoint', () => {
    afterAll(async () => {
        await pool.end();
    });

    afterEach(async () => {
        await UsersTableTestHelper.cleanTable();
        await AuthenticationsTableTestHelper.cleanTable();
        await ThreadsTableTestHelper.cleanTable();
        await CommentsTableTestHelper.cleanTable();
        await RepliesTableTestHelper.cleanTable();
    });

    describe('when POST /threads/{threadId}/comments', () => {
        it('should response 201 status and added comment ', async () => {
            // Arrange
            const requestAddUser = {
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
            };
            const requestCommentPayload = { content: 'A comment' };

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
            const responseAddThread = await server.inject({
                method: 'POST',
                url: '/threads',
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: {
                    id: 'thread-123',
                    title: 'new thread',
                    body: 'new thread body',
                },
            });

            const responseAddThreadJson = JSON.parse(responseAddThread.payload);
            const threadId = responseAddThreadJson.data.addedThread.id;

            const response = await server.inject({
                method: 'POST',
                url: `/threads/${threadId}/comments`,
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: requestCommentPayload,
            });

            // Assert
            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toEqual(201);
            expect(responseJson.status).toEqual('success');
            expect(responseJson.data.addedComment).toBeDefined();
            expect(responseJson.data.addedComment.id).toBeDefined();
            expect(responseJson.data.addedComment.content).toBeDefined();
            expect(responseJson.data.addedComment.owner).toBeDefined();
        });

        it('should response 401 status when the request does not have an authentication', async () => {
            // Arrange
            const server = await createServer(container);

            // Action
            const response = await server.inject({
                method: 'POST',
                url: '/threads/thread-123/comments',
            });

            // Assert
            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toEqual(401);
            expect(responseJson.error).toBeDefined();
            expect(responseJson.message).toEqual('Missing authentication');
        });

        it('should response 400 status when request a bad payload', async () => {
            // Arrange
            const requestAddUser = {
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
            };
            const requestCommentPayload = { content: 400 };

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
            const responseAddThread = await server.inject({
                method: 'POST',
                url: '/threads',
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: {
                    id: 'thread-123',
                    title: 'new thread',
                    body: 'new thread body',
                },
            });

            const responseAddThreadJson = JSON.parse(responseAddThread.payload);
            const threadId = responseAddThreadJson.data.addedThread.id;

            const response = await server.inject({
                method: 'POST',
                url: `/threads/${threadId}/comments`,
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: requestCommentPayload,
            });

            // Assert
            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toEqual(400);
            expect(responseJson.status).toEqual('fail');
            expect(responseJson.message).toEqual('tipe data pada comment tidak valid');
        });

        it('should response 404 status when the thread is not found', async () => {
            // Arrange
            const requestAddUser = {
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
            };
            const commentRequestPayload = { content: 'A comment' };
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
            const responseAddThread = await server.inject({
                method: 'POST',
                url: '/threads',
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: {},
            });

            const responseAddThreadJson = JSON.parse(responseAddThread.payload);
            const threadId = responseAddThreadJson.data;

            const response = await server.inject({
                method: 'POST',
                url: `/threads/${threadId}/comments`,
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: commentRequestPayload,
            });

            // Assert
            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toEqual(404);
            expect(responseJson.status).toEqual('fail');
            expect(responseJson.message).toEqual('thread tidak ditemukan');
        });
    });

    describe('when DELETE /threads/{threadId}/comments/{commentId}', () => {
        it('should response 200 and persisted delete comment ', async () => {
            // Arrange
            const requestAddUser = {
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
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
            const responseAddThread = await server.inject({
                method: 'POST',
                url: '/threads',
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: {
                    id: 'thread-123',
                    title: 'new thread',
                    body: 'new thread body',
                },
            });

            const responseAddThreadJson = JSON.parse(responseAddThread.payload);
            const threadId = responseAddThreadJson.data.addedThread.id;

            const responseAddComment = await server.inject({
                method: 'POST',
                url: `/threads/${threadId}/comments`,
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: {
                    id: 'comment-123',
                    content: 'A comment',
                },
            });

            // Assert
            const responseAddCommentJson = JSON.parse(responseAddComment.payload);
            const commentId = responseAddCommentJson.data.addedComment.id;

            // Action
            const response = await server.inject({
                method: 'DELETE',
                url: `/threads/${threadId}/comments/${commentId}`,
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
            });

            // Assert
            const responseJson = JSON.parse(response.payload);
            expect(response.statusCode).toEqual(200);
            expect(responseJson.status).toEqual('success');
        });

        it('should response 401 when the request does not have an authentication', async () => {
            // Arrange
            const requestAddUser = {
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
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
            const responseAddThread = await server.inject({
                method: 'POST',
                url: '/threads',
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: {
                    id: 'thread-123',
                    title: 'new thread',
                    body: 'new thread body',
                },
            });

            const responseAddThreadJson = JSON.parse(responseAddThread.payload);
            const threadId = responseAddThreadJson.data.addedThread.id;

            const responseAddComment = await server.inject({
                method: 'POST',
                url: `/threads/${threadId}/comments`,
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: {
                    id: 'comment-123',
                    content: 'A comment',
                },
            });

            // Assert
            const responseAddCommentJson = JSON.parse(responseAddComment.payload);
            const commentId = responseAddCommentJson.data.addedComment.id;

            const response = await server.inject({
                method: 'DELETE',
                url: `/threads/${threadId}/comments/${commentId}`,
            });

            // Assert
            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toEqual(401);
            expect(responseJson.error).toBeDefined();
            expect(responseJson.message).toEqual('Missing authentication');
        });

        it('should response 403 when user request not authorized', async () => {
            // Arrange
            const requestAddUser = {
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
            };

            const requestAuthPayload1 = {
                username: requestAddUser.username,
                password: requestAddUser.password,
            };

            const requestAuthPayload2 = {
                username: 'dicoding2',
                password: 'supersecret',
            };

            const requestAddThreadPayload = {
                title: 'new thread',
                body: 'new thread body',
            };

            const requestAddCommentPayload = {
                content: 'new comment',
            };

            const server = await createServer(container);

            /* adding user */
            await server.inject({
                method: 'POST',
                url: '/users',
                payload: {
                    username: requestAuthPayload1.username,
                    password: requestAuthPayload1.password,
                    fullname: requestAddUser.fullname,
                },
            });

            await server.inject({
                method: 'POST',
                url: '/users',
                payload: {
                    username: requestAuthPayload2.username,
                    password: requestAuthPayload2.password,
                    fullname: requestAddUser.fullname,
                },
            });

            // Action
            /* login */
            const responseAuth1 = await server.inject({
                method: 'POST',
                url: '/authentications',
                payload: requestAuthPayload1,
            });
            const responseAuthJson1 = JSON.parse(responseAuth1.payload);
            const accessToken1 = responseAuthJson1.data.accessToken;

            const responseAuth2 = await server.inject({
                method: 'POST',
                url: '/authentications',
                payload: requestAuthPayload2,
            });
            const responseAuthJson2 = JSON.parse(responseAuth2.payload);
            const accessToken2 = responseAuthJson2.data.accessToken;

            /* adding thread */
            const responseAddThread = await server.inject({
                method: 'POST',
                url: '/threads',
                payload: requestAddThreadPayload,
                headers: {
                    Authorization: `Bearer ${accessToken1}`,
                },
            });
            const responseAddThreadJson = JSON.parse(responseAddThread.payload);
            const threadId = responseAddThreadJson.data.addedThread.id;

            /* adding comment */
            const responseAddComment = await server.inject({
                method: 'POST',
                url: `/threads/${threadId}/comments`,
                payload: requestAddCommentPayload,
                headers: {
                    // use access token from user 1
                    Authorization: `Bearer ${accessToken1}`,
                },
            });

            const responseAddCommentJson = JSON.parse(responseAddComment.payload);
            const commentId = responseAddCommentJson.data.addedComment.id;

            /* deleting comment */
            const response = await server.inject({
                method: 'DELETE',
                url: `/threads/${threadId}/comments/${commentId}`,
                headers: {
                    // use access token from user 2
                    Authorization: `Bearer ${accessToken2}`,
                },
            });

            // Assert
            const responseJson = JSON.parse(response.payload);
            expect(response.statusCode).toEqual(403);
            expect(responseJson.status).toEqual('fail');
            expect(responseJson.message).toEqual('anda tidak dapat mengakses resource ini');
        });

        it('should respond 404 status when the thread is not found', async () => {
            // Arrange
            const requestAddUser = {
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
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
            const responseAddThread = await server.inject({
                method: 'POST',
                url: '/threads',
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: {},
            });

            const responseAddThreadJson = JSON.parse(responseAddThread.payload);
            const threadId = responseAddThreadJson.data;

            const responseAddComment = await server.inject({
                method: 'POST',
                url: `/threads/${threadId}/comments`,
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: {},
            });

            // Assert
            const responseAddCommentJson = JSON.parse(responseAddComment.payload);
            const commentId = responseAddCommentJson.data;

            const response = await server.inject({
                method: 'DELETE',
                url: `/threads/${threadId}/comments/${commentId}`,
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
            });

            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toEqual(404);
            expect(responseJson.status).toEqual('fail');
            expect(responseJson.message).toEqual('komentar tidak ditemukan');
        });
    });
});
