package id.my.mrgsrylm.githubuser.data.remote.retrofit

import id.my.mrgsrylm.githubuser.data.remote.response.GithubResponse
import id.my.mrgsrylm.githubuser.data.remote.response.UserDetail
import id.my.mrgsrylm.githubuser.data.remote.response.UserItem
import retrofit2.Call
import retrofit2.http.GET
import retrofit2.http.Path
import retrofit2.http.Query

interface ApiService {
    // @Headers("Authorization: token ${BuildConfig.API_KEY}")
    @GET("search/users")
    fun search(
        @Query("q") username: String?
    ): Call<GithubResponse>

    @GET("users/{username}")
    fun getDetailUser(
        @Path("username") username: String
    ): Call<UserDetail>

    @GET("users/{username}/followers")
    fun getListFollowers(
        @Path("username") username: String
    ): Call<List<UserItem>>

    @GET("users/{username}/following")
    fun getListFollowing(
        @Path("username") username: String
    ): Call<List<UserItem>>
}