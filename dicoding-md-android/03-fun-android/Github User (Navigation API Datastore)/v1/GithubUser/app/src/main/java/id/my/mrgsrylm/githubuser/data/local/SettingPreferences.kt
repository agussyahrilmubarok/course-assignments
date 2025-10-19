package id.my.mrgsrylm.githubuser.data.local

import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.booleanPreferencesKey
import androidx.datastore.preferences.core.edit
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

class SettingPreferences private constructor(
    private val dataStore: DataStore<Preferences>
) {

    private val _themeKey = booleanPreferencesKey("theme")

    companion object {
        @Volatile
        private var INSTANCE: SettingPreferences? = null

        fun getInstance(dataStore: DataStore<Preferences>): SettingPreferences {
            return INSTANCE ?: synchronized(this) {
                val instance = SettingPreferences(dataStore)
                INSTANCE = instance
                instance
            }
        }
    }

    fun getTheme(): Flow<Boolean> {
        return dataStore.data.map { preferences ->
            preferences[_themeKey] ?: false
        }
    }

    suspend fun saveTheme(isDarkMode: Boolean) {
        dataStore.edit { preferences ->
            preferences[_themeKey] = isDarkMode
        }
    }
}