package com.dicoding.storyapp.features.story

import android.os.Bundle
import android.view.View
import androidx.activity.viewModels
import androidx.appcompat.app.AppCompatActivity
import com.bumptech.glide.Glide
import com.bumptech.glide.load.engine.DiskCacheStrategy
import com.bumptech.glide.request.RequestOptions
import com.dicoding.storyapp.common.ResultState
import com.dicoding.storyapp.databinding.ActivityDetailStoryBinding
import com.dicoding.storyapp.features.ViewModelFactory

class DetailStoryActivity : AppCompatActivity() {

    private lateinit var binding: ActivityDetailStoryBinding

    private val viewModel by viewModels<StoryViewModel> {
        ViewModelFactory.getInstance(applicationContext)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityDetailStoryBinding.inflate(layoutInflater)
        setContentView(binding.root)
        setView()
    }

    private fun setView() {
        val storyId = intent.getStringExtra(EXTRA_STORY_ID)
        if (storyId != null) {
            viewModel.getAStory(storyId).observe(this) {
                when (it) {
                    is ResultState.Loading -> {
                        showLoading(true)
                    }

                    is ResultState.Error -> {
                        showLoading(false)
                        com.dicoding.storyapp.common.showDialog(
                            this,
                            "Peringatan",
                            "Story tidak ditemukan."
                        )
                    }

                    is ResultState.Success -> {
                        showLoading(false)
                        binding.tvNameProfile.text = it.data.story.name
                        binding.tvStoryDescription.text = it.data.story.description

                        Glide.with(this)
                            .load(it.data.story.photoUrl)
                            .apply(RequestOptions.diskCacheStrategyOf(DiskCacheStrategy.ALL))
                            .into(binding.ivStoryImage)
                    }
                }
            }
        }
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
    }

    companion object {
        const val EXTRA_STORY_ID = "extra_story_id"
    }
}