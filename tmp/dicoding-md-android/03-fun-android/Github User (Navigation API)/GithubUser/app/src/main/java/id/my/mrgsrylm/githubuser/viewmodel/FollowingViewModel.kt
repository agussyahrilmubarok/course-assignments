package id.my.mrgsrylm.githubuser.viewmodel

import android.util.Log
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import id.my.mrgsrylm.githubuser.data.remote.response.UserItemResponse
import id.my.mrgsrylm.githubuser.data.remote.retrofit.ApiConfig
import id.my.mrgsrylm.githubuser.view.fragment.FollowingFragment
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class FollowingViewModel : ViewModel() {

    private val _following = MutableLiveData<List<UserItemResponse>>()
    val following: LiveData<List<UserItemResponse>> = _following

    private val _isLoading = MutableLiveData<Boolean>()
    val isLoading: LiveData<Boolean> = _isLoading

    internal fun findFollowing(login: String) {
        _isLoading.value = true
        val client = ApiConfig.getApiService().getFollowing(login)
        client.enqueue(object : Callback<List<UserItemResponse>> {
            override fun onResponse(
                call: Call<List<UserItemResponse>>,
                response: Response<List<UserItemResponse>>
            ) {
                _isLoading.value = false
                if (response.isSuccessful) {
                    val responseBody = response.body()
                    if (responseBody != null) {
                        _following.value = response.body()
                    }
                } else {
                    Log.e(FollowingFragment.TAG, "onFailure: ${response.message()}")
                }
            }

            override fun onFailure(call: Call<List<UserItemResponse>>, t: Throwable) {
                _isLoading.value = false
                Log.e(FollowingFragment.TAG, "onFailure: ${t.message}")
            }
        })
    }
}