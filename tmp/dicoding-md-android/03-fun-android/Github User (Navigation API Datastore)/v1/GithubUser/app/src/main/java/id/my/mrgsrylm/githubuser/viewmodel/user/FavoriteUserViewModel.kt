package id.my.mrgsrylm.githubuser.viewmodel.user

import android.app.Application
import androidx.lifecycle.LiveData
import androidx.lifecycle.ViewModel
import id.my.mrgsrylm.githubuser.data.FavoriteUserRepository
import id.my.mrgsrylm.githubuser.data.local.entity.FavoriteUser

class FavoriteUserViewModel(application: Application) : ViewModel() {

    private val favoriteUserRepository: FavoriteUserRepository = FavoriteUserRepository(application)

    fun getUser(): LiveData<List<FavoriteUser>> = favoriteUserRepository.getAll()

}