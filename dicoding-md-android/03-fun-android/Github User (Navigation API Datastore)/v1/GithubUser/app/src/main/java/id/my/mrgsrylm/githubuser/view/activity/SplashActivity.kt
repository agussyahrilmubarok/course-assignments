package id.my.mrgsrylm.githubuser.view.activity

import android.content.Intent
import android.os.Bundle
import android.os.Handler
import android.os.Looper
import androidx.appcompat.app.AppCompatActivity
import id.my.mrgsrylm.githubuser.R

class SplashActivity : AppCompatActivity() {

    private val splashDuration: Long = 3000

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_splash)

        supportActionBar?.hide()

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