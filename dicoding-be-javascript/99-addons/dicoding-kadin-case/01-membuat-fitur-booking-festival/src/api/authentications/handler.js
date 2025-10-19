// This class is used to handle all requests from client to server
// related to authentication resources.
// in here we defined business logic to handle the request.
class AuthenticationsHandler {
  constructor(authenticationsService, usersService, jwtTokenize, validator) {
    this._authenticationsService = authenticationsService;
    this._usersService = usersService;
    this._jwtTokenize = jwtTokenize;
    this._validator = validator;
  }

  // This method will be used to handle POST request to /authentications resource.
  // This method will be used to handle authentication request (login).
  async postAuthenticationHandler(request, h) {
    // validate the payload.
    // the validator will throw an error if the payload is not valid.
    // if the payload is valid, then the method will return the payload.
    const { email, password } = this._validator.validatePostAuthenticationPayload(request.payload);

    // verify user credential
    // this method will throw an error if the credential is not valid
    // if the credential is valid, then the method will return the user id
    const id = await this._usersService.verifyUserCredential(email, password);

    // generate access token and refresh token
    const accessToken = this._jwtTokenize.generateAccessToken({ id });
    const refreshToken = this._jwtTokenize.generateRefreshToken({ id });

    // persist refresh token to database
    await this._authenticationsService.persistRefreshToken(refreshToken);

    // return response with access token and refresh token
    return h.response({
      status: 'success',
      message: 'authentication berhasil',
      data: {
        accessToken,
        refreshToken,
      },
    }).code(201);
  }

  // This method will be used to handle PUT request to /authentications resource.
  // This method will be used to handle refresh token request.
  async putAuthenticationHandler(request) {
    // validate the payload.
    // Same principle as the previous method.
    const { refreshToken } = this._validator.validatePutAuthenticationPayload(request.payload);

    // verify refresh token
    // this method will throw an error if the refresh token is not valid
    await this._authenticationsService.verifyRefreshToken(refreshToken);

    // get user id from refresh token for generating new access token
    const { id } = this._jwtTokenize.verifyRefreshToken(refreshToken);

    // generate new access token
    const accessToken = this._jwtTokenize.generateAccessToken({ id });

    // return response with new access token
    return {
      status: 'success',
      message: 'authentication berhasil diperbarui',
      data: {
        accessToken,
      },
    };
  }

  // This method will be used to handle DELETE request to /authentications resource.
  // This method will be used to handle logout request.
  async deleteAuthenticationHandler(request) {
    // validate the payload.
    // Same principle as the previous method.
    const { refreshToken } = this._validator.validateDeleteAuthenticationPayload(request.payload);

    // verify refresh token
    // this method will throw an error if the refresh token is not valid
    await this._authenticationsService.verifyRefreshToken(refreshToken);

    // delete refresh token from database
    await this._authenticationsService.deleteRefreshToken(refreshToken);

    // return response
    return {
      status: 'success',
      message: 'authentication berhasil dihapus',
      data: {},
    };
  }
}

module.exports = AuthenticationsHandler;
