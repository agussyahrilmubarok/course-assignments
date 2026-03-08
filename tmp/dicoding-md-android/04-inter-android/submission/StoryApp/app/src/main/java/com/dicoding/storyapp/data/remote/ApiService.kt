package com.dicoding.storyapp.data.remote

import com.dicoding.storyapp.data.model.DetailStoryResponse
import com.dicoding.storyapp.data.model.ListStoryResponse
import com.dicoding.storyapp.data.model.LoginRequest
import com.dicoding.storyapp.data.model.LoginResponse
import com.dicoding.storyapp.data.model.MessageResponse
import com.dicoding.storyapp.data.model.RegisterRequest
import com.dicoding.storyapp.data.model.RegisterResponse
import okhttp3.MultipartBody
import okhttp3.RequestBody
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.Multipart
import retrofit2.http.POST
import retrofit2.http.Part
import retrofit2.http.Path
import retrofit2.http.Query

interface ApiService {
    @POST("register")
    suspend fun register(@Body param: RegisterRequest): RegisterResponse

    @POST("login")
    suspend fun login(@Body param: LoginRequest): LoginResponse

    @POST("stories")
    @Multipart
    suspend fun addStory(
        @Header("Authorization") authorizationHeader: String,
        @Part("description") description: RequestBody,
        @Part photo: MultipartBody.Part,
        @Part("lat") lat: RequestBody,
        @Part("lon") lon: RequestBody,
    ): MessageResponse

    @GET("stories")
    suspend fun getAllStories(
        @Header("Authorization") authorizationHeader: String,
        @Query("page") page: Int? = 1,
        @Query("size") size: Int? = 10,
        @Query("location") location: Int? = 0
    ): ListStoryResponse

    @GET("stories/{id}")
    suspend fun getAStory(
        @Header("Authorization") authorizationHeader: String,
        @Path("id") id: String
    ): DetailStoryResponse
}