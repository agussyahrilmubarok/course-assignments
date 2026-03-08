package id.my.mrgsrylm.githubuser.viewmodel

import android.content.ContentValues
import android.util.Log
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import id.my.mrgsrylm.githubuser.data.remote.response.UserDetailResponse
import id.my.mrgsrylm.githubuser.data.remote.retrofit.ApiConfig
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class UserDetailViewModel : ViewModel() {

    private val _detail = MutableLiveData<UserDetailResponse>()
    val detail: LiveData<UserDetailResponse> = _detail

    private val _isLoading = MutableLiveData<Boolean>()
    val isLoading: LiveData<Boolean> = _isLoading

    private val _snackBarText = MutableLiveData<String>()
    val snackBarText: LiveData<String> = _snackBarText

    internal fun getDetail(login: String) {
        _isLoading.value = true

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