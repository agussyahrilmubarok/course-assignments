package id.my.mrgsrylm.githubuser.view.adapter

import android.content.Intent
import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.DiffUtil
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import id.my.mrgsrylm.githubuser.data.local.entity.FavoriteUser
import id.my.mrgsrylm.githubuser.databinding.ItemUserTileBinding
import id.my.mrgsrylm.githubuser.di.DiffCallback
import id.my.mrgsrylm.githubuser.view.activity.UserDetailActivity

class FavoriteUserAdapter : RecyclerView.Adapter<FavoriteUserAdapter.MyViewHolder>() {

    private val listFavorites = ArrayList<FavoriteUser>()

    fun setListFavorites(favorite: List<FavoriteUser>) {
        val diffCallback = DiffCallback(this.listFavorites, favorite)
        val diffResult = DiffUtil.calculateDiff(diffCallback)
        this.listFavorites.clear()
        this.listFavorites.addAll(favorite)
        diffResult.dispatchUpdatesTo(this)
    }

    inner class MyViewHolder(private val binding: ItemUserTileBinding) :
        RecyclerView.ViewHolder(binding.root) {
        fun bind(user: FavoriteUser) {
            with(binding) {
                Glide.with(itemView.context)
                    .load(user.avatarUrl)
                    .into(binding.ivUserPhoto)
                tvUserName.text = user.username
                tvUserLink.text = user.htmlUrl
                itemView.setOnClickListener {
                    val intent = Intent(itemView.context, UserDetailActivity::class.java)
                    intent.putExtra(UserDetailActivity.EXTRA_USER, user.username)
                    intent.putExtra(UserDetailActivity.EXTRA_ACTIONBAR, user.username)
                    itemView.context.startActivity(intent)
                }
            }

        }
    }

    override fun getItemCount(): Int = listFavorites.size

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): MyViewHolder {
        val binding =
            ItemUserTileBinding.inflate(LayoutInflater.from(parent.context), parent, false)
        return MyViewHolder(binding)
    }

    override fun onBindViewHolder(holder: MyViewHolder, position: Int) {
        val favoriteUser = listFavorites[position]
        holder.bind(favoriteUser)
    }
}