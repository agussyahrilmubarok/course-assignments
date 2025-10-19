package id.my.mrgsrylm.githubuser.viewmodel.user

import android.app.Application
import android.content.ContentValues
import android.util.Log
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import id.my.mrgsrylm.githubuser.data.FavoriteUserRepository
import id.my.mrgsrylm.githubuser.data.local.entity.FavoriteUser
import id.my.mrgsrylm.githubuser.data.remote.response.UserDetailResponse
import id.my.mrgsrylm.githubuser.data.remote.retrofit.ApiConfig
import id.my.mrgsrylm.githubuser.di.Idling
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class UserDetailViewModel(
    application: Application
) : ViewModel() {

    private val favoriteUserRepository: FavoriteUserRepository = FavoriteUserRepository(application)

    private val _detail = MutableLiveData<UserDetailResponse>()
    val detail: LiveData<UserDetailResponse> = _detail

    private val _isLoading = MutableLiveData<Boolean>()
    val isLoading: LiveData<Boolean> = _isLoading

    private val _snackBarText = MutableLiveData<String>()
    val snackBarText: LiveData<String> = _snackBarText

    fun getFavoriteUser(username: String): LiveData<FavoriteUser> =
        favoriteUserRepository.getFavoriteUser(username)

    fun insert(user: FavoriteUser) {
        favoriteUserRepository.insert(user)
        Log.d("FavoriteAddViewModel", "${user.username}; ${user.avatarUrl} added")
    }

    fun delete(user: String) {
        favoriteUserRepository.delete(user)
    }

    internal fun getDetail(login: String) {
        _isLoading.value = true
        Idling.increment()

        val client = ApiConfig.getApiService().getDetailUser(login)
        client.enqueue(object : Callback<UserDetailResponse> {
            override fun onResponse(
                call: Call<UserDetailResponse>,
                response: Response<UserDetailResponse>
            ) {
                _isLoading.value = false
                if (response.isSuccessful) {
                    val responseBody = response.body()
                    if (responseBody != null) {
                        _detail.value = response.body()
                        if (!Idling.idlingResource.isIdleNow) {
                            Idling.decrement()
                        }
                    }
                } else {
                    Log.e(ContentValues.TAG, "onFailure: ${response.message()}")
                    _snackBarText.value = "An error has occurred!"
                }
            }

            override fun onFailure(call: Call<UserDetailResponse>, t: Throwable) {
                _isLoading.value = false
                Log.e(ContentValues.TAG, "onFailure: ${t.message}")
                _snackBarText.value = "An error has occurred! Please try again later."
            }
        })
    }
}