const pool = require('../../database/postgres/pool');

const AuthenticationsTableTestHelper = require('../../../../tests/AuthenticationsTableTestHelper');
const UsersTableTestHelper = require('../../../../tests/UsersTableTestHelper');

const ThreadsTableTestHelper = require('../../../../tests/ThreadsTableTestHelper');
const CommentsTableTestHelper = require('../../../../tests/CommentsTableTestHelper');
const CommentLikesTableTestHelper = require('../../../../tests/CommentLikesTableTestHelper');

const container = require('../../container');
const createServer = require('../createServer');

describe('/like endpoint', () => {
    afterEach(async () => {
        await UsersTableTestHelper.cleanTable();
        await AuthenticationsTableTestHelper.cleanTable();
        await ThreadsTableTestHelper.cleanTable();
        await CommentsTableTestHelper.cleanTable();
        await CommentLikesTableTestHelper.cleanTable();
    });

    afterAll(async () => {
        await pool.end();
    });

    describe('when POST /threads/{threadId}/comments/{commentId}/likes', () => {
        it('should response 201 and persisted like', async () => {
            // Arrange
            const requestAddUser = {
                id: 'user-123',
                username: 'dicoding',
                password: 'secret',
                fullname: 'Dicoding Indonesia',
            };

            const requestThreadPayload = {
                title: 'A thread',
                body: 'A thread body',
            };

            const requestCommentPayload = {
                id: 'comment-123',
                content: 'A comment',
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
                payload: {
                    id: 'thread-123',
                    title: requestThreadPayload.title,
                    body: requestThreadPayload.body,
                },
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
            });

            const responseAddThreadJson = JSON.parse(responseAddThread.payload);
            const threadId = responseAddThreadJson.data.addedThread.id;

            /* adding comment */
            const responseAddComment = await server.inject({
                method: 'POST',
                url: `/threads/${threadId}/comments`,
                payload: requestCommentPayload,
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
            });

            const responseAddCommentJson = JSON.parse(responseAddComment.payload);
            const commentId = responseAddCommentJson.data.addedComment.id;

            /* adding like */
            const response = await server.inject({
                method: 'PUT',
                url: `/threads/${threadId}/comments/${commentId}/likes`,
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
            const server = await createServer(container);

            // Action
            const response = await server.inject({
                method: 'PUT',
                url: '/threads/thread-123/comments/comment-123/likes',
            });

            // Assert
            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toEqual(401);
            expect(responseJson.error).toEqual('Unauthorized');
            expect(responseJson.message).toEqual('Missing authentication');
        });

        it('should response 404 when comment not found', async () => {
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

            /* adding comment */
            const responseAddComment = await server.inject({
                method: 'POST',
                url: `/threads/${threadId}/comments`,
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                payload: {},
            });

            const responseAddCommentJson = JSON.parse(responseAddComment.payload);
            const commentId = responseAddCommentJson.data;

            /* adding like */
            const response = await server.inject({
                method: 'PUT',
                url: `/threads/${threadId}/comments/${commentId}/likes`,
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
            });

            // Assert
            const responseJson = JSON.parse(response.payload);

            expect(response.statusCode).toEqual(404);
            expect(responseJson.status).toEqual('fail');
            expect(responseJson.message).toEqual('thread tidak ditemukan');
        });
    });
});
