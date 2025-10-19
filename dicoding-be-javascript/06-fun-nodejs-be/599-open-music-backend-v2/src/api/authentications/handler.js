class AuthenticationsHandler {
    constructor(authenticationService, userService, tokenManager, validator) {
        this._authenticationService = authenticationService;
        this._userService = userService;
        this._tokenManager = tokenManager;
        this._validator = validator;
    }

    async postAuthenticationHandler(request, h) {
        const userValidated = await this._validator.validatePostAuthenticationPayload(request.payload);

        const id = await this._userService.verifyUserCredential(userValidated);

        const accessToken = await this._tokenManager.generateAccessToken({id});
        const refreshToken = await this._tokenManager.generateRefreshToken({id});

        await this._authenticationService.addRefreshToken(refreshToken);

        const response = h.response({
            status: 'success',
            message: 'Authentication berhasil ditambahkan',
            data: {
                accessToken,
                refreshToken,
            },
        });
        response.code(201);
        return response;
    }

    async putAuthenticationHandler(request, h) {
        const refreshTokenValidated = this._validator.validatePutAuthenticationPayload(request.payload);

        await this._authenticationService.verifyRefreshToken(refreshTokenValidated.refreshToken);
        
        const {id} = this._tokenManager.verifyRefreshToken(refreshTokenValidated.refreshToken);

        const accessToken = this._tokenManager.generateRefreshToken({id});

        const response = h.response({
            status: 'success',
            message: 'Access Token berhasil diperbarui',
            data: {
                accessToken,
            },
        });
        return response;
    }

    async deleteAuthenticationHandler(request, h) {
        const refreshTokenValidated = this._validator.validateDeleteAuthenticationPayload(request.payload);

        await this._authenticationService.verifyRefreshToken(refreshTokenValidated.refreshToken);
        await this._authenticationService.deleteRefreshToken(refreshTokenValidated.refreshToken);

        const response = h.response({
            status: 'success',
            message: 'Refresh token berhasil dihapus',
        });
        return response;
    }
}

module.exports = AuthenticationsHandler;
