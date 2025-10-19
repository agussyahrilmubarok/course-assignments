package id.my.mrgsrylm.githubuser.data.remote.retrofit

import id.my.mrgsrylm.githubuser.data.remote.response.GithubResponse
import id.my.mrgsrylm.githubuser.data.remote.response.UserDetailResponse
import id.my.mrgsrylm.githubuser.data.remote.response.UserItemResponse
import retrofit2.Call
import retrofit2.http.GET
import retrofit2.http.Path
import retrofit2.http.Query

interface ApiService {
    @GET("/search/users")
    fun searchUser(@Query("q") login: String?): Call<GithubResponse>

    @GET("/users/{login}")
    fun getDetailUser(@Path("login") login: String?): Call<UserDetailResponse>

    @GET("/users/{login}/followers")
    fun getFollowers(@Path("login") login: String?): Call<List<UserItemResponse>>

    @GET("/users/{login}/following")
    fun getFollowing(@Path("login") login: String?): Call<List<UserItemResponse>>
}