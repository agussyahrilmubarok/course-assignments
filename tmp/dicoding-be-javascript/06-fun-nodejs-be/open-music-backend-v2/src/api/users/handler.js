class UserHandler {
    constructor(service, validator) {
        this._service = service;
        this._validator = validator;
    }

    async postUserHandler(request, h) {
        const userValidated = await this._validator.validateUserPayload(request.payload);

        const userId = await this._service.addUser(userValidated);

        const response = h.response({
            status: 'success',
            data: {
                userId: userId,
            },
        });
        response.code(201);
        return response;
    }
}

module.exports = UserHandler;
