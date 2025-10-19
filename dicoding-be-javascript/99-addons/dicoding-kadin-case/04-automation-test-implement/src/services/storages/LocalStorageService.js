const fs = require('fs');
const path = require('path');

// this class is used to handle file upload to local storage
class LocalStorageService {
  constructor() {
    // set the base upload folder
    this._baseUploadFolder = path.join(process.cwd(), 'src', 'public', 'uploads');

    // create the base upload folder if it doesn't exist
    if (!fs.existsSync(this._baseUploadFolder)) {
      fs.mkdirSync(this._baseUploadFolder, { recursive: true });
    }
  }

  async uploadUserIdCard(userId, file, meta) {
    // create the filename in the format of card-{userId}.{fileExtension}
    // for example: card-user-123.png
    // this is done to prevent filename collision between users
    const filename = `card-${userId}.${meta.filename.split('.').pop()}`;

    // create the output path
    const outputPath = path.join(this._baseUploadFolder, filename);
    // create the file stream and use output path as the destination
    const fileStream = fs.createWriteStream(outputPath);

    // return a promise that will be resolved when the file is successfully uploaded
    return new Promise((resolve, reject) => {
      // if there's an error, reject the promise
      fileStream.on('error', (error) => reject(error));
      // pipe the file to the file stream
      file.pipe(fileStream);
      // when the file is successfully piped, resolve the promise
      file.on('end', () => resolve(filename));
    });
  }
}

module.exports = LocalStorageService;
