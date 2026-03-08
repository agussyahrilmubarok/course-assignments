package id.my.mrgsrylm.githubuser.data.remote.response

import com.google.gson.annotations.SerializedName

data class GithubResponse(
    @field:SerializedName("total_count")
    val totalCount: Int,

    @field:SerializedName("incomplete_results")
    val incompleteResults: Boolean,

    @field:SerializedName("items")
    val items: List<UserItem>
)

data class UserItem(
    @field:SerializedName("login")
    val login: String,

    @field:SerializedName("id")
    val id: Int,

    @field:SerializedName("node_id")
    val nodeId: String,

    @field:SerializedName("avatar_url")
    val avatarUrl: String,

    @field:SerializedName("gravatar_id")
    val gravatarId: String,

    @field:SerializedName("url")
    val url: String,

    @field:SerializedName("html_url")
    val htmlUrl: String,

    @field:SerializedName("followers_url")
    val followersUrl: String,

    @field:SerializedName("following_url")
    val followingUrl: String,

    @field:SerializedName("gists_url")
    val gistsUrl: String,

    @field:SerializedName("starred_url")
    val starredUrl: String,

    @field:SerializedName("subscriptions_url")
    val subscriptionsUrl: String,

    @field:SerializedName("organizations_url")
    val organizationsUrl: String,

    @field:SerializedName("repos_url")
    val reposUrl: String,

    @field:SerializedName("events_url")
    val eventsUrl: String,

    @field:SerializedName("received_events_url")
    val receivedEventsUrl: String,

    @field:SerializedName("type")
    val type: String,

    @field:SerializedName("site_admin")
    val siteAdmin: Boolean
)

data class UserDetail(
    @field:SerializedName("login")
    val login: String,

    @field:SerializedName("id")
    val id: Int,

    @field:SerializedName("node_id")
    val nodeId: String,

    @field:SerializedName("avatar_url")
    val avatarUrl: String,

    @field:SerializedName("gravatar_id")
    val gravatarId: String,

    @field:SerializedName("url")
    val url: String,

    @field:SerializedName("html_url")
    val htmlUrl: String,

    @field:SerializedName("followers_url")
    val followersUrl: String,

    @field:SerializedName("following_url")
    val followingUrl: String,

    @field:SerializedName("gists_url")
    val gistsUrl: String,

    @field:SerializedName("starred_url")
    val starredUrl: String,

    @field:SerializedName("subscriptions_url")
    val subscriptionsUrl: String,

    @field:SerializedName("organizations_url")
    val organizationsUrl: String,

    @field:SerializedName("repos_url")
    val reposUrl: String,

    @field:SerializedName("events_url")
    val eventsUrl: String,

    @field:SerializedName("received_events_url")
    val receivedEventsUrl: String,

    @field:SerializedName("type")
    val type: String,

    @field:SerializedName("site_admin")
    val siteAdmin: Boolean,

    @field:SerializedName("name")
    val name: String,

    @field:SerializedName("company")
    val company: String,

    @field:SerializedName("blog")
    val blog: String,

    @field:SerializedName("location")
    val location: String,

    @field:SerializedName("email")
    val email: String,

    @field:SerializedName("hireable")
    val hireable: String,

    @field:SerializedName("bio")
    val bio: String,

    @field:SerializedName("twitter_username")
    val twitterUsername: String,

    @field:SerializedName("public_repos")
    val publicRepos: Int,

    @field:SerializedName("public_gists")
    val publicGists: Int,

    @field:SerializedName("followers")
    val followers: Int,

    @field:SerializedName("following")
    val following: Int,

    @field:SerializedName("created_at")
    val createdAt: String,

    @field:SerializedName("updated_at")
    val updatedAt: String,
)