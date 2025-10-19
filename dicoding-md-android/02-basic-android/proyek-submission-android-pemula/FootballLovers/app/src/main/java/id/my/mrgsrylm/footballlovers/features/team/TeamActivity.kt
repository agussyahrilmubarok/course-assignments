package id.my.mrgsrylm.footballlovers.features.team

import android.os.Bundle
import android.widget.ImageView
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity
import com.bumptech.glide.Glide
import com.bumptech.glide.request.RequestOptions
import id.my.mrgsrylm.footballlovers.R

class TeamActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_team)

        val actionBar = supportActionBar
        actionBar!!.title = getString(R.string.team)
        actionBar.setDisplayHomeAsUpEnabled(true)

        val tvName: TextView = findViewById(R.id.tv_set_team_name)
        val ivImageUrl: ImageView = findViewById(R.id.iv_set_team_image_url)
        val tvLeague: TextView = findViewById(R.id.tv_set_team_league)
        val tvCountry: TextView = findViewById(R.id.tv_set_team_country)
        val tvFoundedYear: TextView = findViewById(R.id.tv_set_team_founded_year)
        val tvStadium: TextView = findViewById(R.id.tv_set_team_stadium)
        val tvCoach: TextView = findViewById(R.id.tv_set_team_coach)
        val tvDescription: TextView = findViewById(R.id.tv_set_team_description)

        val tName = intent.getStringExtra(EXTRA_NAME)
        val tImage = intent.getStringExtra(EXTRA_IMAGE_URL)
        val tDescription = intent.getStringExtra(EXTRA_DESCRIPTION)
        val tStadium = intent.getStringExtra(EXTRA_STADIUM)
        val tFoundedYear = intent.getStringExtra(EXTRA_FOUNDED_YEAR)
        val tCoach = intent.getStringExtra(EXTRA_COACH)
        val tCountry = intent.getStringExtra(EXTRA_COUNTRY)
        val tLeague = intent.getStringExtra(EXTRA_LEAGUE)

        Glide.with(this).load(tImage).apply(RequestOptions()).into(ivImageUrl)
        tvName.text = tName
        tvCountry.text = tCountry
        tvStadium.text = tStadium
        tvLeague.text = tLeague
        tvCoach.text = tCoach
        tvFoundedYear.text = tFoundedYear
        tvDescription.text = tDescription


    }

    override fun onSupportNavigateUp(): Boolean {
        onBackPressed()
        return true
    }

    companion object {
        const val EXTRA_NAME = "extra_name"
        const val EXTRA_IMAGE_URL = "extra_image_url"
        const val EXTRA_DESCRIPTION = "extra_description"
        const val EXTRA_STADIUM = "extra_stadium"
        const val EXTRA_FOUNDED_YEAR = "extra_founded_year"
        const val EXTRA_COACH = "extra_coach"
        const val EXTRA_COUNTRY = "extra_country"
        const val EXTRA_LEAGUE = "extra_league"
    }
}