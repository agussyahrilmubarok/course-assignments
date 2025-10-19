const { nanoid } = require('nanoid');
const festivals = require('../dummy-data.json');

exports.up = (pgm) => {
  pgm.sql(`
  INSERT INTO festivals(id, title, description, facilities, refundable, banner, price) VALUES ${festivals.map((festival) => `('festival-${nanoid(16)}', '${festival.title}', '${festival.description}', '{${festival.facilities.join(',')}}', ${festival.refundable}, '${festival.banner}', ${festival.price})`).join(',\n')};`);
};

exports.down = (pgm) => {
  pgm.sql('DELETE FROM festivals WHERE id LIKE \'festival-%\';');
};
