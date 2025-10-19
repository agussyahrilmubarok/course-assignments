const InvariantError = require('../../exceptions/InvariantError');

// this handler will be used to handle the request from client related to users feature
class UsersHandler {
  constructor(usersService, validator) {
    this._usersService = usersService;
    this._validator = validator;
  }

  // this handler will be used to handle POST /users
  async postUserHandler(request, h) {
    const user = this._validator.validatePostUserPayload(request.payload);

    // checking availability of email
    const isEmailAvailable = await this._usersService.isEmailAvailable(user.email);

    // if the email is not available, then throw an error
    if (!isEmailAvailable) {
      throw new InvariantError('email sudah digunakan');
    }

    // persist user to database
    const createdUser = await this._usersService.persistUsers(user);

    // response with persisted user
    return h.response({
      status: 'success',
      message: 'user berhasil ditambahkan',
      data: {
        createdUser,
      },
    }).code(201);
  }
}

module.exports = UsersHandler;
