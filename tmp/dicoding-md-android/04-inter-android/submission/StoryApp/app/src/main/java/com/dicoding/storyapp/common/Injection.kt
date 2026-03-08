package com.dicoding.storyapp.common

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import com.dicoding.storyapp.data.AuthRepository
import com.dicoding.storyapp.data.CommonRepository
import com.dicoding.storyapp.data.StoryRepository
import com.dicoding.storyapp.data.local.UserPreferences
import com.dicoding.storyapp.data.remote.ApiConfig

object Injection {
    fun provideAuthRepository(): AuthRepository {
        return AuthRepository.getInstance(ApiConfig.getApiService())
    }

    fun provideStoryRepository(): StoryRepository {
        return StoryRepository.getInstance(ApiConfig.getApiService())
    }

    fun provideUserPreference(dataStore: DataStore<Preferences>): UserPreferences {
        return UserPreferences.getInstance(dataStore)
    }

    fun provideCommonRepository(context: Context): CommonRepository {
        return CommonRepository(context)
    }
}