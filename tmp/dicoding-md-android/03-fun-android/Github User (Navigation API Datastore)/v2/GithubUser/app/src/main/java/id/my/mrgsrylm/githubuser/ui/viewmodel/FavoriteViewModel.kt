package id.my.mrgsrylm.githubuser.ui.viewmodel

import android.app.Application
import androidx.lifecycle.LiveData
import androidx.lifecycle.ViewModel
import id.my.mrgsrylm.githubuser.data.FavoriteUserRepository
import id.my.mrgsrylm.githubuser.data.local.entity.FavoriteUser

class FavoriteViewModel(application: Application) : ViewModel() {

    private val mFavoriteUserRepository: FavoriteUserRepository =
        FavoriteUserRepository(application)

    fun getAll(): LiveData<List<FavoriteUser>> = mFavoriteUserRepository.getAll()
}