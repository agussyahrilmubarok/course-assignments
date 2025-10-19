const { createPool } = require('./pool');
const NotFoundError = require('../../exceptions/NotFoundError');

// this class will be used to handle all the database operations related to festivals table
class FestivalsService {
  constructor() {
    // create database connection pool
    this._pool = createPool();
  }

  /**
   * this method will be used to get all festivals from database
   */
  async getFestivals() {
    const result = await this._pool.query('SELECT * FROM festivals');
    return result.rows.map((festival) => ({ ...festival, price: Number(festival.price) }));
  }

  /**
   * this method will be used to get festival by id from database
   */
  async getFestival(id) {
    const query = {
      text: 'SELECT * FROM festivals WHERE id = $1',
      values: [id],
    };

    const result = await this._pool.query(query);

    if (!result.rows.length) {
      throw new NotFoundError('festival tidak ditemukan');
    }

    const [festival] = result.rows;

    return {
      ...festival,
      price: Number(festival.price),
    };
  }

  /**
   * this method will be used to verify festival availability
   */
  async isFestivalAvailable(id) {
    const query = {
      text: 'SELECT id FROM festivals WHERE id = $1',
      values: [id],
    };

    const { rows } = await this._pool.query(query);

    return rows.length > 0;
  }
}

module.exports = FestivalsService;
