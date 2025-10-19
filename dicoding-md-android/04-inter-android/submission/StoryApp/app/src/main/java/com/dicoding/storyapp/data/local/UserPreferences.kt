package com.dicoding.storyapp.data.local

import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.booleanPreferencesKey
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import com.dicoding.storyapp.data.model.UserModel
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

class UserPreferences private constructor(
    private val dataStore: DataStore<Preferences>
) {
    suspend fun saveUser(user: UserModel) {
        dataStore.edit { pref ->
            pref[ID_KEY] = user.userId
            pref[NAME_KEY] = user.name
            pref[EMAIL_KEY] = user.email
            pref[PASSWORD_KEY] = user.password
            pref[TOKEN_KEY] = user.token
            pref[STATE_KEY] = user.isLogged
        }
    }

    fun getUser(): Flow<UserModel> {
        return dataStore.data.map { pref ->
            UserModel(
                pref[ID_KEY] ?: "",
                pref[NAME_KEY] ?: "",
                pref[EMAIL_KEY] ?: "",
                pref[PASSWORD_KEY] ?: "",
                pref[TOKEN_KEY] ?: "",
                pref[STATE_KEY] ?: false
            )
        }
    }

    fun getToken(): Flow<String> {
        return dataStore.data.map { pref ->
            pref[TOKEN_KEY] ?: ""
        }
    }

    suspend fun setLogout() {
        dataStore.edit { pref ->
            pref[ID_KEY] = ""
            pref[NAME_KEY] = ""
            pref[EMAIL_KEY] = ""
            pref[PASSWORD_KEY] = ""
            pref[TOKEN_KEY] = ""
            pref[STATE_KEY] = false
        }
    }

    companion object {
        @Volatile
        private var INSTANCE: UserPreferences? = null

        private val ID_KEY = stringPreferencesKey("user_id")
        private val NAME_KEY = stringPreferencesKey("user_name")
        private val EMAIL_KEY = stringPreferencesKey("user_email")
        private val PASSWORD_KEY = stringPreferencesKey("user_password")
        private val TOKEN_KEY = stringPreferencesKey("user_token")
        private val STATE_KEY = booleanPreferencesKey("user_state")

        fun getInstance(dataStore: DataStore<Preferences>): UserPreferences {
            return INSTANCE ?: synchronized(this) {
                val instance = UserPreferences(dataStore)
                INSTANCE = instance
                instance
            }
        }
    }
}