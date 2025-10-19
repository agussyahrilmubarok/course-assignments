package id.my.mrgsrylm.githubuser.data.local.room

import androidx.lifecycle.LiveData
import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import id.my.mrgsrylm.githubuser.data.local.entity.FavoriteUser

@Dao
interface FavoriteUserDao {
    @Insert(onConflict = OnConflictStrategy.IGNORE)
    fun insert(favoriteUser: FavoriteUser)

    @Query("SELECT * from favoriteUser ORDER BY username ASC")
    fun getAll(): LiveData<List<FavoriteUser>>

    @Query("SELECT * FROM favoriteUser WHERE username = :username")
    fun getByUsername(username: String): LiveData<FavoriteUser>

    @Query("DELETE FROM favoriteUser WHERE username = :username")
    fun deleteByUsername(username: String)
}