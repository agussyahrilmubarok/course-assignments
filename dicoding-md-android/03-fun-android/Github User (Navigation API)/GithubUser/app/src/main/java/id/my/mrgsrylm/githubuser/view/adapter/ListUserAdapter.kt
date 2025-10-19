package id.my.mrgsrylm.githubuser.view.adapter

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import id.my.mrgsrylm.githubuser.data.remote.response.UserItemResponse
import id.my.mrgsrylm.githubuser.databinding.ItemUserTileBinding

class ListUserAdapter(
    private val listUser: List<UserItemResponse>
) : RecyclerView.Adapter<ListUserAdapter.ViewHolder>() {

    private lateinit var onItemClickCallback: OnItemClickCallback

    interface OnItemClickCallback {
        fun onItemClicked(data: UserItemResponse)
    }

    fun setOnItemClickCallback(onItemClickCallback: OnItemClickCallback) {
        this.onItemClickCallback = onItemClickCallback
    }

    override fun onCreateViewHolder(viewGroup: ViewGroup, viewType: Int): ViewHolder {
        val viewBinding =
            ItemUserTileBinding.inflate(LayoutInflater.from(viewGroup.context), viewGroup, false)
        return ViewHolder(viewBinding)
    }

    override fun onBindViewHolder(viewHolder: ViewHolder, position: Int) {
        val user = listUser[position]
        viewHolder.binding.tvUserName.text = user.login
        viewHolder.binding.tvUserLink.text = user.htmlUrl

        Glide.with(viewHolder.itemView.context)
            .load(user.avatarUrl)
            .into(viewHolder.binding.ivUserPhoto)

        viewHolder.itemView.setOnClickListener {
            onItemClickCallback.onItemClicked(listUser[viewHolder.adapterPosition])
        }
    }

    override fun getItemCount() = listUser.size

    class ViewHolder(var binding: ItemUserTileBinding) : RecyclerView.ViewHolder(binding.root)
}