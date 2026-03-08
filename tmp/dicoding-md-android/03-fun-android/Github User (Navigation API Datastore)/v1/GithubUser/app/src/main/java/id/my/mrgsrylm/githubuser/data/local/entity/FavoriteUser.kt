package id.my.mrgsrylm.githubuser.data.local.entity

import android.os.Parcelable
import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey
import kotlinx.parcelize.Parcelize

@Parcelize
@Entity
class FavoriteUser(
    @field:PrimaryKey(autoGenerate = false)
    @field:ColumnInfo(name = "login")
    var username: String = "",

    @field:ColumnInfo(name = "htmlUrl")
    var htmlUrl: String? = "",

    @field:ColumnInfo(name = "avatarUrl")
    var avatarUrl: String? = null
) : Parcelable