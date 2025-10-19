package id.my.mrgsrylm.githubuser.viewmodel.preference

import androidx.lifecycle.LiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.asLiveData
import androidx.lifecycle.viewModelScope
import id.my.mrgsrylm.githubuser.data.local.SettingPreferences
import kotlinx.coroutines.launch

class PreferencesViewModel(
    private val pref: SettingPreferences
) : ViewModel() {
    fun getThemeSetting(): LiveData<Boolean> {
        return pref.getTheme().asLiveData()
    }

    fun saveThemeSetting(darkMode: Boolean) {
        viewModelScope.launch {
            pref.saveTheme(darkMode)
        }
    }
}