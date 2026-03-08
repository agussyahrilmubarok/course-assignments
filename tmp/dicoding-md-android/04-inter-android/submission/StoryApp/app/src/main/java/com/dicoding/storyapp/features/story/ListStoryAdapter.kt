package com.dicoding.storyapp.features.story

import android.content.Intent
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.LinearLayout
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import com.bumptech.glide.load.engine.DiskCacheStrategy
import com.bumptech.glide.request.RequestOptions
import com.dicoding.storyapp.R
import com.dicoding.storyapp.data.model.StoryResponse

class ListStoryAdapter(
    private val stories: List<StoryResponse>,
) : RecyclerView.Adapter<ListStoryAdapter.ViewHolder>() {

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ViewHolder {
        val itemView = LayoutInflater.from(parent.context)
            .inflate(R.layout.item_story, parent, false)
        return ViewHolder(itemView)
    }

    override fun onBindViewHolder(holder: ViewHolder, position: Int) {
        val data = stories[position]

        holder.tvNameProfile.text = data.name
        holder.tvStoryDesc.text = data.description

        Glide.with(holder.itemView.context)
            .load(data.photoUrl)
            .apply(RequestOptions.diskCacheStrategyOf(DiskCacheStrategy.ALL))
            .into(holder.ivStoryImage)

        holder.layoutItem.setOnClickListener {
            val intent = Intent(holder.itemView.context, DetailStoryActivity::class.java)
            intent.putExtra(DetailStoryActivity.EXTRA_STORY_ID, data.id)
            holder.itemView.context.startActivity(intent)
        }
    }

    override fun getItemCount(): Int {
        return stories.size
    }

    inner class ViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        val tvNameProfile: TextView = itemView.findViewById(R.id.tvNameProfile)
        val tvStoryDesc: TextView = itemView.findViewById(R.id.tvStoryDescription)
        val ivStoryImage: ImageView = itemView.findViewById(R.id.ivStoryImage)
        val layoutItem: LinearLayout = itemView.findViewById(R.id.layoutStoryItem)
    }
}