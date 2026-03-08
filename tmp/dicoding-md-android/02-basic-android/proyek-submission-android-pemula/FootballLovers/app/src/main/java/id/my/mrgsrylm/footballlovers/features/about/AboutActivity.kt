package id.my.mrgsrylm.footballlovers.features.about

import android.os.Bundle
import android.widget.ImageView
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity
import com.bumptech.glide.Glide
import com.bumptech.glide.request.RequestOptions
import id.my.mrgsrylm.footballlovers.R
import id.my.mrgsrylm.footballlovers.data.dummy.usersDummy
import id.my.mrgsrylm.footballlovers.data.model.User

class AboutActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_about)

        val actionBar = supportActionBar
        actionBar!!.title = getString(R.string.about)

        val user: User = usersDummy[0]

        val tvName: TextView = findViewById(R.id.tv_set_user_name)
        val tvEmail: TextView = findViewById(R.id.tv_set_user_email)
        val tvDescription: TextView = findViewById(R.id.tv_set_user_description)
        val ivImage: ImageView = findViewById(R.id.iv_set_user_image_url)

        tvName.text = user.name
        tvEmail.text = user.email
        tvDescription.text = user.description
        Glide.with(this).load(user.imageUrl).apply(RequestOptions()).into(ivImage)
    }

    override fun onSupportNavigateUp(): Boolean {
        onBackPressed()
        return true
    }
}