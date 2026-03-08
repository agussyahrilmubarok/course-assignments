package id.my.mrgsrylm.githubuser.ui.activity

import android.content.Intent
import android.os.Bundle
import android.os.Handler
import android.os.Looper
import androidx.appcompat.app.AppCompatActivity
import androidx.appcompat.app.AppCompatDelegate
import androidx.lifecycle.ViewModelProvider
import id.my.mrgsrylm.githubuser.R
import id.my.mrgsrylm.githubuser.data.local.SettingPreferences
import id.my.mrgsrylm.githubuser.data.local.dataStore
import id.my.mrgsrylm.githubuser.databinding.ActivitySplashBinding
import id.my.mrgsrylm.githubuser.helper.PreferencesViewModelFactory
import id.my.mrgsrylm.githubuser.ui.viewmodel.MainViewModel

class SplashActivity : AppCompatActivity() {

    private val splashDuration: Long = 3000
    private lateinit var binding: ActivitySplashBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivitySplashBinding.inflate(layoutInflater)
        setContentView(binding.root)

        val pref = SettingPreferences.getInstance(application.dataStore)
        val mainViewModel = ViewModelProvider(this, PreferencesViewModelFactory(pref)).get(
            MainViewModel::class.java
        )

        mainViewModel.getThemeSettings().observe(this) { isDarkModeActive: Boolean ->
            if (isDarkModeActive) {
                binding.ivLogo.setImageResource(R.drawable.logo_github_white)
                AppCompatDelegate.setDefaultNightMode(AppCompatDelegate.MODE_NIGHT_YES)
            } else {
                binding.ivLogo.setImageResource(R.drawable.logo_github)
                AppCompatDelegate.setDefaultNightMode(AppCompatDelegate.MODE_NIGHT_NO)
            }
        }

        Handler(Looper.getMainLooper()).postDelayed(splashRunnable, splashDuration)
    }

    override fun onBackPressed() {
        Handler(Looper.getMainLooper()).removeCallbacks(splashRunnable)
        super.onBackPressed()
    }

    private val splashRunnable = Runnable {
        if (!isFinishing) {
            val intent = Intent(this@SplashActivity, MainActivity::class.java)
            startActivity(intent)
            finish()
        }
    }
}