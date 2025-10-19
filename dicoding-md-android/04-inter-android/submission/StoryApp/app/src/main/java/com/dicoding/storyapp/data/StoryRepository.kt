package com.dicoding.storyapp.data

import android.util.Log
import androidx.lifecycle.liveData
import com.dicoding.storyapp.common.ResultState
import com.dicoding.storyapp.data.model.ErrorResponse
import com.dicoding.storyapp.data.remote.ApiService
import com.google.gson.Gson
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.RequestBody.Companion.asRequestBody
import okhttp3.RequestBody.Companion.toRequestBody
import retrofit2.HttpException
import java.io.File

class StoryRepository private constructor(
    private val service: ApiService
) {

    fun addStory(
        token: String,
        description: String,
        photo: File,
        lat: Float,
        lon: Float
    ) = liveData {
        emit(ResultState.Loading)
        try {
            val response = service.addStory(
                authorizationHeader = token,
                description = description.toRequestBody("text/plain".toMediaType()),
                photo = MultipartBody.Part.createFormData(
                    "photo",
                    photo.name,
                    photo.asRequestBody("image/jpeg".toMediaType())
                ),
                lat = lat.toString().toRequestBody("text/plain".toMediaType()),
                lon = lon.toString().toRequestBody("text/plain".toMediaType())
            )
            Log.d("repo-addStory", "$response")
            emit(ResultState.Success(response))
        } catch (e: HttpException) {
            val errorBody = e.response()?.errorBody()?.string()
            val errorResponse = Gson().fromJson(errorBody, ErrorResponse::class.java)
            Log.e("repo-addStory", errorBody.toString())
            emit(ResultState.Error(errorResponse.message))
        }
    }

    fun getAllStories(token: String) = liveData {
        emit(ResultState.Loading)
        try {

            val response = service.getAllStories(token)
            Log.d("repo-getAllStories", "$response")
            emit(ResultState.Success(response))
        } catch (e: HttpException) {
            val errorBody = e.response()?.errorBody()?.string()
            val errorResponse = Gson().fromJson(errorBody, ErrorResponse::class.java)
            Log.e("repo-getAllStories", errorBody.toString())
            emit(ResultState.Error(errorResponse.message))
        }
    }

    fun getAStory(token: String, id: String) = liveData {
        emit(ResultState.Loading)
        try {
            val response = service.getAStory(token, id)
            Log.d("repo-getAStory", "$response")
            emit(ResultState.Success(response))
        } catch (e: HttpException) {
            val errorBody = e.response()?.errorBody()?.string()
            val errorResponse = Gson().fromJson(errorBody, ErrorResponse::class.java)
            Log.e("repo-getAStory", errorBody.toString())
            emit(ResultState.Error(errorResponse.message))
        }
    }

    companion object {
        @Volatile
        private var instance: StoryRepository? = null
        fun getInstance(apiService: ApiService) =
            instance ?: synchronized(this) {
                instance ?: StoryRepository(apiService)
            }.also { instance = it }
    }

}