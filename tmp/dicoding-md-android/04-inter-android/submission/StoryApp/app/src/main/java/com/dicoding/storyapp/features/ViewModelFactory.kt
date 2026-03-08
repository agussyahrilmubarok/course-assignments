package com.dicoding.storyapp.features

import android.content.Context
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import com.dicoding.storyapp.common.Injection
import com.dicoding.storyapp.common.dataStore
import com.dicoding.storyapp.data.AuthRepository
import com.dicoding.storyapp.data.CommonRepository
import com.dicoding.storyapp.data.StoryRepository
import com.dicoding.storyapp.data.local.UserPreferences
import com.dicoding.storyapp.features.home.HomeViewModel
import com.dicoding.storyapp.features.login.LoginViewModel
import com.dicoding.storyapp.features.register.RegisterViewModel
import com.dicoding.storyapp.features.story.StoryViewModel
import com.dicoding.storyapp.features.welcome.WelcomeViewModel

class ViewModelFactory(
    private val authRepo: AuthRepository,
    private val storyRepo: StoryRepository,
    private val userPref: UserPreferences,
    private val common: CommonRepository
) : ViewModelProvider.NewInstanceFactory() {

    @Suppress("UNCHECKED_CAST")
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        return when {
            modelClass.isAssignableFrom(WelcomeViewModel::class.java) -> {
                WelcomeViewModel(userPref, common) as T
            }

            modelClass.isAssignableFrom(RegisterViewModel::class.java) -> {
                RegisterViewModel(authRepo) as T
            }

            modelClass.isAssignableFrom(LoginViewModel::class.java) -> {
                LoginViewModel(authRepo, userPref) as T
            }

            modelClass.isAssignableFrom(HomeViewModel::class.java) -> {
                HomeViewModel(storyRepo, userPref) as T
            }

            modelClass.isAssignableFrom(StoryViewModel::class.java) -> {
                StoryViewModel(storyRepo, userPref) as T
            }

            else -> throw IllegalArgumentException("Unknown ViewModel class: " + modelClass.name)
        }
    }

    companion object {
        @Volatile
        private var instance: ViewModelFactory? = null

        @JvmStatic
        fun getInstance(context: Context) =
            instance ?: synchronized(this) {
                instance ?: ViewModelFactory(
                    Injection.provideAuthRepository(),
                    Injection.provideStoryRepository(),
                    Injection.provideUserPreference(context.dataStore),
                    Injection.provideCommonRepository(context)
                )
            }.also { instance = it }
    }
}