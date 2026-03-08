package id.my.mrgsrylm.githubuser.data

import android.app.Application
import androidx.lifecycle.LiveData
import id.my.mrgsrylm.githubuser.data.local.entity.FavoriteUser
import id.my.mrgsrylm.githubuser.data.local.room.FavoriteUserDao
import id.my.mrgsrylm.githubuser.data.local.room.FavoriteUserRoomDatabase
import java.util.concurrent.ExecutorService
import java.util.concurrent.Executors

class FavoriteUserRepository(
    application: Application
) {

    private val favoriteUserDao: FavoriteUserDao
    private val executorService: ExecutorService = Executors.newSingleThreadExecutor()

    init {
        val db = FavoriteUserRoomDatabase.getInstance(application)
        favoriteUserDao = db.favoriteUserDao()
    }

    fun getAll(): LiveData<List<FavoriteUser>> = favoriteUserDao.getAll()

    fun getFavoriteUser(username: String): LiveData<FavoriteUser> =
        favoriteUserDao.getFavoriteUser(username)

    fun insert(favoriteUser: FavoriteUser) {
        executorService.execute { favoriteUserDao.insert(favoriteUser) }
    }

    fun delete(favoriteUser: String) {
        executorService.execute { favoriteUserDao.deleteFavoriteUser(favoriteUser) }
    }
}