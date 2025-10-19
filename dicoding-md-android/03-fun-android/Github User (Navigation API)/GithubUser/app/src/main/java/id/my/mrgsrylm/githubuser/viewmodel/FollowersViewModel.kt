package id.my.mrgsrylm.githubuser.viewmodel

import android.util.Log
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import id.my.mrgsrylm.githubuser.data.remote.response.UserItemResponse
import id.my.mrgsrylm.githubuser.data.remote.retrofit.ApiConfig
import id.my.mrgsrylm.githubuser.view.fragment.FollowersFragment
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class FollowersViewModel : ViewModel() {

    private val _followers = MutableLiveData<List<UserItemResponse>>()
    val followers: LiveData<List<UserItemResponse>> = _followers

    private val _isLoading = MutableLiveData<Boolean>()
    val isLoading: LiveData<Boolean> = _isLoading

    internal fun findFollowers(login: String) {
        _isLoading.value = true
        val client = ApiConfig.getApiService().getFollowers(login)
        client.enqueue(object : Callback<List<UserItemResponse>> {
            override fun onResponse(
                call: Call<List<UserItemResponse>>,
                response: Response<List<UserItemResponse>>
            ) {
                _isLoading.value = false
                if (response.isSuccessful) {
                    val responseBody = response.body()
                    if (responseBody != null) {
                        _followers.value = response.body()
                    }
                } else {
                    Log.e(FollowersFragment.TAG, "onFailure: ${response.message()}")
                }
            }

            override fun onFailure(call: Call<List<UserItemResponse>>, t: Throwable) {
                _isLoading.value = false
                Log.e(FollowersFragment.TAG, "onFailure: ${t.message}")
            }
        })
    }
}