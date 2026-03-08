package id.my.mrgsrylm.githubuser.di

import androidx.recyclerview.widget.DiffUtil
import id.my.mrgsrylm.githubuser.data.local.entity.FavoriteUser

class DiffCallback(
    private val oldListUser: List<FavoriteUser>,
    private val newListUser: List<FavoriteUser>
) : DiffUtil.Callback() {
    
    override fun getOldListSize(): Int = oldListUser.size

    override fun getNewListSize(): Int = newListUser.size

    override fun areItemsTheSame(oldItemPosition: Int, newItemPosition: Int): Boolean {
        return oldListUser[oldItemPosition].username == newListUser[newItemPosition].username
    }

    override fun areContentsTheSame(oldItemPosition: Int, newItemPosition: Int): Boolean {
        val oldEmployee = oldListUser[oldItemPosition]
        val newEmployee = oldListUser[newItemPosition]
        return oldEmployee.avatarUrl == newEmployee.avatarUrl
    }

}