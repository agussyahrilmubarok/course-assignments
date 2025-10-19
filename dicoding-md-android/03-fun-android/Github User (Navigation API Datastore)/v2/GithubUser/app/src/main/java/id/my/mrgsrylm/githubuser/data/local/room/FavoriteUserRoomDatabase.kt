package id.my.mrgsrylm.githubuser.data.local.room

import android.content.Context
import androidx.room.Database
import androidx.room.Room
import androidx.room.RoomDatabase
import id.my.mrgsrylm.githubuser.data.local.entity.FavoriteUser

@Database(entities = [FavoriteUser::class], version = 1, exportSchema = true)
abstract class FavoriteUserRoomDatabase : RoomDatabase() {

    abstract fun favoriteUserDao(): FavoriteUserDao

    companion object {
        @Volatile
        private var INSTANCE: FavoriteUserRoomDatabase? = null

        @JvmStatic
        fun getDatabase(context: Context): FavoriteUserRoomDatabase {
            if (INSTANCE == null) {
                synchronized(FavoriteUserRoomDatabase::class.java) {
                    INSTANCE = Room.databaseBuilder(
                        context.applicationContext,
                        FavoriteUserRoomDatabase::class.java,
                        "favoriteUser_database"
                    )
                        .allowMainThreadQueries()
                        .build()
                }
            }
            return INSTANCE as FavoriteUserRoomDatabase
        }
    }
}