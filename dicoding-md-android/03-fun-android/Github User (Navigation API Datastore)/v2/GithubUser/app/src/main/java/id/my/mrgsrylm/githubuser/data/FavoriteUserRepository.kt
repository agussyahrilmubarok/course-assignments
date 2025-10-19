package id.my.mrgsrylm.githubuser.data

import android.app.Application
import androidx.lifecycle.LiveData
import id.my.mrgsrylm.githubuser.data.local.entity.FavoriteUser
import id.my.mrgsrylm.githubuser.data.local.room.FavoriteUserDao
import id.my.mrgsrylm.githubuser.data.local.room.FavoriteUserRoomDatabase
import java.util.concurrent.ExecutorService
import java.util.concurrent.Executors

class FavoriteUserRepository(application: Application) {

    private val mFavoriteUserDao: FavoriteUserDao
    private val executorService: ExecutorService = Executors.newSingleThreadExecutor()

    init {
        val db = FavoriteUserRoomDatabase.getDatabase(application)
        mFavoriteUserDao = db.favoriteUserDao()
    }

    fun insert(favoriteUser: FavoriteUser) {
        executorService.execute { mFavoriteUserDao.insert(favoriteUser) }
    }

    fun getAll(): LiveData<List<FavoriteUser>> {
        return mFavoriteUserDao.getAll()
    }

    fun getByUsername(username: String): LiveData<FavoriteUser> {
        return mFavoriteUserDao.getByUsername(username)
    }

    fun deleteByUsername(username: String) {
        executorService.execute { mFavoriteUserDao.deleteByUsername(username) }
    }
}