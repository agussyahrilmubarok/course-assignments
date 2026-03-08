const AddThreadUseCase = require('../../../../Applications/use_case/threads/AddThreadUseCase');
const GetThreadDetailUseCase = require('../../../../Applications/use_case/threads/GetThreadDetailUseCase');

const AddCommentUseCase = require('../../../../Applications/use_case/comments/AddCommentUseCase');
const DeleteCommentUseCase = require('../../../../Applications/use_case/comments/DeleteCommentUseCase');

const AddReplyUseCase = require('../../../../Applications/use_case/replies/AddReplyUseCase');
const DeleteReplyUseCase = require('../../../../Applications/use_case/replies/DeleteReplyUseCase');

const LikeUnlikeUseCase = require('../../../../Applications/use_case/likes/LikeUnlikeUseCase');

class ThreadsHandler {
    constructor(container) {
        this._container = container;

        this.postThreadHandler = this.postThreadHandler.bind(this);
        this.getThreadDetailsHandler = this.getThreadDetailsHandler.bind(this);
        this.postCommentHandler = this.postCommentHandler.bind(this);
        this.deleteCommentHandler = this.deleteCommentHandler.bind(this);
        this.postReplyHandler = this.postReplyHandler.bind(this);
        this.deleteReplyHandler = this.deleteReplyHandler.bind(this);
        this.putLikeUnlikeHandler = this.putLikeUnlikeHandler.bind(this);
    }

    async postThreadHandler(request, h) {
        const { id } = request.auth.credentials;
        const addThreadUseCase = this._container.getInstance(AddThreadUseCase.name);
        const addedThread = await addThreadUseCase.execute(request.payload, id);

        const response = h.response({
            status: 'success',
            data: { addedThread },
        });
        response.code(201);
        return response;
    }

    async getThreadDetailsHandler(request, h) {
        const { threadId } = request.params;

        const getThreadDetails = this._container
            .getInstance(GetThreadDetailUseCase.name);

        const thread = await getThreadDetails.execute(threadId);

        return h.response({
            status: 'success',
            data: { thread },
        }).code(200);
    }

    async postCommentHandler(request, h) {
        const { id } = request.auth.credentials;
        const { threadId } = request.params;

        const addCommentUseCase = this._container
            .getInstance(AddCommentUseCase.name);

        const addedComment = await addCommentUseCase
            .execute(threadId, request.payload, id);

        const response = h.response({
            status: 'success',
            data: { addedComment },
        });
        response.code(201);
        return response;
    }

    async deleteCommentHandler(request, h) {
        const { threadId, commentId } = request.params;
        const { id: credentialId } = request.auth.credentials;

        const deleteCommentUseCase = this._container
            .getInstance(DeleteCommentUseCase.name);

        // await deleteCommentUseCase.execute(threadId, commentId, id);
        await deleteCommentUseCase.execute({
            threadId,
            commentId,
            userId: credentialId,
        });

        const response = h.response({
            status: 'success',
        });
        response.code(200);
        return response;
    }

    async postReplyHandler(request, h) {
        const { id } = request.auth.credentials;
        const { threadId, commentId } = request.params;

        const addReplyUseCase = this._container
            .getInstance(AddReplyUseCase.name);

        const addedReply = await addReplyUseCase
            .execute(threadId, commentId, request.payload, id);

        const response = h.response({
            status: 'success',
            data: { addedReply },
        });
        response.code(201);
        return response;
    }

    async deleteReplyHandler(request, h) {
        const { threadId, commentId, replyId } = request.params;
        const { id: credentialId } = request.auth.credentials;

        const deleteReplyUseCase = this._container
            .getInstance(DeleteReplyUseCase.name);

        await deleteReplyUseCase.execute({
            threadId,
            commentId,
            replyId,
            userId: credentialId,
        });

        const response = h.response({
            status: 'success',
        });
        response.code(200);
        return response;
    }

    async putLikeUnlikeHandler(request, h) {
        const { threadId, commentId } = request.params;
        const { id: credentialId } = request.auth.credentials;

        const likeUnLikeUseCase = this._container.getInstance(LikeUnlikeUseCase.name);

        const useCasePayload = {
            threadId,
            commentId,
            userId: credentialId,
        };

        await likeUnLikeUseCase.execute(useCasePayload);

        const response = h.response({
            status: 'success',
        });
        response.code(200);
        return response;
    }
}

module.exports = ThreadsHandler;
