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

    @Query("SELECT * FROM favoriteUser ORDER BY login ASC")
    fun getAll(): LiveData<List<FavoriteUser>>

    @Query("SELECT * FROM favoriteUser WHERE login = :favoriteUser")
    fun getFavoriteUser(favoriteUser: String): LiveData<FavoriteUser>

    @Query("DELETE FROM favoriteUser WHERE login = :favoriteUser")
    fun deleteFavoriteUser(favoriteUser: String)
}