package id.my.mrgsrylm.githubuser.view.activity

import android.content.Context
import android.os.Bundle
import android.view.MenuItem
import android.widget.CompoundButton
import androidx.appcompat.app.ActionBar
import androidx.appcompat.app.AppCompatActivity
import androidx.appcompat.app.AppCompatDelegate
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.preferencesDataStore
import androidx.lifecycle.ViewModelProvider
import id.my.mrgsrylm.githubuser.data.local.SettingPreferences
import id.my.mrgsrylm.githubuser.databinding.ActivityPreferencesBinding
import id.my.mrgsrylm.githubuser.viewmodel.PreferencesViewModelFactory
import id.my.mrgsrylm.githubuser.viewmodel.preference.PreferencesViewModel

private val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "theme")

class PreferencesActivity : AppCompatActivity() {

    private lateinit var binding: ActivityPreferencesBinding
    private lateinit var actionBar: ActionBar

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityPreferencesBinding.inflate(layoutInflater)
        setContentView(binding.root)

        actionBar = supportActionBar!!
        actionBar.title = "Setting"
        actionBar.setDisplayHomeAsUpEnabled(true)

        val pref = SettingPreferences.getInstance(dataStore)
        val prefViewModel = ViewModelProvider(
            this,
            PreferencesViewModelFactory(pref)
        )[PreferencesViewModel::class.java]

        prefViewModel.getThemeSetting().observe(this) { darkMode: Boolean ->
            if (darkMode) {
                AppCompatDelegate.setDefaultNightMode(AppCompatDelegate.MODE_NIGHT_YES)
                binding.switchTheme.isChecked = true
            } else {
                AppCompatDelegate.setDefaultNightMode(AppCompatDelegate.MODE_NIGHT_NO)
                binding.switchTheme.isChecked = false
            }
        }

        binding.switchTheme.setOnCheckedChangeListener { _: CompoundButton?, isChecked: Boolean ->
            prefViewModel.saveThemeSetting(isChecked)
        }
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        when (item.itemId) {
            android.R.id.home -> {
                finish()
                return true
            }
        }
        return super.onOptionsItemSelected(item)
    }
}