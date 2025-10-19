package id.my.mrgsrylm.githubuser.view.adapter

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import id.my.mrgsrylm.githubuser.data.remote.response.UserItemResponse
import id.my.mrgsrylm.githubuser.databinding.ItemUserTileBinding

class ListFollowersAdapter(
    private val listUser: List<UserItemResponse>
) : RecyclerView.Adapter<ListFollowersAdapter.ViewHolder>() {

    override fun onCreateViewHolder(viewGroup: ViewGroup, viewType: Int): ViewHolder {
        val binding =
            ItemUserTileBinding.inflate(LayoutInflater.from(viewGroup.context), viewGroup, false)
        return ViewHolder(binding)
    }

    override fun onBindViewHolder(viewHolder: ViewHolder, position: Int) {
        val user = listUser[position]
        viewHolder.binding.tvUserName.text = user.login
        viewHolder.binding.tvUserLink.text = user.htmlUrl

        Glide.with(viewHolder.itemView.context)
            .load(user.avatarUrl)
            .into(viewHolder.binding.ivUserPhoto)
    }

    override fun getItemCount() = listUser.size

    class ViewHolder(var binding: ItemUserTileBinding) : RecyclerView.ViewHolder(binding.root)
}