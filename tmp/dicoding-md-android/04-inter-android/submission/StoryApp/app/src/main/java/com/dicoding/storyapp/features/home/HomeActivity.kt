package com.dicoding.storyapp.features.home

import android.content.Intent
import android.os.Bundle
import android.view.View
import androidx.activity.viewModels
import androidx.appcompat.app.AppCompatActivity
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.dicoding.storyapp.common.ResultState
import com.dicoding.storyapp.common.showDialog
import com.dicoding.storyapp.databinding.ActivityMainBinding
import com.dicoding.storyapp.features.ViewModelFactory
import com.dicoding.storyapp.features.login.LoginActivity
import com.dicoding.storyapp.features.story.AddStoryActivity
import com.dicoding.storyapp.features.story.ListStoryAdapter


class HomeActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding

    private val viewModel by viewModels<HomeViewModel> {
        ViewModelFactory.getInstance(applicationContext)
    }

    private var isExpanded = false

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)
        setAction()
        setView()
    }

    private fun setView() {
        getAllStories()
    }

    private fun setAction() {
        binding.fabMenu.setOnClickListener {
            toggleFabMenu()
        }

        binding.fabAddStory.setOnClickListener {
            val moveIntent = Intent(this, AddStoryActivity::class.java)
            startActivity(moveIntent)
            finish()
        }

        binding.fabLeave.setOnClickListener {
            viewModel.logout()
            val moveIntent = Intent(this, LoginActivity::class.java)
            startActivity(moveIntent)
            finish()
        }
    }

    private fun getAllStories() {
        viewModel.getAllStories().observe(this) {
            when (it) {
                is ResultState.Loading -> {
                    showLoading(true)
                }

                is ResultState.Error -> {
                    showLoading(false)
                    showDialog(this, "Peringatan", "Tidak ada story baru.")
                }

                is ResultState.Success -> {
                    showLoading(false)
                    val rv: RecyclerView = binding.rvStories
                    val adapter =
                        ListStoryAdapter(it.data.listStory.sortedByDescending { it.createdAt })
                    rv.layoutManager = LinearLayoutManager(this)
                    rv.adapter = adapter
                }
            }
        }
    }

    private fun toggleFabMenu() {
        isExpanded = !isExpanded

        if (isExpanded) {
            expandFabMenu()
        } else {
            collapseFabMenu()
        }
    }

    private fun expandFabMenu() {
        binding.fabAddStory.visibility = View.VISIBLE
        binding.fabLeave.visibility = View.VISIBLE
    }

    private fun collapseFabMenu() {
        binding.fabAddStory.visibility = View.GONE
        binding.fabLeave.visibility = View.GONE
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
    }
}