const path = require('path');
const InvariantError = require('../../exceptions/InvariantError');

// this handler will be used to handle the request from client related to users feature
class UsersHandler {
  constructor(usersService, storageService, validator) {
    this._usersService = usersService;
    this._storageService = storageService;
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

  // this handler will be used to handle GET /users/{id}
  async patchUsersMeIdCardHandler(request) {
    const { id: userId } = request.auth.credentials;

    // validate the payload
    // if the payload is not valid, then the validator will throw an error
    // if the payload is valid, then the validator will return the payload
    const { images } = this._validator.validatePatchUserPayload(request.payload);

    // get information about the image
    const { hapi: meta } = images;

    // upload image to storage
    const filename = await this._storageService.uploadUserIdCard(userId, images, meta);

    // update id card filename in database
    await this._usersService.updateIdCardUser(userId, filename);

    // return response
    return {
      status: 'success',
      message: 'id card berhasil diperbarui',
      data: {},
    };
  }

  // this handler will be used to handle GET /users/me/id-card
  async getUsersMeIdCardHandler(request, h) {
    // get user id from request.auth.credentials
    const { id: userId } = request.auth.credentials;

    // get id card filename from database
    const idCard = await this._usersService.getUserIdCard(userId);

    // return response with file that has been uploaded
    return h.file(path.join(process.cwd(), 'src', 'public', 'uploads', idCard));
  }
}

module.exports = UsersHandler;
