package com.dicoding.storyapp.data

import android.util.Log
import androidx.lifecycle.liveData
import com.dicoding.storyapp.common.ResultState
import com.dicoding.storyapp.data.model.ErrorResponse
import com.dicoding.storyapp.data.model.LoginRequest
import com.dicoding.storyapp.data.model.RegisterRequest
import com.dicoding.storyapp.data.remote.ApiService
import com.google.gson.Gson
import retrofit2.HttpException

class AuthRepository(
    private val service: ApiService
) {
    fun register(param: RegisterRequest) = liveData {
        emit(ResultState.Loading)
        try {
            val response = service.register(param)
            Log.d("repo-register", "$response")
            emit(ResultState.Success(response))
        } catch (e: HttpException) {
            val errorBody = e.response()?.errorBody()?.string()
            val errorResponse = Gson().fromJson(errorBody, ErrorResponse::class.java)
            Log.e("repo-register", errorBody.toString())
            emit(ResultState.Error(errorResponse.message))
        }
    }

    fun login(param: LoginRequest) = liveData {
        emit(ResultState.Loading)
        try {
            val response = service.login(param)
            Log.d("repo-login", "$response")
            emit(ResultState.Success(response))
        } catch (e: HttpException) {
            val errorBody = e.response()?.errorBody()?.string()
            val errorResponse = Gson().fromJson(errorBody, ErrorResponse::class.java)
            Log.e("repo-login", errorBody.toString())
            emit(ResultState.Error(errorResponse.message))
        }
    }

    companion object {
        @Volatile
        private var instance: AuthRepository? = null
        fun getInstance(apiService: ApiService) =
            instance ?: synchronized(this) {
                instance ?: AuthRepository(apiService)
            }.also { instance = it }
    }
}