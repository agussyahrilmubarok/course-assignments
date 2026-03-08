package id.my.mrgsrylm.githubuser.helper

import androidx.recyclerview.widget.DiffUtil
import id.my.mrgsrylm.githubuser.data.local.entity.FavoriteUser

class FavoriteUserDiffCallback(
    private val mOldFavoriteUserList: List<FavoriteUser>,
    private val mNewFavoriteUserList: List<FavoriteUser>
) : DiffUtil.Callback() {

    override fun getOldListSize(): Int = mOldFavoriteUserList.size
    override fun getNewListSize(): Int = mNewFavoriteUserList.size

    override fun areItemsTheSame(oldItemPosition: Int, newItemPosition: Int): Boolean {
        return mOldFavoriteUserList[oldItemPosition].username == mNewFavoriteUserList[newItemPosition].username
    }

    override fun areContentsTheSame(oldItemPosition: Int, newItemPosition: Int): Boolean {
        val old = mOldFavoriteUserList[oldItemPosition]
        val new = mNewFavoriteUserList[newItemPosition]
        return old.avatarUrl == new.avatarUrl
    }
}