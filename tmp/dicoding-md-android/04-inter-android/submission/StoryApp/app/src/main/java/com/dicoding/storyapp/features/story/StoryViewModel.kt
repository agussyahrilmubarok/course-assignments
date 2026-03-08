package com.dicoding.storyapp.features.story

import androidx.lifecycle.ViewModel
import com.dicoding.storyapp.data.StoryRepository
import com.dicoding.storyapp.data.local.UserPreferences
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import java.io.File

class StoryViewModel(
    private val storyRepo: StoryRepository,
    private val userPref: UserPreferences
) : ViewModel() {

    fun addStory(
        description: String,
        photo: File,
        lat: Float,
        lon: Float
    ) = storyRepo.addStory(getToken(), description, photo, lat, lon)

    fun getAStory(id: String) = storyRepo.getAStory(getToken(), id)

    private fun getToken(): String {
        val token = runBlocking {
            userPref.getToken().first()
        }
        return "Bearer $token"
    }
}